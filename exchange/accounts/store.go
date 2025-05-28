package accounts

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/dispatch"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

// store holds ticker information for each individual exchange
type store struct {
	exchangeAccounts exchangeMap
	mu               sync.Mutex
	mux              *dispatch.Mux
}

type exchangeMap map[exchange]*Accounts

type exchange interface {
	GetName() string
}

var global atomic.Pointer[store]

// NewStore returns a new store with the default global dispatcher mux
func NewStore() *store {
	return &store{
		exchangeAccounts: make(exchangeMap),
		mux:              dispatch.GetNewMux(nil),
	}
}

// GetStore returns the singleton accounts store for global use; Initialising if necessary
func GetStore() *store {
	if s := global.Load(); s != nil {
		return s
	}
	_ = global.CompareAndSwap(nil, NewStore())
	return global.Load()
}

// CollectBalances converts a map of sub-account balances into a slice
func CollectBalances(accountBalances map[string][]Balance, assetType asset.Item) (accounts []SubAccount, err error) {
	if err := common.NilGuard(accountBalances); err != nil {
		return nil, err
	}

	if !assetType.IsValid() {
		return nil, fmt.Errorf("%s, %w", assetType, asset.ErrNotSupported)
	}

	accounts = make([]SubAccount, 0, len(accountBalances))
	for accountID, balances := range accountBalances {
		accounts = append(accounts, SubAccount{
			ID:         accountID,
			AssetType:  assetType,
			Currencies: balances,
		})
	}
	return
}

// GetExchangeAccounts returns accounts for a specific exchange
func (s *store) GetExchangeAccounts(e exchange) (a *Accounts, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.exchangeAccounts[e]
	if !ok {
		if a, err = s.registerExchange(e); err != nil {
			return nil, fmt.Errorf("error subscribing to `%s` exchange account: %w", e, err)
		}
	}
	return a, nil
}

// registerExchange adds a new empty shared account accounts entry for an exchange
// must be called with s.mu locked
func (s *store) registerExchange(e exchange) (*Accounts, error) {
	if _, ok := s.exchangeAccounts[e]; ok {
		return nil, errExchangeAlreadyExists
	}
	a, err := NewAccounts(e, s.mux)
	s.exchangeAccounts[e] = a
	return a, err
}
