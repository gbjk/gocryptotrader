package accounts

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/key"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/dispatch"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

type MockEx struct{}

// GetName returns exchange name
func (_ *MockEx) GetName() string {
	return "mocky"
}

var (
	happyCredentials = &Credentials{Key: "AAAAA"}
	m                = &MockEx{}
)

func TestCollectBalances(t *testing.T) {
	t.Parallel()

	_, err := CollectBalances(nil, asset.Spot)
	require.ErrorIs(t, err, common.ErrNilPointer)

	accounts, err := CollectBalances(map[string][]Balance{}, asset.Spot)
	require.NoError(t, err)
	assert.Empty(t, accounts, "CollectBalances should return empty when given an empty map")

	_, err = CollectBalances(map[string][]Balance{}, asset.Empty)
	require.ErrorIs(t, err, asset.ErrNotSupported)

	accounts, err = CollectBalances(map[string][]Balance{
		"someAccountID": {{Currency: currency.BTC, Total: 40000, Hold: 1}},
	}, asset.Spot)

	require.NoError(t, err)
	require.Equal(t, 1, len(accounts), "CollectBalances must return one sub-account")

	subAccount := accounts[0]
	require.Equal(t, 1, len(subAccount.Currencies), "CollectBalances must return one Currency Balance")
	balance := subAccount.Currencies[0]

	assert.Equal(t, "someAccountID", subAccount.ID, "subAccountID should be correct")
	assert.Equal(t, asset.Spot, subAccount.AssetType, "AssetType should be correct")
	assert.Equal(t, currency.BTC, balance.Currency, "Currency should be correct")
	assert.Equal(t, 40000.0, balance.Total, "Total should be correct")
	assert.Equal(t, 1.0, balance.Hold, "Hold should be correct")
}

func TestGetHoldings(t *testing.T) {
	t.Parallel()

	a := MustNewAccounts(m, dispatch.GetNewMux(nil))

	require.NoError(t, a.Save([]SubAccount{}, happyCredentials), "Save must not error")

	_, err := a.GetHoldings(happyCredentials, asset.Options)
	require.ErrorIs(t, err, ErrExchangeHoldingsNotFound)

	s := []SubAccount{
		{
			AssetType: asset.UpsideProfitContract,
			ID:        "1337",
		},
		{
			AssetType: asset.Spot,
			ID:        "1337",
			Currencies: []Balance{
				{
					Currency: currency.BTC,
					Total:    100,
					Hold:     20,
				},
			},
		},
	}
	require.NoError(t, a.Save(s, happyCredentials), "Save must not error")

	_, err = a.GetHoldings(nil, asset.Spot)
	assert.ErrorIs(t, err, errCredentialsEmpty)

	_, err = a.GetHoldings(happyCredentials, asset.Empty)
	assert.ErrorIs(t, err, asset.ErrNotSupported)

	_, err = a.GetHoldings(&Credentials{Key: "BBBBB"}, asset.Spot)
	assert.ErrorIs(t, err, errNoCredentialBalances)

	u, err := a.GetHoldings(happyCredentials, asset.Spot)
	require.NoError(t, err)

	require.Len(t, u, 1)
	assert.Equal(t, "1337", u[0].ID)
	assert.Equal(t, asset.Spot, u[0].AssetType)
	require.Len(t, u[0].Currencies, 1)
	assert.Equal(t, currency.BTC, u[0].Currencies[0].Currency)
	assert.Equal(t, 100.0, u[0].Currencies[0].Total)
	assert.Equal(t, 20.0, u[0].Currencies[0].Hold)
}

func TestGetBalance(t *testing.T) {
	t.Parallel()

	a := MustNewAccounts(m, dispatch.GetNewMux(nil))

	_, err := a.GetBalance("", nil, asset.Empty, currency.Code{})
	assert.ErrorIs(t, err, asset.ErrNotSupported)

	_, err = a.GetBalance("", nil, asset.Spot, currency.Code{})
	assert.ErrorIs(t, err, errCredentialsEmpty)

	_, err = a.GetBalance("", happyCredentials, asset.Spot, currency.Code{})
	assert.ErrorIs(t, err, currency.ErrCurrencyCodeEmpty)

	err = a.Save([]SubAccount{
		{
			AssetType: asset.Spot,
			ID:        "1337",
		},
	}, happyCredentials)
	require.NoError(t, err, "Save must not error")

	_, err = a.GetBalance("1336", &Credentials{Key: "BBBBB"}, asset.Spot, currency.BTC)
	assert.ErrorIs(t, err, errNoCredentialBalances)

	_, err = a.GetBalance("1336", happyCredentials, asset.Spot, currency.BTC)
	assert.ErrorIs(t, err, errNoExchangeSubAccountBalances)

	_, err = a.GetBalance("1337", happyCredentials, asset.Futures, currency.BTC)
	assert.ErrorIs(t, err, errNoExchangeSubAccountBalances)

	err = a.Save([]SubAccount{
		{
			AssetType: asset.Spot,
			ID:        "1337",
			Currencies: []Balance{
				{
					Currency: currency.BTC,
					Total:    2,
					Hold:     1,
				},
			},
		},
	}, happyCredentials)
	require.NoError(t, err, "Save must not error")

	bal, err := a.GetBalance("1337", happyCredentials, asset.Spot, currency.BTC)
	require.NoError(t, err, "get balance must not error")

	bal.m.Lock()
	assert.Equal(t, 2.0, bal.total)
	assert.Equal(t, 1.0, bal.hold)
	bal.m.Unlock()
}

func TestBalanceInternalWait(t *testing.T) {
	t.Parallel()
	_, _, err := (*ProtectedBalance)(nil).Wait(0)
	require.ErrorIs(t, err, common.ErrNilPointer)

	b := &ProtectedBalance{}
	waiter, _, err := b.Wait(time.Nanosecond)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		_, ok := <-waiter
		return ok
	}, 100*time.Millisecond, time.Millisecond, "Wait must publish within a millisecond")

	waiter, _, err = b.Wait(0)
	require.NoError(t, err)

	b.notice.Alert()
	assert.False(t, <-waiter, "Alert should change Waiter to return false")
}

func TestBalanceInternalLoad(t *testing.T) {
	t.Parallel()
	bi := &ProtectedBalance{}
	err := bi.load(nil)
	assert.ErrorIs(t, err, common.ErrNilPointer, "should error nil pointer correctly")

	err = bi.load(&Balance{Total: 1, Hold: 2, Free: 3, AvailableWithoutBorrow: 4, Borrowed: 5})
	assert.ErrorIs(t, err, errUpdatedAtIsZero, "should error correctly when updatedAt is not set")

	now := time.Now()
	err = bi.load(&Balance{UpdatedAt: now, Total: 1, Hold: 2, Free: 3, AvailableWithoutBorrow: 4, Borrowed: 5})
	require.NoError(t, err)

	bi.m.Lock()
	assert.Equal(t, now, bi.updatedAt)
	assert.Equal(t, 1.0, bi.total)
	assert.Equal(t, 2.0, bi.hold)
	assert.Equal(t, 3.0, bi.free)
	assert.Equal(t, 4.0, bi.availableWithoutBorrow)
	assert.Equal(t, 5.0, bi.borrowed)
	bi.m.Unlock()

	assert.Equal(t, 3.0, bi.GetFree())

	err = bi.load(&Balance{UpdatedAt: now, Total: 2, Hold: 3, Free: 4, AvailableWithoutBorrow: 5, Borrowed: 6})
	assert.ErrorIs(t, err, errOutOfSequence, "should error correctly with same UpdatedAt")

	err = bi.load(&Balance{UpdatedAt: now.Add(time.Second), Total: 2, Hold: 3, Free: 4, AvailableWithoutBorrow: 5, Borrowed: 6})
	assert.NoError(t, err)
}

func TestGetFree(t *testing.T) {
	t.Parallel()
	assert.Zero(t, (*ProtectedBalance)(nil))
	assert.Equal(t, 1.0, (&ProtectedBalance{free: 1}).GetFree())
}

func TestSave(t *testing.T) {
	t.Parallel()

	s := []SubAccount{}
	err := (*Accounts)(nil).Save(s, nil)
	assert.ErrorIs(t, err, common.ErrNilPointer)
	assert.ErrorContains(t, err, "*account.Accounts")

	err = new(Accounts).Save(s, nil)
	assert.ErrorIs(t, err, errNoExchangeSubAccountBalances)

	a := MustNewAccounts(m, dispatch.GetNewMux(nil))

	err = a.Save(s, nil)
	assert.ErrorIs(t, err, errCredentialsEmpty)

	err = dispatch.Start(dispatch.DefaultMaxWorkers, dispatch.DefaultJobsLimit)
	require.NoError(t, common.ExcludeError(err, dispatch.ErrDispatcherAlreadyRunning), "dispatch.Start must not error")

	p, err := a.Subscribe()
	require.NoError(t, err)
	require.NotNil(t, p, "Subscribe must return a pipe")
	require.Empty(t, p.Channel(), "Pipe must be empty before Saving anything")

	s = []SubAccount{
		{AssetType: 6969, ID: "1337", Currencies: []Balance{{Currency: currency.BTC, Total: 100, Hold: 20}}},
		{ID: "1338", AssetType: asset.Options},
	}

	assert.ErrorIs(t, a.Save(s, nil), errCredentialsEmpty)
	assert.ErrorIs(t, a.Save(s, happyCredentials), asset.ErrNotSupported)

	s[0].AssetType = asset.Spot
	require.NoError(t, a.Save(s, happyCredentials))

	updates := map[asset.Item]SubAccount{}
	require.Eventually(t, func() bool {
		if uAny, ok := <-p.Channel(); ok {
			if update, ok := uAny.(SubAccount); ok {
				updates[update.AssetType] = update
			}
		}
		return len(updates) == 2
	}, time.Second, time.Millisecond*10, "Save must publish 2 saves through dispatch channel to subscriber")

	require.Contains(t, updates, asset.Spot, "Save must publish Spot asset update")
	require.Equal(t, s[0], updates[asset.Spot], "Save published Spot update must be correct")
	require.Contains(t, updates, asset.Options, "Save must publish Options asset update")
	require.Equal(t, s[1], updates[asset.Options], "Save published Options update must be correct")
	require.NoError(t, p.Release(), "Releasing the subscription must not error")

	assets, ok := a.subAccounts[*happyCredentials][key.SubAccountAsset{
		SubAccount: "1337",
		Asset:      asset.Spot,
	}]
	require.True(t, ok)

	b, ok := assets[currency.BTC.Item]
	require.True(t, ok)

	assert.NotEmpty(t, b.updatedAt)
	assert.Equal(t, 100.0, b.total)
	assert.Equal(t, 20.0, b.hold)

	s[0].Currencies[0] = Balance{Currency: currency.ETH, Total: 80, Hold: 20}
	require.NoError(t, a.Save(s, happyCredentials))

	b, ok = assets[currency.BTC.Item]
	require.True(t, ok)
	assert.NotEmpty(t, b.updatedAt)
	assert.Zero(t, b.total)
	assert.Zero(t, b.hold)

	e, ok := assets[currency.ETH.Item]
	require.True(t, ok)
	assert.NotEmpty(t, e.updatedAt)
	assert.Equal(t, 80.0, e.total)
	assert.Equal(t, 20.0, e.hold)

	s[0].Currencies[0].UpdatedAt = time.Now().Add(-time.Hour)
	err = a.Save(s, happyCredentials)
	assert.ErrorIs(t, err, errOutOfSequence)

	a.mux = nil
	s[0].Currencies[0].UpdatedAt = time.Now()
	err = a.Save(s, happyCredentials)
	assert.ErrorIs(t, err, common.ErrNilPointer)
	assert.ErrorContains(t, err, "*dispatch.Mux")
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	err := dispatch.Start(dispatch.DefaultMaxWorkers, dispatch.DefaultJobsLimit)
	require.NoError(t, common.ExcludeError(err, dispatch.ErrDispatcherAlreadyRunning), "dispatch.Start must not error")

	a := MustNewAccounts(m, dispatch.GetNewMux(nil))

	require.ErrorIs(t, (*Accounts)(nil).Update(nil, nil), common.ErrNilPointer)
	require.ErrorIs(t, new(Accounts).Update(nil, nil), common.ErrNilPointer)

	err = a.Update(nil, nil)
	assert.ErrorIs(t, err, errCredentialsEmpty)

	err = a.Update([]Change{{AssetType: 6969}}, happyCredentials)
	assert.ErrorIs(t, err, asset.ErrNotSupported)

	p, err := a.Subscribe()
	require.NoError(t, err)
	require.NotNil(t, p, "Subscribe must return a pipe")
	require.Empty(t, p.Channel(), "Pipe must be empty before Saving anything")

	now := time.Now()
	c := []Change{{
		AssetType: asset.Spot,
		Account:   "1337",
		Balance:   &Balance{Currency: currency.BTC, Total: 100, Free: 80, UpdatedAt: now},
	}, {
		AssetType: asset.Options,
		Account:   "1337",
		Balance:   &Balance{Currency: currency.USDT, Total: 20, UpdatedAt: now},
	}}
	err = a.Update(c, happyCredentials)

	require.NoError(t, err)

	updates := map[asset.Item]Change{}
	require.Eventually(t, func() bool {
		if uAny, ok := <-p.Channel(); ok {
			if update, ok := uAny.(Change); ok {
				updates[update.AssetType] = update
			}
		}
		return len(updates) == 2
	}, 2*time.Second, time.Millisecond*10, "Update must publish updates through dispatch channel to subscriber")

	require.Contains(t, updates, asset.Spot, "Update must publish Spot asset update")
	require.Equal(t, c[0], updates[asset.Spot], "Update published Spot update must be correct")
	require.Contains(t, updates, asset.Options, "Update must publish Options asset update")
	require.Equal(t, c[1], updates[asset.Options], "Update published Options update must be correct")
	require.NoError(t, p.Release(), "Releasing the subscription must not error")

	assets, ok := a.subAccounts[*happyCredentials][key.SubAccountAsset{
		SubAccount: "1337",
		Asset:      asset.Spot,
	}]
	require.True(t, ok, "Update must add subAccount for the credentials")

	b, ok := assets[currency.BTC.Item]
	require.True(t, ok, "Update must add currency to the subAccount")

	assert.Equal(t, 100.0, b.total, "Update should set total correctly")
	assert.Equal(t, 80.0, b.free, "Update should set free correctly")
	assert.Equal(t, now, b.updatedAt, "Update should set updatedAt correctly")

	err = a.Update([]Change{{AssetType: asset.Spot, Account: "1337"}}, happyCredentials)
	assert.ErrorIs(t, err, common.ErrNilPointer)
	assert.ErrorContains(t, err, "*account.Balance")

	err = a.Update([]Change{
		{
			AssetType: asset.Spot,
			Account:   "1337",
			Balance: &Balance{
				Currency:  currency.BTC,
				Total:     100,
				Free:      100,
				UpdatedAt: now.Add(-1 * time.Second),
			},
		},
	}, happyCredentials)
	assert.ErrorIs(t, err, errOutOfSequence)

	err = a.Update([]Change{
		{
			AssetType: asset.Spot,
			Account:   "1337",
			Balance: &Balance{
				Currency:  currency.BTC,
				Total:     100,
				Free:      100,
				UpdatedAt: now.Add(1 * time.Second),
			},
		},
	}, happyCredentials)
	require.NoError(t, err)

	assert.Equal(t, 100.0, b.total)
	assert.Equal(t, 100.0, b.free)
	assert.Equal(t, now.Add(1*time.Second), b.updatedAt)
}

func TestTrackNewAccounts(t *testing.T) {
	t.Parallel()

	s := NewStore()

	s.mu.Lock()
	_, err := s.registerExchange(m)
	s.mu.Unlock()
	require.NoError(t, err)

	s.mu.Lock()
	_, err = s.registerExchange(m)
	s.mu.Unlock()
	assert.ErrorIs(t, err, errExchangeAlreadyExists)
}

// TestSubscribe ensures that Subscribe returns a subscription channel
// See TestSave and TestUpdate for exercising publish to subscribers
func TestSubscribe(t *testing.T) {
	t.Parallel()

	err := dispatch.Start(dispatch.DefaultMaxWorkers, dispatch.DefaultJobsLimit)
	require.NoError(t, common.ExcludeError(err, dispatch.ErrDispatcherAlreadyRunning), "dispatch.Start must not error")

	p, err := MustNewAccounts(m, dispatch.GetNewMux(nil)).Subscribe()
	require.NoError(t, err)
	require.NotNil(t, p, "Subscribe must return a pipe")
	require.Empty(t, p.Channel(), "Pipe must be empty before Saving anything")
}
