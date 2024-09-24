package versions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/thrasher-corp/gocryptotrader/config"
)

type Version interface {
	Upgrade(context.Context, *config.Config)
	Downgrade(context.Context, *config.Config)
}

var versions = []Version{}

func Manage(ctx context.Context, c *config.Config) error {
	target := latest()
	current := c.Version
	for current != target {
		// A panic in out of bounds situations is acceptable error handling
		if target > current {
			patch := versions[current+1]
			patch.Upgrade(ctx, c)
			current = current + 1
		} else {
			patch := versions[current-1]
			patch.Downgrade(ctx, c)
			current = current - 1
		}
	}
	return nil
}
func RegisterVersion(v Version) {
	ver, err := strconv.Atoi(strings.TrimPrefix(fmt.Sprintf("%T", v), "*versions.Version"))
	if err != nil {
		panic(err.Error())
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
