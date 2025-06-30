package accounts

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/key"
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
	assert.Equal(t, "mocky", a.Exchange.GetName(), "Exchange should set correctly")
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
	a := MustNewAccounts(&mockEx{"mocky"})
	require.NotNil(t, a)
	require.Panics(t, func() { _ = MustNewAccounts(nil) })
}

func TestNewSubAccount(t *testing.T) {
	t.Parallel()
	a := NewSubAccount(asset.Spot, "")
	require.NotNil(t, a, "must not return nil with no id")
	assert.Equal(t, asset.Spot, a.AssetType, "AssetType should correct")
	assert.Empty(t, a.ID, "ID should not default to anything")
	a = NewSubAccount(asset.Spot, "42")
	assert.Equal(t, "42", a.ID, "ID should be correct")
}

func TestSubscribe(t *testing.T) {
	t.Parallel()
	err := dispatch.Start(dispatch.DefaultMaxWorkers, dispatch.DefaultJobsLimit)
	require.NoError(t, common.ExcludeError(err, dispatch.ErrDispatcherAlreadyRunning), "dispatch.Start must not error")
	p, err := MustNewAccounts(&mockEx{}).Subscribe()
	require.NoError(t, err)
	require.NotNil(t, p, "Subscribe must return a pipe")
	require.Empty(t, p.Channel(), "Pipe must be empty before Saving anything")
}

// TestAccountsCurrencyBalancesPrivateensures that currencyBalances initializes maps, and returns a reference to the same entry
func TestAccountsCurrencyBalancesPrivate(t *testing.T) {
	t.Parallel()

	a := &Accounts{subAccounts: make(credSubAccounts)}
	c := &Credentials{Key: "a"}
	b := a.currencyBalances(c, "", asset.Spot)
	r1 := a.subAccounts[*c]
	// Using reflect since assert.Same cannot be used on maps to ensure same underlying pointer
	assert.Equal(t,
		reflect.ValueOf(b).UnsafePointer(),
		reflect.ValueOf(r1[key.SubAccountAsset{Asset: asset.Spot}]).UnsafePointer(),
		"should make and return the same map")
	assert.Equal(t,
		reflect.ValueOf(b).UnsafePointer(),
		reflect.ValueOf(a.currencyBalances(c, "", asset.Spot)).UnsafePointer(),
		"should return the same map on subsequent calls")
	b = a.currencyBalances(c, "", asset.Futures)
	assert.Equal(t,
		reflect.ValueOf(r1).UnsafePointer(),
		reflect.ValueOf(a.subAccounts[*c]).UnsafePointer(),
		"should not make a new cred key")
	assert.Equal(t,
		reflect.ValueOf(b).UnsafePointer(),
		reflect.ValueOf(r1[key.SubAccountAsset{Asset: asset.Futures}]).UnsafePointer(),
		"should make and return the same map")
}

/*
	require.NotNil(
	// Test creds don't exist
	if _, ok := a.subAccounts[*c]; !ok {

	// Test asset key doesn't exist

	require.NotNil(t, bals)
	b := bals.balance(currency.BTC)
	updated, err := b.update(Balance{Total: 4.0, UpdatedAt: time.Now()})
	require.NoError(t, err, "update must not error")
	require.True(t, updated, "update must return true")
	ref := a.subAccounts[*c][key.SubAccountAsset{Asset: asset.Spot}][currency.BTC].internal
	assert.Equal(t, currency.BTC, ref.Currency)
	assert.Equal(t, 4.0, ref.Total)

	bals = a.currencyBalances(c, "", asset.Spot)
	require.NotNil(t, bals)
	b = bals.balance(currency.ETH)
	updated, err = b.update(Balance{Currency: currency.LTC, Total: 5.0, UpdatedAt: time.Now()})
	require.NoError(t, err)
	require.True(t, updated)
	ref = a.subAccounts[*c][key.SubAccountAsset{Asset: asset.Spot}][currency.LTC].internal
	assert.Equal(t, currency.LTC, ref.Currency)
	assert.Equal(t, 5.0, ref.Total)
	_, ok := a.subAccounts[*c][key.SubAccountAsset{Asset: asset.Spot}][currency.ETH]
	assert.False(t, ok)

	bals = a.currencyBalances(c, "", asset.Spot)
	b = bals.balance(currency.DOGE)
	updated, err = b.update(Balance{Currency: currency.XRP, Total: 100.0, UpdatedAt: time.Now()})
	require.NoError(t, err)
	require.True(t, updated)
	internalBalance := a.subAccounts[*c][key.SubAccountAsset{Asset: asset.Spot}][currency.DOGE].internal
	assert.Equal(t, currency.XRP, internalBalance.Currency)
	assert.Equal(t, 100.0, internalBalance.Total)
}

func TestCurrencyBalances(t *testing.T) {
	t.Parallel()
	n := time.Now()
	creds1 := Credentials{Key: "cred1"}
	creds2 := Credentials{Key: "cred2"}
	mockExInst := &mockEx{name: "TestExchange"}
	a := MustNewAccounts(mockExInst)

	require.NoError(t, a.Save(SubAccounts{
		{ID: "a", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Currency: currency.BTC, Total: 1, Free: 1, UpdatedAt: n}, currency.ETH: {Currency: currency.ETH, Total: 10, Free: 10, UpdatedAt: n}}},
		{ID: "b", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Currency: currency.BTC, Total: 2, Free: 2, UpdatedAt: n}, currency.LTC: {Currency: currency.LTC, Total: 5, Free: 5, UpdatedAt: n}}},
		{ID: "a", AssetType: asset.Margin, Balances: CurrencyBalances{currency.USDT: {Currency: currency.USDT, Total: 100, Free: 100, UpdatedAt: n}, currency.BTC:  {Currency: currency.BTC, Total: 0.5, Free: 0.5, UpdatedAt: n}}},
	}, &creds1))
	require.NoError(t, a.Save(SubAccounts{{ID: "c", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Currency: currency.BTC, Total: 3, UpdatedAt: n}}}}, &creds2))

	t.Run("all credentials, all assets", func(t *testing.T) {
		expected := map[currency.Code]float64{currency.BTC: 6.5, currency.ETH: 10, currency.LTC: 5, currency.USDT: 100}
		res, err := a.CurrencyBalances(nil, asset.All); require.NoError(t, err); require.Len(t, res, len(expected))
		for c, total := range expected { bal, ok := res[c]; require.True(t, ok); assert.Equal(t, total, bal.Total) }
	})
	t.Run("specific credentials, all assets", func(t *testing.T) {
		expected := map[currency.Code]float64{currency.BTC: 3.5, currency.ETH: 10, currency.LTC: 5, currency.USDT: 100}
		res, err := a.CurrencyBalances(&creds1, asset.All); require.NoError(t, err); require.Len(t, res, len(expected))
		for c, total := range expected { bal, ok := res[c]; require.True(t, ok); assert.Equal(t, total, bal.Total) }
	})
	t.Run("all credentials, specific asset", func(t *testing.T) {
		expected := map[currency.Code]float64{currency.BTC: 6, currency.ETH: 10, currency.LTC: 5}
		res, err := a.CurrencyBalances(nil, asset.Spot); require.NoError(t, err); require.Len(t, res, len(expected))
		for c, total := range expected { bal, ok := res[c]; require.True(t, ok); assert.Equal(t, total, bal.Total) }
	})
	t.Run("specific credentials, specific asset", func(t *testing.T) {
		expected := map[currency.Code]float64{currency.BTC: 0.5, currency.USDT: 100}
		res, err := a.CurrencyBalances(&creds1, asset.Margin); require.NoError(t, err); require.Len(t, res, len(expected))
		for c, total := range expected { bal, ok := res[c]; require.True(t, ok); assert.Equal(t, total, bal.Total) }
	})
	t.Run("no subaccounts", func(t *testing.T) {
		_, err := MustNewAccounts(&mockEx{name: "empty"}).CurrencyBalances(nil, asset.All)
		assert.ErrorIs(t, err, ErrBalancesNotFound)
	})
	t.Run("no balances for criteria", func(t *testing.T) {
		_, err := a.CurrencyBalances(&Credentials{Key: "nonexistent"}, asset.Spot); assert.ErrorIs(t, err, ErrBalancesNotFound)
		_, err = a.CurrencyBalances(&creds1, asset.Futures); assert.ErrorIs(t, err, ErrBalancesNotFound)
	})
	t.Run("invalid asset type", func(t *testing.T) {
		_, err := a.CurrencyBalances(nil, testInvalidAsset); assert.ErrorContains(t, err, asset.ErrNotSupported.Error())
	})
	t.Run("nil receiver", func(t *testing.T) {
		var nilAccs *Accounts; _, err := nilAccs.CurrencyBalances(nil, asset.All); assert.ErrorIs(t, err, common.ErrNilPointer)
	})
}

func TestAccountsSave(t *testing.T) {
	t.Parallel()
	requrie.NoError(t, dispatch.EnsureRunning(dispatch.DefaultMaxWorkers, dispatch.DefaultJobsLimit))
	mockCreds := Credentials{Key: "saveCreds"}
	var emptyCreds Credentials
	mockExInst := &mockEx{name: "TestSaveExchange"}
	initialBTC := Balance{Currency: currency.BTC, Total: 1, UpdatedAt: time.Unix(100, 0)}
	initialETH := Balance{Currency: currency.ETH, Total: 10, UpdatedAt: time.Unix(100, 0)}
	updatedBTC := Balance{Currency: currency.BTC, Total: 2, UpdatedAt: time.Unix(200, 0)}
	newLTC := Balance{Currency: currency.LTC, Total: 5, UpdatedAt: time.Unix(200, 0)}

	tests := []struct{
		name string
		creds *Credentials
		toSave SubAccounts
		initState func(e accountsExchangeDependency) *Accounts
		errContains string
		verify func(*testing.T, *Accounts, error)
	}{
		{"Save_New", &mockCreds, SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, MustNewAccounts, "", func(t *testing.T, accs *Accounts, err error) { require.NoError(t, err); b, e := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.BTC); require.NoError(t, e); assert.Equal(t, initialBTC.Total, b.Total) }},
		{"Save_Update_Add_Remove", &mockCreds, SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: updatedBTC, currency.LTC: newLTC}}}, func(e accountsExchangeDependency) *Accounts { accs := MustNewAccounts(e); require.NoError(t, accs.Save(SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC, currency.ETH: initialETH}}}, &mockCreds)); return accs }, "", func(t *testing.T, accs *Accounts, err error) { require.NoError(t, err); bBTC, eBTC := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.BTC); require.NoError(t, eBTC); assert.Equal(t, updatedBTC.Total, bBTC.Total); _, eETH := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.ETH); assert.ErrorContains(t, eETH, errNoExchangeSubAccountBalances.Error()); bLTC, eLTC := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.LTC); require.NoError(t, eLTC); assert.Equal(t, newLTC.Total, bLTC.Total) }},
		{"Save_EmptyCreds", &emptyCreds, SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, MustNewAccounts, errCredentialsEmpty.Error(), nil },
		{"Save_InvalidAsset", &mockCreds, SubAccounts{{ID: "s1", AssetType: testInvalidAsset, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, MustNewAccounts, asset.ErrNotSupported.Error(), nil },
		{"Save_OutOfSequence", &mockCreds, SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Total: 0.5, UpdatedAt: time.Unix(50,0)}}}}, func(e accountsExchangeDependency) *Accounts { accs := MustNewAccounts(e); require.NoError(t, accs.Save(SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, &mockCreds)); return accs }, errOutOfSequence.Error(), nil },
		{"Save_ZeroUpdateAt", &mockCreds, SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Total: 1, UpdatedAt: time.Time{}}}}}, MustNewAccounts, "", func(t *testing.T, accs *Accounts, err error) { require.NoError(t, err); b, e := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.BTC); require.NoError(t, e); assert.False(t, b.UpdatedAt.IsZero()) }},
		{"Save_NilReceiver", &mockCreds, SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, func(e accountsExchangeDependency) *Accounts { return nil }, common.ErrNilPointer.Error(), nil },
		{"Save_NilSubAccountsMap", &mockCreds, SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, func(e accountsExchangeDependency) *Accounts { mux := dispatch.GetNewMux(nil); id, _ := mux.GetID(); return &Accounts{Exchange: e, mux: mux, routingID: id, subAccounts: nil} }, common.ErrNilPointer.Error(), nil },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accs := tt.initState(mockExInst)
			err := accs.Save(tt.toSave, tt.creds)
			if tt.errContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
			}
			if tt.verify != nil {
				tt.verify(t, accs, err)
			}
		})
	}
}

func TestAccountsGetBalance(t *testing.T) {
	t.Parallel()
	mockCreds := Credentials{Key: "mainAcc"}
	mockExInst := &mockEx{name: "TestGetBalanceExchange"}
	accs := MustNewAccounts(mockExInst)
	spotBTC := Balance{Currency: currency.BTC, Total: 1, UpdatedAt: time.Now()}
	spotETH := Balance{Currency: currency.ETH, Total: 10, UpdatedAt: time.Now()}
	require.NoError(t, accs.Save(SubAccounts{{ID: "sS1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: spotBTC, currency.ETH: spotETH}}, {ID: "sS2", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Total: 5}}}}, &mockCreds))
	tests := []struct{ name, sID string; creds *Credentials; aType asset.Item; c currency.Code; wantBal Balance; wantErr bool; errContains string }{
		{"OK", "sS1", &mockCreds, asset.Spot, currency.BTC, spotBTC, false, ""},
		{"NoCurrency", "sS1", &mockCreds, asset.Spot, currency.LTC, Balance{}, true, errNoExchangeSubAccountBalances.Error()},
		{"NoAsset", "sS1", &mockCreds, asset.Margin, currency.BTC, Balance{}, true, errNoExchangeSubAccountBalances.Error()},
		{"NoCreds", "sS1", &Credentials{Key: "unknown"}, asset.Spot, currency.BTC, Balance{}, true, ErrBalancesNotFound.Error()},
		{"InvalidAsset", "sS1", &mockCreds, testInvalidAsset, currency.BTC, Balance{}, true, asset.ErrNotSupported.Error()},
		{"EmptyCreds", "sS1", &Credentials{}, asset.Spot, currency.BTC, Balance{}, true, errCredentialsEmpty.Error()},
		{"EmptyCurrency", "sS1", &mockCreds, asset.Spot, currency.EMPTYCODE, Balance{}, true, currency.ErrCurrencyCodeEmpty.Error()},
		{"NoSubAccountID", "nonExistent", &mockCreds, asset.Spot, currency.BTC, Balance{}, true, errNoExchangeSubAccountBalances.Error()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bal, err := accs.GetBalance(tt.sID, tt.creds, tt.aType, tt.c)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.ErrorContains(t, err, tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantBal.Currency, bal.Currency)
				assert.Equal(t, tt.wantBal.Total, bal.Total)
			}
		})
	}
}

func TestAccountsBalances(t *testing.T) {
	t.Parallel()
	c1, c2 := Credentials{Key: "cA"}, Credentials{Key: "cB"}
	mockExInst := &mockEx{name: "TestBalancesExchange"}
	accs := MustNewAccounts(mockExInst)
	sBTC1 := Balance{Currency: currency.BTC, Total: 1, UpdatedAt: time.Now()}
	sETH1 := Balance{Currency: currency.ETH, Total: 10, UpdatedAt: time.Now()}
	mUSDT1 := Balance{Currency: currency.USDT, Total: 1000, UpdatedAt: time.Now()}
	sBTC2 := Balance{Currency: currency.BTC, Total: 2, UpdatedAt: time.Now()}
	require.NoError(t, accs.Save(SubAccounts{{ID: "sA", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sBTC1, currency.ETH: sETH1}}, {ID: "mA", AssetType: asset.Margin, Balances: CurrencyBalances{currency.USDT: mUSDT1}}}, &c1))
	require.NoError(t, accs.Save(SubAccounts{{ID: "sB", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sBTC2}}}, &c2))
	tests := []struct{ name string; creds *Credentials; aType asset.Item; wantCount int; wantErr bool; errContains string; verify func(*testing.T, SubAccounts); initState func(e accountsExchangeDependency) *Accounts }{
		{"AllCredsAllAssets", nil, asset.All, 3, false, "", func(t *testing.T, r SubAccounts) { assert.Len(t, r, 3) }, nil},
		{"SpecificCredsAllAssets", &c1, asset.All, 2, false, "", nil, nil},
		{"AllCredsSpecificAsset", nil, asset.Spot, 2, false, "", func(t *testing.T, r SubAccounts) { for _, sa := range r { assert.Equal(t, asset.Spot, sa.AssetType) } }, nil},
		{"SpecificCredsSpecificAsset", &c1, asset.Margin, 1, false, "", func(t *testing.T, r SubAccounts) { assert.Equal(t, "mA", r[0].ID); assert.Contains(t, r[0].Balances, currency.USDT) }, nil},
		{"NoMatchingCreds", &Credentials{Key: "unknown"}, asset.All, 0, true, ErrBalancesNotFound.Error(), nil, nil},
		{"NoMatchingAsset", &c1, asset.Futures, 0, true, ErrBalancesNotFound.Error(), nil, nil},
		{"InvalidAsset", nil, testInvalidAsset, 0, true, asset.ErrNotSupported.Error(), nil, nil},
		{"NilReceiver", nil, asset.All, 0, true, common.ErrNilPointer.Error(), nil, func(e accountsExchangeDependency) *Accounts { return nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentAccs := accs
			if tt.initState != nil {
				currentAccs = tt.initState(mockExInst)
			}
			r, err := currentAccs.Balances(tt.creds, tt.aType)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.ErrorContains(t, err, tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Len(t, r, tt.wantCount)
				if tt.verify != nil {
					tt.verify(t, r)
				}
			}
		})
	}
}

func TestAccountsUpdate(t *testing.T) {
	t.Parallel()
	require.NoError(t, common.ExcludeError(dispatch.Start(dispatch.DefaultMaxWorkers, dispatch.DefaultJobsLimit), dispatch.ErrDispatcherAlreadyRunning))
	mockCreds := Credentials{Key: "updateCreds"}
	mockExInst := &mockEx{name: "TestUpdateExchange"}
	initialBTC := Balance{Currency: currency.BTC, Total: 10, UpdatedAt: time.Unix(100, 0)}
	initialETH := Balance{Currency: currency.ETH, Total: 20, UpdatedAt: time.Unix(100, 0)}

	changedBTC := Balance{Currency: currency.BTC, Total: 15, UpdatedAt: time.Unix(150,0)}
	addedLTC := Balance{Currency: currency.LTC, Total: 5, UpdatedAt: time.Unix(150,0)}
	newMarginUSDT := Balance{Currency: currency.USDT, Total: 100, UpdatedAt: time.Unix(150,0)}
	outOfSeqBTC := Balance{Currency: currency.BTC, Total: 5, UpdatedAt: time.Unix(50,0)}
	updatedETH := Balance{Currency: currency.ETH, Total: 25, UpdatedAt: time.Unix(150,0)}
	anotherLTC := Balance{Currency: currency.LTC, Total: 30, UpdatedAt: time.Unix(160,0)}
	anotherBTC := Balance{Currency: currency.BTC, Total: 1, UpdatedAt: time.Now()}

	tests := []struct{ name string; creds *Credentials; changes []Change; initState func(e accountsExchangeDependency) *Accounts; errContains string; verify func(*testing.T, *Accounts, error) }{
		{"UpdateExisting", &mockCreds, []Change{{Account: "s1", AssetType: asset.Spot, Balance: &changedBTC}}, func(e accountsExchangeDependency) *Accounts { accs := MustNewAccounts(e); require.NoError(t, accs.Save(SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, &mockCreds)); return accs }, "", func(t *testing.T, accs *Accounts, err error) { require.NoError(t, err); b, _ := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.BTC); assert.Equal(t, 15.0, b.Total) }},
		{"AddNewCurrency", &mockCreds, []Change{{Account: "s1", AssetType: asset.Spot, Balance: &addedLTC}}, func(e accountsExchangeDependency) *Accounts { accs := MustNewAccounts(e); require.NoError(t, accs.Save(SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, &mockCreds)); return accs }, "", func(t *testing.T, accs *Accounts, err error) { require.NoError(t, err); b, e := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.LTC); require.NoError(t, e); assert.Equal(t, 5.0, b.Total) }},
		{"NewSubAccountAsset", &mockCreds, []Change{{Account: "mN", AssetType: asset.Margin, Balance: &newMarginUSDT}}, MustNewAccounts, "", func(t *testing.T, accs *Accounts, err error) { require.NoError(t, err); b, e := accs.GetBalance("mN", &mockCreds, asset.Margin, currency.USDT); require.NoError(t, e); assert.Equal(t, 100.0, b.Total) }},
		{"EmptyCreds", &Credentials{}, []Change{{Account: "s1", AssetType: asset.Spot, Balance: &initialBTC}}, MustNewAccounts, errCredentialsEmpty.Error(), nil },
		{"InvalidAsset", &mockCreds, []Change{{Account: "s1", AssetType: testInvalidAsset, Balance: &initialBTC}}, MustNewAccounts, asset.ErrNotSupported.Error(), nil },
		{"NilBalance", &mockCreds, []Change{{Account: "s1", AssetType: asset.Spot, Balance: nil}}, MustNewAccounts, common.ErrNilPointer.Error(), nil },
		{"OutOfSequence", &mockCreds, []Change{{Account: "s1", AssetType: asset.Spot, Balance: &outOfSeqBTC}}, func(e accountsExchangeDependency) *Accounts { accs := MustNewAccounts(e); require.NoError(t, accs.Save(SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC}}}, &mockCreds)); return accs }, errOutOfSequence.Error(), nil },
		{"MultipleChangesMixed", &mockCreds, []Change{{Account: "s1", AssetType: asset.Spot, Balance: &updatedETH}, {Account: "s1", AssetType: asset.Spot, Balance: &outOfSeqBTC}, {Account: "s1", AssetType: asset.Spot, Balance: &anotherLTC}, {Account: "s2", AssetType: testInvalidAsset, Balance: &anotherBTC}}, func(e accountsExchangeDependency) *Accounts { accs := MustNewAccounts(e); require.NoError(t, accs.Save(SubAccounts{{ID: "s1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: initialBTC, currency.ETH: initialETH}}}, &mockCreds)); return accs }, errOutOfSequence.Error(), func(t *testing.T, accs *Accounts, err error) { require.Error(t, err); assert.ErrorContains(t, err, errOutOfSequence.Error()); assert.ErrorContains(t, err, asset.ErrNotSupported.Error()); bETH, _ := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.ETH); assert.Equal(t, 25.0, bETH.Total); bBTC, _ := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.BTC); assert.Equal(t, initialBTC.Total, bBTC.Total); bLTC, eLTC := accs.GetBalance("s1", &mockCreds, asset.Spot, currency.LTC); require.NoError(t, eLTC); assert.Equal(t, 30.0, bLTC.Total); _, eS2 := accs.GetBalance("s2", &mockCreds, testInvalidAsset, currency.BTC); assert.Error(t, eS2) }},
		{"NilReceiver", &mockCreds, []Change{{Account: "s1", AssetType: asset.Spot, Balance: &initialBTC}}, func(e accountsExchangeDependency) *Accounts{return nil}, common.ErrNilPointer.Error(), nil },
		{"NilSubAccountsMap", &mockCreds, []Change{{Account: "s1", AssetType: asset.Spot, Balance: &initialBTC}}, func(e accountsExchangeDependency) *Accounts { mux := dispatch.GetNewMux(nil); id, _ := mux.GetID(); return &Accounts{Exchange: e, mux: mux, routingID: id, subAccounts: nil} }, common.ErrNilPointer.Error(), nil },
	}
	for _, tt := range tests { t.Run(tt.name, func(t *testing.T) { accs := tt.initState(mockExInst); err := accs.Update(tt.changes, tt.creds); if tt.errContains != "" { require.Error(t, err); assert.ErrorContains(t, err, tt.errContains) } else { require.NoError(t, err) } if tt.verify != nil { tt.verify(t, accs, err) }})}
}

func TestSubAccountsMerge(t *testing.T) {
	t.Parallel()
	sa1BTC1 := Balance{Currency: currency.BTC, Total: 1, Free: 1, UpdatedAt: time.Now()}
	sa1ETH1 := Balance{Currency: currency.ETH, Total: 10, Free: 10, UpdatedAt: time.Now()}
	sa1USDT1 := Balance{Currency: currency.USDT, Total: 100, Free: 100, UpdatedAt: time.Now()}
	sa2LTC1 := Balance{Currency: currency.LTC, Total: 20, Free: 20, UpdatedAt: time.Now()}
	sa1BTC2 := Balance{Currency: currency.BTC, Total: 0.5, Free: 0.5, Borrowed: 0.1, UpdatedAt: time.Now().Add(time.Second)}

	tests := []struct{ name string; initial SubAccounts; toMerge *SubAccount; expected SubAccounts }{
		{"EmptyList", SubAccounts{}, &SubAccount{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}, SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}},
		{"NewID", SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}, &SubAccount{ID: "a2", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Total: 2, UpdatedAt: time.Now()}}}, SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}, {ID: "a2", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Total: 2, UpdatedAt: time.Now()}}}}},
		{"NewAsset", SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}, &SubAccount{ID: "a1", AssetType: asset.Margin, Balances: CurrencyBalances{currency.USDT: sa1USDT1}}, SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}, {ID: "a1", AssetType: asset.Margin, Balances: CurrencyBalances{currency.USDT: sa1USDT1}}}},
		{"MergeExisting", SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}, &SubAccount{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.ETH: sa1ETH1, currency.BTC: sa1BTC2}}, SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: {Currency: sa1BTC2.Currency, Total: 1.5, Free: 0.5, Borrowed: 0.1, UpdatedAt: sa1BTC2.UpdatedAt}, currency.ETH: sa1ETH1}}}},
		{"MergeWithEmptyBalances", SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}, &SubAccount{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.LTC: sa2LTC1}}, SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1, currency.LTC: sa2LTC1}}}},
		{"MergeTotallyEmptyBalances", SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}, &SubAccount{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{}}, SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}},
		{"MergeNil", SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}, nil, SubAccounts{{ID: "a1", AssetType: asset.Spot, Balances: CurrencyBalances{currency.BTC: sa1BTC1}}}},
	}
	for _, tt := range tests { t.Run(tt.name, func(t *testing.T) {
		listCopy := make(SubAccounts, len(tt.initial))
		for i, sa := range tt.initial { balancesCopy := make(CurrencyBalances); for c, b := range sa.Balances { balancesCopy[c] = b }; listCopy[i] = &SubAccount{ID: sa.ID, AssetType: sa.AssetType, Balances: balancesCopy} }
		result := listCopy.Merge(tt.toMerge)
		require.Len(t, result, len(tt.expected))
		expectedMap := make(map[string]*SubAccount); for _, sa := range tt.expected { expectedMap[sa.ID+sa.AssetType.String()] = sa }
		resultMap := make(map[string]*SubAccount); for _, sa := range result { resultMap[sa.ID+sa.AssetType.String()] = sa }
		require.Equal(t, len(expectedMap), len(resultMap))
		for k, expSA := range expectedMap {
			actSA, ok := resultMap[k]; require.True(t, ok)
			assert.Equal(t, expSA.ID, actSA.ID); assert.Equal(t, expSA.AssetType, actSA.AssetType)
			require.Equal(t, len(expSA.Balances), len(actSA.Balances))
			for curr, expBal := range expSA.Balances {
				actBal, currOk := actSA.Balances[curr]; require.True(t, currOk)
				assert.Equal(t, expBal.Total, actBal.Total)
				assert.Equal(t, expBal.Free, actBal.Free)
				assert.Equal(t, expBal.Borrowed, actBal.Borrowed)
				if !expBal.UpdatedAt.IsZero() { assert.WithinDuration(t, expBal.UpdatedAt, actBal.UpdatedAt, time.Millisecond) }
			}
		}
	})}
}
*/
