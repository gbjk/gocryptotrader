package v8

import (
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

type Subscription struct {
	Enable        bool           `json:"enabled"`
	Channel       string         `json:"channel,omitempty"`
	Pairs         any            `json:"pairs,omitempty"`
	Asset         asset.Item     `json:"asset,omitempty"`
	Params        map[string]any `json:"params,omitempty"`
	Interval      any            `json:"interval,omitempty"`
	Levels        int            `json:"levels,omitempty"`
	Authenticated bool           `json:"authenticated,omitempty"`
}
