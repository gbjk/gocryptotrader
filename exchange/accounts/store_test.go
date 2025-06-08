package accounts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	t.Parallel()
	s := NewStore()
	require.NotNil(t, s, "NewStore must return a store")
	require.NotNil(t, s.mux, "NewStore must set mux")
	require.NotNil(t, s.exchangeAccounts, "NewStore must set exchangeAccounts")
}

func TestGetStore(t *testing.T) {
	t.Parallel()
	// Initialize global in case of -count=N+; No other tests should be relying on it
	global.Store(nil)
	s := GetStore()
	require.NotNil(t, s)
	require.Same(t, global.Load(), s, "GetStore must initialize the store")
	require.Same(t, s, GetStore(), "GetStore must return the global store on second call")
}
