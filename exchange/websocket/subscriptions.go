package websocket

import (
	"errors"
	"fmt"
	"slices"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/exchange/subscription"
)

// Public subscription errors
var (
	ErrSubscriptionFailure     = errors.New("subscription failure")
	ErrSubscriptionsNotAdded   = errors.New("subscriptions not added")
	ErrSubscriptionsNotRemoved = errors.New("subscriptions not removed")
)

// Public subscription errors
var (
	errSubscriptionsExceedsLimit = errors.New("subscriptions exceeds limit")
)

// UnsubscribeChannels unsubscribes from a list of websocket channel
func (m *Manager) UnsubscribeChannels(channels subscription.List) error {
	if len(channels) == 0 {
		return nil // No channels to unsubscribe from is not an error
	}

	if err := common.NilGuard(m); err != nil {
		return err
	}
	if err := common.NilGuard(m.Unsubscriber); err != nil {
		return err
	}

	connSubs := map[*Connection]subscription.List{}
	for i, s := range channels {
		if m.subscriptions.Get(s) == nil {
			return fmt.Errorf("%w: %s", subscription.ErrNotFound, s)
		}
		connSubs[i] = s
	}

	return m.Unsubscriber(subs)
}

// ResubscribeToChannel resubscribes to channel
// Sets state to Resubscribing, and exchanges which want to maintain a lock on it can respect this state and not RemoveSubscription
// Errors if subscription is already subscribing
func (m *Manager) ResubscribeToChannel(s *subscription.Subscription) error {
	l := subscription.List{s}
	if err := s.SetState(subscription.ResubscribingState); err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	if err := m.UnsubscribeChannels(l); err != nil {
		return err
	}
	return m.SubscribeToChannels(l)
}

// SubscribeToChannels subscribes to websocket channels using the exchange specific Subscriber method
// Errors are returned for duplicates or exceeding max Subscriptions
func (m *Manager) SubscribeToChannels(subs subscription.List) error {
	if slices.Contains(subs, nil) {
		return fmt.Errorf("%w: List parameter contains an nil element", common.ErrNilPointer)
	}

	if err := m.checkSubscriptions(nil, subs); err != nil {
		return err
	}

	if m.Subscriber == nil {
		return fmt.Errorf("%w: Global Subscriber not set", common.ErrNilPointer)
	}

	if err := m.Subscriber(conn, subs); err != nil {
		return fmt.Errorf("%w: %w", ErrSubscriptionFailure, err)
	}
	return nil
}

// AddSubscriptions adds subscriptions to the subscription store
// Sets state to Subscribing unless the state is already set
func (m *Manager) AddSubscriptions(conn Connection, subs ...*subscription.Subscription) error {
	if m == nil {
		return fmt.Errorf("%w: AddSubscriptions called on nil Websocket", common.ErrNilPointer)
	}

	// TODO: GBJK Please tell me we don't need to do this
	// We definitely don't - Because we're not adding Subscriptions to a Connection
	var subscriptionStore **subscription.Store
	/*
		    * TODO: GBJK - Need to work out what to add it to
			if wrapper, ok := m.connections[conn]; ok && conn != nil {
				subscriptionStore = &wrapper.subscriptions
			} else {
				subscriptionStore = &m.subscriptions
			}
	*/

	if *subscriptionStore == nil {
		*subscriptionStore = subscription.NewStore()
	}
	var errs error
	for _, s := range subs {
		if s.State() == subscription.InactiveState {
			if err := s.SetState(subscription.SubscribingState); err != nil {
				errs = common.AppendError(errs, fmt.Errorf("%w: %s", err, s))
			}
		}
		if err := (*subscriptionStore).Add(s); err != nil {
			errs = common.AppendError(errs, err)
		}
	}
	return errs
}

// AddSuccessfulSubscriptions marks subscriptions as subscribed and adds them to the subscription store
func (m *Manager) AddSuccessfulSubscriptions(conn Connection, subs ...*subscription.Subscription) error {
	if m == nil {
		return fmt.Errorf("%w: AddSuccessfulSubscriptions called on nil Websocket", common.ErrNilPointer)
	}
	var errs error
	/*
	    TODO: GBJK - This doesn't feel right - Manager should be doing this internally, right?
	   Should an exchange even get to call this?
	   I mean, it gives us flexibility, but at the cost of control

	   	var subscriptionStore **subscription.Store
	   	if wrapper, ok := m.connections[conn]; ok && conn != nil {
	   		subscriptionStore = &wrapper.subscriptions
	   	} else {
	   		subscriptionStore = &m.subscriptions
	   	}

	   	if *subscriptionStore == nil {
	   		*subscriptionStore = subscription.NewStore()
	   	}

	   	for _, s := range subs {
	   		if err := s.SetState(subscription.SubscribedState); err != nil {
	   			errs = common.AppendError(errs, fmt.Errorf("%w: %s", err, s))
	   		}
	   		if err := (*subscriptionStore).Add(s); err != nil {
	   			errs = common.AppendError(errs, err)
	   		}
	   	}
	*/
	return errs
}

// RemoveSubscriptions removes subscriptions from the subscription list and sets the status to Unsubscribed
func (m *Manager) RemoveSubscriptions(conn Connection, subs ...*subscription.Subscription) error {
	if m == nil {
		return fmt.Errorf("%w: RemoveSubscriptions called on nil Websocket", common.ErrNilPointer)
	}

	if m.subscriptions == nil {
		return fmt.Errorf("%w: RemoveSubscriptions called on uninitialised Websocket", common.ErrNilPointer)
	}

	var errs error
	for _, s := range subs {
		if err := s.SetState(subscription.UnsubscribedState); err != nil {
			errs = common.AppendError(errs, fmt.Errorf("%w: %s", err, s))
		}
		if err := m.subscriptions.Remove(s); err != nil {
			errs = common.AppendError(errs, err)
		}
	}
	return errs
}

// GetSubscription returns a subscription at the key provided
// returns nil if no subscription is at that key or the key is nil
// Keys can implement subscription.MatchableKey in order to provide custom matching logic
func (m *Manager) GetSubscription(key any) *subscription.Subscription {
	if m == nil || key == nil {
		return nil
	}
	for _, c := range m.connectionConfigs {
		if c.subscriptions == nil {
			continue
		}
		sub := c.subscriptions.Get(key)
		if sub != nil {
			return sub
		}
	}
	if m.subscriptions == nil {
		return nil
	}
	return m.subscriptions.Get(key)
}

// GetSubscriptions returns a new slice of the subscriptions
func (m *Manager) GetSubscriptions() subscription.List {
	if m == nil {
		return nil
	}
	var subs subscription.List
	for _, c := range m.connectionConfigs {
		if c.subscriptions != nil {
			subs = append(subs, c.subscriptions.List()...)
		}
	}
	if m.subscriptions != nil {
		subs = append(subs, m.subscriptions.List()...)
	}
	return subs
}

// TODO: GBJK: This needs rethinking
// checkSubscriptions checks subscriptions against the max subscription limit and if the subscription already exists
// The subscription state is not considered when counting existing subscriptions
func (m *Manager) checkSubscriptions(conn Connection, subs subscription.List) error {
	var subscriptionStore *subscription.Store
	if wrapper, ok := m.connections[conn]; ok && conn != nil {
		subscriptionStore = wrapper.subscriptions
	} else {
		subscriptionStore = m.subscriptions
	}
	if subscriptionStore == nil {
		return fmt.Errorf("%w: Websocket.subscriptions", common.ErrNilPointer)
	}

	existing := subscriptionStore.Len()
	if m.MaxSubscriptionsPerConnection > 0 && existing+len(subs) > m.MaxSubscriptionsPerConnection {
		return fmt.Errorf("%w: current subscriptions: %v, incoming subscriptions: %v, max subscriptions per connection: %v - please reduce enabled pairs",
			errSubscriptionsExceedsLimit,
			existing,
			len(subs),
			m.MaxSubscriptionsPerConnection)
	}

	for _, s := range subs {
		if s.State() == subscription.ResubscribingState {
			continue
		}
		if found := subscriptionStore.Get(s); found != nil {
			return fmt.Errorf("%w: %s", subscription.ErrDuplicate, s)
		}
	}

	return nil
}

// SyncSubscriptions flushes channel subscriptions when there is a pair/asset change
// TODO: GBJK Add window for max subscriptions per connection, to spawn new connections if needed.
func (m *Manager) SyncSubscriptions() error {
	if !m.IsEnabled() {
		return fmt.Errorf("%s %w", m.exchangeName, ErrWebsocketNotEnabled)
	}

	if !m.IsConnected() {
		return fmt.Errorf("%s %w", m.exchangeName, ErrNotConnected)
	}

	// If the exchange does not support subscribing and or unsubscribing the full connection needs to be flushed to maintain consistency
	// TODO: GBJK Any examples of this?
	// TODO: GBJK - This isn't sane in an N+ Connections situation
	if !m.features.Subscribe || !m.features.Unsubscribe {
		m.m.Lock()
		defer m.m.Unlock()
		if err := m.shutdown(); err != nil {
			return err
		}
		return m.connect()
	}

	/*
		    * TODO: GBJK - Generate subscriptions first. Then find the conns for them
			for x := range m.connectionConfigs {
				newSubs, err := m.connectionConfigs[x].setup.GenerateSubscriptions()
				if err != nil {
					return err
				}

				// Case if there is nothing to unsubscribe from and the connection is nil
				if len(newSubs) == 0 && m.connectionConfigs[x].connection == nil {
					continue
				}

				// If there are subscriptions to subscribe to but no connection to subscribe to, establish a new connection.
				if m.connectionConfigs[x].connection == nil {
					conn := m.getConnectionFromSetup(m.connectionConfigs[x].setup)
					if err := m.connectionConfigs[x].setup.Connector(context.TODO(), conn); err != nil {
						return err
					}
					m.Wg.Add(1)
					go m.Reader(context.TODO(), conn, m.connectionConfigs[x].setup.Handler)
					m.connections[conn] = m.connectionConfigs[x]
					m.connectionConfigs[x].connection = conn
				}

				err = m.updateChannelSubscriptions(m.connectionConfigs[x].connection, m.connectionConfigs[x].subscriptions, newSubs)
				if err != nil {
					return err
				}

				// If there are no subscriptions to subscribe to, close the connection as it is no longer needed.
				if m.connectionConfigs[x].subscriptions.Len() == 0 {
					delete(m.connections, m.connectionConfigs[x].connection) // Remove from lookup map
					if err := m.connectionConfigs[x].connection.Shutdown(); err != nil {
						log.Warnf(log.WebsocketMgr, "%v websocket: failed to shutdown connection: %v", m.exchangeName, err)
					}
					m.connectionConfigs[x].connection = nil
				}
			}
	*/
	return nil
}

// updateChannelSubscriptions subscribes or unsubscribes from channels and checks that the correct number of channels
// have been subscribed to or unsubscribed from.
func (m *Manager) updateChannelSubscriptions(c Connection, store *subscription.Store, incoming subscription.List) error {
	subs, unsubs := store.Diff(incoming)
	if len(unsubs) != 0 {
		if err := m.UnsubscribeChannels(c, unsubs); err != nil {
			return err
		}

		if contained := store.Contained(unsubs); len(contained) > 0 {
			return fmt.Errorf("%v %w %q", m.exchangeName, ErrSubscriptionsNotRemoved, contained)
		}
	}
	if len(subs) != 0 {
		if err := m.SubscribeToChannels(c, subs); err != nil {
			return err
		}

		if missing := store.Missing(subs); len(missing) > 0 {
			return fmt.Errorf("%v %w %q", m.exchangeName, ErrSubscriptionsNotAdded, missing)
		}
	}
	return nil
}
