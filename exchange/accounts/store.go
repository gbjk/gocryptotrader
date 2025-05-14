package accounts

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/thrasher-corp/gocryptotrader/dispatch"
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

type exchangeWrapper interface {
	GetBase() exchange
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

// GetExchangeAccounts returns accounts for a specific exchange
func (s *store) GetExchangeAccounts(e exchange) (a *Accounts, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if w, ok := e.(exchangeWrapper); ok {
		// Because SetupDefualts is called on Base, it's easiest to just use the Base pointer as the key
		e = w.GetBase()
	}
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
