package websocket

import (
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
)

// PingHandler container for ping handler settings
type PingHandler struct {
	UseGorillaHandler bool
	MessageType       int
	Message           []byte
	Delay             time.Duration
}

// FundingData defines funding data
type FundingData struct {
	Timestamp    time.Time
	CurrencyPair currency.Pair
	AssetType    asset.Item
	Exchange     string
	Amount       float64
	Rate         float64
	Period       int64
	Side         order.Side
}

// UnhandledMessageWarning defines a container for unhandled message warnings
type UnhandledMessageWarning struct {
	Message string
}

// Reporter interface groups observability functionality over Websocket request latency.
type Reporter interface {
	Latency(name string, message []byte, t time.Duration)
}

// ConnectionDown reports a connection that has stopped delivering. URL and
// MessageFilter identify the connection setup, and Subscriptions lists what it
// carried when it went down. Err is why it went down.
type ConnectionDown struct {
	URL           string
	MessageFilter any
	Subscriptions subscription.List
	Err           error
}

// ConnectionUp reports a connection that is delivering again. URL and
// MessageFilter identify the connection setup, and Subscriptions lists what it
// carries now. Downtime is how long it was down.
type ConnectionUp struct {
	URL           string
	MessageFilter any
	Subscriptions subscription.List
	Downtime      time.Duration
}
