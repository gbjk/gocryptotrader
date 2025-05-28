package accounts

import (
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/alert"
)

// Balance contains an exchange currency balance
type Balance struct {
	Currency               currency.Code
	Total                  float64
	Hold                   float64
	Free                   float64
	AvailableWithoutBorrow float64
	Borrowed               float64
}

// LiveBalance contains a balance with live updates
type LiveBalance struct {
	Currency  currency.Code
	b         Balance
	updatedAt time.Time
	m         sync.RWMutex
	notifier  alert.Notice
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
