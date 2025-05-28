package accounts

import (
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/alert"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

// ProtectedBalance stores the full balance information for that specific asset
type ProtectedBalance struct {
	total                  float64
	hold                   float64
	free                   float64
	availableWithoutBorrow float64
	borrowed               float64
	m                      sync.Mutex
	updatedAt              time.Time

	// notice alerts for when the balance changes for strategy inspection and usage
	notice alert.Notice
}

// SubAccount defines a singular account type with associated currency balances
type SubAccount struct {
	Credentials Protected
	ID          string
	AssetType   asset.Item
	Currencies  []Balance
}

// Balance is a sub-type to store currency name and individual totals
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

// Protected limits the access to the underlying credentials outside of this package
type Protected struct {
	creds Credentials
}
