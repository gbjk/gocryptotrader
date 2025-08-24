package websocket

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/exchange/subscription"
	"github.com/thrasher-corp/gocryptotrader/exchange/websocket/buffer"
	"github.com/thrasher-corp/gocryptotrader/exchanges/fill"
	"github.com/thrasher-corp/gocryptotrader/exchanges/protocol"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/exchanges/trade"
	"github.com/thrasher-corp/gocryptotrader/log"
)

// Public websocket errors
var (
	ErrWebsocketNotEnabled  = errors.New("websocket manager not enabled")
	ErrAlreadyDisabled      = errors.New("websocket manager already disabled")
	ErrNotRunning           = errors.New("websocket manager is not running")
	ErrNotConnected         = errors.New("websocket is not connected")
	ErrSignatureTimeout     = errors.New("websocket timeout waiting for response with signature")
	ErrRequestRouteNotFound = errors.New("request route not found")
	ErrSignatureNotSet      = errors.New("signature not set")
	ErrInvalidMessageFilter = errors.New("invalid message filter")
)

// Private websocket errors
var (
	errWebsocketAlreadyInitialised          = errors.New("websocket already initialised")
	errWebsocketAlreadyEnabled              = errors.New("websocket already enabled")
	errAlreadyRunning                       = errors.New("websocket already running")
	errInvalidWebsocketURL                  = errors.New("invalid websocket url")
	errExchangeConfigNameEmpty              = errors.New("exchange config name empty")
	errInvalidTrafficTimeout                = errors.New("invalid traffic timeout")
	errTrafficAlertNil                      = errors.New("traffic alert is nil")
	errNoSubscriber                         = errors.New("websocket subscriber function needs to be set")
	errNoUnsubscriber                       = errors.New("websocket unsubscriber functionality allowed but unsubscriber function not set")
	errNoDataHandler                        = errors.New("websocket data handler not set")
	errReadMessageErrorsNil                 = errors.New("read message errors is nil")
	errWebsocketSubscriptionsGeneratorUnset = errors.New("websocket subscriptions generator function needs to be set")
	errInvalidMaxSubscriptions              = errors.New("max subscriptions cannot be less than 0")
	errSameProxyAddress                     = errors.New("cannot set proxy address to the same address")
	errNoConnectFunc                        = errors.New("websocket manager connect func not set")
	errAlreadyConnected                     = errors.New("websocket already connected")
	errCannotShutdown                       = errors.New("websocket cannot shutdown")
	errAlreadyReconnecting                  = errors.New("websocket in the process of reconnection")
	errConnSetup                            = errors.New("error in connection setup")
	errNoPendingConnections                 = errors.New("no pending connections, call SetupNewConnection first")
	errDuplicateConnectionSetup             = errors.New("duplicate connection setup")
	errCannotChangeConnectionURL            = errors.New("cannot change connection URL when using multi connection management")
	errCannotObtainOutboundConnection       = errors.New("cannot obtain outbound connection")
)

// Websocket functionality list and state consts
const (
	UnhandledMessage = " - Unhandled websocket message: "
	jobBuffer        = 5000
)

const (
	uninitialisedState uint32 = iota
	dormantState
	runningState
)

// Manager provides connection and subscription management and routing
type Manager struct {
	enabled                      atomic.Bool
	state                        atomic.Uint32
	verbose                      bool
	canUseAuthenticatedEndpoints atomic.Bool // TODO: GBJK - Do we keep and maintain this flag, or work it out dynamically
	connectionMonitorRunning     atomic.Bool
	trafficTimeout               time.Duration // TODO: GBJK - This needs to move to connection
	connectionMonitorDelay       time.Duration
	proxyAddr                    string
	defaultURL                   string
	defaultURLAuth               string
	runningURL                   string
	runningURLAuth               string
	exchangeName                 string
	features                     *protocol.Features
	m                            sync.Mutex
	connectionConfigs            []*ConnectionSetup
	connections                  []*Connection // TODO: GBJK Not sure if there will be a use for a plain slice like this
	subscriptions                *subscription.Store
	connector                    func(context.Context, *Connection) error
	rateLimitDefinitions         request.RateLimitDefinitions // rate limiters shared between Websocket and REST connections
	Subscriber                   func(subscription.Connection, subscription.List) error
	Unsubscriber                 func(subscription.Connection, subscription.List) error
	GenerateSubs                 func() (subscription.List, error)
	Authenticate                 func(ctx context.Context, conn *Connection) error // TODO: GBJK - This concept is blown, isn't it? Or can we fan out inside the exchanges
	DataHandler                  chan any
	ToRoutine                    chan any
	Match                        *Match // TODO: GBJK - I think we don't have a global Match available
	ShutdownC                    chan struct{}
	Wg                           sync.WaitGroup
	Orderbook                    buffer.Orderbook
	Trade                        trade.Trade // Trade is a notifier for trades
	Fills                        fill.Fills  // Fills is a notifier for fills
	TrafficAlert                 chan struct{}
	ReadMessageErrors            chan error
	ExchangeLevelReporter        Reporter // Latency reporter
	// TODO: GBJK - Think this is a connection setup thing, because Auth and non-auth might have different values
	// That does imply duplicating it a lot, though
	MaxSubscriptionsPerConnection int
}

// ManagerSetup defines variables for setting up a websocket manager
type ManagerSetup struct {
	ExchangeConfig *config.Exchange
	Handler        MessageHandler

	Connector             func() error
	Subscriber            func(subscription.List) error
	Unsubscriber          func(subscription.List) error
	GenerateSubscriptions func() (subscription.List, error)
	Features              *protocol.Features
	OrderbookBufferConfig buffer.Config

	TradeFeed bool
	FillsFeed bool

	MaxWebsocketSubscriptionsPerConnection int

	// RateLimitDefinitions contains the rate limiters shared between WebSocket and REST connections for all endpoints.
	// These rate limits take precedence over any rate limits specified in individual connection configurations.
	// If no connection-specific rate limit is provided and the endpoint does not match any of these definitions,
	// an error will be returned. However, if a connection configuration includes its own rate limit,
	// it will fall back to that configuration’s rate limit without raising an error.
	RateLimitDefinitions request.RateLimitDefinitions
}

var globalReporter Reporter

// SetupGlobalReporter sets a reporter interface to be used
// for all exchange requests
func SetupGlobalReporter(r Reporter) {
	globalReporter = r
}

// NewManager initialises the websocket struct
func NewManager() *Manager {
	return &Manager{
		DataHandler:  make(chan any, jobBuffer),
		ToRoutine:    make(chan any, jobBuffer),
		ShutdownC:    make(chan struct{}),
		TrafficAlert: make(chan struct{}, 1),
		// ReadMessageErrors is buffered for an edge case when `Connect` fails
		// after subscriptions are made but before the connectionMonitor has
		// started. This allows the error to be read and handled in the
		// connectionMonitor and start a connection cycle again.
		ReadMessageErrors: make(chan error, 1),
		Match:             NewMatch(),
		subscriptions:     subscription.NewStore(),
		features:          &protocol.Features{},
		Orderbook:         buffer.Orderbook{},
		connections:       []*Connection{},
	}
}

// Setup sets main variables for websocket connection
func (m *Manager) Setup(s *ManagerSetup) error {
	if err := common.NilGuard(m, s); err != nil {
		return err
	}

	m.m.Lock()
	defer m.m.Unlock()

	if m.IsInitialised() {
		return fmt.Errorf("%s %w", m.exchangeName, errWebsocketAlreadyInitialised)
	}

	if err := s.validate(); err != nil {
		return fmt.Errorf("%s %w", m.exchangeName, err)
	}

	m.exchangeName = s.ExchangeConfig.Name
	m.verbose = s.ExchangeConfig.Verbose
	m.features = s.Features
	m.rateLimitDefinitions = s.RateLimitDefinitions
	m.connectionMonitorDelay = s.ExchangeConfig.ConnectionMonitorDelay
	m.trafficTimeout = s.ExchangeConfig.WebsocketTrafficTimeout

	if m.connectionMonitorDelay <= 0 {
		m.connectionMonitorDelay = config.DefaultConnectionMonitorDelay
	}

	if m.trafficTimeout < time.Second {
		return fmt.Errorf("%s %w cannot be less than %s", m.exchangeName, errInvalidTrafficTimeout, time.Second)
	}

	m.setState(dormantState)
	m.setEnabled(s.ExchangeConfig.Features.Enabled.Websocket)
	m.SetCanUseAuthenticatedEndpoints(s.ExchangeConfig.API.AuthenticatedWebsocketSupport)

	if err := m.Orderbook.Setup(s.ExchangeConfig, &s.OrderbookBufferConfig, m.DataHandler); err != nil {
		return err
	}

	m.Trade.Setup(s.TradeFeed, m.DataHandler)
	m.Fills.Setup(s.FillsFeed, m.DataHandler)

	return nil
}

// SetupNewConnection sets up an auth or unauth streaming connection
func (m *Manager) SetupNewConnection(c *ConnectionSetup) error {
	if err := common.NilGuard(m, c); err != nil {
		return err
	}

	if c.URL == "" {
		return fmt.Errorf("%w: %w", errConnSetup, errInvalidWebsocketURL)
	}

	if c.ConnectionLevelReporter == nil {
		c.ConnectionLevelReporter = m.ExchangeLevelReporter
	}
	if c.ConnectionLevelReporter == nil {
		c.ConnectionLevelReporter = globalReporter
	}

	if c.MessageFilter != nil && !reflect.TypeOf(c.MessageFilter).Comparable() {
		return ErrInvalidMessageFilter
	}

	if slices.ContainsFunc(m.connectionConfigs, func(b *ConnectionSetup) bool { return b.URL == c.URL && b.MessageFilter == c.MessageFilter }) {
		return fmt.Errorf("%w: %w", errConnSetup, errDuplicateConnectionSetup)
	}

	m.connectionConfigs = append(m.connectionConfigs, c)

	return nil
}

// getConnectionFromSetup returns a websocket connection from a setup
// configuration. This is used for setting up new connections on the fly.
func (m *Manager) getConnectionFromSetup(c *ConnectionSetup) *connection {
	connectionURL := m.GetWebsocketURL()
	if c.URL != "" {
		connectionURL = c.URL
	}
	return &connection{
		ExchangeName:         m.exchangeName,
		URL:                  connectionURL,
		ProxyURL:             m.GetProxyAddress(),
		Verbose:              m.verbose,
		ResponseMaxLimit:     c.ResponseMaxLimit,
		Traffic:              m.TrafficAlert,
		readMessageErrors:    m.ReadMessageErrors,
		shutdown:             m.ShutdownC,
		Wg:                   &m.Wg,
		Match:                NewMatch(),
		RateLimit:            c.RateLimit,
		Reporter:             c.ConnectionLevelReporter,
		RateLimitDefinitions: m.rateLimitDefinitions,
	}
}

// Start begins the Manager's process of checking for subscriptions and making new connections
func (m *Manager) Start() error {
	m.m.Lock()
	defer m.m.Unlock()

	if !m.IsEnabled() {
		return ErrWebsocketNotEnabled
	}

	if m.IsRunning() {
		return fmt.Errorf("%v %w", m.exchangeName, errAlreadyRunning)
	}

	m.Wg.Add(2)
	go m.monitorFrame(&m.Wg, m.monitorData)
	go m.monitorFrame(&m.Wg, m.monitorTraffic)

	m.SyncSubscriptions()

	m.setState(runningState)

	if m.connectionMonitorRunning.CompareAndSwap(false, true) {
		// This oversees all connections and does not need to be part of wait group management.
		go m.monitorFrame(nil, m.monitorConnection)
	}

	return nil
}

func (m *Manager) connect(ctx context.Context, s *ConnectionSetup) error {
	conn := m.getConnectionFromSetup(s)

	if err := m.connector(ctx, conn); err != nil {
		return fmt.Errorf("%v Error connecting %w", m.exchangeName, err)
	}

	if !conn.IsConnected() {
		return fmt.Errorf("%s websocket: [conn:%d] [URL:%s] failed to connect", m.exchangeName, i+1, conn.URL)
	}

	m.connectionPool[conn] = w.connectionConfigs[i]
	m.connectionConfigs[i].Connection = conn

	m.Wg.Add(1)
	go m.Reader(context.TODO(), conn, w.connectionConfigs[i].Setup.Handler)

	if m.connectionConfigs[i].Setup.Authenticate != nil && w.CanUseAuthenticatedEndpoints() {
		err = m.connectionConfigs[i].Setup.Authenticate(context.TODO(), conn)
		if err != nil {
			// Opted to not fail entirely here for POC. This should be
			// revisited and handled more gracefully.
			log.Errorf(log.WebsocketMgr, "%s websocket: [conn:%d] [URL:%s] failed to authenticate %v", m.exchangeName, i+1, conn.URL, err)
		}
	}

	err = m.connectionConfigs[i].Setup.Subscriber(context.TODO(), conn, subs)
	if err != nil {
		multiConnectFatalError = fmt.Errorf("%v Error subscribing %w", m.exchangeName, err)
		break
	}

	if m.verbose {
		log.Debugf(log.WebsocketMgr, "%s websocket: [conn:%d] [URL:%s] connected. [Subscribed: %d]",
			m.exchangeName,
			i+1,
			conn.URL,
			len(subs))
	}
	return nil
}

// Disable disables the exchange websocket protocol
// Note that connectionMonitor will be responsible for shutting down the websocket after disabling
func (m *Manager) Disable() error {
	if !m.IsEnabled() {
		return fmt.Errorf("%s %w", m.exchangeName, ErrAlreadyDisabled)
	}

	m.setEnabled(false)
	return nil
}

// Enable enables the exchange websocket protocol
func (m *Manager) Enable() error {
	if m.IsConnected() || m.IsEnabled() {
		return fmt.Errorf("%s %w", m.exchangeName, errWebsocketAlreadyEnabled)
	}

	m.setEnabled(true)
	return m.Connect()
}

// Shutdown attempts to shut down a websocket connection and associated routines
// by using a package defined shutdown function
func (m *Manager) Shutdown() error {
	m.m.Lock()
	defer m.m.Unlock()
	return m.shutdown()
}

func (m *Manager) shutdown() error {
	if !m.IsConnected() {
		return fmt.Errorf("%v %w: %w", m.exchangeName, errCannotShutdown, ErrNotConnected)
	}

	// TODO: Interrupt connection and or close connection when it is re-established.
	if m.IsConnecting() {
		return fmt.Errorf("%v %w: %w ", m.exchangeName, errCannotShutdown, errAlreadyReconnecting)
	}

	if m.verbose {
		log.Debugf(log.WebsocketMgr, "%v websocket: shutting down websocket", m.exchangeName)
	}

	defer m.Orderbook.FlushBuffer()

	// During the shutdown process, all errors are treated as non-fatal to avoid issues when the connection has already
	// been closed. In such cases, attempting to close the connection may result in a
	// "failed to send closeNotify alert (but connection was closed anyway)" error. Treating these errors as non-fatal
	// prevents the shutdown process from being interrupted, which could otherwise trigger a continuous traffic monitor
	// cycle and potentially block the initiation of a new connection.
	var nonFatalCloseConnectionErrors error

	// Shutdown managed connections
	for x := range m.connectionConfigs {
		if m.connectionConfigs[x].connection != nil {
			if err := m.connectionConfigs[x].connection.Shutdown(); err != nil {
				nonFatalCloseConnectionErrors = common.AppendError(nonFatalCloseConnectionErrors, err)
			}
			m.connectionConfigs[x].connection = nil
			// Flush any subscriptions from last connection across any managed connections
			m.connectionConfigs[x].subscriptions.Clear()
		}
	}
	// Clean map of old connections
	clear(m.connections)

	if m.Conn != nil {
		if err := m.Conn.Shutdown(); err != nil {
			nonFatalCloseConnectionErrors = common.AppendError(nonFatalCloseConnectionErrors, err)
		}
	}
	if m.AuthConn != nil {
		if err := m.AuthConn.Shutdown(); err != nil {
			nonFatalCloseConnectionErrors = common.AppendError(nonFatalCloseConnectionErrors, err)
		}
	}
	// flush any subscriptions from last connection if needed
	m.subscriptions.Clear()

	m.setState(disconnectedState)

	close(m.ShutdownC)
	m.Wg.Wait()
	m.ShutdownC = make(chan struct{})
	if m.verbose {
		log.Debugf(log.WebsocketMgr, "%v websocket: completed websocket shutdown", m.exchangeName)
	}

	// Drain residual error in the single buffered channel, this mitigates
	// the cycle when `Connect` is called again and the connectionMonitor
	// starts but there is an old error in the channel.
	drain(m.ReadMessageErrors)

	if nonFatalCloseConnectionErrors != nil {
		log.Warnf(log.WebsocketMgr, "%v websocket: shutdown error: %v", m.exchangeName, nonFatalCloseConnectionErrors)
	}

	return nil
}

func (m *Manager) setState(s uint32) {
	m.state.Store(s)
}

// IsInitialised returns whether the manager has been Setup() already
func (m *Manager) IsInitialised() bool {
	return m.state.Load() != uninitialisedState
}

// IsConnected returns whether the manager is Running
func (m *Manager) IsRunning() bool {
	return m.state.Load() == runningState
}

func (m *Manager) setEnabled(b bool) {
	m.enabled.Store(b)
}

// IsEnabled returns whether the websocket is enabled
func (m *Manager) IsEnabled() bool {
	return m.enabled.Load()
}

// CanUseAuthenticatedWebsocketForWrapper Handles a common check to
// verify whether a wrapper can use an authenticated websocket endpoint
func (m *Manager) CanUseAuthenticatedWebsocketForWrapper() bool {
	if m.IsConnected() {
		if m.CanUseAuthenticatedEndpoints() {
			return true
		}
		log.Infof(log.WebsocketMgr, "%v - Websocket not authenticated, using REST\n", m.exchangeName)
	}
	return false
}

// SetWebsocketURL sets websocket URL and can refresh underlying connections
func (m *Manager) SetWebsocketURL(url string, auth, reconnect bool) error {
	if m.useMultiConnectionManagement {
		// TODO: Add functionality for multi-connection management to change URL
		return fmt.Errorf("%s: %w", m.exchangeName, errCannotChangeConnectionURL)
	}
	defaultVals := url == "" || url == config.WebsocketURLNonDefaultMessage
	if auth {
		if defaultVals {
			url = m.defaultURLAuth
		}

		err := checkWebsocketURL(url)
		if err != nil {
			return err
		}
		m.runningURLAuth = url

		if m.verbose {
			log.Debugf(log.WebsocketMgr, "%s websocket: setting authenticated websocket URL: %s\n", m.exchangeName, url)
		}

		if m.AuthConn != nil {
			m.AuthConn.SetURL(url)
		}
	} else {
		if defaultVals {
			url = m.defaultURL
		}
		err := checkWebsocketURL(url)
		if err != nil {
			return err
		}
		m.runningURL = url

		if m.verbose {
			log.Debugf(log.WebsocketMgr, "%s websocket: setting unauthenticated websocket URL: %s\n", m.exchangeName, url)
		}

		if m.Conn != nil {
			m.Conn.SetURL(url)
		}
	}

	if m.IsConnected() && reconnect {
		log.Debugf(log.WebsocketMgr, "%s websocket: flushing websocket connection to %s\n", m.exchangeName, url)
		return m.Shutdown()
	}
	return nil
}

// GetWebsocketURL returns the running websocket URL
func (m *Manager) GetWebsocketURL() string {
	return m.runningURL
}

// SetProxyAddress sets websocket proxy address
func (m *Manager) SetProxyAddress(proxyAddr string) error {
	m.m.Lock()
	defer m.m.Unlock()
	if proxyAddr != "" {
		if _, err := url.ParseRequestURI(proxyAddr); err != nil {
			return fmt.Errorf("%v websocket: cannot set proxy address: %w", m.exchangeName, err)
		}

		if m.proxyAddr == proxyAddr {
			return fmt.Errorf("%v websocket: %w '%v'", m.exchangeName, errSameProxyAddress, m.proxyAddr)
		}

		log.Debugf(log.ExchangeSys, "%s websocket: setting websocket proxy: %s", m.exchangeName, proxyAddr)
	} else {
		log.Debugf(log.ExchangeSys, "%s websocket: removing websocket proxy", m.exchangeName)
	}

	for _, wrapper := range m.connectionConfigs {
		if wrapper.connection != nil {
			wrapper.connection.SetProxy(proxyAddr)
		}
	}
	if m.Conn != nil {
		m.Conn.SetProxy(proxyAddr)
	}
	if m.AuthConn != nil {
		m.AuthConn.SetProxy(proxyAddr)
	}

	m.proxyAddr = proxyAddr

	if !m.IsConnected() {
		return nil
	}
	if err := m.shutdown(); err != nil {
		return err
	}
	return m.connect()
}

// GetProxyAddress returns the current websocket proxy
func (m *Manager) GetProxyAddress() string {
	return m.proxyAddr
}

// GetName returns exchange name
func (m *Manager) GetName() string {
	return m.exchangeName
}

// SetCanUseAuthenticatedEndpoints sets canUseAuthenticatedEndpoints val in a thread safe manner
func (m *Manager) SetCanUseAuthenticatedEndpoints(b bool) {
	m.canUseAuthenticatedEndpoints.Store(b)
}

// CanUseAuthenticatedEndpoints gets canUseAuthenticatedEndpoints val in a thread safe manner
func (m *Manager) CanUseAuthenticatedEndpoints() bool {
	return m.canUseAuthenticatedEndpoints.Load()
}

// checkWebsocketURL checks for a valid websocket url
func checkWebsocketURL(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}
	if u.Scheme != "ws" && u.Scheme != "wss" {
		return fmt.Errorf("cannot set %w %s", errInvalidWebsocketURL, s)
	}
	return nil
}

// Reader reads and handles data from a specific connection
func (m *Manager) Reader(ctx context.Context, conn Connection, handler func(ctx context.Context, conn Connection, message []byte) error) {
	defer m.Wg.Done()
	for {
		resp := conn.ReadMessage()
		if resp.Raw == nil {
			return // Connection has been closed
		}
		if err := handler(ctx, conn, resp.Raw); err != nil {
			m.DataHandler <- fmt.Errorf("connection URL:[%v] error: %w", conn.GetURL(), err)
		}
	}
}

func drain(ch <-chan error) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

// ClosureFrame is a closure function that wraps monitoring variables with observer, if the return is true the frame will exit
type ClosureFrame func() func() bool

// monitorFrame monitors a specific websocket component or critical system. It will exit if the observer returns true
// This is used for monitoring data throughput, connection status and other critical websocket components. The waitgroup
// is optional and is used to signal when the monitor has finished.
func (m *Manager) monitorFrame(wg *sync.WaitGroup, fn ClosureFrame) {
	if wg != nil {
		defer m.Wg.Done()
	}
	observe := fn()
	for {
		if observe() {
			return
		}
	}
}

// monitorData monitors data throughput and logs if there is a back log of data
func (m *Manager) monitorData() func() bool {
	dropped := 0
	return func() bool { return m.observeData(&dropped) }
}

// observeData observes data throughput and logs if there is a back log of data
func (m *Manager) observeData(dropped *int) (exit bool) {
	select {
	case <-m.ShutdownC:
		return true
	case d := <-m.DataHandler:
		select {
		case m.ToRoutine <- d:
			if *dropped != 0 {
				log.Infof(log.WebsocketMgr, "%s exchange websocket ToRoutine channel buffer recovered; %d messages were dropped", m.exchangeName, dropped)
				*dropped = 0
			}
		default:
			if *dropped == 0 {
				// If this becomes prone to flapping we could drain the buffer, but that's extreme and we'd like to avoid it if possible
				log.Warnf(log.WebsocketMgr, "%s exchange websocket ToRoutine channel buffer full; dropping messages", m.exchangeName)
			}
			*dropped++
		}
		return false
	}
}

// monitorConnection monitors the connection and attempts to reconnect if the connection is lost
func (m *Manager) monitorConnection() func() bool {
	timer := time.NewTimer(m.connectionMonitorDelay)
	return func() bool { return m.observeConnection(timer) }
}

// observeConnection observes the connection and attempts to reconnect if the connection is lost
func (m *Manager) observeConnection(t *time.Timer) (exit bool) {
	select {
	case err := <-m.ReadMessageErrors:
		if errors.Is(err, errConnectionFault) {
			log.Warnf(log.WebsocketMgr, "%v websocket has been disconnected. Reason: %v", m.exchangeName, err)
			if m.IsConnected() {
				if shutdownErr := m.Shutdown(); shutdownErr != nil {
					log.Errorf(log.WebsocketMgr, "%v websocket: connectionMonitor shutdown err: %s", m.exchangeName, shutdownErr)
				}
			}
		}
		// Speedier reconnection, instead of waiting for the next cycle.
		if m.IsEnabled() && (!m.IsConnected() && !m.IsConnecting()) {
			if connectErr := m.Connect(); connectErr != nil {
				log.Errorln(log.WebsocketMgr, connectErr)
			}
		}
		m.DataHandler <- err // hand over the error to the data handler (shutdown and reconnection is priority)
	case <-t.C:
		if m.verbose {
			log.Debugf(log.WebsocketMgr, "%v websocket: running connection monitor cycle", m.exchangeName)
		}
		if !m.IsEnabled() {
			if m.verbose {
				log.Debugf(log.WebsocketMgr, "%v websocket: connectionMonitor - websocket disabled, shutting down", m.exchangeName)
			}
			if m.IsConnected() {
				if err := m.Shutdown(); err != nil {
					log.Errorln(log.WebsocketMgr, err)
				}
			}
			if m.verbose {
				log.Debugf(log.WebsocketMgr, "%v websocket: connection monitor exiting", m.exchangeName)
			}
			t.Stop()
			m.connectionMonitorRunning.Store(false)
			return true
		}
		if !m.IsConnecting() && !m.IsConnected() {
			err := m.Connect()
			if err != nil {
				log.Errorln(log.WebsocketMgr, err)
			}
		}
		t.Reset(m.connectionMonitorDelay)
	}
	return false
}

// monitorTraffic monitors to see if there has been traffic within the trafficTimeout time window. If there is no traffic
// the connection is shutdown and will be reconnected by the connectionMonitor routine.
func (m *Manager) monitorTraffic() func() bool {
	timer := time.NewTimer(m.trafficTimeout)
	return func() bool { return m.observeTraffic(timer) }
}

func (m *Manager) observeTraffic(t *time.Timer) bool {
	select {
	case <-m.ShutdownC:
		if m.verbose {
			log.Debugf(log.WebsocketMgr, "%v websocket: trafficMonitor shutdown message received", m.exchangeName)
		}
	case <-t.C:
		if m.IsConnecting() || signalReceived(m.TrafficAlert) {
			t.Reset(m.trafficTimeout)
			return false
		}
		if m.verbose {
			log.Warnf(log.WebsocketMgr, "%v websocket: has not received a traffic alert in %v. Reconnecting", m.exchangeName, m.trafficTimeout)
		}
		if m.IsConnected() {
			go func() { // Without this the m.Shutdown() call below will deadlock
				if err := m.Shutdown(); err != nil {
					log.Errorf(log.WebsocketMgr, "%v websocket: trafficMonitor shutdown err: %s", m.exchangeName, err)
				}
			}()
		}
	}
	t.Stop()
	return true
}

// signalReceived checks if a signal has been received, this also clears the signal.
func signalReceived(ch chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

// GetConnection returns a connection by message filter (defined in exchange package _wrapper.go websocket connection)
// for request and response handling in a multi connection context.
func (m *Manager) GetConnection(messageFilter any) (Connection, error) {
	if err := common.NilGuard(m); err != nil {
		return nil, err
	}
	if messageFilter == nil {
		return nil, fmt.Errorf("%w: messageFilter", common.ErrNilPointer)
	}

	m.m.Lock()
	defer m.m.Unlock()

	if !m.useMultiConnectionManagement {
		return nil, fmt.Errorf("%s: multi connection management not enabled %w please use exported Conn and AuthConn fields", m.exchangeName, errCannotObtainOutboundConnection)
	}

	if !m.IsConnected() {
		return nil, ErrNotConnected
	}

	for _, wrapper := range m.connectionConfigs {
		if wrapper.setup.MessageFilter == messageFilter {
			if wrapper.connection == nil {
				return nil, fmt.Errorf("%s: %s %w associated with message filter: '%v'", m.exchangeName, wrapper.setup.URL, ErrNotConnected, messageFilter)
			}
			return wrapper.connection, nil
		}
	}

	return nil, fmt.Errorf("%s: %w associated with message filter: '%v'", m.exchangeName, ErrRequestRouteNotFound, messageFilter)
}

func (s *ManagerSetup) validate() error {
	if s.Connector == nil {
		return errNoConnectFunc
	}
	if s.Handler == nil {
		return errNoDataHandler
	}
	if s.Subscriber == nil {
		return errNoSubscriber
	}
	if s.Subscriber == nil {
		return errNoSubscriber
	}
	if s.ExchangeConfig == nil {
		return fmt.Errorf("%w: ManagerSetup.ExchangeConfig", common.ErrNilPointer)
	}
	if s.ExchangeConfig.Name == "" {
		return errExchangeConfigNameEmpty
	}
	if s.Features == nil {
		return fmt.Errorf("%w: ManagerSetup.Features", common.ErrNilPointer)
	}
	if s.ExchangeConfig.Features == nil {
		return fmt.Errorf("%w: ManagerSetup.ExchangeConfig.Features", common.ErrNilPointer)
	}
	if s.ExchangeConfig.WebsocketTrafficTimeout < time.Second {
		return fmt.Errorf("%w: cannot be less than %s", errInvalidTrafficTimeout, time.Second)
	}
	return nil
}
