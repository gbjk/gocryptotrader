package versions

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/buger/jsonparser"
	v2 "github.com/thrasher-corp/gocryptotrader/config/versions/v2"
)

// Version5 is an ExchangeVersion to split GateIO futures into CoinM and USDT margined futures assets
type Version5 struct{}

func init() {
	Manager.registerVersion(5, &Version5{})
}

// Exchanges returns just GateIO
func (v *Version5) Exchanges() []string { return []string{"GateIO"} }

// UpgradeExchange split GateIO futures into CoinM and USDT margined futures assets
func (v *Version5) UpgradeExchange(_ context.Context, e []byte) ([]byte, error) {
	fs := v2.FullStore{"coinmarginedfutures": {}, "usdtmarginedfutures": {}}
	fsJSON, _, _, err := jsonparser.Get(e, "currencyPairs", "pairs")
	if err != nil {
		return e, err
	}
	if err := json.Unmarshal(fsJSON, &fs); err != nil {
		return e, err
	}
	if _, ok := fs["futures"]; !ok {
		// Not our job to add CoinM, USDT. Only to split them
		return e, nil
	}
	for _, p := range strings.Split(fs["futures"].Available, ",") {
		where := "usdtmarginedfutures"
		if strings.HasSuffix(p, "USD") {
			where = "coinmarginedfutures"
		}
		if fs[where].Available != "" {
			fs[where].Available += ","
		}
		fs[where].Available += p
	}
	for _, p := range strings.Split(fs["futures"].Enabled, ",") {
		where := "usdtmarginedfutures"
		if strings.HasSuffix(p, "USD") {
			where = "coinmarginedfutures"
		}
		if fs[where].Enabled != "" {
			fs[where].Enabled += ","
		}
		fs[where].Enabled += p
	}
	fs["usdtmarginedfutures"].AssetEnabled = fs["futures"].AssetEnabled
	fs["coinmarginedfutures"].AssetEnabled = fs["futures"].AssetEnabled
	delete(fs, "futures")
	val, err := json.Marshal(fs)
	if err == nil {
		e, err = jsonparser.Set(e, val, "currencyPairs", "pairs")
	}
	return e, err
}

// DowngradeExchange will merge GateIO CoinM and USDT margined futures assets into futures
func (v *Version5) DowngradeExchange(_ context.Context, e []byte) ([]byte, error) {
	fs := v2.FullStore{"futures": {}, "coinmarginedfutures": {}, "usdtmarginedfutures": {}}
	fsJSON, _, _, err := jsonparser.Get(e, "currencyPairs", "pairs")
	if err != nil {
		return e, err
	}
	if err := json.Unmarshal(fsJSON, &fs); err != nil {
		return e, err
	}
	fs["futures"].Enabled = fs["coinmarginedfutures"].Enabled
	if fs["futures"].Enabled != "" {
		fs["futures"].Enabled += ","
	}
	fs["futures"].Enabled += fs["usdtmarginedfutures"].Enabled
	fs["futures"].Available = fs["coinmarginedfutures"].Available
	if fs["futures"].Available != "" {
		fs["futures"].Available += ","
	}
	fs["futures"].Available += fs["usdtmarginedfutures"].Available
	fs["futures"].AssetEnabled = fs["usdtmarginedfutures"].AssetEnabled || fs["coinmarginedfutures"].AssetEnabled
	delete(fs, "coinmarginedfutures")
	delete(fs, "usdtmarginedfutures")
	val, err := json.Marshal(fs)
	if err == nil {
		e, err = jsonparser.Set(e, val, "currencyPairs", "pairs")
	}
	return e, err
}
