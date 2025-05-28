package accounts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStore(t *testing.T) {
	t.Parallel()
	require.NotNil(t, NewStore())
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
