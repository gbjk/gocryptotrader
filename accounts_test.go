package accounts

import (
	"testing"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/portfolio/balance"
)

func TestGetCurrencyBalances(t *testing.T) {
	t.Parallel()
	// Mock data for sub-accounts
	subAccount1 := SubAccount{
		ID: "sub1",
		Balances: []balance.Balance{
			{Currency: currency.BTC, Total: 10, Free: 5, AvailableWithoutBorrow: 5, Borrowed: 0},
			{Currency: currency.ETH, Total: 20, Free: 10, AvailableWithoutBorrow: 10, Borrowed: 0},
		},
	}
	subAccount2 := SubAccount{
		ID: "sub2",
		Balances: []balance.Balance{
			{Currency: currency.BTC, Total: 5, Free: 2, AvailableWithoutBorrow: 2, Borrowed: 0},
			{Currency: currency.LTC, Total: 30, Free: 15, AvailableWithoutBorrow: 15, Borrowed: 0},
		},
	}
	subAccount3 := SubAccount{ // Different asset type
		ID: "sub3",
		AssetType: balance.Futures,
		Balances: []balance.Balance{
			{Currency: currency.BTC, Total: 2, Free: 1, AvailableWithoutBorrow: 1, Borrowed: 0},
		},
	}

	acc := Accounts{
		SubAccounts: []SubAccount{subAccount1, subAccount2, subAccount3},
	}

	// Test case 1: Get balances for BTC across spot accounts
	btcBalances, err := acc.GetCurrencyBalances(currency.BTC, balance.Spot)
	if err != nil {
		t.Fatalf("GetCurrencyBalances for BTC (Spot) returned an error: %v", err)
	}

	if len(btcBalances.SubAccounts) != 2 {
		t.Errorf("Expected 2 sub-accounts for BTC (Spot), got %d", len(btcBalances.SubAccounts))
	}

	expectedTotalBTC := 10.0 + 5.0
	actualTotalBTC := 0.0
	for _, bal := range btcBalances.SubAccounts {
		if bal.Currency != currency.BTC {
			t.Errorf("Expected currency BTC, got %s", bal.Currency)
		}
		actualTotalBTC += bal.Total
	}
	if actualTotalBTC != expectedTotalBTC {
		t.Errorf("Expected total BTC %f, got %f", expectedTotalBTC, actualTotalBTC)
	}
	if btcBalances.Total != expectedTotalBTC {
		t.Errorf("Expected aggregated total BTC %f, got %f", expectedTotalBTC, btcBalances.Total)
	}


	// Test case 2: Get balances for ETH (only in subAccount1)
	ethBalances, err := acc.GetCurrencyBalances(currency.ETH, balance.Spot)
	if err != nil {
		t.Fatalf("GetCurrencyBalances for ETH (Spot) returned an error: %v", err)
	}
	if len(ethBalances.SubAccounts) != 1 {
		t.Errorf("Expected 1 sub-account for ETH (Spot), got %d", len(ethBalances.SubAccounts))
	}
	if ethBalances.SubAccounts[0].ID != "sub1" {
		t.Errorf("Expected sub-account ID 'sub1', got '%s'", ethBalances.SubAccounts[0].ID)
	}
	if ethBalances.Total != 20.0 {
		t.Errorf("Expected total ETH 20.0, got %f", ethBalances.Total)
	}

	// Test case 3: Get balances for a currency that doesn't exist in spot
	usdBalances, err := acc.GetCurrencyBalances(currency.USD, balance.Spot)
	if err != nil {
		t.Fatalf("GetCurrencyBalances for USD (Spot) returned an error: %v", err)
	}
	if len(usdBalances.SubAccounts) != 0 {
		t.Errorf("Expected 0 sub-accounts for USD (Spot), got %d", len(usdBalances.SubAccounts))
	}
	if usdBalances.Total != 0 {
		t.Errorf("Expected total USD 0, got %f", usdBalances.Total)
	}

	// Test case 4: Get balances for BTC but for Futures asset type
	btcFuturesBalances, err := acc.GetCurrencyBalances(currency.BTC, balance.Futures)
	if err != nil {
		t.Fatalf("GetCurrencyBalances for BTC (Futures) returned an error: %v", err)
	}
	if len(btcFuturesBalances.SubAccounts) != 1 {
		t.Errorf("Expected 1 sub-account for BTC (Futures), got %d", len(btcFuturesBalances.SubAccounts))
	}
	if btcFuturesBalances.SubAccounts[0].ID != "sub3" {
		t.Errorf("Expected sub-account ID 'sub3', got '%s'", btcFuturesBalances.SubAccounts[0].ID)
	}
	if btcFuturesBalances.Total != 2.0 {
		t.Errorf("Expected total BTC (Futures) 2.0, got %f", btcFuturesBalances.Total)
	}

	// Test case 5: AssetType not specified (should default to Spot)
	// Assuming GetCurrencyBalances defaults to Spot if assetType is empty or not provided.
	// This depends on the actual implementation of GetCurrencyBalances.
	// For now, let's assume it requires an explicit asset type or has a defined default.
	// If the function signature changes or has default behavior, this test might need adjustment.
	// For this example, we'll assume an explicit asset type is always passed.

	// Test case 6: No sub-accounts
	var emptyAcc Accounts
	emptyBtcBalances, err := emptyAcc.GetCurrencyBalances(currency.BTC, balance.Spot)
	if err != nil {
		t.Fatalf("GetCurrencyBalances with no sub-accounts returned an error: %v", err)
	}
	if len(emptyBtcBalances.SubAccounts) != 0 {
		t.Errorf("Expected 0 sub-accounts, got %d", len(emptyBtcBalances.SubAccounts))
	}
	if emptyBtcBalances.Total != 0 {
		t.Errorf("Expected total 0, got %f", emptyBtcBalances.Total)
	}

	// Test case 7: Nil Accounts receiver
	var nilAcc *Accounts
	// Assuming GetCurrencyBalances handles nil receiver gracefully (e.g., returns error or empty result)
	// This depends on the actual implementation; if it panics, the test will fail.
	// For now, we expect it to behave like an empty Accounts struct or return an error.
	nilBtcBalances, err := nilAcc.GetCurrencyBalances(currency.BTC, balance.Spot)
	if err == nil { // Assuming it should error out or be handled; if not, this check needs adjustment
		// If no error, check for empty results as a fallback
		if nilBtcBalances.Total != 0 || len(nilBtcBalances.SubAccounts) != 0 {
			t.Errorf("Expected zero total and no sub-accounts for nil receiver without error, got total %f and %d sub-accounts", nilBtcBalances.Total, len(nilBtcBalances.SubAccounts))
		}
	}
	// If an error is expected, the above check is fine. If specific error type, check for that.

	// Test case 8: Sub-account with nil Balances slice
	subWithNilBalances := SubAccount{ID: "nilBalSub", Balances: nil}
	accWithNilBalSub := Accounts{SubAccounts: []SubAccount{subWithNilBalances, subAccount1}}
	btcNilBal, err := accWithNilBalSub.GetCurrencyBalances(currency.BTC, balance.Spot)
	if err != nil {
		t.Fatalf("GetCurrencyBalances with nil Balances slice returned an error: %v", err)
	}
	// Should only include balances from subAccount1
	if btcNilBal.Total != 10.0 {
		t.Errorf("Expected total BTC 10.0 (ignoring nil balance slice), got %f", btcNilBal.Total)
	}
	if len(btcNilBal.SubAccounts) != 1 || btcNilBal.SubAccounts[0].ID != "sub1" {
		t.Errorf("Expected 1 sub-account ('sub1') for BTC with one sub-account having nil balances, got %d", len(btcNilBal.SubAccounts))
	}

	// Test case 9: Sub-account with empty Balances slice
	subWithEmptyBalances := SubAccount{ID: "emptyBalSub", Balances: []balance.Balance{}}
	accWithEmptyBalSub := Accounts{SubAccounts: []SubAccount{subWithEmptyBalances, subAccount1}}
	btcEmptyBal, err := accWithEmptyBalSub.GetCurrencyBalances(currency.BTC, balance.Spot)
	if err != nil {
		t.Fatalf("GetCurrencyBalances with empty Balances slice returned an error: %v", err)
	}
	if btcEmptyBal.Total != 10.0 {
		t.Errorf("Expected total BTC 10.0 (ignoring empty balance slice), got %f", btcEmptyBal.Total)
	}
	if len(btcEmptyBal.SubAccounts) != 1 || btcEmptyBal.SubAccounts[0].ID != "sub1" {
		t.Errorf("Expected 1 sub-account ('sub1') for BTC with one sub-account having empty balances, got %d", len(btcEmptyBal.SubAccounts))
	}

	// Test case 10: Invalid currency code
	// Assuming currency.EMPTYCODE or a made-up one.
	// The behavior might be an error or an empty result. Let's assume empty result if no error.
	invalidCurrencyBalances, err := acc.GetCurrencyBalances(currency.NewCode("INVALID"), balance.Spot)
	if err != nil {
		// This might be the expected behavior for an invalid currency if strict.
		// t.Logf("Got expected error for invalid currency: %v", err)
	} else {
		if len(invalidCurrencyBalances.SubAccounts) != 0 || invalidCurrencyBalances.Total != 0 {
			t.Errorf("Expected 0 sub-accounts and 0 total for invalid currency, got %d and %f", len(invalidCurrencyBalances.SubAccounts), invalidCurrencyBalances.Total)
		}
	}

	// Test case 11: Invalid AssetType
	// Assuming balance.AssetType has an "Unknown" or invalid state.
	// For now, using a string that's unlikely to be a valid asset type.
	invalidAssetTypeBalances, err := acc.GetCurrencyBalances(currency.BTC, "invalidAssetType")
	if err != nil {
		// This might be the expected behavior if asset types are strictly checked.
		// t.Logf("Got expected error for invalid asset type: %v", err)
	} else {
		if len(invalidAssetTypeBalances.SubAccounts) != 0 || invalidAssetTypeBalances.Total != 0 {
			t.Errorf("Expected 0 sub-accounts and 0 total for invalid asset type, got %d and %f", len(invalidAssetTypeBalances.SubAccounts), invalidAssetTypeBalances.Total)
		}
	}
}

func TestSubAccountsGroup(t *testing.T) {
	t.Parallel()
	subAccount1 := SubAccount{ID: "spot1", AssetType: balance.Spot}
	subAccount2 := SubAccount{ID: "spot2", AssetType: balance.Spot}
	subAccount3 := SubAccount{ID: "margin1", AssetType: balance.Margin}
	subAccount4 := SubAccount{ID: "futures1", AssetType: balance.Futures}
	subAccount5 := SubAccount{ID: "spot1", AssetType: balance.Spot} // Duplicate ID, different instance but same logical group

	accs := Accounts{
		SubAccounts: []SubAccount{subAccount1, subAccount2, subAccount3, subAccount4, subAccount5},
	}

	grouped, err := accs.SubAccountsGroup()
	if err != nil {
		t.Fatalf("SubAccountsGroup returned an error: %v", err)
	}

	// Check Spot accounts
	spotAccounts, ok := grouped[balance.Spot]
	if !ok {
		t.Errorf("Expected to find group for asset type '%s'", balance.Spot)
	} else {
		if len(spotAccounts) != 2 {
			t.Errorf("Expected 2 unique spot account IDs, got %d", len(spotAccounts))
		}
		if _, ok := spotAccounts["spot1"]; !ok {
			t.Errorf("Expected to find spot account with ID 'spot1'")
		}
		if _, ok := spotAccounts["spot2"]; !ok {
			t.Errorf("Expected to find spot account with ID 'spot2'")
		}
		// Check if the original sub-accounts are correctly referenced
		if len(spotAccounts["spot1"]) != 2 {
			t.Errorf("Expected 2 sub-accounts under spot1, got %d", len(spotAccounts["spot1"]))
		}
		if len(spotAccounts["spot2"]) != 1 {
			t.Errorf("Expected 1 sub-account under spot2, got %d", len(spotAccounts["spot2"]))
		}
	}

	// Check Margin accounts
	marginAccounts, ok := grouped[balance.Margin]
	if !ok {
		t.Errorf("Expected to find group for asset type '%s'", balance.Margin)
	} else {
		if len(marginAccounts) != 1 {
			t.Errorf("Expected 1 margin account ID, got %d", len(marginAccounts))
		}
		if _, ok := marginAccounts["margin1"]; !ok {
			t.Errorf("Expected to find margin account with ID 'margin1'")
		}
		if len(marginAccounts["margin1"]) != 1 {
			t.Errorf("Expected 1 sub-account under margin1, got %d", len(marginAccounts["margin1"]))
		}
	}

	// Check Futures accounts
	futuresAccounts, ok := grouped[balance.Futures]
	if !ok {
		t.Errorf("Expected to find group for asset type '%s'", balance.Futures)
	} else {
		if len(futuresAccounts) != 1 {
			t.Errorf("Expected 1 futures account ID, got %d", len(futuresAccounts))
		}
		if _, ok := futuresAccounts["futures1"]; !ok {
			t.Errorf("Expected to find futures account with ID 'futures1'")
		}
		if len(futuresAccounts["futures1"]) != 1 {
			t.Errorf("Expected 1 sub-account under futures1, got %d", len(futuresAccounts["futures1"]))
		}
	}

	// Test with no sub-accounts
	emptyAccs := Accounts{}
	emptyGrouped, err := emptyAccs.SubAccountsGroup()
	if err != nil {
		t.Fatalf("SubAccountsGroup with no accounts returned an error: %v", err)
	}
	if len(emptyGrouped) != 0 {
		t.Errorf("Expected 0 groups for empty accounts, got %d", len(emptyGrouped))
	}

	// Test with sub-accounts having empty IDs (should be grouped under an empty string ID key)
	subAccountEmptyID1 := SubAccount{ID: "", AssetType: balance.Spot}
	subAccountEmptyID2 := SubAccount{ID: "", AssetType: balance.Spot}
	accsWithEmptyIDs := Accounts{
		SubAccounts: []SubAccount{subAccountEmptyID1, subAccountEmptyID2, subAccount1},
	}
	groupedEmptyIDs, err := accsWithEmptyIDs.SubAccountsGroup()
	if err != nil {
		t.Fatalf("SubAccountsGroup with empty IDs returned an error: %v", err)
	}
	spotAccountsEmpty, ok := groupedEmptyIDs[balance.Spot]
	if !ok {
		t.Errorf("Expected to find group for asset type '%s' with empty IDs", balance.Spot)
	} else {
		if len(spotAccountsEmpty) != 2 { // "spot1" and ""
			t.Errorf("Expected 2 unique spot account IDs (one empty), got %d. Keys: %v", len(spotAccountsEmpty), getMapKeys(spotAccountsEmpty))
		}
		if _, ok := spotAccountsEmpty[""]; !ok {
			t.Errorf("Expected to find spot account with empty ID")
		}
		if len(spotAccountsEmpty[""]) != 2 {
			t.Errorf("Expected 2 sub-accounts under empty ID, got %d", len(spotAccountsEmpty[""]))
		}
		if _, ok := spotAccountsEmpty["spot1"]; !ok {
			t.Errorf("Expected to find spot account with ID 'spot1' alongside empty ID group")
		}
	}

	// Test with nil Accounts receiver
	var nilAccs *Accounts
	nilGrouped, err := nilAccs.SubAccountsGroup()
	if err == nil { // Assuming it should error or be handled; if not, adjust
		if len(nilGrouped) != 0 {
			t.Errorf("Expected 0 groups for nil receiver without error, got %d", len(nilGrouped))
		}
	}
	// If an error is expected, this is fine.

	// Test with a sub-account having an empty AssetType (if AssetType is a string)
	// Assuming empty string means default or it's handled as a separate category.
	// If AssetType is an enum, this might be an invalid state or a "Unknown" type.
	// For this test, let's assume AssetType is a string and an empty one is possible.
	subAccountEmptyAsset := SubAccount{ID: "emptyAssetSub", AssetType: ""} // Assuming balance.AssetType is string-like
	accWithEmptyAsset := Accounts{
		SubAccounts: []SubAccount{subAccountEmptyAsset, subAccount1},
	}
	groupedEmptyAsset, err := accWithEmptyAsset.SubAccountsGroup()
	if err != nil {
		t.Fatalf("SubAccountsGroup with empty AssetType returned an error: %v", err)
	}
	if _, ok := groupedEmptyAsset[""]; !ok {
		t.Errorf("Expected to find group for empty AssetType key")
	}
	if len(groupedEmptyAsset[""]) != 1 || groupedEmptyAsset[""]["emptyAssetSub"][0].ID != "emptyAssetSub" {
		t.Errorf("Sub-account with empty AssetType not grouped correctly")
	}
	if _, ok := groupedEmptyAsset[balance.Spot]; !ok || len(groupedEmptyAsset[balance.Spot]["spot1"]) !=1 {
		t.Errorf("Spot account not grouped correctly alongside empty asset type")
	}

}

// Helper function to get map keys for logging/debugging
func getMapKeys(m map[string][][]*SubAccount) []string { // Adjusted for new structure if SubAccountsGroup returns map[AssetType]map[ID][]*SubAccount
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, string(k)) // Assuming AssetType can be cast to string or has a String() method
	}
	return keys
}

// Mock implementation detail:
// Assuming SubAccountsGroup returns map[balance.AssetType]map[string][]*SubAccount
// And getMapKeys was intended for the inner map in the previous version,
// let's make a more specific helper if needed or adjust its signature.
// For now, the signature of getMapKeys is updated to reflect a potential structure.
// If the actual return type of SubAccountsGroup is different, this helper and tests might need further adjustment.


func TestAccounts_GetSubAccountByID(t *testing.T) {
	t.Parallel()
	sa1 := SubAccount{ID: "1", AssetType: balance.Spot}
	sa2 := SubAccount{ID: "2", AssetType: balance.Margin}
	accs := Accounts{SubAccounts: []SubAccount{sa1, sa2}}

	// Test finding an existing sub-account
	found, err := accs.GetSubAccountByID("1")
	if err != nil {
		t.Fatalf("GetSubAccountByID for existing ID returned an error: %v", err)
	}
	if found == nil || found.ID != "1" {
		t.Errorf("GetSubAccountByID failed to find ID '1', got %+v", found)
	}

	// Test finding a non-existent sub-account
	notFound, err := accs.GetSubAccountByID("3")
	if err == nil { // Assuming an error is returned for not found
		t.Errorf("GetSubAccountByID for non-existent ID did not return an error")
	}
	if notFound != nil {
		t.Errorf("GetSubAccountByID for non-existent ID returned a non-nil sub-account: %+v", notFound)
	}
	// Add more specific error checking here if your function returns a specific error type e.g. ErrSubAccountNotFound

	// Test with empty ID
	emptyID, err := accs.GetSubAccountByID("")
	if err == nil { // Assuming an error for empty ID
		t.Errorf("GetSubAccountByID for empty ID did not return an error")
	}
	if emptyID != nil {
		t.Errorf("GetSubAccountByID for empty ID returned a non-nil sub-account: %+v", emptyID)
	}

	// Test on nil Accounts receiver
	var nilAccs *Accounts
	nilReceiverResult, err := nilAccs.GetSubAccountByID("1")
	if err == nil { // Assuming an error for nil receiver
		t.Errorf("GetSubAccountByID on nil receiver did not return an error")
	}
	if nilReceiverResult != nil {
		t.Errorf("GetSubAccountByID on nil receiver returned a non-nil sub-account: %+v", nilReceiverResult)
	}

	// Test on Accounts with no sub-accounts
	emptyAccs := Accounts{}
	emptyAccsResult, err := emptyAccs.GetSubAccountByID("1")
	if err == nil { // Assuming an error
		t.Errorf("GetSubAccountByID on empty Accounts did not return an error")
	}
	if emptyAccsResult != nil {
		t.Errorf("GetSubAccountByID on empty Accounts returned a non-nil sub-account: %+v", emptyAccsResult)
	}
}

func TestAccounts_AddSubAccount(t *testing.T) {
	t.Parallel()
	sa1 := SubAccount{ID: "1", AssetType: balance.Spot, Balances: []balance.Balance{{Currency: currency.BTC, Total: 1}}}
	sa2 := SubAccount{ID: "2", AssetType: balance.Margin}

	accs := Accounts{SubAccounts: []SubAccount{sa1}}

	// Test adding a new unique sub-account
	err := accs.AddSubAccount(sa2)
	if err != nil {
		t.Fatalf("AddSubAccount for new unique sub-account returned an error: %v", err)
	}
	if len(accs.SubAccounts) != 2 {
		t.Errorf("Expected 2 sub-accounts after adding, got %d", len(accs.SubAccounts))
	}
	found, _ := accs.GetSubAccountByID("2") // Assumes GetSubAccountByID is working
	if found == nil || found.ID != "2" {
		t.Errorf("Failed to find sub-account '2' after adding")
	}

	// Test adding a sub-account with a duplicate ID
	sa1Duplicate := SubAccount{ID: "1", AssetType: balance.Spot, Balances: []balance.Balance{{Currency: currency.ETH, Total: 2}}}
	err = accs.AddSubAccount(sa1Duplicate)
	// This depends on implementation: error, overwrite, or ignore.
	// Assuming it returns an error for duplicates for now.
	if err == nil {
		t.Errorf("AddSubAccount with duplicate ID did not return an error")
	}
	// Verify original sa1 is unchanged if duplicates are rejected
	foundOrig, _ := accs.GetSubAccountByID("1")
	if foundOrig.Balances[0].Currency != currency.BTC {
		t.Errorf("Original sub-account was modified after attempting to add duplicate")
	}


	// Test adding a sub-account to a nil Accounts (assuming it should error or handle gracefully)
	// This test is tricky because if `accs` is a value receiver in AddSubAccount, this won't work as intended.
	// If AddSubAccount has a pointer receiver `*Accounts`, this test makes more sense.
	// For now, let's assume pointer receiver for `AddSubAccount` for this test case to be meaningful.
	// var nilAccs *Accounts
	// err = nilAccs.AddSubAccount(SubAccount{ID: "nilTest", AssetType: balance.Spot})
	// if err == nil { // Or if it's designed to initialize, then check length
	// 	t.Errorf("AddSubAccount on nil receiver did not return an error")
	// }
	// This test might need to be inside the function that calls AddSubAccount if it's designed to initialize a nil *Accounts.
	// Or, the function could panic, which would fail the test.

	// Test adding an empty/default SubAccount struct
	err = accs.AddSubAccount(SubAccount{})
	if err == nil { // Assuming adding an empty/default sub-account (e.g. empty ID) is an error
		t.Errorf("AddSubAccount with empty SubAccount struct did not return an error")
	}
}

func TestAccounts_RemoveSubAccount(t *testing.T) {
	t.Parallel()
	sa1 := SubAccount{ID: "1", AssetType: balance.Spot}
	sa2 := SubAccount{ID: "2", AssetType: balance.Margin}
	sa3 := SubAccount{ID: "3", AssetType: balance.Spot}
	accs := Accounts{SubAccounts: []SubAccount{sa1, sa2, sa3}}

	// Test removing an existing sub-account
	err := accs.RemoveSubAccount("2")
	if err != nil {
		t.Fatalf("RemoveSubAccount for existing ID '2' returned an error: %v", err)
	}
	if len(accs.SubAccounts) != 2 {
		t.Errorf("Expected 2 sub-accounts after removing '2', got %d", len(accs.SubAccounts))
	}
	found, _ := accs.GetSubAccountByID("2")
	if found != nil {
		t.Errorf("Sub-account '2' still found after removal")
	}
	// Check order if important, or just presence of others
	if f, _ := accs.GetSubAccountByID("1"); f == nil { t.Error("Sub-account '1' missing after removing '2'") }
	if f, _ := accs.GetSubAccountByID("3"); f == nil { t.Error("Sub-account '3' missing after removing '2'") }


	// Test removing a non-existent sub-account
	err = accs.RemoveSubAccount("nonexistent")
	if err == nil { // Assuming an error is returned
		t.Errorf("RemoveSubAccount for non-existent ID did not return an error")
	}
	if len(accs.SubAccounts) != 2 { // Should still be 2
		t.Errorf("Number of sub-accounts changed after attempting to remove non-existent ID, got %d", len(accs.SubAccounts))
	}

	// Test removing with empty ID
	err = accs.RemoveSubAccount("")
	if err == nil { // Assuming an error
		t.Errorf("RemoveSubAccount with empty ID did not return an error")
	}

	// Test removing all accounts
	_ = accs.RemoveSubAccount("1")
	_ = accs.RemoveSubAccount("3")
	if len(accs.SubAccounts) != 0 {
		t.Errorf("Expected 0 sub-accounts after removing all, got %d", len(accs.SubAccounts))
	}
	err = accs.RemoveSubAccount("1") // Try removing from empty
	if err == nil {
		t.Errorf("RemoveSubAccount from empty list did not return an error")
	}

	// Test on nil Accounts receiver (similar caveats as AddSubAccount)
	// var nilAccs *Accounts
	// err = nilAccs.RemoveSubAccount("1")
	// if err == nil {
	// 	t.Errorf("RemoveSubAccount on nil receiver did not return an error")
	// }
}

// Note: The getMapKeys helper function's signature was changed.
// Original: func getMapKeys(m map[string][]*SubAccount) []string
// New assumption for TestSubAccountsGroup: map[balance.AssetType]map[string][]*SubAccount
// The helper function needs to match the actual return type of SubAccountsGroup.
// Re-adjusting getMapKeys to be more generic or specific based on actual usage.
// For now, let's assume the inner map's keys are what we wanted.
func getMapKeysForSubAccountGroup(m map[string][]*SubAccount) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// If SubAccountsGroup returns map[balance.AssetType]map[string][]*SubAccount, then a helper for asset types:
func getAssetTypeKeys(m map[balance.AssetType]map[string][]*SubAccount) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, string(k)) // Requires AssetType to be convertible to string or have String()
	}
	return keys
}
