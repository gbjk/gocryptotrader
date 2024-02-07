package subscriptionstest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
)

func Equal(tb testing.TB, a, b subscription.List) {
	tb.Helper()
	added, missing := subscription.ListToMap(a).Diff(subscription.ListToMap(b))
	if len(added) > 0 || len(missing) > 0 {
		fail := "Differences:"
		if len(added) > 0 {
			fail = fail + "\n + " + strings.Join(added.Strings(), "\n + ")
		}
		if len(missing) > 0 {
			fail = fail + "\n - " + strings.Join(missing.Strings(), "\n - ")
		}
		assert.Fail(tb, fail, "Subscriptions should be equal")
	}
}
