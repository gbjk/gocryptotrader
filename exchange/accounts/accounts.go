package accounts

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/key"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/dispatch"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

// Public errors
var (
	ErrBalancesNotFound = errors.New("balances not found")
)

var (
	errExchangeNameUnset            = errors.New("exchange name unset")
	errNoExchangeSubAccountBalances = errors.New("no exchange sub account balances")
	errCredentialsEmpty             = errors.New("credentials are nil")
	errOutOfSequence                = errors.New("out of sequence")
	errUpdatedAtIsZero              = errors.New("updatedAt may not be zero")
	errLoadingBalance               = errors.New("error loading balance")
	errExchangeAlreadyExists        = errors.New("exchange already exists")
	errCannotUpdateBalance          = errors.New("cannot update balance")
)

// Accounts holds a stream ID and a map to the exchange holdings
type Accounts struct {
	Exchange    exchange
	routingID   uuid.UUID // GCT internal routing mux id
	subAccounts credSubAccounts
	mu          sync.RWMutex
	mux         *dispatch.Mux
}

type (
	credSubAccounts map[Credentials]subAccounts
	subAccounts     map[key.SubAccountAsset]currencyBalances
)

// SubAccount contains a account for an asset type and its balances
// The subAccount may be the main account depending on exchange structure
type SubAccount struct {
	ID        string
	AssetType asset.Item
	Balances  CurrencyBalances
}

// SubAccounts contains a list of public SubAccounts
type SubAccounts []*SubAccount

// MustNewAccounts returns an initialized Accounts store for use in isolation from a global exchange accounts store
// mux is set to dispatch.GetNewMux(nil)
// Any errors in mux ID generation will panic, so users should balance risk vs utility accordingly depending on use-case
func MustNewAccounts(e exchange) *Accounts {
	a, err := NewAccounts(e, dispatch.GetNewMux(nil))
	if err != nil {
		panic(err)
	}
	return a
}

// NewAccounts returns an initialized Accounts store for use in isolation from a global exchange accounts store
func NewAccounts(e exchange, mux *dispatch.Mux) (*Accounts, error) {
	if err := common.NilGuard(e); err != nil {
		return nil, err
	}
	id, err := mux.GetID()
	if err != nil {
		return nil, err
	}
	return &Accounts{
		Exchange:    e,
		subAccounts: make(credSubAccounts),
		routingID:   id,
		mux:         mux,
	}, nil
}

// Subscribe subscribes to your exchange accounts
func (a *Accounts) Subscribe() (dispatch.Pipe, error) {
	return a.mux.Subscribe(a.routingID)
}

// CurrencyBalances returns the balances for the Accounts grouped by currency
// If creds is nil, all credential SubAccounts will be collated
// If assetType is asset.All, all assets will be collated
func (a *Accounts) CurrencyBalances(creds *Credentials, assetType asset.Item) (CurrencyBalances, error) {
	if err := common.NilGuard(a); err != nil {
		return nil, err
	}
	if !assetType.IsValid() && assetType != asset.All {
		return nil, fmt.Errorf("%s %s %w", a.Exchange.GetName(), assetType, asset.ErrNotSupported)
	}
	currs := CurrencyBalances{}
	for credsKey, subAccountsForCreds := range a.subAccounts {
		if !creds.IsEmpty() && *creds != credsKey {
			continue
		}
		for subAcctKey, balances := range subAccountsForCreds {
			if assetType != asset.All && assetType != subAcctKey.Asset {
				continue
			}
			for curr, bal := range balances {
				currs.Add(curr, bal.Balance())
			}
		}
	}
	if len(currs) == 0 {
		return nil, fmt.Errorf("%w for %s credentials %s asset %s", ErrBalancesNotFound, a.Exchange.GetName(), creds, assetType)
	}
	return currs, nil
}

// Balances returns the public SubAccounts and their balances
// If creds is nil, all credential SubAccounts will be returned
// If assetType is asset.All, all assets will be returned
func (a *Accounts) Balances(creds *Credentials, assetType asset.Item) (SubAccounts, error) {
	if err := common.NilGuard(a); err != nil {
		return nil, err
	}

	if !assetType.IsValid() && assetType != asset.All {
		return nil, fmt.Errorf("%s %s %w", a.Exchange.GetName(), assetType, asset.ErrNotSupported)
	}

	var subAccts SubAccounts
	for credsKey, subAccountsForCreds := range a.subAccounts {
		if !creds.IsEmpty() && *creds != credsKey {
			continue
		}
		for subAcctKey, balances := range subAccountsForCreds {
			if assetType != asset.All && assetType != subAcctKey.Asset {
				continue
			}
			subAccts = append(subAccts, &SubAccount{
				ID:        subAcctKey.SubAccount,
				AssetType: subAcctKey.Asset,
				Balances:  balances.Public(),
			})
		}
	}

	if len(subAccts) == 0 {
		return nil, fmt.Errorf("%w for %s credentials %s asset %s", ErrBalancesNotFound, a.Exchange.GetName(), creds, assetType)
	}
	return subAccts, nil
}

// GetBalance returns a copy of the balance for that asset item
func (a *Accounts) GetBalance(subAccount string, creds *Credentials, aType asset.Item, c currency.Code) (Balance, error) {
	if !aType.IsValid() {
		return Balance{}, fmt.Errorf("cannot get balance: %w: %q", asset.ErrNotSupported, aType)
	}

	if creds.IsEmpty() {
		return Balance{}, fmt.Errorf("cannot get balance: %w", errCredentialsEmpty)
	}

	if c.IsEmpty() {
		return Balance{}, fmt.Errorf("cannot get balance: %w", currency.ErrCurrencyCodeEmpty)
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	subAccounts, ok := a.subAccounts[*creds]
	if !ok {
		return Balance{}, fmt.Errorf("%w for %s", ErrBalancesNotFound, creds)
	}

	assets, ok := subAccounts[key.SubAccountAsset{
		SubAccount: subAccount,
		Asset:      aType,
	}]
	if !ok {
		return Balance{}, fmt.Errorf("%w for %s SubAccount %q %s %s", errNoExchangeSubAccountBalances, a.Exchange.GetName(), subAccount, aType, c)
	}
	b, ok := assets[c]
	if !ok {
		return Balance{}, fmt.Errorf("%w for %s SubAccount %q %s %s", errNoExchangeSubAccountBalances, a.Exchange.GetName(), subAccount, aType, c)
	}
	return b.Balance(), nil
}

// Save saves the holdings with a new snapshot of account balances; Any missing currencies will be removed
// Each SubAccount change is Published if it changes the balance for that currency
func (a *Accounts) Save(subAccts SubAccounts, creds *Credentials) error {
	if err := common.NilGuard(a); err != nil {
		return fmt.Errorf("cannot save holdings: %w", err)
	}
	if err := common.NilGuard(a.subAccounts); err != nil {
		return fmt.Errorf("cannot save holdings: %w", err)
	}

	if creds.IsEmpty() {
		return fmt.Errorf("cannot save holdings: %w", errCredentialsEmpty)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	var errs error
	for _, s := range subAccts {
		if !s.AssetType.IsValid() {
			errs = common.AppendError(errs, fmt.Errorf("error loading %s[%s] SubAccount holdings: %w", s.ID, s.AssetType, asset.ErrNotSupported))
			continue
		}

		accBalances := a.currencyBalances(creds, s.ID, s.AssetType)

		updated := false
		missing := maps.Clone(accBalances)
		for curr, newBal := range s.Balances {
			delete(missing, curr)
			if newBal.UpdatedAt.IsZero() {
				newBal.UpdatedAt = time.Now()
			}
			if newBal.Currency == currency.EMPTYCODE {
				newBal.Currency = curr
			}
			b := accBalances.balance(curr)
			if u, err := b.update(newBal); err != nil {
				errs = common.AppendError(errs, fmt.Errorf("%w for account ID `%s` [%s %s]: %w", errLoadingBalance, s.ID, s.AssetType, curr, err))
			} else if u {
				updated = true
			}
		}
		for cur := range missing {
			delete(accBalances, cur)
			updated = true
		}
		if updated {
			if err := a.mux.Publish(s, a.routingID); err != nil {
				errs = common.AppendError(errs, fmt.Errorf("cannot publish load for %s %w", a.Exchange, err))
			}
		}
	}

	return errs
}

// Update updates the account balances
func (a *Accounts) Update(changes []Change, creds *Credentials) error {
	if err := common.NilGuard(a); err != nil {
		return fmt.Errorf("cannot save holdings: %w", err)
	}
	if err := common.NilGuard(a.subAccounts); err != nil {
		return fmt.Errorf("cannot save holdings: %w", err)
	}
	if creds.IsEmpty() {
		return fmt.Errorf("%w: %w", errCannotUpdateBalance, errCredentialsEmpty)
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	var errs error
	for _, change := range changes {
		if !change.AssetType.IsValid() {
			errs = common.AppendError(errs, fmt.Errorf("%w for %s.%s %w", errCannotUpdateBalance, change.Account, change.AssetType, asset.ErrNotSupported))
			continue
		}
		if err := common.NilGuard(change.Balance); err != nil {
			errs = common.AppendError(errs, fmt.Errorf("%w for %s.%s %w", errCannotUpdateBalance, change.Account, change.AssetType, err))
			continue
		}

		accBalances := a.currencyBalances(creds, change.Account, change.AssetType)
		if u, err := accBalances.balance(change.Balance.Currency).update(change.Balance); err != nil {
			errs = common.AppendError(errs, fmt.Errorf("%w for %s.%s.%s %w", errCannotUpdateBalance, change.Account, change.AssetType, change.Balance.Currency, err))
		} else if u {
			if err := a.mux.Publish(change, a.routingID); err != nil {
				errs = common.AppendError(errs, fmt.Errorf("cannot publish update balance for %s: %w", a.Exchange, err))
			}
		}
	}
	return errs
}

// NewSUbAccount returns a new sub account
// id may be empty
func NewSubAccount(a asset.Item, id string) *SubAccount {
	return &SubAccount{
		AssetType: a,
		ID:        id,
		Balances:  CurrencyBalances{},
	}
}

// Merge adds CurrencyBalances in s to the SubAccount in l with a matching AssetType and ID
// If no SubAccount matches, s is appended
// Duplicate Currency Balances are added together
func (l SubAccounts) Merge(s *SubAccount) SubAccounts {
	i := slices.IndexFunc(l, func(b *SubAccount) bool { return s.AssetType == b.AssetType && s.ID == b.ID })
	if i == -1 {
		return append(l, s)
	}
	for curr, newBal := range s.Balances {
		l[i].Balances[curr] = newBal.Add(l[i].Balances[curr])
	}
	return l
}

func (a *Accounts) currencyBalances(c *Credentials, subAcct string, aType asset.Item) currencyBalances {
	k := key.SubAccountAsset{SubAccount: subAcct, Asset: aType}
	if _, ok := a.subAccounts[*c]; !ok {
		a.subAccounts[*c] = make(subAccounts)
	}
	if _, ok := a.subAccounts[*c][k]; !ok {
		a.subAccounts[*c][k] = make(currencyBalances)
	}
	return a.subAccounts[*c][k]
}

func (s currencyBalances) balance(c currency.Code) *balance {
	if _, ok := s[c]; !ok {
		s[c] = &balance{}
	}
	return s[c]
}
