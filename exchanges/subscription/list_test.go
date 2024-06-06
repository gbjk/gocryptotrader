package subscription

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
)

// TestListStrings exercises List.Strings()
func TestListStrings(t *testing.T) {
	t.Parallel()
	l := List{
		&Subscription{
			Channel: TickerChannel,
			Asset:   asset.Spot,
			Pairs:   currency.Pairs{ethusdcPair, btcusdtPair},
		},
		&Subscription{
			Channel: OrderbookChannel,
			Pairs:   currency.Pairs{ethusdcPair},
		},
	}
	exp := []string{"orderbook  ETH/USDC", "ticker spot ETH/USDC,BTC/USDT"}
	assert.ElementsMatch(t, exp, l.Strings(), "String must return correct sorted list")
}

// TestListGroupPairs exercises List.GroupPairs()
func TestListGroupPairs(t *testing.T) {
	t.Parallel()
	l := List{
		{Asset: asset.Spot, Channel: TickerChannel, Pairs: currency.Pairs{ethusdcPair, btcusdtPair}},
	}
	for _, c := range []string{TickerChannel, OrderbookChannel} {
		for _, p := range []currency.Pair{ethusdcPair, btcusdtPair} {
			l = append(l, &Subscription{
				Channel: c,
				Asset:   asset.Spot,
				Pairs:   currency.Pairs{p},
			})
		}
	}
	n := l.GroupPairs()
	assert.Len(t, l, 5, "Orig list should not be changed")
	assert.Len(t, n, 2, "New list should be grouped")
	exp := []string{"ticker spot ETH/USDC,BTC/USDT", "orderbook spot ETH/USDC,BTC/USDT"}
	assert.ElementsMatch(t, exp, n.Strings(), "String must return correct sorted list")
}

type mockEx struct{}

func (m *mockEx) GetAssetTypes(e bool) asset.Items { return asset.Items{asset.Spot, asset.Futures} }

func (m *mockEx) GetEnabledPairs(a asset.Item) (currency.Pairs, error) {
	return currency.Pairs{btcusdtPair, ethusdcPair}, nil
}

func (m *mockEx) GetPairFormat(a asset.Item, r bool) (currency.PairFormat, error) {
	return currency.PairFormat{Uppercase: true}, nil
}

func (m *mockEx) GetSubscriptionTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"assetName": func(a asset.Item) string {
			if a == asset.Futures {
				return "future"
			}
			return a.String()
		},
	}
}

// TestQualifiedChannels exercises QualifiedChannels
func TestQualifiedChannels(t *testing.T) {
	t.Parallel()
	l := &List{
		{Channel: "asset.{{$asset}}.{{$s.Interval.Short}}",
			Asset:    asset.All,
			Pairs:    currency.Pairs{btcusdtPair, ethusdcPair},
			Interval: kline.FifteenMin},
		{Channel: "pair.{{$pair}}.{{$s.Interval.Short}}",
			Asset:    asset.Spot,
			Pairs:    currency.Pairs{btcusdtPair, ethusdcPair},
			Interval: kline.FifteenMin},
		{Channel: "candles.{{assetName $asset}}.{{$pair.Swap.String}}.{{if eq $pair.String `BTCUSDT`}}{{$s.Params.color}}{{else}}red{{end}}.{{$s.Interval.Short}}",
			Asset: asset.All,
			Pairs: currency.Pairs{btcusdtPair, ethusdcPair},
			Params: map[string]any{
				"color": "green",
			},
			Interval: kline.FifteenMin},
	}
	got, err := l.QualifiedChannels2(&mockEx{})
	require.NoError(t, err, "QualifiedChannels must not error")
	exp := List{
		{Channel: "asset.spot.15m", Asset: asset.Spot, Pairs: currency.Pairs{btcusdtPair, ethusdcPair}, Interval: kline.FifteenMin},
		{Channel: "asset.futures.15m", Asset: asset.Futures, Pairs: currency.Pairs{btcusdtPair, ethusdcPair}, Interval: kline.FifteenMin},
		{Channel: "pair.BTCUSDT.15m", Asset: asset.Spot, Pairs: currency.Pairs{btcusdtPair}, Interval: kline.FifteenMin},
		{Channel: "pair.ETHUSDC.15m", Asset: asset.Spot, Pairs: currency.Pairs{ethusdcPair}, Interval: kline.FifteenMin},
		{Channel: "candles.spot.USDTBTC.green.15m", Asset: asset.Spot, Pairs: currency.Pairs{btcusdtPair}, Interval: kline.FifteenMin},
		{Channel: "candles.spot.USDCETH.red.15m", Asset: asset.Spot, Pairs: currency.Pairs{ethusdcPair}, Interval: kline.FifteenMin},
		{Channel: "candles.future.USDTBTC.green.15m", Asset: asset.Futures, Pairs: currency.Pairs{btcusdtPair}, Interval: kline.FifteenMin},
		{Channel: "candles.future.USDCETH.red.15m", Asset: asset.Futures, Pairs: currency.Pairs{ethusdcPair}, Interval: kline.FifteenMin},
	}
	equalLists(t, exp, got)

	// Test funcmap
}
