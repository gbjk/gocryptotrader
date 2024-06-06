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
		"assetName": func(s string) string {
			a, err := asset.New(s)
			if err != nil {
				return "unknown"
			}
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
		{Channel: "plain.$pair.{{$s.Interval.Short}}",
			Asset:    asset.Spot,
			Pairs:    currency.Pairs{btcusdtPair, ethusdcPair},
			Interval: kline.FifteenMin},
		{Channel: "plain-pipeline.{{`$pair`}}.{{$s.Interval.Short}}",
			Asset:    asset.Spot,
			Pairs:    currency.Pairs{btcusdtPair, ethusdcPair},
			Interval: kline.FifteenMin},
		{Channel: "plain.{{assetName `$asset`}}.$pair.{{$s.Params.color}}.{{$s.Interval.Short}}",
			Asset: asset.All,
			Pairs: currency.Pairs{btcusdtPair, ethusdcPair},
			Params: map[string]any{
				"color": "green",
			},
			Interval: kline.FifteenMin},
	}
	got, err := l.QualifiedChannels(&mockEx{})
	require.NoError(t, err, "QualifiedChannels must not error")
	exp := List{
		{Channel: "plain.BTCUSDT.15m", Asset: asset.Spot, Pairs: currency.Pairs{btcusdtPair}, Interval: kline.FifteenMin},
		{Channel: "plain.ETHUSDC.15m", Asset: asset.Spot, Pairs: currency.Pairs{ethusdcPair}, Interval: kline.FifteenMin},
		{Channel: "plain-pipeline.BTCUSDT.15m", Asset: asset.Spot, Pairs: currency.Pairs{btcusdtPair}, Interval: kline.FifteenMin},
		{Channel: "plain-pipeline.ETHUSDC.15m", Asset: asset.Spot, Pairs: currency.Pairs{ethusdcPair}, Interval: kline.FifteenMin},
		{Channel: "plain.spot.BTCUSDT.15m", Asset: asset.Spot, Pairs: currency.Pairs{btcusdtPair}, Interval: kline.FifteenMin},
		{Channel: "plain.spot.ETHUSDC.15m", Asset: asset.Spot, Pairs: currency.Pairs{ethusdcPair}, Interval: kline.FifteenMin},
		{Channel: "plain.future.BTCUSDT.15m", Asset: asset.Futures, Pairs: currency.Pairs{btcusdtPair}, Interval: kline.FifteenMin},
		{Channel: "plain.future.ETHUSDC.15m", Asset: asset.Futures, Pairs: currency.Pairs{ethusdcPair}, Interval: kline.FifteenMin},
	}
	equalLists(t, exp, got)

	// Test funcmap
}
