package v8

import (
	"context"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/encoding/json"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

// Version is an Exchange upgrade to ensure kucoin websocket subscriptions are not missing for assets
// This handles users who upgraded only once after #1394 and before #1579
type Version struct{}

// Exchanges returns just kucoin
func (v *Version) Exchanges() []string { return []string{"kucoin"} }

/*
UpgradeExchange adds assets to Kucoin subscriptions, attempting to preserve any other changes where possible
  - a single allTrades assetless subscription is replaced with separate spot and margin subscritions
  - /contractMarket/level2Depth50 is replaced by subsctiption.Orderbook for asset.All
  - /contractMarket/tickerV2 is replaced by subscription.Ticker for asset.All
  - /margin/fundingBook is deprecated and removed
*/
func (v *Version) UpgradeExchange(_ context.Context, e []byte) ([]byte, error) {
	var errs error
	newSubs := [][]byte{}
	fn := func(subStr []byte, _ jsonparser.ValueType, _ int, _ error) {
		var s Subscription
		if err := json.Unmarshal(subStr, &s); err != nil {
			errs = common.AppendError(errs, err)
			return
		}
		if s.Asset != asset.Empty {
			newSubs = append(newSubs, subStr)
			return
		}
		s.Channel = strings.TrimSuffix(s.Channel, ":%s")
		switch s.Channel {
		case "ticker", "orderbook":
			s.Asset = asset.All
		case "allTrades":
			newSubs = append(newSubs,
				[]byte(`{"enabled":true,asset:"spot","channel":"allTrades"}`),
				[]byte(`{"enabled":true,asset:"margin","channel":"allTrades"}`),
			)
			return
		case "/contractMarket/tradeOrders", "/contractMarket/advancedOrders", "/contractAccount/wallet":
			s.Asset = asset.Futures
		case "/margin/position", "/margin/loan":
			s.Asset = asset.Margin
		}
		if newSub, err := json.Marshal(s); err != nil {
			errs = common.AppendError(errs, err)
		} else {
			newSubs = append(newSubs, newSub)
		}
	}

	_, err := jsonparser.ArrayEach(e, fn, "features", "subscriptions")
	return e, err
}

// DowngradeExchange performs no-operation, since v7 is ensuring an upgrade which was already in place
func (v *Version) DowngradeExchange(_ context.Context, e []byte) ([]byte, error) {
	return e, nil
}
