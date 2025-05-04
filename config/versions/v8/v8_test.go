package v8_test

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/config"
	v8 "github.com/thrasher-corp/gocryptotrader/config/versions/v8"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
)

func TestUpgradeConfig(t *testing.T) {
	t.Parallel()

	/*
		- Things with assets are left alone
		- Orderbook and Ticker get all asset added
		- If allTrades exists we add  spot and margin channels
						[]byte(`{"enabled":true,asset:"spot","channel":"allTrades"}`),
						[]byte(`{"enabled":true,asset:"margin","channel":"allTrades"}`),
		- case "/contractMarket/tradeOrders", "/contractMarket/advancedOrders", "/contractAccount/wallet":
		gets futures
		- case "/margin/position", "/margin/loan":
		gets margin

	*/

	user := config.Exchange{
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
	}
	userJSON, err := json.Marshal(user)
	t.Log(string(userJSON))
	require.NoError(t, err, "Must not error Marshalling input")
	v := new(v8.Version)
	outConfig, err := v.UpgradeExchange(t.Context(), userJSON)
	require.NoError(t, err)
	out := &config.Exchange{}
	err = json.Unmarshal(outConfig, out)
	require.NoError(t, err, "Must not error unmarshalling result")
	spew.Dump(out.Features.Subscriptions)
}

func TestDowngradeConfig(t *testing.T) {
	t.Parallel()
}
