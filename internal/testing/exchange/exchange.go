package exchange

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/thrasher-corp/gocryptotrader/config"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/mock"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
	"github.com/thrasher-corp/gocryptotrader/exchanges/sharedtestvalues"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
)

func TestInstance(e exchange.IBotExchange) error {
	cfg := config.GetConfig()
	err := cfg.LoadConfig("../../testdata/configtest.json", true)
	if err != nil {
		return fmt.Errorf("LoadConfig() error: %w", err)
	}
	parts := strings.Split(fmt.Sprintf("%T", e), ".")
	if len(parts) != 2 {
		return errors.New("Unexpected parts splitting exchange type name")
	}
	eName := parts[1]
	exchConf, err := cfg.GetExchangeConfig(eName)
	if err != nil {
		return fmt.Errorf("GetExchangeConfig(`%s`) error: %w", eName, err)
	}
	e.SetDefaults()
	b := e.GetBase()
	b.Websocket = sharedtestvalues.NewTestWebsocket()
	err = e.Setup(exchConf)
	if err != nil {
		return fmt.Errorf("Setup() error: %w", err)
	}
	return nil
}

// httpMockFile is a consistent path under each exchange to find the mock server definitions
const httpMockFile = "testdata/http.json"

// mockHTTPInstance takes an existing Exchange instance and attaches it to a new http server
// It is expected to be run once,  since http requests do not often tangle with each other
func MockHTTPInstance(e exchange.IBotExchange) error {
	serverDetails, newClient, err := mock.NewVCRServer(httpMockFile)
	if err != nil {
		return fmt.Errorf("Mock server error %s", err)
	}
	b := e.GetBase()
	b.SkipAuthCheck = true
	err = b.SetHTTPClient(newClient)
	if err != nil {
		return fmt.Errorf("Mock server error %s", err)
	}
	endpointMap := b.API.Endpoints.GetURLMap()
	for k := range endpointMap {
		err = b.API.Endpoints.SetRunning(k, serverDetails)
		if err != nil {
			return fmt.Errorf("Mock server error %s", err)
		}
	}
	request.MaxRequestJobs = 100
	log.Printf(sharedtestvalues.MockTesting, e.GetName())

	return nil
}

var upgrader = websocket.Upgrader{}

type wsMockFunc func(msg []byte, w *websocket.Conn) error

// mockWSInstancecreates a new Exchange instance, attaches a new HTTP server to it, and a new mock WS instance
// It is expected to be run from any WS tests which need a specific response function
func MockWSInstance[T any, PT interface {
	*T
	exchange.IBotExchange
}](tb testing.TB, m wsMockFunc) *T {
	tb.Helper()

	e := PT(new(T))

	err := TestInstance(e)
	if !assert.NoError(tb, err, "testInstance should not error") {
		tb.FailNow()
	}

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { wsMockWrapper(tb, w, r, m) }))
	if !assert.NoError(tb, err, "testInstance should not error") {
		tb.FailNow()
	}

	b := e.GetBase()
	b.SkipAuthCheck = true
	err = b.API.Endpoints.SetRunning("RestSpotURL", s.URL)
	assert.NoError(tb, err, "SetRunning should not error")
	for _, auth := range []bool{true, false} {
		err = b.Websocket.SetWebsocketURL("ws"+strings.TrimPrefix(s.URL, "http"), auth, true)
		assert.NoErrorf(tb, err, "SetWebsocketURL should not error for auth: %v", auth)
	}

	b.Features.Subscriptions = []*subscription.Subscription{}
	err = b.Websocket.Connect()
	assert.NoError(tb, err, "Connect should not error")

	return e
}

func wsMockWrapper(tb testing.TB, w http.ResponseWriter, r *http.Request, m wsMockFunc) {
	tb.Helper()
	if strings.Contains(r.URL.Path, "GetWebSocketsToken") {
		_, err := w.Write([]byte(`{"result":{"token":"mockAuth"}}`))
		assert.NoError(tb, err, "Write should not error")
		return
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if !assert.NoError(tb, err, "Upgrade connection should not error") {
		return
	}
	defer c.Close()
	for {
		_, p, err := c.ReadMessage()
		if !assert.NoError(tb, err, "ReadMessage should not error") {
			return
		}

		err = m(p, c)
		if !assert.NoError(tb, err, "WS Mock Function should not error") {
			return
		}
	}
}

var setupWsMutex sync.Mutex
var setupWsOnce = make(map[exchange.IBotExchange]bool)

// SetupWs is a helper function to connect both auth and normal websockets
// It will skip the test if websockets are not enabled
// It's up to the test to skip if it requires creds, though
func SetupWs(tb testing.TB, e exchange.IBotExchange) {
	tb.Helper()

	setupWsMutex.Lock()
	defer setupWsMutex.Unlock()

	if _, ok := setupWsOnce[e]; ok {
		return
	}

	b := e.GetBase()
	if !b.Websocket.IsEnabled() {
		tb.Skip("Websocket not enabled")
	}
	if b.Websocket.IsConnected() {
		return
	}
	err := b.Websocket.Connect()
	if !assert.NoError(tb, err, "WsConnect should not error") {
		tb.FailNow()
	}

	setupWsOnce[e] = true
}
