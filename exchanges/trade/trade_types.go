package trade

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
)

const DefaultSaveInterval = time.Second * 15

var (
	processor Processor
	// ErrNoTradesSupplied is returned when an attempt is made to process trades, but is an empty slice
	ErrNoTradesSupplied = errors.New("no trades supplied")
)

// Trade defines trade data
type Trade struct {
	ID           uuid.UUID `json:"ID,omitempty"`
	TID          string
	Exchange     string
	CurrencyPair currency.Pair
	AssetType    asset.Item
	Side         order.Side
	Price        float64
	Amount       float64
	Timestamp    time.Time
}

// Processor used for processing trade data in batches and saving them to the database
type Processor struct {
	started                 atomic.Bool
	bufferProcessorInterval time.Duration
	queue                   chan *Trade
	buffer                  []*Trade
}

// ByDate sorts trades by date ascending
type ByDate []Trade

func (b ByDate) Len() int {
	return len(b)
}

func (b ByDate) Less(i, j int) bool {
	return b[i].Timestamp.Before(b[j].Timestamp)
}

func (b ByDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
