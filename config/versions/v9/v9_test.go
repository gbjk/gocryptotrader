package v9_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v9 "github.com/thrasher-corp/gocryptotrader/config/versions/v9"
)

func TestExchanges(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []string{"Binance"}, new(v9.Version).Exchanges())
}

func TestUpgradeExchange(t *testing.T) {
	t.Parallel()
	t.Error("G didn't implement")
}

func TestDowngradeExchange(t *testing.T) {
	t.Parallel()
	t.Error("G didn't implement")
}
