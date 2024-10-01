package versions

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersion3Upgrade(t *testing.T) {
	t.Parallel()

	v := &Version3{}
	in := []byte(`{"name":"Wibble","currencyPairs":{"pairs":{"spot":{"enabled":"LTC_BTC","available":"ETH_BTC,BTC_USD"}}}}`)
	exp := []byte(`{"name":"Wibble","currencyPairs":{"pairs":{"spot":{"enabled":"LTC_BTC","available":"ETH_BTC,BTC_USD","assetEnabled":true}}}}`)
	out, err := v.UpgradeExchange(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, string(exp), string(out), "Exchanges without assetEnabled should have it added")

	in = []byte(`{"name":"Wibble","currencyPairs":{"pairs":{"spot":{"enabled":"LTC_BTC","available":"ETH_BTC,BTC_USD","assetEnabled":false}}}}`)
	out, err = v.UpgradeExchange(context.Background(), bytes.Clone(in))
	require.NoError(t, err)
	assert.Equal(t, string(in), string(out), "Exchanges with assetEnabled false should be left alone")

	in = []byte(`{"name":"Wibble","currencyPairs":{"pairs":{"spot":{"enabled":"LTC_BTC","available":"ETH_BTC,BTC_USD","assetEnabled":null}}}}`)
	out, err = v.UpgradeExchange(context.Background(), in)
	require.NoError(t, err)
	assert.Equal(t, string(exp), string(out), "Exchanges with assetEnabled null should enabled")
}

func TestVersion3Downgrade(t *testing.T) {
	t.Parallel()
	v := &Version3{}
	in := []byte(`{"name":"Wibble","currencyPairs":{"pairs":{"spot":{"enabled":"LTC_BTC","available":"ETH_BTC,BTC_USD","assetEnabled":true}}}}`)
	out, err := v.UpgradeExchange(context.Background(), bytes.Clone(in))
	require.NoError(t, err)
	assert.Equal(t, string(in), string(out), "Downgrade should leave the config alone")
}
