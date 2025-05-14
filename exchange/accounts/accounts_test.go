package accounts

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/key"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/dispatch"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

type mockEx struct {
	name string
}

func (m *mockEx) GetName() string {
	return m.name
}

func TestNewAccounts(t *testing.T) {
	t.Parallel()
	a, err := NewAccounts(&mockEx{"mocky"}, dispatch.GetNewMux(nil))
	require.NoError(t, err)
	require.NotNil(t, a)
	assert.Equal(t, "mocky", a.Exchange.GetName(), "Exchande should set correctly")
	assert.NotNil(t, a.subAccounts, "subAccounts should be initialized")
	assert.NotEmpty(t, a.routingID, "ID should not be empty")
	assert.NotNil(t, a.mux, "mux should be set correctly")
	_, err = NewAccounts(nil, dispatch.GetNewMux(nil))
	assert.ErrorIs(t, err, common.ErrNilPointer)
	_, err = NewAccounts(&mockEx{"mocky"}, nil)
	assert.ErrorContains(t, err, "nil pointer: *dispatch.Mux")
}

func TestMustNewAccounts(t *testing.T) {
	t.Parallel()
	a := MustNewAccounts(&mockEx{"mocky"}, dispatch.GetNewMux(nil))
	require.NotNil(t, a)
	require.Panics(t, func() { _ = MustNewAccounts(nil, nil) })
}

// TestSubscribe ensures that Subscribe returns a subscription channel
// See TestSave and TestUpdate for exercising publish to subscribers
func TestSubscribe(t *testing.T) {
	t.Parallel()

	err := dispatch.Start(dispatch.DefaultMaxWorkers, dispatch.DefaultJobsLimit)
	require.NoError(t, common.ExcludeError(err, dispatch.ErrDispatcherAlreadyRunning), "dispatch.Start must not error")

	p, err := MustNewAccounts(&mockEx{}, dispatch.GetNewMux(nil)).Subscribe()
	require.NoError(t, err)
	require.NotNil(t, p, "Subscribe must return a pipe")
	require.Empty(t, p.Channel(), "Pipe must be empty before Saving anything")
}

func TestCurrencyBalancesPublic(t *testing.T) {
	t.Parallel()

	a := MustNewAccounts(&mockEx{}, dispatch.GetNewMux(nil))
	c := &Credentials{Key: "Creds"}
	bals := a.currencyBalances(c, "", asset.Spot)
	require.NotNil(t, bals)
	b := bals.balance(currency.BTC)
	updated, err := b.update(Balance{Total: 4.0, UpdatedAt: time.Now()})
	require.NoError(t, err, "update must not error")
	require.True(t, updated, "update must return true")
	spew.Dump(a.subAccounts[*c][key.SubAccountAsset{Asset: asset.Spot}], b)
}

/*
func TestCurrencyBalances(t *testing.T) {
	t.Parallel()

	n := time.Now()
	creds := Credentials{Key: "cred1"}
	a := Accounts{
		subAccounts: map[Credentials]map[key.SubAccountAsset]CurrencyBalances{
			creds: {
				{"a", asset.Spot}: {
					currency.BTC.Item: &balance{internal: Balance{Currency: currency.BTC, Total: 1, Free: 1, UpdatedAt: n}},
					currency.ETH.Item: &balance{internal: Balance{Currency: currency.ETH, Total: 10, Free: 10, UpdatedAt: n}},
				},
				{"b", asset.Spot}: {
					currency.BTC.Item: &balance{internal: Balance{Currency: currency.BTC, Total: 2, Free: 2, UpdatedAt: n}},
					currency.LTC.Item: &balance{internal: Balance{Currency: currency.LTC, Total: 5, Free: 5, UpdatedAt: n}},
				},
				{"a", asset.Margin}: {
					currency.USDT.Item: &balance{internal: Balance{Currency: currency.USDT, Total: 100, Free: 100, UpdatedAt: n}},
					currency.BTC.Item:  &balance{internal: Balance{Currency: currency.BTC, Total: 0.5, Free: 0.5, UpdatedAt: n}},
				},
			},
		},
	}

		expectedBalances := map[currency.Code]Balance{
			currency.BTC:  {Currency: currency.BTC, Total: 3.5, Free: 3.5},
			currency.ETH:  {Currency: currency.ETH, Total: 10, Free: 10},
			currency.LTC:  {Currency: currency.LTC, Total: 5, Free: 5},
			currency.USDT: {Currency: currency.USDT, Total: 100, Free: 100},
		}

	b := a.CurrencyBalances()
	require.Equal(t, 4, len(b), "Must get 4 balances")
			if len(result) != len(expectedBalances) {
				t.Errorf("Expected %d aggregated currencies, got %d. Result: %+v", len(expectedBalances), len(result), result)
			}
			for c, expectedBal := range expectedBalances {
				actualBal, ok := result[c]
				if !ok {
					t.Errorf("Expected currency %s missing from result", c)
					continue
				}
				if actualBal.Total != expectedBal.Total {
					t.Errorf("For currency %s, expected Total %f, got %f", c, expectedBal.Total, actualBal.Total)
				}
				if actualBal.Free != expectedBal.Free {
					t.Errorf("For currency %s, expected Free %f, got %f", c, expectedBal.Free, actualBal.Free)
				}
				if actualBal.Hold != expectedBal.Hold {
					t.Errorf("For currency %s, expected Hold %f, got %f", c, expectedBal.Hold, actualBal.Hold)
				}
			}
		})

		t.Run("no subaccounts", func(t *testing.T) {
			emptyAccs := Accounts{
				subAccounts: make(map[Credentials]map[key.SubAccountAsset]CurrencyBalances),
			}
			result := emptyAccs.CurrencyBalances()
			if len(result) != 0 {
				t.Errorf("Expected 0 balances for empty accounts, got %d", len(result))
			}
		})

		t.Run("subaccounts exist but no balances for some currencies", func(t *testing.T) {
			complexAccs := Accounts{
				subAccounts: make(map[Credentials]map[key.SubAccountAsset]CurrencyBalances),
			}
			complexAccs.subAccounts[creds] = make(map[key.SubAccountAsset]CurrencyBalances)
			saKeyTest := key.SubAccountAsset{SubAccount: "testSub", Asset: asset.Spot}
			complexAccs.subAccounts[creds][saKeyTest] = CurrencyBalances{
				currency.BTC.Item: &balance{internal: Balance{Currency: currency.BTC, Total: 1, UpdatedAt: time.Now()}},
			}
			saKeyEmpty := key.SubAccountAsset{SubAccount: "emptySub", Asset: asset.Spot}
			complexAccs.subAccounts[creds][saKeyEmpty] = make(CurrencyBalances)

			result := complexAccs.CurrencyBalances()
			if len(result) != 1 {
				t.Errorf("Expected 1 currency (BTC), got %d. Result: %+v", len(result), result)
			}
			if _, ok := result[currency.BTC]; !ok {
				t.Error("Expected BTC to be present in the result.")
			}
		})
}

// TestAccounts_Save tests the Save method of the Accounts struct.
func TestAccounts_Save(t *testing.T) {
	t.Parallel()

	mockCreds := Credentials{Key: "saveCreds"}
	var emptyCredsInstance Credentials
	var invalidAssetType asset.Item

	mockEx := &mockEx{name: "TestSaveExchange"}

	initialBTCHoldings := Balance{Currency: currency.BTC, Total: 1, Free: 1, UpdatedAt: time.Unix(100, 0)}
	initialETHHoldings := Balance{Currency: currency.ETH, Total: 10, Free: 10, UpdatedAt: time.Unix(100, 0)}
	updatedBTCHoldings := Balance{Currency: currency.BTC, Total: 2, Free: 2, UpdatedAt: time.Unix(200, 0)}
	newLTCHoldings := Balance{Currency: currency.LTC, Total: 5, Free: 5, UpdatedAt: time.Unix(200, 0)}

	tests := []struct {
		name                string
		credsToSave         *Credentials
		subAccountsToSave   []SubAccount
		initialAccountState func() *Accounts
		expectedErr         bool
		verifyState         func(t *testing.T, accs *Accounts, errEncountered error)
	}{
		{
			name:        "Save_NewSubAccountData_NewCredentials",
			credsToSave: &mockCreds,
			subAccountsToSave: []SubAccount{
				{ID: "spot1", AssetType: asset.Spot, Currencies: []Balance{initialBTCHoldings}},
			},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				accs, _ := NewAccounts(mockEx, mux)
				return accs
			},
			expectedErr: false,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered != nil {
					t.Fatalf("Unexpected error: %v", errEncountered)
				}
				b, err := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.BTC)
				if err != nil {
					t.Fatalf("Failed to get balance for verification: %v", err)
				}
				if b.Total != initialBTCHoldings.Total {
					t.Errorf("Expected BTC total %f, got %f", initialBTCHoldings.Total, b.Total)
				}
			},
		},
		{
			name:        "Save_UpdateExisting_AddNewCurrency_RemoveOldOne",
			credsToSave: &mockCreds,
			subAccountsToSave: []SubAccount{
				{ID: "spot1", AssetType: asset.Spot, Currencies: []Balance{updatedBTCHoldings, newLTCHoldings}},
			},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				accs, _ := NewAccounts(mockEx, mux)
				accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
					{SubAccount: "spot1", Asset: asset.Spot}: {
						currency.BTC.Item: &balance{internal: initialBTCHoldings},
						currency.ETH.Item: &balance{internal: initialETHHoldings},
					},
				}
				return accs
			},
			expectedErr: false,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered != nil {
					t.Fatalf("Unexpected error: %v", errEncountered)
				}
				btcBal, _ := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.BTC)
				if btcBal.Total != updatedBTCHoldings.Total {
					t.Errorf("Expected updated BTC total %f, got %f", updatedBTCHoldings.Total, btcBal.Total)
				}
				_, errEth := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.ETH)
				if errEth == nil {
					t.Error("Expected ETH to be removed, but GetBalance found it.")
				}
				ltcBal, _ := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.LTC)
				if ltcBal.Total != newLTCHoldings.Total {
					t.Errorf("Expected new LTC total %f, got %f", newLTCHoldings.Total, ltcBal.Total)
				}
			},
		},
		{
			name:                "Save_EmptyCredentials_ErrorExpected",
			credsToSave:         &emptyCredsInstance,
			subAccountsToSave:   []SubAccount{{ID: "spot1", AssetType: asset.Spot, Currencies: []Balance{initialBTCHoldings}}},
			initialAccountState: func() *Accounts { mux := dispatch.GetNewMux(nil); accs, _ := NewAccounts(mockEx, mux); return accs },
			expectedErr:         true,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered == nil {
					t.Error("Expected errCredentialsEmpty, got nil")
				}
			},
		},
		{
			name:                "Save_InvalidAssetTypeInSubAccount_ErrorExpected",
			credsToSave:         &mockCreds,
			subAccountsToSave:   []SubAccount{{ID: "spot1", AssetType: invalidAssetType, Currencies: []Balance{initialBTCHoldings}}},
			initialAccountState: func() *Accounts { mux := dispatch.GetNewMux(nil); accs, _ := NewAccounts(mockEx, mux); return accs },
			expectedErr:         true,
		},
		{
			name:        "Save_BalanceUpdateError_OutOfSequence_ErrorExpected",
			credsToSave: &mockCreds,
			subAccountsToSave: []SubAccount{
				{ID: "spot1", AssetType: asset.Spot, Currencies: []Balance{{Currency: currency.BTC, Total: 0.5, UpdatedAt: time.Unix(50, 0)}}},
			},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				accs, _ := NewAccounts(mockEx, mux)
				accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
					{SubAccount: "spot1", Asset: asset.Spot}: {
						currency.BTC.Item: &balance{internal: initialBTCHoldings},
					},
				}
				return accs
			},
			expectedErr: true,
		},
		{
			name:        "Save_ZeroUpdatedAt_GetsTimestamped_Success",
			credsToSave: &mockCreds,
			subAccountsToSave: []SubAccount{
				{ID: "spot1", AssetType: asset.Spot, Currencies: []Balance{{Currency: currency.BTC, Total: 1, UpdatedAt: time.Time{}}}},
			},
			initialAccountState: func() *Accounts { mux := dispatch.GetNewMux(nil); accs, _ := NewAccounts(mockEx, mux); return accs },
			expectedErr:         false,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered != nil {
					t.Fatalf("Unexpected error: %v", errEncountered)
				}
				b, err := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.BTC)
				if err != nil {
					t.Fatalf("Failed to get balance: %v", err)
				}
				if b.UpdatedAt.IsZero() {
					t.Error("Expected UpdatedAt to be set, but it's still zero.")
				}
			},
		},
		{
			name:                "Save_NilReceiver_ErrorExpected",
			credsToSave:         &mockCreds,
			subAccountsToSave:   []SubAccount{{ID: "spot1", AssetType: asset.Spot, Currencies: []Balance{initialBTCHoldings}}},
			initialAccountState: func() *Accounts { return nil },
			expectedErr:         true,
		},
		{
			name:              "Save_NilSubAccountsMapInReceiver_ErrorExpected",
			credsToSave:       &mockCreds,
			subAccountsToSave: []SubAccount{{ID: "spot1", AssetType: asset.Spot, Currencies: []Balance{initialBTCHoldings}}},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				return &Accounts{Exchange: mockEx, mux: mux, ID: uuid.Must(uuid.NewV4()), subAccounts: nil}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accs := tt.initialAccountState()
			err := accs.Save(tt.subAccountsToSave, tt.credsToSave)
			if (err != nil) != tt.expectedErr {
				t.Errorf("Save() error = %v, wantErr %v. Test: %s", err, tt.expectedErr, tt.name)
			}
			if tt.verifyState != nil {
				tt.verifyState(t, accs, err)
			}
		})
	}
}

// TestAccounts_GetBalance tests the GetBalance method of the Accounts struct.
func TestAccounts_GetBalance(t *testing.T) {
	t.Parallel()

	mockCreds := Credentials{Key: "mainAcc"}
	mockEx := &mockEx{name: "TestGetBalanceExchange"}
	mux := dispatch.GetNewMux(nil)
	accs, _ := NewAccounts(mockEx, mux)

	spotBTC := Balance{Currency: currency.BTC, Total: 1, Free: 1, UpdatedAt: time.Now()}
	spotETH := Balance{Currency: currency.ETH, Total: 10, Free: 10, UpdatedAt: time.Now()}

	accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
		{SubAccount: "spotSub1", Asset: asset.Spot}: {
			currency.BTC.Item: &balance{internal: spotBTC},
			currency.ETH.Item: &balance{internal: spotETH},
		},
		{SubAccount: "spotSub2", Asset: asset.Spot}: {
			currency.BTC.Item: &balance{internal: Balance{Currency: currency.BTC, Total: 5, Free: 5}},
		},
	}

	emptyCreds := Credentials{}
	var invalidAssetTypeForGetBalance asset.Item
	emptyCurrency := currency.Code{}

	tests := []struct {
		name            string
		subAccountID    string
		creds           *Credentials
		assetType       asset.Item
		currency        currency.Code
		expectedBalance Balance
		expectedErr     bool
	}{
		{
			name:            "GetBalance_Success",
			subAccountID:    "spotSub1",
			creds:           &mockCreds,
			assetType:       asset.Spot,
			currency:        currency.BTC,
			expectedBalance: spotBTC,
			expectedErr:     false,
		},
		{
			name:            "GetBalance_CurrencyNotFoundInSubAccount",
			subAccountID:    "spotSub1",
			creds:           &mockCreds,
			assetType:       asset.Spot,
			currency:        currency.LTC,
			expectedBalance: Balance{},
			expectedErr:     true,
		},
		{
			name:            "GetBalance_SubAccountAssetKeyNotFound",
			subAccountID:    "spotSub1",
			creds:           &mockCreds,
			assetType:       asset.Margin,
			expectedBalance: Balance{},
			expectedErr:     true,
		},
		{
			name:            "GetBalance_CredentialsNotFound",
			subAccountID:    "spotSub1",
			creds:           &Credentials{Key: "unknown"},
			assetType:       asset.Spot,
			currency:        currency.BTC,
			expectedBalance: Balance{},
			expectedErr:     true,
		},
		{
			name:            "GetBalance_InvalidAssetType",
			subAccountID:    "spotSub1",
			creds:           &mockCreds,
			assetType:       invalidAssetTypeForGetBalance,
			currency:        currency.BTC,
			expectedBalance: Balance{},
			expectedErr:     true,
		},
		{
			name:            "GetBalance_EmptyCredentials",
			subAccountID:    "spotSub1",
			creds:           &emptyCreds,
			assetType:       asset.Spot,
			currency:        currency.BTC,
			expectedBalance: Balance{},
			expectedErr:     true,
		},
		{
			name:            "GetBalance_EmptyCurrencyCode",
			subAccountID:    "spotSub1",
			creds:           &mockCreds,
			assetType:       asset.Spot,
			currency:        emptyCurrency,
			expectedBalance: Balance{},
			expectedErr:     true,
		},
		{
			name:            "GetBalance_SubAccountIDNotFound",
			subAccountID:    "nonExistentSubID",
			creds:           &mockCreds,
			assetType:       asset.Spot,
			currency:        currency.BTC,
			expectedBalance: Balance{},
			expectedErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bal, err := accs.GetBalance(tt.subAccountID, tt.creds, tt.assetType, tt.currency)
			if (err != nil) != tt.expectedErr {
				t.Errorf("GetBalance() error = %v, wantErr %v. Test: %s", err, tt.expectedErr, tt.name)
				return
			}
			if !tt.expectedErr {
				if bal.Currency != tt.expectedBalance.Currency ||
					bal.Total != tt.expectedBalance.Total ||
					bal.Free != tt.expectedBalance.Free ||
					bal.Hold != tt.expectedBalance.Hold {
					t.Errorf("GetBalance() balance = %+v, want %+v. Test: %s", bal, tt.expectedBalance, tt.name)
				}
			}
		})
	}
}

// TestAccounts_GetHoldings tests the GetHoldings method of the Accounts struct.
func TestAccounts_GetHoldings(t *testing.T) {
	t.Parallel()

	mockCreds := Credentials{Key: "mainAcc"}
	mockEx := &mockEx{name: "TestHoldingsExchange"}
	mux := dispatch.GetNewMux(nil)
	accs, _ := NewAccounts(mockEx, mux)

	spotBTC := Balance{Currency: currency.BTC, Total: 1, Free: 1, UpdatedAt: time.Now()}
	spotETH := Balance{Currency: currency.ETH, Total: 10, Free: 10, UpdatedAt: time.Now()}
	marginUSDT := Balance{Currency: currency.USDT, Total: 1000, Free: 1000, UpdatedAt: time.Now()}

	accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
		{SubAccount: "spotA", Asset: asset.Spot}: {
			currency.BTC.Item: &balance{internal: spotBTC},
			currency.ETH.Item: &balance{internal: spotETH},
		},
		{SubAccount: "marginA", Asset: asset.Margin}: {
			currency.USDT.Item: &balance{internal: marginUSDT},
		},
	}

	emptyCreds := Credentials{}
	var invalidAssetTypeForGetHoldings asset.Item

	tests := []struct {
		name          string
		creds         *Credentials
		assetType     asset.Item
		expectedCount int
		expectedErr   bool
	}{
		{
			name:          "GetHoldings_Spot_ActuallyReturnsNonSpot",
			creds:         &mockCreds,
			assetType:     asset.Spot,
			expectedCount: 1,
			expectedErr:   false,
		},
		{
			name:          "GetHoldings_Margin_ActuallyReturnsNonMargin",
			creds:         &mockCreds,
			assetType:     asset.Margin,
			expectedCount: 2,
			expectedErr:   false,
		},
		{
			name:          "GetHoldings_NoMatchingCredentials",
			creds:         &Credentials{Key: "unknownCreds"},
			assetType:     asset.Spot,
			expectedCount: 0,
			expectedErr:   true,
		},
		{
			name:          "GetHoldings_EmptyCredentials",
			creds:         &emptyCreds,
			assetType:     asset.Spot,
			expectedCount: 0,
			expectedErr:   true,
		},
		{
			name:          "GetHoldings_InvalidAssetType",
			creds:         &mockCreds,
			assetType:     invalidAssetTypeForGetHoldings,
			expectedCount: 0,
			expectedErr:   true,
		},
		{
			name:          "GetHoldings_AssetTypeWithNoOtherHoldings",
			creds:         &mockCreds,
			assetType:     asset.Futures,
			expectedCount: 3,
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			holdings, err := accs.GetHoldings(tt.creds, tt.assetType)
			if (err != nil) != tt.expectedErr {
				t.Errorf("GetHoldings() error = %v, wantErr %v. Test: %s", err, tt.expectedErr, tt.name)
				return
			}
			if err == nil && len(holdings) != tt.expectedCount {
				t.Errorf("GetHoldings() count = %d, wantCount %d. Holdings: %+v. Test: %s", len(holdings), tt.expectedCount, holdings, tt.name)
			}
			if !tt.expectedErr && tt.expectedCount > 0 {
				if len(holdings) == 0 && tt.expectedCount > 0 {
					t.Errorf("Expected some holdings but got none for test: %s", tt.name)
				}
			}
		})
	}
}

// TestAccounts_Update tests the Update method of the Accounts struct.
func TestAccounts_Update(t *testing.T) {
	t.Parallel()

	mockCreds := Credentials{Key: "updateCreds"}
	var emptyCredsInstanceForUpdate Credentials
	var invalidAssetTypeForUpdate asset.Item
	mockEx := &mockEx{name: "TestUpdateExchange"}

	initialBTC := Balance{Currency: currency.BTC, Total: 10, Free: 10, UpdatedAt: time.Unix(100, 0)}
	initialETH := Balance{Currency: currency.ETH, Total: 20, Free: 20, UpdatedAt: time.Unix(100, 0)}

	tests := []struct {
		name                string
		credsToUpdate       *Credentials
		changesToApply      []Change
		initialAccountState func() *Accounts
		expectedErr         bool
		verifyState         func(t *testing.T, accs *Accounts, errEncountered error)
	}{
		{
			name:          "Update_ExistingBalance_Success",
			credsToUpdate: &mockCreds,
			changesToApply: []Change{
				{Account: "spot1", AssetType: asset.Spot, Balance: &Balance{Currency: currency.BTC, Total: 15, Free: 15, UpdatedAt: time.Unix(150, 0)}},
			},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				accs, _ := NewAccounts(mockEx, mux)
				accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
					{SubAccount: "spot1", Asset: asset.Spot}: {currency.BTC.Item: &balance{internal: initialBTC}},
				}
				return accs
			},
			expectedErr: false,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered != nil {
					t.Fatalf("Unexpected error: %v", errEncountered)
				}
				b, _ := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.BTC)
				if b.Total != 15 {
					t.Errorf("Expected BTC total 15, got %f", b.Total)
				}
			},
		},
		{
			name:          "Update_AddNewCurrencyToExistingSubAccount_Success",
			credsToUpdate: &mockCreds,
			changesToApply: []Change{
				{Account: "spot1", AssetType: asset.Spot, Balance: &Balance{Currency: currency.LTC, Total: 5, Free: 5, UpdatedAt: time.Unix(150, 0)}},
			},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				accs, _ := NewAccounts(mockEx, mux)
				accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
					{SubAccount: "spot1", Asset: asset.Spot}: {currency.BTC.Item: &balance{internal: initialBTC}},
				}
				return accs
			},
			expectedErr: false,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered != nil {
					t.Fatalf("Unexpected error: %v", errEncountered)
				}
				bLTC, errLTC := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.LTC)
				if errLTC != nil {
					t.Fatalf("Failed to get LTC: %v", errLTC)
				}
				if bLTC.Total != 5 {
					t.Errorf("Expected LTC total 5, got %f", bLTC.Total)
				}
			},
		},
		{
			name:          "Update_NewSubAccountAndAssetViaChange_Success",
			credsToUpdate: &mockCreds,
			changesToApply: []Change{
				{Account: "marginNew", AssetType: asset.Margin, Balance: &Balance{Currency: currency.USDT, Total: 100, Free: 100, UpdatedAt: time.Unix(150, 0)}},
			},
			initialAccountState: func() *Accounts { mux := dispatch.GetNewMux(nil); accs, _ := NewAccounts(mockEx, mux); return accs },
			expectedErr:         false,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered != nil {
					t.Fatalf("Unexpected error: %v", errEncountered)
				}
				b, err := accs.GetBalance("marginNew", &mockCreds, asset.Margin, currency.USDT)
				if err != nil {
					t.Fatalf("Failed to get USDT: %v", err)
				}
				if b.Total != 100 {
					t.Errorf("Expected USDT total 100, got %f", b.Total)
				}
			},
		},
		{
			name:                "Update_EmptyCredentials_ErrorExpected",
			credsToUpdate:       &emptyCredsInstanceForUpdate,
			changesToApply:      []Change{{Account: "spot1", AssetType: asset.Spot, Balance: &initialBTC}},
			initialAccountState: func() *Accounts { mux := dispatch.GetNewMux(nil); accs, _ := NewAccounts(mockEx, mux); return accs },
			expectedErr:         true,
		},
		{
			name:                "Update_InvalidAssetTypeInChange_ErrorExpected",
			credsToUpdate:       &mockCreds,
			changesToApply:      []Change{{Account: "spot1", AssetType: invalidAssetTypeForUpdate, Balance: &initialBTC}},
			initialAccountState: func() *Accounts { mux := dispatch.GetNewMux(nil); accs, _ := NewAccounts(mockEx, mux); return accs },
			expectedErr:         true,
		},
		{
			name:                "Update_NilBalanceInChange_ErrorExpected",
			credsToUpdate:       &mockCreds,
			changesToApply:      []Change{{Account: "spot1", AssetType: asset.Spot, Balance: nil}},
			initialAccountState: func() *Accounts { mux := dispatch.GetNewMux(nil); accs, _ := NewAccounts(mockEx, mux); return accs },
			expectedErr:         true,
		},
		{
			name:          "Update_BalanceUpdateError_OutOfSequence_ErrorExpected",
			credsToUpdate: &mockCreds,
			changesToApply: []Change{
				{Account: "spot1", AssetType: asset.Spot, Balance: &Balance{Currency: currency.BTC, Total: 5, UpdatedAt: time.Unix(50, 0)}},
			},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				accs, _ := NewAccounts(mockEx, mux)
				accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
					{SubAccount: "spot1", Asset: asset.Spot}: {currency.BTC.Item: &balance{internal: initialBTC}},
				}
				return accs
			},
			expectedErr: true,
		},
		{
			name:          "Update_MultipleChanges_MixedSuccessAndFailures_ExpectPartialUpdateAndError",
			credsToUpdate: &mockCreds,
			changesToApply: []Change{
				{Account: "spot1", AssetType: asset.Spot, Balance: &Balance{Currency: currency.ETH, Total: 25, UpdatedAt: time.Unix(150, 0)}},
				{Account: "spot1", AssetType: asset.Spot, Balance: &Balance{Currency: currency.BTC, Total: 5, UpdatedAt: time.Unix(50, 0)}},
				{Account: "spot1", AssetType: asset.Spot, Balance: &Balance{Currency: currency.LTC, Total: 30, UpdatedAt: time.Unix(160, 0)}},
			},
			initialAccountState: func() *Accounts {
				mux := dispatch.GetNewMux(nil)
				accs, _ := NewAccounts(mockEx, mux)
				accs.subAccounts[mockCreds] = map[key.SubAccountAsset]CurrencyBalances{
					{SubAccount: "spot1", Asset: asset.Spot}: {
						currency.BTC.Item: &balance{internal: initialBTC},
						currency.ETH.Item: &balance{internal: initialETH},
					},
				}
				return accs
			},
			expectedErr: true,
			verifyState: func(t *testing.T, accs *Accounts, errEncountered error) {
				if errEncountered == nil {
					t.Fatal("Expected an error, got nil")
				}
				ethBal, _ := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.ETH)
				if ethBal.Total != 25 {
					t.Errorf("Expected ETH total 25, got %f", ethBal.Total)
				}
				btcBal, _ := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.BTC)
				if btcBal.Total != initialBTC.Total {
					t.Errorf("Expected BTC total %f, got %f", initialBTC.Total, btcBal.Total)
				}
				ltcBal, _ := accs.GetBalance("spot1", &mockCreds, asset.Spot, currency.LTC)
				if ltcBal.Total != 30 {
					t.Errorf("Expected LTC total 30, got %f", ltcBal.Total)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accs := tt.initialAccountState()
			err := accs.Update(tt.changesToApply, tt.credsToUpdate)

			if (err != nil) != tt.expectedErr {
				t.Errorf("Update() error = %v, wantErr %v. Test: %s", err, tt.expectedErr, tt.name)
			}
			if tt.verifyState != nil {
				tt.verifyState(t, accs, err)
			}
		})
	}
}

// TestSubAccounts_Group tests the Group method of the SubAccounts slice type.
func TestSubAccounts_Group(t *testing.T) {
	t.Parallel()

	sa1BTC1 := Balance{Currency: currency.BTC, Total: 1, Free: 1}
	sa1ETH1 := Balance{Currency: currency.ETH, Total: 10, Free: 10}
	sa1USDT1 := Balance{Currency: currency.USDT, Total: 100, Free: 100}

	sa2BTC1 := Balance{Currency: currency.BTC, Total: 2, Free: 2}
	sa2LTC1 := Balance{Currency: currency.LTC, Total: 20, Free: 20}

	sa1BTC2 := Balance{Currency: currency.BTC, Total: 0.5, Free: 0.5} // More BTC for Spot account "acc1"

	tests := []struct {
		name     string
		input    SubAccounts // This is accounts.SubAccounts type (i.e., []SubAccount)
		expected SubAccounts
	}{
		{
			name: "Group basic different IDs and Assets",
			input: SubAccounts{
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1, sa1ETH1}},
				{ID: "acc2", AssetType: asset.Spot, Currencies: []Balance{sa2BTC1, sa2LTC1}},
				{ID: "acc1", AssetType: asset.Margin, Currencies: []Balance{sa1USDT1}},
			},
			// Expected output is sorted by AssetType then ID
			expected: SubAccounts{
				{ID: "acc1", AssetType: asset.Margin, Currencies: []Balance{sa1USDT1}},
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1, sa1ETH1}},
				{ID: "acc2", AssetType: asset.Spot, Currencies: []Balance{sa2BTC1, sa2LTC1}},
			},
		},
		{
			name: "Group and concatenate currencies for same ID and AssetType",
			input: SubAccounts{
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1}},
				{ID: "acc2", AssetType: asset.Spot, Currencies: []Balance{sa2BTC1}},
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1ETH1}}, // acc1, Spot again
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC2}}, // acc1, Spot, more BTC
			},
			expected: SubAccounts{
				// Expected: acc1 Spot (BTC, ETH, BTC) - order of concatenated currencies is preserved from input order
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1, sa1ETH1, sa1BTC2}},
				{ID: "acc2", AssetType: asset.Spot, Currencies: []Balance{sa2BTC1}},
			},
		},
		{
			name:     "Empty input",
			input:    SubAccounts{},
			expected: SubAccounts{},
		},
		{
			name: "Single element",
			input: SubAccounts{
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1}},
			},
			expected: SubAccounts{
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1}},
			},
		},
		{
			name: "All same ID and AssetType",
			input: SubAccounts{
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1}},
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1ETH1}},
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1USDT1}},
			},
			expected: SubAccounts{
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1, sa1ETH1, sa1USDT1}},
			},
		},
		{
			name: "Mixed with empty currencies list initially",
			input: SubAccounts{
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1}},
				{ID: "acc2", AssetType: asset.Spot, Currencies: []Balance{}}, // acc2 starts with empty
				{ID: "acc1", AssetType: asset.Margin, Currencies: []Balance{sa1USDT1}},
				{ID: "acc2", AssetType: asset.Spot, Currencies: []Balance{sa2LTC1}}, // acc2 gets currencies
			},
			expected: SubAccounts{
				{ID: "acc1", AssetType: asset.Margin, Currencies: []Balance{sa1USDT1}},
				{ID: "acc1", AssetType: asset.Spot, Currencies: []Balance{sa1BTC1}},
				// acc2, Spot had Currencies:[] then Currencies:[LTC]. The Group func concatenates.
				{ID: "acc2", AssetType: asset.Spot, Currencies: []Balance{sa2LTC1}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Group()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SubAccounts.Group() for test '%s':\nGot:      %+v\nExpected: %+v", tt.name, result, tt.expected)
			}
		})
	}
}
*/
