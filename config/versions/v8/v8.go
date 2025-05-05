package v8

import (
	"context"
	"encoding/json"
	"slices"
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

// UpgradeExchange adds assets to Kucoin's subscriptions; see summary below
func (v *Version) UpgradeExchange(_ context.Context, e []byte) ([]byte, error) {
	const summary = `
* %s placeholders are removed from channel names
* a single allTrades assetless subscription is replaced with separate spot and margin subscriptions
* /contractMarket/level2Depth50 is replaced by subsctiption.Orderbook for asset.All
* /contractMarket/tickerV2 is replaced by subscription.Ticker for asset.All
* /margin/fundingBook is deprecated and removed`

	oldDefaults := subscription.List{
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
	}

	newDefaults := subscription.List{
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
	}

	return subutils.Upgrade(e, oldDefaults, newDefaults, mutateSubs, summary)
}

func mutateSubs(subs subscription.List) (newSubs subscription.List, err error) {
	for _, s := range subs {
		if s.Asset == asset.Empty {
		s.Channel = strings.TrimSuffix(s.Channel, ":%s")
		switch s.Channel {
		case subscription.TickerChannel, subscription.OrderbookChannel:
			s.Asset = asset.All
		case subscription.AllTradesChannel:
			spotSub, marginSub := s.Clone(), s
			spotSub.Asset, marginSub.Asset = asset.Spot, asset.Margin
			newSubs = append(newSubs, spotSub, marginSub)
			continue
		case "/contractMarket/tradeOrders", "/contractMarket/advancedOrders", "/contractAccount/wallet":
			s.Asset = asset.Futures
		case "/margin/position", "/margin/loan":
			s.Asset = asset.Margin
		case "/contractMarket/level2Depth50", // Replaced by subsctiption.Orderbook for asset.All
			"/contractMarket/tickerV2", // Replaced by subscription.Ticker for asset.All
			"/margin/fundingBook":      // Deprecated and removed
			continue
		}
		newSubs = append(newSubs, s)
	}
	return newSubs, nil
}

// DowngradeExchange performs no-operation, since v7 is ensuring an upgrade which was already in place
func (v *Version) DowngradeExchange(_ context.Context, e []byte) ([]byte, error) {
	return e, nil
}
