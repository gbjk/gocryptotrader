package versions

import (
	"context"
	"errors"
	"fmt"

	"github.com/thrasher-corp/gocryptotrader/common/convert"
	"github.com/thrasher-corp/gocryptotrader/config"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/log"
)

// Version1 is an Exchange config to fix assets structure
type Version1 struct {
}

var _ ExchangeVersion = &Version1{}

func init() {
	RegisterVersion(&Version1{})
}

// Exchanges returns just Huobi
func (v *Version1) Exchanges() []string { return []string{"*"} }

// UpgradeExchange will upgrade pair configuration
func (v *Version1) UpgradeExchange(ctx context.Context, e *config.Exchange) error {
	// Check if see if the new currency pairs format is empty and flesh it out if so
	if e.CurrencyPairs == nil {
		e.CurrencyPairs = new(currency.PairsManager)
		e.CurrencyPairs.Pairs = make(map[asset.Item]*currency.PairStore)

		if e.PairsLastUpdated != nil {
			e.CurrencyPairs.LastUpdated = *e.PairsLastUpdated
		}

		e.CurrencyPairs.ConfigFormat = e.ConfigCurrencyPairFormat
		e.CurrencyPairs.RequestFormat = e.RequestCurrencyPairFormat

		var availPairs, enabledPairs currency.Pairs
		if e.AvailablePairs != nil {
			availPairs = *e.AvailablePairs
		}

		if e.EnabledPairs != nil {
			enabledPairs = *e.EnabledPairs
		}

		e.CurrencyPairs.UseGlobalFormat = true
		err := e.CurrencyPairs.Store(asset.Spot, &currency.PairStore{
			AssetEnabled: convert.BoolPtr(true),
			Available:    availPairs,
			Enabled:      enabledPairs,
		})
		if err != nil {
			return err
		}

		// flush old values
		e.PairsLastUpdated = nil
		e.ConfigCurrencyPairFormat = nil
		e.RequestCurrencyPairFormat = nil
		e.AssetTypes = nil
		e.AvailablePairs = nil
		e.EnabledPairs = nil
	} else {
		if err := e.CurrencyPairs.SetDelimitersFromConfig(); err != nil {
			return fmt.Errorf("%s: %w", e.Name, err)
		}

		assets := e.CurrencyPairs.GetAssetTypes(false)
		if len(assets) == 0 {
			e.Enabled = false
			log.Warnf(log.ConfigMgr, "%s no assets found, disabling...", e.Name)
			return nil
		}

		var atLeastOne bool
		for index := range assets {
			err := e.CurrencyPairs.IsAssetEnabled(assets[index])
			if err != nil {
				if errors.Is(err, currency.ErrAssetIsNil) {
					// Checks if we have an old config without the ability to
					// enable disable the entire asset
					log.Warnf(log.ConfigMgr,
						"Exchange %s: upgrading config for asset type %s and setting enabled.\n",
						e.Name,
						assets[index])
					err = e.CurrencyPairs.SetAssetEnabled(assets[index], true)
					if err != nil {
						return err
					}
					atLeastOne = true
				}
				continue
			}
			atLeastOne = true
		}

		if !atLeastOne {
			// turn on an asset if all disabled
			log.Warnf(log.ConfigMgr,
				"%s assets disabled, turning on asset %s",
				e.Name,
				assets[0])

			err := e.CurrencyPairs.SetAssetEnabled(assets[0], true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Version1) DowngradeExchange(ctx context.Context, e *config.Exchange) error {
	return nil
}
