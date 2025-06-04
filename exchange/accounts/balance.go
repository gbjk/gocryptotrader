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

// LiveBalance contains a balance with live updates
type LiveBalance struct {
	b        Balance
	m        sync.RWMutex
	notifier alert.Notice
}

// balance returns a snapshot copy of the Balance
// Does not enjoy protection from locking
func (l *LiveBalance) Balance() Balance {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.b
}

// Total returns the total in this balance
func (l *LiveBalance) Total() float64 {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.b.Total
}

// Hold returns the amount on hold in this balance
func (l *LiveBalance) Hold() float64 {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.b.Hold
}

// Free returns the free amount in this balance
func (l *LiveBalance) Free() float64 {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.b.Free
}

// AvailableWithoutBorrow returns the available without borrowing  in this balance
func (l *LiveBalance) AvailableWithoutBorrow() float64 {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.b.AvailableWithoutBorrow
}

// Borrowed returns the amount borrowed in this balance
func (l *LiveBalance) Borrowed() float64 {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.b.Borrowed
}

// UpdatedAt returns the time this balance was last updated
func (l *LiveBalance) UpdatedAt() time.Time {
	l.m.RLock()
	defer l.m.RUnlock()
	return l.b.UpdatedAt
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
