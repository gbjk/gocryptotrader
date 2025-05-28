package accounts

import (
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

// Change defines incoming balance change on currency holdings
type Change struct {
	Account   string
	AssetType asset.Item
	Balance   *BalanceUpdate
}
