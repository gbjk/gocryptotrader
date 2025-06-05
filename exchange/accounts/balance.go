package accounts

import (
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/alert"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

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
	Balance   *Balance
}

// balance contains a balance with live updates
type balance struct {
	internal Balance
	m        sync.RWMutex
	notifier alert.Notice
}

// balance returns a snapshot copy of the Balance
// Does not enjoy protection from locking
func (l *balance) Balance() Balance {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.internal
}

// Add returns a new Balance adding together a and b
// UpdatedAt is the later of the two Balances
func (b Balance) Add(a Balance) Balance {
	c := Balance{
		Total:                  b.Total + a.Total,
		Hold:                   b.Hold + a.Hold,
		Free:                   b.Free + a.Free,
		AvailableWithoutBorrow: b.AvailableWithoutBorrow + a.AvailableWithoutBorrow,
		Borrowed:               b.Borrowed + a.Borrowed,
		UpdatedAt:              b.UpdatedAt,
	}
	if a.UpdatedAt.After(b.UpdatedAt) {
		c.UpdatedAt = a.UpdatedAt
	}
	return c
}
