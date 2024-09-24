package versions

import (
	"context"

	"github.com/thrasher-corp/gocryptotrader/config"
)

// Version0 is a baseline version with no changes, so we can downgrade back to nothing
type Version0 struct {
}

func init() {
	RegisterVersion(&Version0{})
}

func (v *Version0) Upgrade(ctx context.Context, c *config.Config) {
	return
}

func (v *Version0) Downgrade(ctx context.Context, c *config.Config) {
	return
}
