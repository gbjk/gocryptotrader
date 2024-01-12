//go:build mock_test_off

// This will build if build tag mock_test_off is parsed and will do live testing
// using all tests in (exchange)_test.go
package bitstamp

import (
	"log"
	"os"
	"testing"

	"github.com/thrasher-corp/gocryptotrader/exchanges/sharedtestvalues"
	testexch "github.com/thrasher-corp/gocryptotrader/internal/testing/exchange"
)

var mockTests = false

func TestMain(m *testing.M) {
	b = new(Bitstamp)
	if err := testexch.TestInstance(b); err != nil {
		log.Fatal(err)
	}

	if apiKey != "" && apiSecret != "" {
		b.API.CredentialsValidator.RequiresBase64DecodeSecret = false
		b.SetCredentials(apiKey, apiSecret, customerID, "", "", "")
	}

	log.Printf(sharedtestvalues.LiveTesting, b.Name)
	os.Exit(m.Run())
}
