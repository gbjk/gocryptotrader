package subscription

// MatchableKey interface should be implemented by Key types which want a more complex matching than a simple key equality check
// The Subscription method allows keys to compare against keys of other types
type MatchableKey interface {
	Match(MatchableKey) bool
	GetSubscription() *Subscription
}

// SubsetPairsKey is a key type for subscriptions where a subset of the pairs may match
// With empty pairs then only a sub without pairs will match
// Otherwise a sub which contains all the pairs will match
// Most likely usage is searching for a subscriptions to many pairs given just one of those pairs
type SubsetPairsKey struct {
	*Subscription
}

var _ MatchableKey = SubsetPairsKey{} // Enforce SubsetPairsKey must implement MatchableKey

// Subscription returns the underlying subscription
func (k SubsetPairsKey) GetSubscription() *Subscription {
	return k.Subscription
}

// Match implements MatchableKey
// Returns true if the key fields exactly matches the subscription, including all Pairs
func (k SubsetPairsKey) Match(eachKey MatchableKey) bool {
	eachSub := eachKey.GetSubscription()
	switch {
	case eachSub.Channel != k.Channel,
		eachSub.Asset != k.Asset,
		// len(eachSub.Pairs) == 0 && len(s.Pairs) == 0: Okay; continue to next non-pairs check
		len(eachSub.Pairs) == 0 && len(k.Pairs) != 0,
		len(eachSub.Pairs) != 0 && len(k.Pairs) == 0,
		len(k.Pairs) != 0 && eachSub.Pairs.ContainsAll(k.Pairs, true) != nil,
		eachSub.Levels != k.Levels,
		eachSub.Interval != k.Interval:
		return false
	}

	return true
}

// ExactKey is key type for subscriptions where all the pairs in a Subscription must match exactly
type ExactKey struct {
	*Subscription
}

var _ MatchableKey = ExactKey{} // Enforce ExactKey must implement MatchableKey

// Subscription returns the underlying subscription
func (k ExactKey) GetSubscription() *Subscription {
	return k.Subscription
}

// Match implements MatchableKey
// Returns true if the key fields exactly matches the subscription, including all Pairs
func (k ExactKey) Match(eachKey MatchableKey) bool {
	eachSub := eachKey.GetSubscription()
	switch {
	case eachSub.Channel != k.Channel,
		eachSub.Asset != k.Asset,
		eachSub.Pairs.ContainsAll(k.Pairs, true) != nil,
		k.Pairs.ContainsAll(eachSub.Pairs, true) != nil,
		eachSub.Levels != k.Levels,
		eachSub.Interval != k.Interval:
		return false
	}

	return true
}

// IdentityKey is key type for subscriptions where we know the exact sub we want to remove, so we use the pointer itself
// This is useful during Unsubscribe, when you might have concurrently reduced 2 subscriptions to having no Pairs
type IdentityKey struct {
	*Subscription
}

var _ MatchableKey = IdentityKey{} // Enforce IdentityKey must implement MatchableKey

// Subscription returns the underlying subscription
func (k IdentityKey) GetSubscription() *Subscription {
	return k.Subscription
}

// Match implements MatchableKey
// Returns true if the key is the same pointer address
func (k IdentityKey) Match(eachKey MatchableKey) bool {
	eachSub := eachKey.GetSubscription()
	return eachSub == k.Subscription
}

// ImmutableKey is key type for subscriptions where we don't want the key to change as the subscription changes
// So we want to consider the Key of each subscription, not the subscription itself
// Unlike other keys, this key will only consider subs with an immutable key when used as a key itself
type ImmutableKey struct {
	Key Subscription
	*Subscription
}

var _ MatchableKey = ImmutableKey{} // Enforce ImmutableKey must implement MatchableKey

// Subscription returns the underlying subscription
func (k ImmutableKey) GetSubscription() *Subscription {
	return k.Subscription
}

// Match implements MatchableKey
// Returns true if the key is the same pointer address
func (k ImmutableKey) Match(eachSubKey MatchableKey) bool {
	i, ok := eachSubKey.(*ImmutableKey)
	if !ok {
		return false
	}

	eachKey := i.Key

	switch {
	case eachKey.Channel != k.Channel,
		eachKey.Asset != k.Asset,
		eachKey.Pairs.ContainsAll(k.Pairs, true) != nil,
		k.Pairs.ContainsAll(eachKey.Pairs, true) != nil,
		eachKey.Levels != k.Levels,
		eachKey.Interval != k.Interval:
		return false
	}

	return true
}
