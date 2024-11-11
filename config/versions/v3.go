package versions

import (
	"bytes"
	"context"
	"fmt"

	"github.com/buger/jsonparser"
	"github.com/thrasher-corp/gocryptotrader/common"
)

// Version3 implements ExchangeVersion
type Version3 struct {
}

func init() {
	Manager.registerVersion(3, &Version3{})
}

// Exchanges returns all exchanges: "*"
func (v *Version3) Exchanges() []string { return []string{"*"} }

// UpgradeExchange sets AssetEnabed: true for any exchange missing it
func (v *Version3) UpgradeExchange(ctx context.Context, e []byte) ([]byte, error) {
	name, err := jsonparser.GetString(e, "name")
	if err != nil {
		return e, fmt.Errorf("%w `name`: %w", common.ErrGettingField, err)
	}
	cb := func(k []byte, v []byte, _ jsonparser.ValueType, _ int) error {
		if _, err := jsonparser.GetBoolean(v, "assetEnabled"); err != nil {
			fmt.Printf("Exchange %s: Setting asset %s enabled\n", name, k)
			e, err = jsonparser.Set(e, []byte(`true`), "currencyPairs", "pairs", string(k), "assetEnabled")
			return err
		}
		return nil
	}
	err = jsonparser.ObjectEach(bytes.Clone(e), cb, "currencyPairs", "pairs")
	return e, err
}

// DowngradeExchange doesn't do anything for this version, because it's a lossy downgrade to disable all assets
func (v *Version3) DowngradeExchange(ctx context.Context, e []byte) ([]byte, error) {
	return e, nil
}
