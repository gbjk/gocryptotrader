package subscription

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

var (
	errRecordSeparator         = errors.New("subscription template must not contain the RecordSeparator character")
	errInvalidAssetExpandPairs = errors.New("subscription template containing $pair with must contain either specific Asset or $asset with asset.All")
	errTemplateLines           = errors.New("subscription template did not generate the expected number of lines")
)

// List is a container of subscription pointers
type List []*Subscription

// Strings returns a sorted slice of subscriptions
func (l List) Strings() []string {
	s := make([]string, len(l))
	for i := range l {
		s[i] = l[i].String()
	}
	slices.Sort(s)
	return s
}

// GroupPairs groups subscriptions which are identical apart from the Pairs
// The returned List contains cloned Subscriptions, and the original Subscriptions are left alone
func (l List) GroupPairs() (n List) {
	s := NewStore()
	for _, sub := range l {
		if found := s.match(&IgnoringPairsKey{sub}); found == nil {
			s.unsafeAdd(sub.Clone())
		} else {
			found.AddPairs(sub.Pairs...)
		}
	}
	return s.List()
}

type assetPairs map[asset.Item]currency.Pairs
type iExchange interface {
	GetAssetTypes(enabled bool) asset.Items
	GetEnabledPairs(asset.Item) (currency.Pairs, error)
	GetPairFormat(asset.Item, bool) (currency.PairFormat, error)
	GetSubscriptionTemplateFuncs() template.FuncMap
}

func fillAssetPairs(ap assetPairs, a asset.Item, e iExchange) error {
	p, err := e.GetEnabledPairs(a)
	if err != nil {
		return err
	}
	f, err := e.GetPairFormat(a, true)
	if err != nil {
		return err
	}
	ap[a] = p.Format(f)
	return nil
}

// AssetPairs returns a map of enabled pairs for the subscriptions in the list, formatted for the asset
func (l List) AssetPairs(e iExchange) (assetPairs, error) { //nolint:revive // unexported-return is acceptable since it's a undecorated primitive
	at := e.GetAssetTypes(true)
	ap := assetPairs{}
	for _, s := range l {
		switch s.Asset {
		case asset.Empty:
			// Nothing to do
		case asset.All:
			for _, a := range at {
				if err := fillAssetPairs(ap, a, e); err != nil {
					return nil, err
				}
			}
		default:
			if slices.Contains(at, s.Asset) {
				if err := fillAssetPairs(ap, s.Asset, e); err != nil {
					return nil, err
				}
			}
		}
	}
	return ap, nil
}

type tplCtx struct {
	Sub        *Subscription
	AssetPairs assetPairs
	Assets     asset.Items
}

// QualifiedChannels returns subscriptions with Channel expanded with any placeholders replaced with subscription fields and parameters
// Format of the Channel should be text/template compatible.
// Template Variables:
// * $s is the subscription; e.g. {{$s.Interval}} {{$s.Params.freq}}
// * $asset will fan out the enabled assets
//   sub.Asset must be All; Otherwise just hardcode the asset the subscription is for
// * $pair will fan out the pairs for the assets, formatted for request
//   Must be used in cojunction with $asset when Asset is All, otherwise we don't know what pairs to use
//   May not not be used when Asset is Empty
// Calls e.GetSubscriptionTemplateFuncs for a template.FuncMap for flexibility in pipelines, e.g. {{ assetName "$asset" }}
// Filters out Authenticated subscriptions if CanUseAuthenticatedEndpoints is false
func (l List) QualifiedChannels(e iExchange) (List, error) {

	if !d.Websocket.CanUseAuthenticatedEndpoints() {
		l = slices.DeleteFunc(l, func(s *Subscription) bool {
			return s.Authenticated
		})
	}

	ap, err := l.AssetPairs(e)
	if err != nil {
		return nil, err
	}

	assets := make([]asset.Item, 0, len(ap))
	for k := range ap {
		assets = append(assets, k)
	}

	baseTpl := template.New("channel")
	if funcs := e.GetSubscriptionTemplateFuncs(); funcs != nil {
		baseTpl = baseTpl.Funcs(funcs)
	}

	subs := List{}

	rs := "\x1E"
	for _, s := range l {
		if strings.Contains(s.Channel, rs) {
			return nil, errRecordSeparator
		}

		subCtx := &tplCtx{
			Sub:        s,
			AssetPairs: ap,
		}

		tpl := s.Channel + rs

		xpandPairs := strings.Contains(s.Channel, "$pair")
		if xpandPairs {
			tpl = "{{range $pair := index $ctx.AssetPairs $asset}}" + tpl + "{{end}}"
		}

		if xpandAssets := strings.Contains(s.Channel, "$asset"); xpandAssets {
			if s.Asset != asset.All {
				return nil, ErrAssetTemplateWithoutAll
			}
			subCtx.Assets = assets
			tpl = "{{range $asset := $ctx.Assets}}" + tpl + "{{end}}"
		} else {
			if xpandPairs && (s.Asset == asset.All || s.Asset == asset.Empty) {
				// We don't currenply support expanding Pairs without expanding Assets for All or Empty assets, but we could; waiting for a use-case
				return nil, errInvalidAssetExpandPairs
			}
			subCtx.Assets = asset.Items{s.Asset}
			tpl = "{{with $asset := $ctx.Sub.Asset}}" + tpl + "{{end}}"
		}

		tpl = "{{with $ctx := .}}{{with $s := $ctx.Sub}}" + tpl + "{{end}}{{end}}"

		t, err := baseTpl.Parse(tpl)
		if err != nil {
			return nil, fmt.Errorf("%w parsing %s", err, tpl)
		}

		buf := &bytes.Buffer{}
		if err := t.Execute(buf, subCtx); err != nil {
			return nil, err
		}

		channels := strings.Split(strings.TrimSuffix(buf.String(), rs), rs)

		i := 0
		line := func(a asset.Item, p currency.Pairs) {
			if i < len(channels) {
				c := s.Clone()
				c.Asset = a
				c.Pairs = p
				c.Channel = channels[i]
				subs = append(subs, c)
			}
			i++ // Trigger errTemplateLines if we go over len(channels)
		}

		for _, a := range subCtx.Assets {
			if xpandPairs {
				for _, p := range ap[a] {
					line(a, currency.Pairs{p})
				}
			} else {
				line(a, ap[a])
			}
		}
		if i != len(channels) {
			return nil, fmt.Errorf("%w: Got %d Expected %d", errTemplateLines, len(channels), i)
		}
	}

	return subs, nil
}
