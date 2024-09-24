package versions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/thrasher-corp/gocryptotrader/config"
)

// DisabledVersion allows authors to rollback changes easily during development
type DisabledVersion interface {
	Disabled() bool
}

// ConfigVersion is a version that affects the general configuration
type ConfigVersion interface {
	UpgradeConfig(context.Context, *config.Config) error
	DowngradeConfig(context.Context, *config.Config) error
}

// ConfigVersion is a version that affects specific exchange configurations
type ExchangeVersion interface {
	Exchanges() []string
	UpgradeExchange(context.Context, *config.Exchange) error
	DowngradeExchange(context.Context, *config.Exchange) error
}

var versions = []any{}

// Manage upgrades or downgrades the config between versions
// panics in out of bounds situations
func Manage(ctx context.Context, c *config.Config) error {
	target := latest()
	current := c.Version
	for current != target {
		next := current + 1
		configMethod := ConfigVersion.UpgradeConfig
		exchMethod := ExchangeVersion.UpgradeExchange
		if target < current {
			next = current - 1
			configMethod = ConfigVersion.DowngradeConfig
			exchMethod = ExchangeVersion.DowngradeExchange
		}

		patch := versions[next]

		if cPatch, ok := patch.(ConfigVersion); ok {
			if err := configMethod(cPatch, ctx, c); err != nil {
				return err
			}
		}

		if ePatch, ok := patch.(ExchangeVersion); ok {
			for _, name := range ePatch.Exchanges() {
				e := c.GetExchange(name)
				if err := exchMethod(ePatch, ctx, e); err != nil {
					return err
				}
			}
		}

		current = next
	}
	return nil
}

func RegisterVersion(v any) {
	ver, err := strconv.Atoi(strings.TrimPrefix(fmt.Sprintf("%T", v), "*versions.Version"))
	if err != nil {
		panic(err.Error())
	}
	switch v.(type) {
	case ExchangeVersion, ConfigVersion:
	default:
		panic(fmt.Sprintf("Version %v does not implement ConfigVersion or ExchangeVersion", ver))
	}
	if len(versions) < ver+1 {
		versions = append(versions, v)
	} else {
		if versions[ver] != nil {
			panic("Duplicate version registered")
		}
		versions[ver] = v
	}
}

// latest returns the highest version number
// May return -1 if something has gone deeply wrong
func latest() int {
	return len(versions) - 1
}
