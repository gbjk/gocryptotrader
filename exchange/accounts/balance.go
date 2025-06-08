package accounts

import (
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
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

// update checks that an incoming change has a valid change, and returns if the balances were changed
func (b *balance) update(change Balance) (bool, error) {
	if err := common.NilGuard(b, change); err != nil {
		return false, err
	}
	if change.UpdatedAt.IsZero() {
		return false, errUpdatedAtIsZero
	}
	b.m.Lock()
	defer b.m.Unlock()
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

// Wait waits for a change in amounts for an asset type. This will pause
// indefinitely if no change ever occurs. Max wait will return true if it failed
// to achieve a state change in the time specified. If Max wait is not specified
// it will default to a minute wait time.
func (b *balance) Wait(maxWait time.Duration) (wait <-chan bool, cancel chan<- struct{}, err error) {
	if err := common.NilGuard(b); err != nil {
		return nil, nil, err
	}

	if maxWait <= 0 {
		maxWait = time.Minute
	}
	ch := make(chan struct{})
	go func(ch chan<- struct{}, until time.Duration) {
		time.Sleep(until)
		close(ch)
	}(ch, maxWait)

	return b.notice.Wait(ch), ch, nil
}

// GetFree returns the current free balance for the exchange
func (b *ProtectedBalance) GetFree() float64 {
	if b == nil {
		return 0
	}
	b.m.Lock()
	defer b.m.Unlock()
	return b.free
}

func (b *ProtectedBalance) reset() {
	b.m.Lock()
	defer b.m.Unlock()

	b.total = 0
	b.hold = 0
	b.free = 0
	b.availableWithoutBorrow = 0
	b.borrowed = 0
	b.updatedAt = time.Now()
	b.notice.Alert()
}
