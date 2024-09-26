package v1

import (
	"context"
	"encoding/json"

	"github.com/buger/jsonparser"
)

// Version is an Exchange config to fix assets structure
type Version struct {
}

// Exchanges returns just Huobi
func (v *Version) Exchanges() []string { return []string{"*"} }

// UpgradeExchange will upgrade currency pair format
func (v *Version) UpgradeExchange(ctx context.Context, c, e []byte, name string) ([]byte, error) {
	if _, d, _, err := jsonparser.Get(e, "currencyPairs"); err == nil && d == jsonparser.Object {
		return e, nil
	}

	d := &exchDeprecated{}
	if err := json.Unmarshal(e, d); err != nil {
		return e, err
	}

	p := &pairsManager{
		UseGlobalFormat: true,
		LastUpdated:     d.PairsLastUpdated,
		ConfigFormat:    d.ConfigCurrencyPairFormat,
		RequestFormat:   d.RequestCurrencyPairFormat,
		Pairs: fullStore{
			"spot": {
				AssetEnabled: true,
				Available:    d.AvailablePairs,
				Enabled:      d.EnabledPairs,
			},
		},
	}
	pJson, err := json.Marshal(p)
	if err != nil {
		return e, err
	}
	for _, f := range []string{"pairsLastUpdated", "configCurrencyPairFormat", "requestCurrencyPairFormat", "assetTypes", "availablePairs", "enabledPairs"} {
		e = jsonparser.Delete(e, f)
	}
	return jsonparser.Set(e, pJson, "currencyPairs")
}

func (v *Version) DowngradeExchange(ctx context.Context, c, e []byte, name string) ([]byte, error) {
	return e, nil
}

type exchDeprecated struct {
	AvailablePairs            string      `json:"availablePairs,omitempty"`
	EnabledPairs              string      `json:"enabledPairs,omitempty"`
	PairsLastUpdated          int64       `json:"pairsLastUpdated,omitempty"`
	ConfigCurrencyPairFormat  *pairFormat `json:"configCurrencyPairFormat,omitempty"`
	RequestCurrencyPairFormat *pairFormat `json:"requestCurrencyPairFormat,omitempty"`
}

type pairsManager struct {
	BypassConfigFormatUpgrades bool        `json:"bypassConfigFormatUpgrades"`
	RequestFormat              *pairFormat `json:"requestFormat,omitempty"`
	ConfigFormat               *pairFormat `json:"configFormat,omitempty"`
	UseGlobalFormat            bool        `json:"useGlobalFormat,omitempty"`
	LastUpdated                int64       `json:"lastUpdated,omitempty"`
	Pairs                      fullStore   `json:"pairs"`
}

type fullStore map[string]struct {
	AssetEnabled  bool        `json:"assetEnabled"`
	Enabled       string      `json:"enabled"`
	Available     string      `json:"available"`
	RequestFormat *pairFormat `json:"requestFormat,omitempty"`
	ConfigFormat  *pairFormat `json:"configFormat,omitempty"`
}

type pairFormat struct {
	Uppercase bool   `json:"uppercase"`
	Delimiter string `json:"delimiter,omitempty"`
	Separator string `json:"separator,omitempty"`
}
