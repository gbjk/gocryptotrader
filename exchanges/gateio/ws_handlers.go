package gateio

import (
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

type wsHandlerFunc func(*Gateio, asset.Item, *WSResponse) error

var wsHandlerFuncs = map[asset.Item]map[string]wsHandlerFunc{
	asset.Spot: {
		tickerChannel: (*Gateio).processTicker,
	},
}
