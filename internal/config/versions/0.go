package versions

import (
	"context"

	"github.com/thrasher-corp/gocryptotrader/config"
)

// Version0 is a baseline version with no changes, so we can downgrade back to nothing
// It does not implement any upgrade interfaces
type Version0 struct {
}

var _ ConfigVersion = &Version0{}

func init() {
	RegisterVersion(&Version0{})
}

// UpgradeConfig is an empty stub for Version0
func (v *Version0) UpgradeConfig(ctx context.Context, c *config.Config) error {
	return nil
}

// DowngradeConfig is an empty stub for Version0
func (v *Version0) DowngradeConfig(ctx context.Context, c *config.Config) error {
	return nil
}
