package v8

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/thrasher-corp/gocryptotrader/common"
	subutils "github.com/thrasher-corp/gocryptotrader/config/versions/utils/subscriptions"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
)

// Version is an Exchange upgrade to ensure kucoin websocket subscriptions are not missing for assets
// This handles users who upgraded only once after #1394 and before #1579
type Version struct{}

// Exchanges returns just kucoin
func (v *Version) Exchanges() []string { return []string{"Kucoin"} }

/*
UpgradeExchange adds assets to Kucoin's subscriptions, attempting to preserve any other changes where possible
*/
func (v *Version) UpgradeExchange(_ context.Context, e []byte) ([]byte, error) {
	return subutils.Upgrade(
		e,
		subscription.List{
			{Enabled: true, Channel: subscription.TickerChannel},                                         // marketTickerChannel
			{Enabled: true, Channel: subscription.AllTradesChannel},                                      // marketMatchChannel
			{Enabled: true, Channel: subscription.OrderbookChannel, Interval: kline.HundredMilliseconds}, // marketOrderbookLevel2Channels
			{Enabled: true, Channel: "/contractMarket/tickerV2:%s"},
			{Enabled: true, Channel: "/contractMarket/level2Depth50:%s"},
			{Enabled: true, Channel: "/margin/fundingBook:%s", Authenticated: true},
			{Enabled: true, Channel: "/account/balance", Authenticated: true},
			{Enabled: true, Channel: "/margin/position", Authenticated: true},
			{Enabled: true, Channel: "/margin/loan:%s", Authenticated: true},
			{Enabled: true, Channel: "/contractMarket/tradeOrders", Authenticated: true},
			{Enabled: true, Channel: "/contractMarket/advancedOrders", Authenticated: true},
			{Enabled: true, Channel: "/contractAccount/wallet", Authenticated: true},
		},
		subscription.List{
			{Enabled: true, Asset: asset.All, Channel: subscription.TickerChannel},
			{Enabled: true, Asset: asset.All, Channel: subscription.OrderbookChannel, Interval: kline.HundredMilliseconds},
			{Enabled: true, Asset: asset.Spot, Channel: subscription.AllTradesChannel},
			{Enabled: true, Asset: asset.Margin, Channel: subscription.AllTradesChannel},
			{Enabled: true, Asset: asset.Futures, Channel: "/contractMarket/tradeOrders", Authenticated: true},
			{Enabled: true, Asset: asset.Futures, Channel: "/contractMarket/advancedOrders", Authenticated: true},
			{Enabled: true, Asset: asset.Futures, Channel: "/contractAccount/wallet", Authenticated: true},
			{Enabled: true, Asset: asset.Margin, Channel: "/margin/position", Authenticated: true},
			{Enabled: true, Asset: asset.Margin, Channel: "/margin/loan", Authenticated: true},
			{Enabled: true, Channel: "/account/balance", Authenticated: true},
		}, mutateSubs, `
* %s placeholders are removed from channel names
* a single allTrades assetless subscription is replaced with separate spot and margin subscriptions
* /contractMarket/level2Depth50 is replaced by subsctiption.Orderbook for asset.All
* /contractMarket/tickerV2 is replaced by subscription.Ticker for asset.All
* /margin/fundingBook is deprecated and removed`,
	)
}

func mutateSubs(subs subscription.List) (subscription.List, error) {
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
		case "/contractMarket/level2Depth50", "/contractMarket/tickerV2", "/margin/fundingBook":
			return
		}
		if newSub, err := json.Marshal(s); err != nil {
			errs = common.AppendError(errs, err)
		} else {
			newSubs = append(newSubs, newSub)
		}
	}

	return subs, errs
}

// DowngradeExchange performs no-operation, since v7 is ensuring an upgrade which was already in place
func (v *Version) DowngradeExchange(_ context.Context, e []byte) ([]byte, error) {
	return e, nil
}
