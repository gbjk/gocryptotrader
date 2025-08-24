package websocket

import (
	"errors"
	"fmt"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/exchange/subscription"
	"github.com/thrasher-corp/gocryptotrader/log"
)

// Public subscription errors
var (
	ErrSubscribe               = errors.New("subscribe to websocket channel failed")
	ErrUnsubscribe             = errors.New("unsubscribe from websocket channel failed")
	ErrSubscriptionsNotAdded   = errors.New("subscriptions not added")
	ErrSubscriptionsNotRemoved = errors.New("subscriptions not removed")
)

// Private subscription errors
var (
	errSubscriptionsExceedsLimit = errors.New("subscriptions exceeds limit")
	errSyncSubscriptions         = errors.New("error syncing websocket manager subscriptions")
)

// Unsubscribe unsubscribes from a list of websocket channels
func (m *Manager) Unsubscribe(subs subscription.List) error {
	if len(subs) == 0 {
		return nil // No channels to unsubscribe from is not an error
	}

	if err := common.NilGuard(m); err != nil {
		return err
	}
	if err := common.NilGuard(m.Unsubscriber); err != nil {
		return err
	}

	connSubs := subs.GroupByConnection()

	errs := common.CollectErrors(len(connSubs))
	for c, subs := range connSubs {
		go func() {
			defer errs.Wg.Done()
			if err := m.Unsubscriber(c, subs); err != nil {
				errs.C <- fmt.Errorf("%w: %w", ErrUnsubscribe, err)
			}
		}()
	}

	return errs.Collect()
}

// ResubscribeToChannel resubscribes to channel
// Sets state to Resubscribing, and exchanges which want to maintain a lock on it can respect this state and not RemoveSubscription
// Errors if subscription is already subscribing
func (m *Manager) ResubscribeToChannel(s *subscription.Subscription) error {
	l := subscription.List{s}
	if err := s.SetState(subscription.ResubscribingState); err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	if err := m.Unsubscribe(l); err != nil {
		return err
	}
	return m.Subscribe(l)
}

// Subscribe subscribes to websocket channels using the exchange specific Subscriber method
// Errors are returned for duplicates or exceeding max Subscriptions
func (m *Manager) Subscribe(subs subscription.List) error {
	if len(subs) == 0 {
		return nil // No channels to unsubscribe from is not an error
	}
	if err := common.NilGuard(m); err != nil {
		return err
	}
	if err := common.NilGuard(m.Subscriber); err != nil {
		return err
	}

	connSubs, err := m.assignSubsToConns(subs)
	if err != nil {
		return err
	}

	errs := common.CollectErrors(len(connSubs))
	for c, subs := range connSubs {
		go func() {
			defer errs.Wg.Done()
			if err := m.Subscriber(c, subs); err != nil {
				errs.C <- fmt.Errorf("%w: %w", ErrSubscribe, err)
			}
		}()
	}

	return errs.Collect()
}

func (m *Manager) assignSubsToConns(subs subscription.List) (map[*Connection]subscription.List, error) {
	return nil, nil
}

// AddSubscriptions adds subscriptions to the subscription store
// Sets state to Subscribing unless the state is already set
// TODO: When and why and how would we do this like this?
// I think anything we're going to add should already be in the store?
func (m *Manager) AddSubscriptions(subs ...*subscription.Subscription) error {
	if err := common.NilGuard(m); err != nil {
		return err
	}

	var errs error
	for _, s := range subs {
		if s.State() == subscription.InactiveState {
			if err := s.SetState(subscription.SubscribingState); err != nil {
				errs = common.AppendError(errs, fmt.Errorf("%w: %s", err, s))
			}
		}
		if err := m.subscriptions.Add(s); err != nil {
			errs = common.AppendError(errs, err)
		}
	}
	return errs
}

// AddSuccessfulSubscriptions marks subscriptions as subscribed and adds them to the subscription store
func (m *Manager) AddSuccessfulSubscriptions(conn Connection, subs ...*subscription.Subscription) error {
	if err := common.NilGuard(m); err != nil {
		return err
	}
	if err := common.NilGuard(m.subscriptions); err != nil {
		return err
	}

	var errs error
	for _, s := range subs {
		if err := s.SetState(subscription.SubscribedState); err != nil {
			errs = common.AppendError(errs, fmt.Errorf("%w: %s", err, s))
		}
		if err := m.subscriptions.Add(s); err != nil {
			errs = common.AppendError(errs, err)
		}
	}
	return errs
}

// RemoveSubscriptions removes subscriptions from the subscription list and sets the status to Unsubscribed
func (m *Manager) RemoveSubscriptions(conn Connection, subs ...*subscription.Subscription) error {
	if err := common.NilGuard(m); err != nil {
		return err
	}
	if err := common.NilGuard(m.subscriptions); err != nil {
		return err
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
	if err := common.NilGuard(m, key); err != nil {
		return nil
	}
	if err := common.NilGuard(m.subscriptions); err != nil {
		return nil
	}
	return m.subscriptions.Get(key)
}

// GetSubscriptions returns a new slice of the subscriptions
func (m *Manager) GetSubscriptions() subscription.List {
	if err := common.NilGuard(m); err != nil {
		return nil
	}
	return m.subscriptions.List()
}

/*
// TODO: GBJK: This needs rethinking
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
*/

// SyncSubscriptions reconciles the existing subscriptions for changes
// Call when after a change affecting subscriptions external to the websocket manager, e.g. changing assets or pairs
// Does not need to be called after calling any other function inside the websocket manager
// Expceted to be called in a goro, so any errors are logged
func (m *Manager) SyncSubscriptions() {
	if err := m.SyncSubscriptions(); err != nil {
		log.Errorf(log.WebsocketMgr, "%s %w: %w", m.exchangeName, ErrSyncSubscriptions, err)
	}
}

func (m *Manager) syncSubscriptions() error {
	if !m.IsEnabled() {
		return ErrWebsocketNotEnabled
	}
	if !m.IsRunning() {
		return ErrNotRunning
	}

	subs, err := m.GenerateSubs()
	if err != nil {
		return err
	}

	subs, unsubs := m.subscriptions.Diff(subs)

	return common.AppendError(
		m.Unsubscribe(unsubs),
		m.subscriptions.Add(subs...), // Add directly so state is Inactiev; manager will find connections
	)
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
