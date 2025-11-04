package data

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thrasher-corp/gocryptotrader/backtester/common"
	"github.com/thrasher-corp/gocryptotrader/common/key"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

// HandlerHolder stores an event handler per exchange asset pair
type HandlerHolder struct {
	m    sync.Mutex
	data map[key.ExchangeAssetPair]Handler
}

// Holder interface dictates what a Data holder is expected to do
type Holder interface {
	SetDataForCurrency(string, asset.Item, currency.Pair, Handler) error
	GetAllData() ([]Handler, error)
	GetDataForCurrency(ev common.Event) (Handler, error)
	Reset() error
}

// Base is the base implementation of some interface functions
// where further specific functions are implemented in DataFromKline
type Base struct {
	m          sync.Mutex
	latest     Event
	stream     []Event
	offset     int64
	isLiveData bool
}

// Handler interface for Loading and Streaming Data
type Handler interface {
	Loader
	Streamer
	GetDetails() (string, asset.Item, currency.Pair, error)
	Reset() error
}

// Loader interface for Loading Data into backtest supported format
type Loader interface {
	Load() error
	AppendStream(s ...Event) error
}

// Streamer interface handles loading, parsing, distributing BackTest Data
type Streamer interface {
	Next() (Event, error)
	GetStream() (Events, error)
	History() (Events, error)
	Latest() (Event, error)
	List() (Events, error)
	IsLastEvent() (bool, error)
	Offset() (int64, error)

	StreamOpen() ([]decimal.Decimal, error)
	StreamHigh() ([]decimal.Decimal, error)
	StreamLow() ([]decimal.Decimal, error)
	StreamClose() ([]decimal.Decimal, error)
	StreamVol() ([]decimal.Decimal, error)

	HasDataAtTime(time.Time) (bool, error)
}

// Event interface used for loading and interacting with Data
type Event interface {
	common.Event
	GetUnderlyingPair() currency.Pair
	GetClosePrice() decimal.Decimal
	GetHighPrice() decimal.Decimal
	GetLowPrice() decimal.Decimal
	GetOpenPrice() decimal.Decimal
	GetVolume() decimal.Decimal
}

// Events allows for some common functions on a slice of events
type Events []Event
