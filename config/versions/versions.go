package versions

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/thrasher-corp/gocryptotrader/common"
)

// DisabledVersion allows authors to rollback changes easily during development
type DisabledVersion interface {
	Disabled() bool
}

// ConfigVersion is a version that affects the general configuration
type ConfigVersion interface {
	UpgradeConfig(context.Context, []byte) ([]byte, error)
	DowngradeConfig(context.Context, []byte) ([]byte, error)
}

// ConfigVersion is a version that affects specific exchange configurations
type ExchangeVersion interface {
	Exchanges() []string // Use `*` for all exchanges
	UpgradeExchange(ctx context.Context, configJSON []byte, exchangeJSON []byte, name string) ([]byte, error)
	DowngradeExchange(ctx context.Context, configJSON []byte, exchangeJSON []byte, names string) ([]byte, error)
}

var versions = []any{}

/*
Deploy upgrades or downgrades the config between versions
Accepts and returns io.ReadSeeker as a convenience for chaining config in a processing pipeline of decryption and upgrading
An important concept about upgrade versioning is the statefullness of the config definition at the time of the upgrade
Each version must use a config definition that matches the type at the time of that version
It should not use a config.* type, since that may be altered in the future, breaking the stateless multi-version upgrades
Therefore it should either define it's own types, or alter the json byte stream directly
It should return a config representation for the next version, never for the latest version
*/
func Deploy(ctx context.Context, j []byte) ([]byte, error) {
	target := latest()
	current, err := jsonparser.GetInt(j, "version")
	switch {
	case err == jsonparser.KeyPathNotFoundError:
		current = -1
	case err != nil:
		return j, err
	case target == current:
		return j, nil
	}

	for current != target {
		next := current + 1
		action := "upgrade"
		configMethod := ConfigVersion.UpgradeConfig
		exchMethod := ExchangeVersion.UpgradeExchange

		if target < current {
			next = current - 1
			action = "downgrade"
			configMethod = ConfigVersion.DowngradeConfig
			exchMethod = ExchangeVersion.DowngradeExchange
		}

		log.Printf("Running %s to config version %v\n", action, next)

		patch := versions[next]

		if cPatch, ok := patch.(ConfigVersion); ok {
			if j, err = configMethod(cPatch, ctx, j); err != nil {
				return j, err
			}
		}

		if ePatch, ok := patch.(ExchangeVersion); ok {
			var errs error
			for _, want := range ePatch.Exchanges() {
				eFunc := func(eOrig []byte, _ jsonparser.ValueType, _ int, err error) {
					name, err := jsonparser.GetString(eOrig, "name")
					if want != "*" && want != name {
						return
					}
					eNew, err := exchMethod(ePatch, ctx, j, eOrig, name)
					if err != nil {
						errs = common.AppendError(errs, err)
					}
					if !bytes.Equal(eNew, eOrig) {
						if j, err = jsonparser.Set(j, eNew, "exchanges", name); err != nil {
							errs = common.AppendError(errs, err)
						}
					}
				}
				if _, err := jsonparser.ArrayEach(bytes.Clone(j), eFunc, "exchanges"); err != nil {
					return j, err
				}
			}
			if errs != nil {
				return j, errs
			}
		}
		current = next
		if j, err = jsonparser.Set(j, []byte(strconv.Itoa(int(current))), "version"); err != nil {
			return j, err
		}
	}

	return j, nil
}

func RegisterVersion(v any) {
	ver, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%T", v), "*v"), ".Version"))
	if err != nil {
		log.Fatal(err.Error())
	}
	switch v.(type) {
	case ExchangeVersion, ConfigVersion:
	default:
		log.Fatalf("Version %v does not implement ConfigVersion or ExchangeVersion", ver)
	}
	if len(versions) != ver {
		log.Fatal("Out of sequence versions registered")
	}
	versions = append(versions, v)
}

// latest returns the highest version number
// May return -1 if something has gone deeply wrong
func latest() int64 {
	return int64(len(versions) - 1)
}
