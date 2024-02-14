package subscription

import (
	"maps"
	"sync"
)

// Store is a container of subscription pointers
type Store struct {
	m  map[any]*Subscription
	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		m: map[any]*Subscription{},
	}
}

// NewStoreFromList creates a Store from a List
func NewStoreFromList(s List) *Store {
	n := new(Store)
	for _, c := range s {
		n.Add(c)
	}
	return n
}

// Add copies a subscription and adds it to the store
// Key can be already set; if ommitted EnsureKeyed will be used
// Errors if it already exists
func (s *Store) Add(sub *Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := sub.EnsureKeyed()
	if e := s.get(key); e != nil {
		return ErrDuplicate
	}

	s.m[key] = sub

	return nil
}

// Get returns a pointer to a subscription or nil if not found
// If key implements MatchableKey then key.Match will be used
func (s *Store) Get(key any) *Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.get(key)
}

// get returns a pointer to subscription or nil if not found
// If key implements MatchableKey then key.Match will be used
// If key is actually a subscription then ensureKeyed will be used and we'll return the
// This method provides no locking protection
// returned subscriptions are implicitly guaranteed to have a key
func (s *Store) get(key any) *Subscription {
	switch v := key.(type) {
	case *Subscription:
		return s.get(v.EnsureKeyed())
	case Subscription:
		return s.get(v.EnsureKeyed())
	case MatchableKey:
		return s.match(v)
	default:
		return s.m[v]
	}
}

// Remove removes a subscription from the store
func (s *Store) Remove(sub *Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if found := s.get(sub); found != nil {
		delete(s.m, found.Key)
	}
}

// List returns a slice of Subscriptions pointers
func (s *Store) List() []*Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	subs := make([]*Subscription, 0, len(s.m))
	for _, s := range s.m {
		subs = append(subs, s)
	}
	return subs
}

// Clear empties the subscription store
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.m)
}

// match returns the first subscription which matches the Key's Asset, Channel and Pairs
// If the key provided has:
// 1) Empty pairs then only Subscriptions without pairs will be considered
// 2) >=1 pairs then Subscriptions which contain all the pairs will be considered
// This method provides no locking protection
func (s *Store) match(key MatchableKey) *Subscription {
	for anyKey, s := range s.m {
		if key.Match(anyKey) {
			return s
		}
	}
	return nil
}

// Diff returns a list of the added and missing subs from a new list
// The store Diff is invoked upon is read-lock protected
// The new store is assumed to be a new instance and enjoys no locking protection
func (s *Store) Diff(compare List) (added, removed List) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	removedMap := maps.Clone(s.m)
	for _, sub := range compare {
		if found := s.get(sub); found != nil {
			delete(removedMap, found.Key)
		} else {
			added = append(added, sub)
		}
	}

	for _, c := range removedMap {
		removed = append(removed, c)
	}

	return
}

// Len returns the number of subscriptions
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.m)
}
