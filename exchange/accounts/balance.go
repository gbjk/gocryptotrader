package accounts

import (
	"errors"
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

var errBalanceCurrencyMismatch = errors.New("balance currency does not match update currency")

// Balance contains an exchange currency balance
type Balance struct {
	Currency               currency.Code
	Total                  float64
	Hold                   float64
	Free                   float64
	AvailableWithoutBorrow float64
	Borrowed               float64
	UpdatedAt              time.Time
}

// Change defines incoming balance change on currency holdings
type Change struct {
	Account   string
	AssetType asset.Item
	Balance   Balance
}

// balance contains a balance with live updates
type balance struct {
	internal Balance
	m        sync.RWMutex
}

// currencyBalances provides a map of currencies to balances
type currencyBalances map[currency.Code]*balance

// CurrencyBalances provides a map of currencies to balances
type CurrencyBalances map[currency.Code]Balance

// Add will Set a currency balance, overwriting any previous Balance
// currency may be a string or a currency.Code
func (c *CurrencyBalances) Set(cAny any, b Balance) {
	var curr currency.Code
	switch v := cAny.(type) {
	case string:
		curr = currency.NewCode(v)
	case currency.Code:
		curr = v
	}
	b.Currency = curr
	(*c)[curr] = b
}

// Add will Add to a currency balance
func (c *CurrencyBalances) Add(cAny any, b Balance) {
	var curr currency.Code
	switch v := cAny.(type) {
	case string:
		curr = currency.NewCode(v)
	case currency.Code:
		curr = v
	}
	if e, ok := (*c)[curr]; !ok {
		b.Currency = curr
		(*c)[curr] = b
	} else {
		(*c)[curr] = e.Add(b)
	}
}

// Balance returns a snapshot copy of the Balance
func (l *balance) Balance() Balance {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.internal
}

// Add returns a new Balance adding together a and b
// UpdatedAt is the later of the two Balances
func (b Balance) Add(a Balance) Balance {
	b.Total += a.Total
	b.Hold += a.Hold
	b.Free += a.Free
	b.AvailableWithoutBorrow += a.AvailableWithoutBorrow
	b.Borrowed += a.Borrowed
	if a.UpdatedAt.After(b.UpdatedAt) {
		b.UpdatedAt = a.UpdatedAt
	}
	return b
}

func (c *currencyBalances) Public() CurrencyBalances {
	n := make(CurrencyBalances, len(*c))
	for curr, bal := range *c {
		n[curr] = bal.Balance()
	}
	return n
}

// update checks that an incoming change has a valid change, and returns if the balances were changed
// If change does not have a Currency set, the existing Currency is preserved
func (b *balance) update(change Balance) (bool, error) {
	if err := common.NilGuard(b, change); err != nil {
		return false, err
	}
	if change.UpdatedAt.IsZero() {
		return false, errUpdatedAtIsZero
	}
	b.m.Lock()
	defer b.m.Unlock()
	if b.internal.Currency != currency.EMPTYCODE {
		switch change.Currency {
		case b.internal.Currency:
			// All good
		case currency.EMPTYCODE:
			change.Currency = b.internal.Currency
		default:
			return false, errBalanceCurrencyMismatch
		}
	}
	if !b.internal.UpdatedAt.Before(change.UpdatedAt) {
		return false, errOutOfSequence
	}
	b.internal.UpdatedAt = change.UpdatedAt // Set just the time, and then can compare easily
	if b.internal == change {
		return false, nil
	}
	b.internal = change
	return true, nil
}

// balance rutens a balance for a currency
func (s currencyBalances) balance(c currency.Code) *balance {
	if _, ok := s[c]; !ok {
		s[c] = &balance{internal: Balance{Currency: c}}
	}
	return s[c]
}
