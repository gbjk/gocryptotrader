package versions

import (
	v0 "github.com/thrasher-corp/gocryptotrader/config/versions/v0"
	v1 "github.com/thrasher-corp/gocryptotrader/config/versions/v1"
)

func init() {
	RegisterVersion(&v0.Version{})
	RegisterVersion(&v1.Version{})
}
