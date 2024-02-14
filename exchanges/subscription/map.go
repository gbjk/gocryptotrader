package subscription

import (
	"maps"
	"sync"
)

// Map is a container of subscription pointers
type Map struct {
	m  map[any]*Subscription
	mu sync.RWMutex
}

// ListToMap creates a Map from a slice of subscriptions
func ListToMap(s List) *Map {
	n := new(Map)
	for _, c := range s {
		n.Add(c)
	}
	return n
}

// Add copies a subscription and adds it to the map
// Key can be already set; if ommitted EnsureKeyed will be used
// Errors if it already exists
func (m *Map) Add(s *Subscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := s.ensureKeyed()
	if e := m.get(key); e != nil {
		return ErrDuplicate
	}

	m.m[key] = s

	return nil
}

// Get returns a pointer to a subscription or nil if not found
// If key implements MatchableKey then key.Match will be used
func (m *Map) Get(key any) *Subscription {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.get(key)
}

// get returns a pointer to subscription or nil if not found
// If key implements MatchableKey then key.Match will be used
// If key is actually a subscription then ensureKeyed will be used and we'll return the
// This method provides no locking protection
// returned subscriptions are implicitly guaranteed to have a key
func (m *Map) get(key any) *Subscription {
	switch v := key.(type) {
	case *Subscription:
		return m.get(v.ensureKeyed())
	case Subscription:
		return m.get(v.ensureKeyed())
	case MatchableKey:
		return m.match(v)
	default:
		return m.m[v]
	}
}

// Remove removes a subscription from the map
func (m *Map) Remove(s *Subscription) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s := m.get(s); s != nil {
		delete(m.m, s.key)
	}
}

// List returns a slice of Subscriptions pointers
func (m *Map) List() []*Subscription {
	m.mu.RLock()
	defer m.mu.RUnlock()
	subs := make([]*Subscription, 0, len(m.m))
	for _, s := range m.m {
		subs = append(subs, s)
	}
	return subs
}

// Clear empties the subscription map
func (m *Map) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.m)
}

// match returns the first subscription which matches the Key's Asset, Channel and Pairs
// If the key provided has:
// 1) Empty pairs then only Subscriptions without pairs will be considered
// 2) >=1 pairs then Subscriptions which contain all the pairs will be considered
// This method provides no locking protection
func (m *Map) match(key MatchableKey) *Subscription {
	for anyKey, s := range m.m {
		if key.Match(anyKey) {
			return s
		}
	}
	return nil
}

// Diff returns a list of the added and missing subs between two maps
func (m *Map) Diff(newSubs *Map) (sub, unsub List) {
	m.mu.RLock()
	defer m.mu.RUnlock() // More efficient to hold lock and call .get
	oldSubs := maps.Clone(m.m)
	for _, s := range newSubs.m {
		if found := m.get(s); found != nil {
			delete(oldSubs, found.key) // If it's in both then we remove it from the unsubscribe list
		} else {
			sub = append(sub, s) // If it's in newSubs but not oldSubs subs we want to subscribe
		}
	}

	for _, c := range oldSubs {
		unsub = append(unsub, c)
	}

	return
}

// Len returns the number of subscriptions
func (m *Map) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}
