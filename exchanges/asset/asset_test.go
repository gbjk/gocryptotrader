package asset

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Parallel()
	a := Spot
	if a.String() != "spot" {
		t.Fatal("TestString returned an unexpected result")
	}

	a = 0
	if a.String() != "" {
		t.Fatal("TestString returned an unexpected result")
	}
}

func TestStrings(t *testing.T) {
	t.Parallel()
	assert.ElementsMatch(t, Items{Spot, Futures}.Strings(), []string{"spot", "futures"})
}

func TestContains(t *testing.T) {
	t.Parallel()
	a := Items{Spot, Futures}
	if a.Contains(666) {
		t.Fatal("TestContains returned an unexpected result")
	}

	if !a.Contains(Spot) {
		t.Fatal("TestContains returned an unexpected result")
	}

	if a.Contains(Binary) {
		t.Fatal("TestContains returned an unexpected result")
	}

	// Every asset should be created and matched with func New so this should
	// not be matched against list
	if a.Contains(0) {
		t.Error("TestContains returned an unexpected result")
	}
}

func TestJoinToString(t *testing.T) {
	t.Parallel()
	a := Items{Spot, Futures}
	if a.JoinToString(",") != "spot,futures" {
		t.Fatal("TestJoinToString returned an unexpected result")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		Input    string
		Expected Item
		Error    error
	}{
		{Input: "Spota", Error: ErrNotSupported},
		{Input: "MARGIN", Expected: Margin},
		{Input: "MARGINFUNDING", Expected: MarginFunding},
		{Input: "INDEX", Expected: Index},
		{Input: "BINARY", Expected: Binary},
		{Input: "PERPETUALCONTRACT", Expected: PerpetualContract},
		{Input: "PERPETUALSWAP", Expected: PerpetualSwap},
		{Input: "FUTURES", Expected: Futures},
		{Input: "UpsideProfitContract", Expected: UpsideProfitContract},
		{Input: "DownsideProfitContract", Expected: DownsideProfitContract},
		{Input: "CoinMarginedFutures", Expected: CoinMarginedFutures},
		{Input: "USDTMarginedFutures", Expected: USDTMarginedFutures},
		{Input: "USDCMarginedFutures", Expected: USDCMarginedFutures},
		{Input: "Options", Expected: Options},
		{Input: "Option", Expected: Options},
		{Input: "Future", Error: ErrNotSupported},
		{Input: "future_combo", Expected: FutureCombo},
		{Input: "option_combo", Expected: OptionCombo},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			t.Parallel()
			returned, err := New(tt.Input)
			if !errors.Is(err, tt.Error) {
				t.Fatalf("received: '%v' but expected: '%v'", err, tt.Error)
			}
			if returned != tt.Expected {
				t.Fatalf("received: '%v' but expected: '%v'", returned, tt.Expected)
			}
		})
	}
}

func TestIsFutures(t *testing.T) {
	t.Parallel()
	for _, a := range []Item{Spot, Margin, MarginFunding, Index, Binary} {
		assert.Falsef(t, a.IsFutures(), "%s should return correctly for IsFutures")
	}
	for _, a := range []Item{PerpetualContract, PerpetualSwap, Futures, UpsideProfitContract, DownsideProfitContract, CoinMarginedFutures, USDTMarginedFutures, USDCMarginedFutures, FutureCombo} {
		assert.Truef(t, a.IsFutures(), "%s should return correctly for IsFutures")
	}
}

func TestIsOptions(t *testing.T) {
	t.Parallel()
	type scenario struct {
		item      Item
		isOptions bool
	}
	scenarios := []scenario{
		{
			item:      Options,
			isOptions: true,
		}, {
			item:      OptionCombo,
			isOptions: true,
		},
		{
			item:      Futures,
			isOptions: false,
		},
		{
			item:      Empty,
			isOptions: false,
		},
	}
	for _, s := range scenarios {
		testScenario := s
		t.Run(testScenario.item.String(), func(t *testing.T) {
			t.Parallel()
			if testScenario.item.IsOptions() != testScenario.isOptions {
				t.Errorf("expected %v isOptions to be %v", testScenario.item, testScenario.isOptions)
			}
		})
	}
}

func TestUnmarshalMarshal(t *testing.T) {
	t.Parallel()
	data, err := json.Marshal(Item(0))
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if string(data) != `""` {
		t.Fatal("unexpected value")
	}

	data, err = json.Marshal(Spot)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if string(data) != `"spot"` {
		t.Fatal("unexpected value")
	}

	var spot Item

	err = json.Unmarshal(data, &spot)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if spot != Spot {
		t.Fatal("unexpected value")
	}

	err = json.Unmarshal([]byte(`"confused"`), &spot)
	if !errors.Is(err, ErrNotSupported) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrNotSupported)
	}

	err = json.Unmarshal([]byte(`""`), &spot)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	err = json.Unmarshal([]byte(`123`), &spot)
	if errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", nil, "an error")
	}
}

func TestUseDefault(t *testing.T) {
	t.Parallel()
	if UseDefault() != Spot {
		t.Fatalf("received: '%v' but expected: '%v'", UseDefault(), Spot)
	}
}
