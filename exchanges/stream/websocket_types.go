package stream

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/exchanges/fill"
	"github.com/thrasher-corp/gocryptotrader/exchanges/protocol"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream/buffer"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
	"github.com/thrasher-corp/gocryptotrader/exchanges/trade"
)

// Websocket functionality list and state consts
const (
	WebsocketNotAuthenticatedUsingRest = "%v - Websocket not authenticated, using REST\n"
	Ping                               = "ping"
	Pong                               = "pong"
	UnhandledMessage                   = " - Unhandled websocket message: "
)

const (
	uninitialisedState uint32 = iota
	disconnectedState
	connectingState
	connectedState
)

// Websocket provides control over websocket connections and subscriptions
// TODO: Rename to manager
type Websocket struct {
	canUseAuthenticatedEndpoints atomic.Bool
	enabled                      atomic.Bool
	state                        atomic.Uint32
	verbose                      bool
	connectionMonitorRunning     atomic.Bool
	trafficTimeout               time.Duration
	connectionMonitorDelay       time.Duration
	proxyAddr                    string
	defaultURL                   string
	defaultURLAuth               string
	runningURL                   string
	runningURLAuth               string
	exchangeName                 string
	m                            sync.Mutex
	connector                    func(context.Context, *WebsocketConnection) error

	// connectionConfigs stores all *potential* connections for the exchange, organised within ConnectionWrapper structs.
	// Each ConnectionWrapper one connection (will be expanded soon) tailored for specific exchange functionalities or asset types.
	// TODO: Expand this to support multiple connections per ConnectionWrapper
	// For example, separate connections can be used for Spot, Margin, and Futures trading. This structure is especially useful
	// for exchanges that differentiate between trading pairs by using different connection endpoints or protocols for various asset classes.
	// If an exchange does not require such differentiation, all connections may be managed under a single ConnectionWrapper.
	connectionConfigs []*ConnectionSetup
	// connections holds a look up table for all connections to their corresponding ConnectionWrapper and subscription holder
	connections []*Connection

	subscriptions *subscription.Store

	// Subscriber function for exchange specific subscribe implementation
	Subscriber func(subscription.List) error
	// Subscriber function for exchange specific unsubscribe implementation
	Unsubscriber func(subscription.List) error
	// GenerateSubs function for exchange specific generating subscriptions from Features.Subscriptions, Pairs and Assets
	GenerateSubs func() (subscription.List, error)
	Authenticate func(ctx context.Context, conn Connection) error

	DataHandler chan interface{}
	ToRoutine   chan interface{}

	Match *Match

	ShutdownC chan struct{} // ShutdownC synchronises shutdown across routines
	Wg        sync.WaitGroup

	Orderbook             buffer.Orderbook // Orderbook is a local buffer of orderbooks
	Trade                 trade.Trade      // Trade is a notifier of occurring trades
	Fills                 fill.Fills       // Fills is a notifier of occurring fills
	TrafficAlert          chan struct{}    // trafficAlert monitors if there is a halt in traffic throughput
	ReadMessageErrors     chan error       // ReadMessageErrors will received all errors from ws.ReadMessage() and verify if its a disconnection
	features              *protocol.Features
	ExchangeLevelReporter Reporter                     // Latency reporter
	rateLimitDefinitions  request.RateLimitDefinitions // rate limiters shared between Websocket and REST connections for all potential endpoints
}

// TODO: Rename to ManagerSetup
type WebsocketSetup struct {
	ExchangeConfig        *config.Exchange
	Handler               MessageHandler
	Connector             func() error
	Subscriber            func(subscription.List) error
	Unsubscriber          func(subscription.List) error
	GenerateSubscriptions func() (subscription.List, error)
	Features              *protocol.Features
	OrderbookBufferConfig buffer.Config // Local orderbook buffer config values
	TradeFeed             bool
	FillsFeed             bool

	// RateLimitDefinitions contains the rate limiters shared between WebSocket and REST connections for all endpoints.
	// These rate limits take precedence over any rate limits specified in individual connection configurations.
	// If no connection-specific rate limit is provided and the endpoint does not match any of these definitions,
	// an error will be returned. However, if a connection configuration includes its own rate limit,
	// it will fall back to that configuration’s rate limit without raising an error.
	RateLimitDefinitions request.RateLimitDefinitions
}

// ConnectionSetup contains per connection configuration
type ConnectionSetup struct {
	URL                     string
	Authenticated           bool
	ConnectionLevelReporter Reporter
	ResponseCheckTimeout    time.Duration
	ResponseMaxLimit        time.Duration
	RateLimit               *request.RateLimiterWithWeight
	MessageFilter           any // MessageFilter defines the criteria used to match messages to a specific connection.
	Features                *protocol.Features
	MaxSubscriptions        int
}

// WebsocketConnection wraps a websocket.Conn with GCT fields required to send and receive websocket messages
type WebsocketConnection struct {
	*ConnectionSetup
	Verbose   bool
	connected int32

	// Gorilla websocket does not allow more than one goroutine to utilise
	// writes methods
	writeControl sync.Mutex

	// RateLimit is a rate limiter for the connection itself
	RateLimit *request.RateLimiterWithWeight
	// RateLimitDefinitions contains the rate limiters shared between WebSocket and REST connections for all
	// potential endpoints.
	RateLimitDefinitions request.RateLimitDefinitions

	ExchangeName         string
	URL                  string
	ProxyURL             string
	Wg                   *sync.WaitGroup
	UnderlyingConnection *websocket.Conn

	// shutdown synchronises shutdown event across routines associated with this connection only e.g. ping handler
	shutdown chan struct{}

	Match             *Match
	ResponseMaxLimit  time.Duration
	Traffic           chan struct{}
	readMessageErrors chan error
	messageFilter     any

	// bespokeGenerateMessageID is a function that returns a unique message ID
	// defined externally. This is used for exchanges that require a unique
	// message ID for each message sent.
	bespokeGenerateMessageID func(highPrecision bool) int64

	Reporter Reporter
}

type MessageHandler func(context.Context, Connection, []byte) error
