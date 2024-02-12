package subscription

import (
	"errors"
	"fmt"
	"maps"
	"slices"

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

// Key is the fallback key for AddSuccessfulSubscriptions
// It provides for matching on one or more keys
type Key struct {
	Channel string
	Pairs   *currency.Pairs
	Asset   asset.Item
}

// Map is a container of subscription pointers
type Map map[any]*Subscription

// List is a container of subscription pointers
type List []Subscription

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
		s.Key = Key{
			Channel: s.Channel,
			Asset:   s.Asset,
			Pairs:   &s.Pairs,
		}
	}
	return s.Key
}

// Match returns the first subscription which matches the Key's Asset, Channel and Pairs
// If the key provided has:
// 1) Empty pairs then only Subscriptions without pairs will be considered
// 2) >=1 pairs then Subscriptions which contain all the pairs will be considered
func (k Key) Match(m Map) *Subscription {
	for anyKey, s := range m {
		candidate, ok := anyKey.(Key)
		if !ok {
			continue
		}
		if k.Channel != candidate.Channel {
			continue
		}
		if k.Asset != candidate.Asset {
			continue
		}
		if k.Pairs == nil || len(*k.Pairs) == 0 {
			if candidate.Pairs == nil || len(*candidate.Pairs) == 0 {
				return s
			}
			continue // Case (1) - key doesn't have any pairs but candidate does
		}
		if err := candidate.Pairs.ContainsAll(*k.Pairs, true); err == nil {
			return s
		}
	}
	return nil
}

// ListToMap creates a Map from a slice of subscriptions
func ListToMap(s List) *Map {
	n := Map{}
	for _, c := range s {
		n[c.EnsureKeyed()] = &c
	}
	return &n
}

// Diff returns a list of the added and missing subs between two maps
func (m *Map) Diff(newSubs *Map) (sub, unsub List) {
	oldSubs := maps.Clone(*m)
	for _, s := range *newSubs {
		key := s.EnsureKeyed()

		var found *Subscription
		if m, ok := key.(MatchableKey); ok {
			found = m.Match(oldSubs)
		} else {
			found = oldSubs[key]
		}

		if found != nil {
			delete(oldSubs, found.Key) // If it's in both then we remove it from the unsubscribe list
		} else {
			sub = append(sub, *s) // If it's in newSubs but not oldSubs subs we want to subscribe
		}
	}

	for _, c := range oldSubs {
		unsub = append(unsub, *c)
	}

	return
}

func (l List) Strings() []string {
	s := make([]string, len(l))
	for i := range l {
		s[i] = l[i].String()
	}
	slices.Sort(s)
	return s
}
