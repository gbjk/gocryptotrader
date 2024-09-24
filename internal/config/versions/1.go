package versions

import (
	"context"

	"github.com/thrasher-corp/gocryptotrader/config"
)

const exch = "Huobi"

// Version1 is an Exchange config upgrade for Huobi assets
type Version1 struct {
}

var _ ExchangeVersion = &Version1{}

func init() {
	RegisterVersion(&Version1{})
}

// Exchanges returns just Huobi

func (v *Version1) Exchanges() []string { return []string{"Huobi"} }

// UpgradeExchange will move Delivery Future contracts from Future asset into Coin-M asset configuratin
func (v *Version1) UpgradeExchange(ctx context.Context, e *config.Exchange) error {
	panic("called")
	return nil
}

func (v *Version1) DowngradeExchange(ctx context.Context, e *config.Exchange) error {
	return nil
}
