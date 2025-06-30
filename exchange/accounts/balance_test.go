package accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/currency"
)

// TestCurrencyBalancesBalance exercises currencyBalances.balance
func TestCurrencyBalancesBalance(t *testing.T) {
	t.Parallel()

	c := currencyBalances{}
	b := c.balance(currency.BTC)
	require.NotNil(t, b)
	assert.Same(t, c[currency.BTC], b, "should make and return the same entry")
	assert.Same(t, b, c.balance(currency.BTC), "should make and return the same entry")
}
