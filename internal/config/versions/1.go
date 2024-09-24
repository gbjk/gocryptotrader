package versions

import (
	"context"

	"github.com/thrasher-corp/gocryptotrader/config"
)

// Version1 is a baseline version with no changes, so we can downgrade back to nothing
type Version1 struct {
}

func init() {
	RegisterVersion(&Version1{})
}

func (v *Version1) Upgrade(ctx context.Context, c *config.Config) {
	panic("called")
	return
}

func (v *Version1) Downgrade(ctx context.Context, c *config.Config) {
	return
}
