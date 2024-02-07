package subscription

import (
	"errors"
	"fmt"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/kline"
)

const (
	UnknownState       State = iota // UnknownState subscription state is not registered, but doesn't imply Inactive
	SubscribingState                // SubscribingState means channel is in the process of subscribing
	SubscribedState                 // SubscribedState means the channel has finished a successful and acknowledged subscription
	UnsubscribingState              // UnsubscribingState means the channel has started to unsubscribe, but not yet confirmed

	TickerChannel    = "ticker"    // TickerChannel Subscription Type
	OrderbookChannel = "orderbook" // OrderbookChannel Subscription Type
	CandlesChannel   = "candles"   // CandlesChannel Subscription Type
	AllOrdersChannel = "allOrders" // AllOrdersChannel Subscription Type
	AllTradesChannel = "allTrades" // AllTradesChannel Subscription Type
	MyTradesChannel  = "myTrades"  // MyTradesChannel Subscription Type
	MyOrdersChannel  = "myOrders"  // MyOrdersChannel Subscription Type
)

var (
	ErrNotSinglePair = errors.New("only single pair subscriptions expected")
)

// State tracks the status of a subscription channel
type State uint8

// MatchableKey interface should be implemented by Key types which want a more complex matching than a simple key equality check
type MatchableKey interface {
	Match(Map) *Subscription
}

// MultiPairKey is the fallback key for AddSuccessfulSubscriptions
// It provides for matching on one or more keys
type MultiPairKey struct {
	Channel string
	Pairs   currency.Pairs
	Asset   asset.Item
}

// SinglePairKey is available as a key type which expects only one pair for a subscription
type SinglePairKey struct {
	Channel string
	Pair    currency.Pair
	Asset   asset.Item
}

// Map is a container of subscription pointers
type Map map[any]*Subscription

type SubscriptionInterface interface {
}

// Subscription container for streaming subscriptions
type Subscription struct {
	Enabled       bool                   `json:"enabled"`
	Key           any                    `json:"-"`
	Channel       string                 `json:"channel,omitempty"`
	Pairs         currency.Pairs         `json:"pairs,omitempty"`
	Asset         asset.Item             `json:"asset,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
	State         State                  `json:"-"`
	Interval      kline.Interval         `json:"interval,omitempty"`
	Levels        int                    `json:"levels,omitempty"`
	Authenticated bool                   `json:"authenticated,omitempty"`
}

// String implements the Stringer interface for Subscription, giving a human representation of the subscription
func (s *Subscription) String() string {
	return fmt.Sprintf("%s %s %s", s.Channel, s.Asset, s.Pairs)
}

// EnsureKeyed sets the default key on a channel if it doesn't have one
// Returns key for convenience
func (s *Subscription) EnsureKeyed() any {
	if s.Key == nil {
		s.Key = MultiPairKey{
			Channel: s.Channel,
			Asset:   s.Asset,
			Pairs:   s.Pairs,
		}
	}
	return s.Key
}

// Match returns the first subscription which matches the MultiPairKey's Asset, Channel and Pairs
// If the key provided has:
// * Empty pairs then only Subscriptions without pairs will be considered
// * >=1 pairs then Subscriptions which contain all the pairs will be considered
func (k *MultiPairKey) Match(m Map) *Subscription {
	for a, v := range m {
		candidate, ok := a.(MultiPairKey)
		if !ok {
			continue
		}
		if k.Channel != candidate.Channel {
			continue
		}
		if k.Asset != candidate.Asset {
			continue
		}
		if len(k.Pairs) == 0 && len(candidate.Pairs) == 0 {
			return v
		}
		if err := candidate.Pairs.ContainsAll(k.Pairs, true); err == nil {
			return v
		}
	}
	return nil
}
