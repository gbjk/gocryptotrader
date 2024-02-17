package subscription

import (
	"errors"
	"fmt"
	"sync"

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

// Public errors
var (
	ErrNotSinglePair  = errors.New("only single pair subscriptions expected")
	ErrInStateAlready = errors.New("subscription already in state")
	ErrInvalidState   = errors.New("invalid subscription state")
	ErrDuplicate      = errors.New("duplicate subscription")
)

// State tracks the status of a subscription channel
type State uint8

// Subscription container for streaming subscriptions
type Subscription struct {
	Enabled       bool                   `json:"enabled"`
	Key           any                    `json:"-"`
	Channel       string                 `json:"channel,omitempty"`
	Pairs         currency.Pairs         `json:"pairs,omitempty"`
	Asset         asset.Item             `json:"asset,omitempty"`
	Params        map[string]interface{} `json:"params,omitempty"`
	Interval      kline.Interval         `json:"interval,omitempty"`
	Levels        int                    `json:"levels,omitempty"`
	Authenticated bool                   `json:"authenticated,omitempty"`
	state         State
	m             sync.RWMutex
}

// MatchableKey interface should be implemented by Key types which want a more complex matching than a simple key equality check
type MatchableKey interface {
	Match(any) bool
}

// String implements the Stringer interface for Subscription, giving a human representation of the subscription
func (s *Subscription) String() string {
	return fmt.Sprintf("%s %s %s", s.Channel, s.Asset, s.Pairs)
}

// State returns the subscription state
func (s *Subscription) State() State {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.state
}

// SetState sets the subscription state
// Errors if already in that state or the new state is not valid
func (s *Subscription) SetState(state State) error {
	s.m.Lock()
	defer s.m.Unlock()
	if state == s.state {
		return ErrInStateAlready
	}
	if state > UnsubscribingState {
		return ErrInvalidState
	}
	s.state = state
	return nil
}

// EnsureKeyed sets the default key on a channel if it doesn't have one
// Returns key for convenience
func (s *Subscription) EnsureKeyed() any {
	if s.Key == nil {
		s.Key = s
	}
	return s.Key
}

// Match returns if the two keys match Channels, Assets, Pairs, Interval and Levels:
// Key Pairs comparison:
// 1) Empty pairs then only Subscriptions without pairs match
// 2) >=1 pairs then Subscriptions which contain all the pairs match
// Such that a subscription for all enabled pairs will be matched when seaching for any one pair
func (s *Subscription) Match(key any) bool {
	b, ok := key.(*Subscription)
	switch {
	case !ok,
		s.Channel != b.Channel,
		s.Asset != b.Asset,
		len(b.Pairs) == 0 && len(s.Pairs) != 0,
		// len(b.Pairs) == 0 && len(s.Pairs) == 0: Okay; continue to next non-pairs check
		len(b.Pairs) != 0 && len(s.Pairs) == 0,
		len(b.Pairs) != 0 && s.Pairs.ContainsAll(b.Pairs, true) != nil,
		s.Levels != b.Levels,
		s.Interval != b.Interval:
		return false
	}

	return true
}
