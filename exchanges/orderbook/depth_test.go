package orderbook

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

var id = uuid.Must(uuid.NewV4())

func TestGetLength(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	err := d.Invalidate(nil)
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Invalidate should error correctly")

	_, err = d.GetAskLength()
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "GetAskLength should error with invalid depth")

	err = d.LoadSnapshot([]Item{{Price: 1337}}, nil, 0, time.Now(), true)
	assert.NoError(t, err, "LoadSnapshot should not error")

	askLen, err := d.GetAskLength()
	assert.NoError(t, err, "GetAskLength should not error")
	assert.Zero(t, askLen, "ask length should be zero")

	d.asks.load([]Item{{Price: 1337}}, d.stack, time.Now())

	askLen, err = d.GetAskLength()
	assert.NoError(t, err, "GetAskLength should not error")
	assert.Equal(t, 1, askLen, "Ask Length should be correct")

	d = NewDepth(id)
	err = d.Invalidate(nil)
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Invalidate should error correctly")

	_, err = d.GetBidLength()
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "GetBidLength should error with invalid depth")

	err = d.LoadSnapshot(nil, []Item{{Price: 1337}}, 0, time.Now(), true)
	assert.NoError(t, err, "LoadSnapshot should not error")

	bidLen, err := d.GetBidLength()
	assert.NoError(t, err, "GetBidLength should not error")
	assert.Zero(t, askLen, "bid length should be zero")

	d.bids.load([]Item{{Price: 1337}}, d.stack, time.Now())

	bidLen, err = d.GetBidLength()
	assert.NoError(t, err, "GetBidLength should not error")
	assert.Equal(t, 1, bidLen, "Bid Length should be correct")
}

func TestRetrieve(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	d.asks.load([]Item{{Price: 1337}}, d.stack, time.Now())
	d.bids.load([]Item{{Price: 1337}}, d.stack, time.Now())
	d.options = options{
		exchange:               "THE BIG ONE!!!!!!",
		pair:                   currency.NewPair(currency.THETA, currency.USD),
		asset:                  asset.DownsideProfitContract,
		lastUpdated:            time.Now(),
		lastUpdateID:           1337,
		priceDuplication:       true,
		isFundingRate:          true,
		VerifyOrderbook:        true,
		restSnapshot:           true,
		idAligned:              true,
		maxDepth:               10,
		checksumStringRequired: true,
	}

	// If we add anymore options to the options struct later this will complain
	// generally want to return a full carbon copy
	mirrored := reflect.Indirect(reflect.ValueOf(d.options))
	for n := 0; n < mirrored.NumField(); n++ {
		structVal := mirrored.Field(n)
		assert.Falsef(t, structVal.IsZero(), "struct field '%s' not tested", mirrored.Type().Field(n).Name)
	}

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")
	assert.Len(t, ob.Asks, 1, "Should have correct Asks")
	assert.Len(t, ob.Bids, 1, "Should have correct Bids")
	assert.Equal(t, 10, ob.MaxDepth, "Should have correct MaxDepth")
}

func TestTotalAmounts(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)

	err := d.Invalidate(nil)
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Invalidate should error correctly")
	_, _, err = d.TotalBidAmounts()
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "TotalBidAmounts should error correctly")

	d.validationError = nil
	liquidity, value, err := d.TotalBidAmounts()
	assert.NoError(t, err, "TotalBidAmounts should not error")
	assert.Zero(t, liquidity, "total bid liquidity should be zero")
	assert.Zero(t, value, "total bid value should be zero")

	err = d.Invalidate(nil)
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Invalidate should error correctly")

	_, _, err = d.TotalAskAmounts()
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "TotalAskAmounts should error correctly")

	d.validationError = nil

	liquidity, value, err = d.TotalAskAmounts()
	assert.NoError(t, err, "TotalBidAmounts should not error")
	assert.Zero(t, liquidity, "total bid liquidity should be zero")
	assert.Zero(t, value, "total bid value should be zero")

	d.asks.load([]Item{{Price: 1337, Amount: 1}}, d.stack, time.Now())
	d.bids.load([]Item{{Price: 1337, Amount: 10}}, d.stack, time.Now())

	liquidity, value, err = d.TotalBidAmounts()
	assert.NoError(t, err, "TotalBidAmounts should not error")
	assert.Equal(t, 10.0, liquidity, "total bid liquidity should be correct")
	assert.Equal(t, 13370.0, value, "total bid value should be correct")

	liquidity, value, err = d.TotalAskAmounts()
	assert.NoError(t, err, "TotalAskAmounts should not error")
	assert.Equal(t, 1.0, liquidity, "total ask liquidity should be correct")
	assert.Equal(t, 1337.0, value, "total ask value should be correct")
}

func TestLoadSnapshot(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	err := d.LoadSnapshot(Items{{Price: 1337, Amount: 1}}, Items{{Price: 1337, Amount: 10}}, 0, time.Time{}, false)
	assert.ErrorIs(t, err, errLastUpdatedNotSet, "LoadSnapshot should error correctly")

	err = d.LoadSnapshot(Items{{Price: 1337, Amount: 2}}, Items{{Price: 1338, Amount: 10}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")

	assert.Equal(t, 1338.0, ob.Asks[0].Price, "Top ask price should be correct")
	assert.Equal(t, 10.0, ob.Asks[0].Amount, "Top ask amount should be correct")
	assert.Equal(t, 1337.0, ob.Bids[0].Price, "Top bid price should be correct")
	assert.Equal(t, 2.0, ob.Bids[0].Amount, "Top bid amount should be correct")
}

func TestInvalidate(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	d.exchange = "testexchange"
	d.pair = currency.NewPair(currency.BTC, currency.WABI)
	d.asset = asset.Spot

	err := d.LoadSnapshot(Items{{Price: 1337, Amount: 1}}, Items{{Price: 1337, Amount: 10}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")
	assert.NotNil(t, ob, "ob should not be nil")

	testReason := errors.New("random reason")

	err = d.Invalidate(testReason)
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Invalidate should error correctly")

	_, err = d.Retrieve()
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Retrieve should error correctly")
	assert.ErrorIs(t, err, testReason, "Invalidate should error correctly")

	d.validationError = nil

	ob, err = d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")

	assert.Empty(t, ob.Asks, "Orderbook Asks should be flushed")
	assert.Empty(t, ob.Bids, "Orderbook Bids should be flushed")
}

func TestUpdateBidAskByPrice(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	err := d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1338, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	err = d.UpdateBidAskByPrice(&Update{})
	assert.ErrorIs(t, err, errLastUpdatedNotSet, "UpdateBidAskByPrice should error correctly")

	err = d.UpdateBidAskByPrice(&Update{UpdateTime: time.Now()})
	assert.NoError(t, err, "UpdateBidAskByPrice should not error")

	updates := &Update{
		Bids:       Items{{Price: 1337, Amount: 2, ID: 1}},
		Asks:       Items{{Price: 1338, Amount: 3, ID: 2}},
		UpdateID:   1,
		UpdateTime: time.Now(),
	}
	err = d.UpdateBidAskByPrice(updates)
	assert.NoError(t, err, "UpdateBidAskByPrice should not error")

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")
	assert.Equal(t, 3.0, ob.Asks[0].Amount, "Asks amount should be correct")
	assert.Equal(t, 2.0, ob.Bids[0].Amount, "Bids amount should be correct")

	updates = &Update{
		Bids:       Items{{Price: 1337, Amount: 0, ID: 1}},
		Asks:       Items{{Price: 1338, Amount: 0, ID: 2}},
		UpdateID:   2,
		UpdateTime: time.Now(),
	}
	err = d.UpdateBidAskByPrice(updates)
	assert.NoError(t, err, "UpdateBidAskByPrice should not error")

	askLen, err := d.GetAskLength()
	assert.NoError(t, err, "GetAskLength should not error")
	assert.Zero(t, askLen, "Ask Length should be correct")

	bidLen, err := d.GetBidLength()
	assert.NoError(t, err, "GetBidLength should not error")
	assert.Zero(t, bidLen, "Bid Length should be correct")
}

func TestDeleteBidAskByID(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	err := d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates := &Update{
		Bids: Items{{Price: 1337, Amount: 2, ID: 1}},
		Asks: Items{{Price: 1337, Amount: 2, ID: 2}},
	}

	err = d.DeleteBidAskByID(updates, false)
	assert.ErrorIs(t, err, errLastUpdatedNotSet, "DeleteBidAskByID should error correctly")

	updates.UpdateTime = time.Now()
	err = d.DeleteBidAskByID(updates, false)
	assert.NoError(t, err, "DeleteBidAskByID should not error")

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")
	assert.Empty(t, ob.Asks, "Asks should be empty")
	assert.Empty(t, ob.Bids, "Bids should be empty")

	updates = &Update{
		Bids:       Items{{Price: 1337, Amount: 2, ID: 1}},
		UpdateTime: time.Now(),
	}
	err = d.DeleteBidAskByID(updates, false)
	assert.ErrorIs(t, err, errIDCannotBeMatched, "DeleteBidAskByID should error correctly")

	updates = &Update{
		Asks:       Items{{Price: 1337, Amount: 2, ID: 2}},
		UpdateTime: time.Now(),
	}
	err = d.DeleteBidAskByID(updates, false)
	assert.ErrorIs(t, err, errIDCannotBeMatched, "DeleteBidAskByID should error correctly")

	updates = &Update{
		Asks:       Items{{Price: 1337, Amount: 2, ID: 2}},
		UpdateTime: time.Now(),
	}
	err = d.DeleteBidAskByID(updates, true)
	assert.NoError(t, err, "DeleteBidAskByID should not error")
}

func TestUpdateBidAskByID(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	err := d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates := &Update{
		Bids: Items{{Price: 1337, Amount: 2, ID: 1}},
		Asks: Items{{Price: 1337, Amount: 2, ID: 2}},
	}

	err = d.UpdateBidAskByID(updates)
	assert.ErrorIs(t, err, errLastUpdatedNotSet, "UpdateBidAskByID should error correctly")

	updates.UpdateTime = time.Now()
	err = d.UpdateBidAskByID(updates)
	assert.NoError(t, err, "UpdateBidAskByID should not error")

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")
	assert.Equal(t, 2.0, ob.Asks[0].Amount, "First ask amount should be correct")
	assert.Equal(t, 2.0, ob.Bids[0].Amount, "Fisrt bid amount should be correct")

	updates = &Update{
		Bids:       Items{{Price: 1337, Amount: 2, ID: 666}},
		UpdateTime: time.Now(),
	}
	// random unmatching IDs
	err = d.UpdateBidAskByID(updates)
	assert.ErrorIs(t, err, errIDCannotBeMatched, "UpdateBidAskByID should error correctly")

	updates = &Update{
		Asks:       Items{{Price: 1337, Amount: 2, ID: 69}},
		UpdateTime: time.Now(),
	}
	err = d.UpdateBidAskByID(updates)
	assert.ErrorIs(t, err, errIDCannotBeMatched, "UpdateBidAskByID should error correctly")
}

func TestInsertBidAskByID(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	err := d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates := &Update{
		Asks: Items{{Price: 1337, Amount: 2, ID: 3}},
	}
	err = d.InsertBidAskByID(updates)
	assert.ErrorIs(t, err, errLastUpdatedNotSet, "InsertBidAskByID should error correctly")

	updates.UpdateTime = time.Now()

	err = d.InsertBidAskByID(updates)
	assert.ErrorIs(t, err, errCollisionDetected, "InsertBidAskByID should error correctly on collision")

	err = d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates = &Update{
		Bids:       Items{{Price: 1337, Amount: 2, ID: 3}},
		UpdateTime: time.Now(),
	}

	err = d.InsertBidAskByID(updates)
	assert.ErrorIs(t, err, errCollisionDetected, "InsertBidAskByID should error correctly on collision")

	err = d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates = &Update{
		Bids:       Items{{Price: 1338, Amount: 2, ID: 3}},
		Asks:       Items{{Price: 1336, Amount: 2, ID: 4}},
		UpdateTime: time.Now(),
	}
	err = d.InsertBidAskByID(updates)
	assert.NoError(t, err, "InsertBidAskByID should not error")

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")
	assert.Len(t, ob.Asks, 2, "Should have correct Asks")
	assert.Len(t, ob.Bids, 2, "Should have correct Bids")
}

func TestUpdateInsertByID(t *testing.T) {
	t.Parallel()
	d := NewDepth(id)
	err := d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates := &Update{
		Bids: Items{{Price: 1338, Amount: 0, ID: 3}},
		Asks: Items{{Price: 1336, Amount: 2, ID: 4}},
	}
	err = d.UpdateInsertByID(updates)
	assert.ErrorIs(t, err, errLastUpdatedNotSet, "UpdateInsertByID should error correctly")

	updates.UpdateTime = time.Now()
	err = d.UpdateInsertByID(updates)
	assert.ErrorIs(t, err, errAmountCannotBeLessOrEqualToZero, "UpdateInsertByID should error correctly")

	// Above will invalidate the book
	_, err = d.Retrieve()
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Retrieve should error correctly")

	err = d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates = &Update{
		Bids:       Items{{Price: 1338, Amount: 2, ID: 3}},
		Asks:       Items{{Price: 1336, Amount: 0, ID: 4}},
		UpdateTime: time.Now(),
	}
	err = d.UpdateInsertByID(updates)
	assert.ErrorIs(t, err, errAmountCannotBeLessOrEqualToZero, "UpdateInsertByID should error correctly")

	// Above will invalidate the book
	_, err = d.Retrieve()
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Retrieve should error correctly")

	err = d.LoadSnapshot(Items{{Price: 1337, Amount: 1, ID: 1}}, Items{{Price: 1337, Amount: 10, ID: 2}}, 0, time.Now(), false)
	assert.NoError(t, err, "LoadSnapshot should not error")

	updates = &Update{
		Bids:       Items{{Price: 1338, Amount: 2, ID: 3}},
		Asks:       Items{{Price: 1336, Amount: 2, ID: 4}},
		UpdateTime: time.Now(),
	}
	err = d.UpdateInsertByID(updates)
	assert.NoError(t, err, "UpdateInsertByID should not error")

	ob, err := d.Retrieve()
	assert.NoError(t, err, "Retrieve should not error")
	assert.Len(t, ob.Asks, 2, "Should have correct Asks")
	assert.Len(t, ob.Bids, 2, "Should have correct Bids")
}

func TestAssignOptions(t *testing.T) {
	t.Parallel()
	d := Depth{}
	cp := currency.NewPair(currency.LINK, currency.BTC)
	tn := time.Now()
	d.AssignOptions(&Base{
		Exchange:         "test",
		Pair:             cp,
		Asset:            asset.Spot,
		LastUpdated:      tn,
		LastUpdateID:     1337,
		PriceDuplication: true,
		IsFundingRate:    true,
		VerifyOrderbook:  true,
		RestSnapshot:     true,
		IDAlignment:      true,
	})

	if d.exchange != "test" ||
		d.pair != cp ||
		d.asset != asset.Spot ||
		d.lastUpdated != tn ||
		d.lastUpdateID != 1337 ||
		!d.priceDuplication ||
		!d.isFundingRate ||
		!d.VerifyOrderbook ||
		!d.restSnapshot ||
		!d.idAligned {
		t.Fatalf("failed to set correctly")
	}
}

func TestGetName(t *testing.T) {
	t.Parallel()
	d := Depth{}
	d.exchange = "test"
	if d.GetName() != "test" {
		t.Fatalf("failed to get correct value")
	}
}

func TestIsRestSnapshot(t *testing.T) {
	t.Parallel()
	d := Depth{}
	d.restSnapshot = true
	err := d.Invalidate(nil)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}
	_, err = d.IsRESTSnapshot()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	d.validationError = nil
	b, err := d.IsRESTSnapshot()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !b {
		t.Fatalf("received: '%v' but expected: '%v'", b, true)
	}
}

func TestLastUpdateID(t *testing.T) {
	t.Parallel()
	d := Depth{}
	err := d.Invalidate(nil)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}
	_, err = d.LastUpdateID()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	d.validationError = nil
	d.lastUpdateID = 1337
	id, err := d.LastUpdateID()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if id != 1337 {
		t.Fatalf("received: '%v' but expected: '%v'", id, 1337)
	}
}

func TestIsFundingRate(t *testing.T) {
	t.Parallel()
	d := Depth{}
	d.isFundingRate = true
	if !d.IsFundingRate() {
		t.Fatalf("failed to get correct value")
	}
}

func TestPublish(t *testing.T) {
	t.Parallel()
	d := Depth{}
	if err := d.Invalidate(nil); !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}
	d.Publish()
	d.validationError = nil
	d.Publish()
}

func TestIsValid(t *testing.T) {
	t.Parallel()
	d := Depth{}
	if !d.IsValid() {
		t.Fatalf("received: '%v' but expected: '%v'", d.IsValid(), true)
	}
	if err := d.Invalidate(nil); !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}
	if d.IsValid() {
		t.Fatalf("received: '%v' but expected: '%v'", d.IsValid(), false)
	}
}

func TestHitTheBidsByNominalSlippage(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsByNominalSlippage(10, 1355.5)
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Should error invalid orderbook")

	depth := NewDepth(id)
	_, err = depth.HitTheBidsByNominalSlippage(10, 1355.5)
	assert.ErrorIs(t, err, errNoLiquidity, "Should error no liquidity")

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	assert.NoError(t, err, "LoadSnapshot should not error")

	// First tranche
	amt, err := depth.HitTheBidsByNominalSlippage(0, 1336)
	assert.NoError(t, err, "HitTheBidsByNominalSlippage should not error")
	assert.Equal(t, 1.0, amt.Sold, "Amount sold should be correct")
	assert.Zero(t, amt.NominalPercentage, "NominalPercentage should be 0")
	assert.Equal(t, 1336.0, amt.StartPrice, "Amount StartPrice should be correct")
	assert.Equal(t, 1336.0, amt.EndPrice, "Amount EndPrice should be correct")
	assert.False(t, amt.FullBookSideConsumed, "Orderbook should not be entirely consumed")

	// First and second price
	amt, err = depth.HitTheBidsByNominalSlippage(0.037425149700598806, 1336)
	assert.NoError(t, err, "HitTheBidsByNominalSlippage should not error")
	assert.Equal(t, 2.0, amt.Sold, "Amount sold should be correct")
	assert.Equal(t, 0.037425149700598806, amt.NominalPercentage, "NominalPercentage should be correct")
	assert.Equal(t, 1336.0, amt.StartPrice, "Amount StartPrice should be correct")
	assert.Equal(t, 1335.0, amt.EndPrice, "Amount EndPrice should be correct")
	assert.False(t, amt.FullBookSideConsumed, "Orderbook should not be entirely consumed")

	// First and half of second tranche
	amt, err = depth.HitTheBidsByNominalSlippage(0.02495009980039353, 1336)
	assert.NoError(t, err, "HitTheBidsByNominalSlippage should not error")
	assert.InDelta(t, 1.5, amt.Sold, 0.00000000001, "Amount sold should be correct")
	assert.Equal(t, 0.02495009980039353, amt.NominalPercentage, "NominalPercentage should be correct")
	assert.Equal(t, 1336.0, amt.StartPrice, "Amount StartPrice should be correct")
	assert.Equal(t, 1335.0, amt.EndPrice, "Amount EndPrice should be correct")
	assert.False(t, amt.FullBookSideConsumed, "Orderbook should not be entirely consumed")

	// All the way up to the last price
	amt, err = depth.HitTheBidsByNominalSlippage(0.7110778443113772, 1336)
	assert.NoError(t, err, "HitTheBidsByNominalSlippage should not error")
	assert.Equal(t, 20.0, amt.Sold, "Amount sold should be correct")
	assert.Equal(t, 0.71107784431137723, amt.NominalPercentage, "NominalPercentage should be correct")
	assert.Equal(t, 1336.0, amt.StartPrice, "Amount StartPrice should be correct")
	assert.Equal(t, 1317.0, amt.EndPrice, "Amount EndPrice should be correct")
	assert.True(t, amt.FullBookSideConsumed, "Orderbook should be entirely consumed")
}

func TestHitTheBidsByNominalSlippageFromMid(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsByNominalSlippageFromMid(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.HitTheBidsByNominalSlippageFromMid(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First price from mid point
	amt, err := depth.HitTheBidsByNominalSlippageFromMid(0.03741114852226)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 1 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 1)
	}

	// All the way up to the last price from mid price
	amt, err = depth.HitTheBidsByNominalSlippageFromMid(0.74822297044519)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	// This exceeds the entire total base available
	if amt.Sold != 20 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 20)
	}
}

func TestHitTheBidsByNominalSlippageFromBest(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsByNominalSlippageFromBest(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.HitTheBidsByNominalSlippageFromBest(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First and second price from best bid
	amt, err := depth.HitTheBidsByNominalSlippageFromBest(0.037425149700599)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 2 {
		t.Fatalf("received: '%+v' but expected: '%+v'", amt, 2)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.HitTheBidsByNominalSlippageFromBest(0.71107784431138)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	// This exceeds the entire total base available
	if amt.Sold != 20 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 20)
	}
}

func TestLiftTheAsksByNominalSlippage(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksByNominalSlippage(10, 1355.5)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.LiftTheAsksByNominalSlippage(10, 1355.5)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First and second price
	amt, err := depth.LiftTheAsksByNominalSlippage(0.037397157816006, 1337)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 2675 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 2675)
	}

	// All the way up to the last price
	amt, err = depth.LiftTheAsksByNominalSlippage(0.71054599850411, 1337)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 26930 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 26930)
	}
}

func TestLiftTheAsksByNominalSlippageFromMid(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksByNominalSlippageFromMid(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.LiftTheAsksByNominalSlippageFromMid(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First price from mid point
	amt, err := depth.LiftTheAsksByNominalSlippageFromMid(0.074822297044519)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 2675 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 2675)
	}

	// All the way up to the last price from mid price
	amt, err = depth.LiftTheAsksByNominalSlippageFromMid(0.74822297044519)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	// This does not match the entire total quote available
	if amt.Sold != 26930 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 26930)
	}
}

func TestLiftTheAsksByNominalSlippageFromBest(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksByNominalSlippageFromBest(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.LiftTheAsksByNominalSlippageFromBest(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First and second price from best bid
	amt, err := depth.LiftTheAsksByNominalSlippageFromBest(0.037397157816006)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 2675 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 2675)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.LiftTheAsksByNominalSlippageFromBest(0.71054599850411)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	// This does not match the entire total quote available
	if amt.Sold != 26930 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 26930)
	}
}

func TestHitTheBidsByImpactSlippage(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsByImpactSlippage(0.7485029940119761, 1336)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First and second price from best bid - price level target 1326 (which should be kept)
	amt, err := depth.HitTheBidsByImpactSlippage(0.7485029940119761, 1336)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 10 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 10)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.HitTheBidsByImpactSlippage(1.4221556886227544, 1336)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	// This does not match the entire total quote available - should be 26930.
	if amt.Sold != 19 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 19)
	}
}

func TestHitTheBidsByImpactSlippageFromMid(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsByImpactSlippageFromMid(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.HitTheBidsByImpactSlippageFromMid(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First and second price from mid - price level target 1326 (which should be kept)
	amt, err := depth.HitTheBidsByImpactSlippageFromMid(0.7485029940119761)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 10 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 10)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.HitTheBidsByImpactSlippageFromMid(1.4221556886227544)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 19 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 19)
	}
}

func TestHitTheBidsByImpactSlippageFromBest(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsByImpactSlippageFromBest(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.HitTheBidsByImpactSlippageFromBest(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First and second price from mid - price level target 1326 (which should be kept)
	amt, err := depth.HitTheBidsByImpactSlippageFromBest(0.7485029940119761)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 10 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 10)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.HitTheBidsByImpactSlippageFromBest(1.4221556886227544)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 19 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 19)
	}
}

func TestLiftTheAsksByImpactSlippage(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksByImpactSlippage(0.7479431563201197, 1337)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	// First and second price from best bid - price level target 1326 (which should be kept)
	amt, err := depth.LiftTheAsksByImpactSlippage(0.7479431563201197, 1337)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 13415 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 13415)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.LiftTheAsksByImpactSlippage(1.4210919970082274, 1337)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 25574 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 25574)
	}
}

func TestLiftTheAsksByImpactSlippageFromMid(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksByImpactSlippageFromMid(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.LiftTheAsksByImpactSlippageFromMid(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	// First and second price from mid - price level target 1326 (which should be kept)
	amt, err := depth.LiftTheAsksByImpactSlippageFromMid(0.7485029940119761)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 13415 {
		t.Fatalf("received: '%+v' but expected: '%v'", amt, 13415)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.LiftTheAsksByImpactSlippageFromMid(1.4221556886227544)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 25574 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 25574)
	}
}

func TestLiftTheAsksByImpactSlippageFromBest(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksByImpactSlippageFromBest(10)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.LiftTheAsksByImpactSlippageFromBest(10)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	// First and second price from mid - price level target 1326 (which should be kept)
	amt, err := depth.LiftTheAsksByImpactSlippageFromBest(0.7479431563201197)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if amt.Sold != 13415 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 13415)
	}

	// All the way up to the last price from best bid price
	amt, err = depth.LiftTheAsksByImpactSlippageFromBest(1.4210919970082274)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	// This goes to price 1356, it will not count that tranches' volume as it
	// is needed to sustain the slippage.
	if amt.Sold != 25574 {
		t.Fatalf("received: '%v' but expected: '%v'", amt, 25574)
	}
}

func TestHitTheBids(t *testing.T) {
	t.Parallel()
	depth := NewDepth(id)
	err := depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.HitTheBids(20.1, 1336, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.HitTheBids(1, 1336, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0)
	}
	if mov.ImpactPercentage != 0.07485029940119761 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.07485029940119761)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.HitTheBids(19.5, 1336, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.692845079072617 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.692845079072617)
	}
	if mov.ImpactPercentage != 1.4221556886227544 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4221556886227544)
	}

	if mov.SlippageCost != 180.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 180.5)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.HitTheBids(20, 1336, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7110778443113772 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7110778443113772)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestHitTheBids_QuotationRequired(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBids(26531, 1336, true)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.HitTheBids(26531, 1336, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.HitTheBids(1336, 1336, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0)
	}
	if mov.ImpactPercentage != 0.07485029940119761 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.07485029940119761)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.HitTheBids(25871.5, 1336, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.692845079072617 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.692845079072617)
	}
	if mov.ImpactPercentage != 1.4221556886227544 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4221556886227544)
	}

	if mov.SlippageCost != 180.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 180.5)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.HitTheBids(26530, 1336, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7110778443113772 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7110778443113772)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestHitTheBidsFromMid(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsFromMid(10, false)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.HitTheBidsFromMid(10, false)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.HitTheBidsFromMid(20.1, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.HitTheBidsFromMid(1, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.03741114852225963 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.03741114852225963)
	}

	if mov.ImpactPercentage != 0.11223344556677892 { // mid price 1336.5 -> 1335
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.11223344556677892)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.HitTheBidsFromMid(19.5, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7299970262933156 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7299970262933156)
	}
	if mov.ImpactPercentage != 1.4590347923681257 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4590347923681257)
	}
	if mov.SlippageCost != 180.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 180.5)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.HitTheBidsFromMid(20, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7482229704451926 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7482229704451926)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestHitTheBidsFromMid_QuotationRequired(t *testing.T) {
	t.Parallel()
	depth := NewDepth(id)
	_, err := depth.HitTheBidsFromMid(10, false)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.HitTheBidsFromMid(26531, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.HitTheBidsFromMid(1336, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.03741114852225963 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.03741114852225963)
	}

	if mov.ImpactPercentage != 0.11223344556677892 { // mid price 1336.5 -> 1335
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.11223344556677892)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.HitTheBidsFromMid(25871.5, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7299970262933156 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7299970262933156)
	}
	if mov.ImpactPercentage != 1.4590347923681257 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4590347923681257)
	}
	if mov.SlippageCost != 180.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 180.5)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.HitTheBidsFromMid(26530, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7482229704451926 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7482229704451926)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestHitTheBidsFromBest(t *testing.T) {
	t.Parallel()
	depth := NewDepth(id)
	_, err := depth.HitTheBidsFromBest(10, false)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.HitTheBidsFromBest(20.1, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.HitTheBidsFromBest(1, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0)
	}
	if mov.ImpactPercentage != 0.07485029940119761 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.07485029940119761)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.HitTheBidsFromBest(19.5, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.692845079072617 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.692845079072617)
	}
	if mov.ImpactPercentage != 1.4221556886227544 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4221556886227544)
	}
	if mov.SlippageCost != 180.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 180.5)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.HitTheBidsFromBest(20, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7110778443113772 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7110778443113772)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestHitTheBidsFromBest_QuotationRequired(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().HitTheBidsFromBest(10, false)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.HitTheBidsFromBest(10, false)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.HitTheBidsFromBest(26531, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.HitTheBidsFromBest(1336, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0)
	}
	if mov.ImpactPercentage != 0.07485029940119761 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.07485029940119761)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.HitTheBidsFromBest(25871.5, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.692845079072617 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.692845079072617)
	}
	if mov.ImpactPercentage != 1.4221556886227544 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4221556886227544)
	}
	if mov.SlippageCost != 180.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 180.5)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.HitTheBidsFromBest(26530, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7110778443113772 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7110778443113772)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestLiftTheAsks(t *testing.T) {
	t.Parallel()
	depth := NewDepth(id)
	err := depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.LiftTheAsks(26931, 1337, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.LiftTheAsks(1337, 1337, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0)
	}
	if mov.ImpactPercentage != 0.07479431563201197 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.07479431563201197)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.LiftTheAsks(26900, 1337, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7097591258590459 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7097591258590459)
	}
	if mov.ImpactPercentage != 1.4210919970082274 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4210919970082274)
	}
	if mov.SlippageCost != 189.57964601770072 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 189.57964601770072)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.LiftTheAsks(26930, 1336, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7859281437125748 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7859281437125748)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestLiftTheAsks_BaseRequired(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsks(21, 1337, true)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.LiftTheAsks(21, 1337, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.LiftTheAsks(1, 1337, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0)
	}
	if mov.ImpactPercentage != 0.07479431563201197 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.07479431563201197)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.LiftTheAsks(19.97787610619469, 1337, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7097591258590288 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7097591258590288)
	}
	if mov.ImpactPercentage != 1.4210919970082274 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4210919970082274)
	}
	if mov.SlippageCost != 189.5796460176971 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 189.5796460176971)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.LiftTheAsks(20, 1336, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7859281437125748 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7859281437125748)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestLiftTheAsksFromMid(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksFromMid(10, false)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.LiftTheAsksFromMid(10, false)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.LiftTheAsksFromMid(26931, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.LiftTheAsksFromMid(1337, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.03741114852225963 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.03741114852225963)
	}
	if mov.ImpactPercentage != 0.11223344556677892 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.11223344556677892)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.LiftTheAsksFromMid(26900, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.747435803422031 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.747435803422031)
	}
	if mov.ImpactPercentage != 1.4590347923681257 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4590347923681257)
	}
	if mov.SlippageCost != 189.57964601770072 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 189.57964601770072)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.LiftTheAsksFromMid(26930, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7482229704451926 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7482229704451926)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestLiftTheAsksFromMid_BaseRequired(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksFromMid(10, false)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.LiftTheAsksFromMid(10, false)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.LiftTheAsksFromMid(21, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.LiftTheAsksFromMid(1, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.03741114852225963 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.03741114852225963)
	}
	if mov.ImpactPercentage != 0.11223344556677892 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.11223344556677892)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.LiftTheAsksFromMid(19.97787610619469, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7474358034220139 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7474358034220139)
	}
	if mov.ImpactPercentage != 1.4590347923681257 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4590347923681257)
	}
	if mov.SlippageCost != 189.5796460176971 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 189.5796460176971)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.LiftTheAsksFromMid(20, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7482229704451926 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7482229704451926)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestLiftTheAsksFromBest(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksFromBest(10, false)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.LiftTheAsksFromBest(10, false)
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mov, err := depth.LiftTheAsksFromBest(26931, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !mov.FullBookSideConsumed {
		t.Fatal("entire side should be consumed by this value")
	}

	mov, err = depth.LiftTheAsksFromBest(1337, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0)
	}
	if mov.ImpactPercentage != 0.07479431563201197 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 0.07479431563201197)
	}
	if mov.SlippageCost != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 0)
	}

	mov, err = depth.LiftTheAsksFromBest(26900, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7097591258590459 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7097591258590459)
	}
	if mov.ImpactPercentage != 1.4210919970082274 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, 1.4210919970082274)
	}
	if mov.SlippageCost != 189.57964601770072 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 189.57964601770072)
	}

	// All the way up to the last price from best bid price
	mov, err = depth.LiftTheAsksFromBest(26930, false)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7105459985041137 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7105459985041137)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestLiftTheAsksFromBest_BaseRequired(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().LiftTheAsksFromBest(10, false)
	assert.ErrorIs(t, err, ErrOrderbookInvalid, "Should error invalid orderbook")

	depth := NewDepth(id)
	_, err = depth.LiftTheAsksFromBest(10, false)
	assert.ErrorIs(t, err, errNoLiquidity, "Should error no liquidity")

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	assert.NoError(t, err, "LoadSnapshot should not error")

	mov, err := depth.LiftTheAsksFromBest(21, true)
	assert.NoError(t, err, "LiftTheAsksFromBest should not error")
	assert.True(t, mov.FullBookSideConsumed, "entire side should be consumed by this lift")

	mov, err = depth.LiftTheAsksFromBest(1, true)
	assert.NoError(t, err, "LiftTheAsksFromBest should not error")
	assert.Zero(t, mov.NominalPercentage, "NominalPercentage should be 0")
	assert.Equal(t, 0.07479431563201197, mov.ImpactPercentage, "ImpactPercentage should be correct")
	assert.Zero(t, mov.SlippageCost, "SlippageCost should be 0")

	mov, err = depth.LiftTheAsksFromBest(19.97787610619469, true)
	assert.NoError(t, err, "LiftTheAsksFromBest should not error")
	assert.Equal(t, 0.7097591258590288, mov.NominalPercentage, "NominalPercentage should be correct")
	assert.Equal(t, 1.4210919970082274, mov.ImpactPercentage, "ImpactPercentage should be correct")
	assert.InDelta(t, 189.57964601769947, mov.SlippageCost, 0.00000000000001, "SlippageCost should be correct")

	// All the way up to the last price from best bid price
	mov, err = depth.LiftTheAsksFromBest(20, true)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mov.NominalPercentage != 0.7105459985041137 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.NominalPercentage, 0.7105459985041137)
	}
	if mov.ImpactPercentage != FullLiquidityExhaustedPercentage {
		t.Fatalf("received: '%v' but expected: '%v'", mov.ImpactPercentage, FullLiquidityExhaustedPercentage)
	}
	if mov.SlippageCost != 190 {
		t.Fatalf("received: '%v' but expected: '%v'", mov.SlippageCost, 190)
	}
}

func TestGetMidPrice_Depth(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().GetMidPrice()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.GetMidPrice()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mid, err := depth.GetMidPrice()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mid != 1336.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mid, 1336.5)
	}
}

func TestGetMidPriceNoLock_Depth(t *testing.T) {
	t.Parallel()
	depth := NewDepth(id)
	_, err := depth.getMidPriceNoLock()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, nil, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	_, err = depth.getMidPriceNoLock()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mid, err := depth.getMidPriceNoLock()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if mid != 1336.5 {
		t.Fatalf("received: '%v' but expected: '%v'", mid, 1336.5)
	}
}

func TestGetBestBidASk_Depth(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().GetBestBid()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	_, err = getInvalidDepth().GetBestAsk()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)
	_, err = depth.GetBestBid()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	_, err = depth.GetBestAsk()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}
	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	mid, err := depth.GetBestBid()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}
	if mid != 1336 {
		t.Fatalf("received: '%v' but expected: '%v'", mid, 1336)
	}
	mid, err = depth.GetBestAsk()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}
	if mid != 1337 {
		t.Fatalf("received: '%v' but expected: '%v'", mid, 1337)
	}
}

func TestGetSpreadAmount(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().GetSpreadAmount()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.GetSpreadAmount()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(nil, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	_, err = depth.GetSpreadAmount()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	spread, err := depth.GetSpreadAmount()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if spread != 1 {
		t.Fatalf("received: '%v' but expected: '%v'", spread, 1)
	}
}

func TestGetSpreadPercentage(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().GetSpreadPercentage()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.GetSpreadPercentage()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(nil, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	_, err = depth.GetSpreadPercentage()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	spread, err := depth.GetSpreadPercentage()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if spread != 0.07479431563201197 {
		t.Fatalf("received: '%v' but expected: '%v'", spread, 0.07479431563201197)
	}
}

func TestGetImbalance_Depth(t *testing.T) {
	t.Parallel()
	_, err := getInvalidDepth().GetImbalance()
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, err = depth.GetImbalance()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(nil, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	_, err = depth.GetImbalance()
	if !errors.Is(err, errNoLiquidity) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errNoLiquidity)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}
	imbalance, err := depth.GetImbalance()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if imbalance != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", imbalance, 0)
	}
}

func TestGetTranches(t *testing.T) {
	t.Parallel()
	_, _, err := getInvalidDepth().GetTranches(0)
	if !errors.Is(err, ErrOrderbookInvalid) {
		t.Fatalf("received: '%v' but expected: '%v'", err, ErrOrderbookInvalid)
	}

	depth := NewDepth(id)

	_, _, err = depth.GetTranches(-1)
	if !errors.Is(err, errInvalidBookDepth) {
		t.Fatalf("received: '%v' but expected: '%v'", err, errInvalidBookDepth)
	}

	askT, bidT, err := depth.GetTranches(0)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if len(askT) != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", len(askT), 0)
	}

	if len(bidT) != 0 {
		t.Fatalf("received: '%v' but expected: '%v'", len(bidT), 0)
	}

	err = depth.LoadSnapshot(bid, ask, 0, time.Now(), true)
	if err != nil {
		t.Fatal(err)
	}

	askT, bidT, err = depth.GetTranches(0)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if len(askT) != 20 {
		t.Fatalf("received: '%v' but expected: '%v'", len(askT), 20)
	}

	if len(bidT) != 20 {
		t.Fatalf("received: '%v' but expected: '%v'", len(bidT), 20)
	}

	askT, bidT, err = depth.GetTranches(5)
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if len(askT) != 5 {
		t.Fatalf("received: '%v' but expected: '%v'", len(askT), 5)
	}

	if len(bidT) != 5 {
		t.Fatalf("received: '%v' but expected: '%v'", len(bidT), 5)
	}
}

func TestGetPair(t *testing.T) {
	t.Parallel()
	depth := NewDepth(id)

	_, err := depth.GetPair()
	if !errors.Is(err, currency.ErrCurrencyPairEmpty) {
		t.Fatalf("received: '%v' but expected: '%v'", err, currency.ErrCurrencyPairEmpty)
	}

	expected := currency.NewPair(currency.BTC, currency.WABI)
	depth.pair = expected

	pair, err := depth.GetPair()
	if !errors.Is(err, nil) {
		t.Fatalf("received: '%v' but expected: '%v'", err, nil)
	}

	if !pair.Equal(expected) {
		t.Fatalf("received: '%v' but expected: '%v'", pair, expected)
	}
}

func getInvalidDepth() *Depth {
	depth := NewDepth(id)
	_ = depth.Invalidate(errors.New("invalid reasoning"))
	return depth
}
