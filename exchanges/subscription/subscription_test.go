package subscription

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
)

// TestEnsureKeyed logic test
func TestEnsureKeyed(t *testing.T) {
	t.Parallel()
	s := &Subscription{
		Channel: "candles",
		Asset:   asset.Spot,
		Pairs:   []currency.Pair{currency.NewPair(currency.BTC, currency.USDT)},
	}
	k1, ok := s.EnsureKeyed().(*Subscription)
	if assert.True(t, ok, "EnsureKeyed should return a *Subscription") {
		assert.Same(t, k1, s, "Key should point to the same struct")
	}
	type platypus string
	s = &Subscription{
		Key:     platypus("Gerald"),
		Channel: "orderbook",
		Asset:   asset.Margin,
		Pairs:   []currency.Pair{currency.NewPair(currency.ETH, currency.USDC)},
	}
	k2, ok := s.EnsureKeyed().(platypus)
	if assert.True(t, ok, "EnsureKeyed should return a platypus") {
		assert.Exactly(t, k2, s.Key, "ensureKeyed should set the same key")
		assert.EqualValues(t, "Gerald", k2, "key should have the correct value")
	}
}

// TestMarshalling logic test
func TestMarshaling(t *testing.T) {
	t.Parallel()
	j, err := json.Marshal(&Subscription{Channel: CandlesChannel})
	assert.NoError(t, err, "Marshalling should not error")
	assert.Equal(t, `{"enabled":false,"channel":"candles"}`, string(j), "Marshalling should be clean and concise")

	j, err = json.Marshal(&Subscription{Enabled: true, Channel: OrderbookChannel, Interval: kline.FiveMin, Levels: 4})
	assert.NoError(t, err, "Marshalling should not error")
	assert.Equal(t, `{"enabled":true,"channel":"orderbook","interval":"5m","levels":4}`, string(j), "Marshalling should be clean and concise")

	j, err = json.Marshal(&Subscription{Enabled: true, Channel: OrderbookChannel, Interval: kline.FiveMin, Levels: 4, Pairs: currency.Pairs{currency.NewPair(currency.BTC, currency.USDT)}})
	assert.NoError(t, err, "Marshalling should not error")
	assert.Equal(t, `{"enabled":true,"channel":"orderbook","pairs":"BTCUSDT","interval":"5m","levels":4}`, string(j), "Marshalling should be clean and concise")

	j, err = json.Marshal(&Subscription{Enabled: true, Channel: MyTradesChannel, Authenticated: true})
	assert.NoError(t, err, "Marshalling should not error")
	assert.Equal(t, `{"enabled":true,"channel":"myTrades","authenticated":true}`, string(j), "Marshalling should be clean and concise")
}

// TestSetState tests Subscription state changes
func TestSetState(t *testing.T) {
	t.Parallel()

	s := &Subscription{Key: 42, Channel: "Gophers"}
	assert.Equal(t, InactiveState, s.State(), "State should start as unknown")
	require.NoError(t, s.SetState(SubscribingState), "SetState should not error")
	assert.Equal(t, SubscribingState, s.State(), "State should be set correctly")
	assert.ErrorIs(t, s.SetState(SubscribingState), ErrInStateAlready, "SetState should error on same state")
	assert.ErrorIs(t, s.SetState(UnsubscribingState+1), ErrInvalidState, "Setting an invalid state should error")
	require.NoError(t, s.SetState(UnsubscribingState), "SetState should not error")
	assert.Equal(t, UnsubscribingState, s.State(), "State should be set correctly")
}
