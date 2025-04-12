package v8_test

import (
	"testing"
)

func TestUpgradeConfig(t *testing.T) {
	t.Parallel()
	/*
		ku := &Kucoin{ //nolint:govet // Intentional shadow to avoid future copy/paste mistakes
			Base: exchange.Base{
				Config: &config.Exchange{
					Features: &config.FeaturesConfig{
						Subscriptions: subscription.List{
							{Enabled: true, Channel: "ticker"},
							{Enabled: true, Channel: "allTrades"},
							{Enabled: true, Channel: "orderbook", Interval: kline.HundredMilliseconds},
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
					},
				},
				Features: exchange.Features{},
			},
		}

		ku.checkSubscriptions()
		testsubs.EqualLists(t, defaultSubscriptions, ku.Features.Subscriptions)
		testsubs.EqualLists(t, defaultSubscriptions, ku.Config.Features.Subscriptions)
	*/
}

func TestDowngradeConfig(t *testing.T) {
	t.Parallel()
}
