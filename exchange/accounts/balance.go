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
