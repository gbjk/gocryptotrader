package subscription

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
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

// QualifiedChannels returns subscriptions with Channel expanded with any placeholders replaced with subscription fields and parameters
// Format of the Channel should be text/template compatible.
// $asset and $pair will be expanded as plain substitutions, meaning they do not need to be wrapped in {{ }} and if they are they should be quoted
// This allows for custom functions like {{ assetName "$asset" }}
// If the channel contains $pair then the template will be ranged over the Pairs with each pair set to $pair
func (l List) QualifiedChannels(e iExchange, funcs template.FuncMap) (List, error) {
	ap, err := l.AssetPairs(e)
	if err != nil {
		return nil, err
	}

	// Break out Assets
	l2 := List{}
	for _, s := range l {
		if strings.Contains(s.Channel, "$asset") {
			if s.Asset != asset.All {
				return nil, ErrAssetTemplateWithoutAll
			}
			for a := range ap {
				tpl := strings.ReplaceAll(s.Channel, "$asset", a.String())
				c := s.Clone()
				c.Asset = a
				c.Channel = tpl
				l2 = append(l2, c)
			}
		} else {
			l2 = append(l2, s)
		}
	}

	// Break out Pairs
	l3 := List{}
	for _, s := range l2 {
		if strings.Contains(s.Channel, "$pair") {
			for a, pairs := range ap {
				if s.Asset != a && s.Asset != asset.All {
					continue
				}
				for _, p := range pairs {
					c := s.Clone()
					c.Pairs = currency.Pairs{p}
					c.Channel = strings.ReplaceAll(s.Channel, "$pair", p.String())
					l3 = append(l3, c)
				}
			}
		} else {
			l3 = append(l3, s)
		}
	}

	for _, s := range l3 {
		tpl := "{{with $s := . }}" + s.Channel + "{{end}}"
		t := template.New("channel")
		if funcs != nil {
			t = t.Funcs(funcs)
		}
		t, err = t.Parse(tpl)
		if err != nil {
			return nil, fmt.Errorf("%w parsing %s", err, tpl)
		}
		buf := &bytes.Buffer{}
		if err := t.Execute(buf, s); err != nil {
			return nil, err
		}
		s.Channel = buf.String()
	}

	return l3, nil
}
