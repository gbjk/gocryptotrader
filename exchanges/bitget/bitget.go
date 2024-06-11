package bitget

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/request"
)

// Bitget is the overarching type across this package
type Bitget struct {
	exchange.Base
}

const (
	bitgetAPIURL = "https://api.bitget.com/api/v2/"

	// Public endpoints
	bitgetPublic              = "public/"
	bitgetAnnouncements       = "annoucements" // sic
	bitgetTime                = "time"
	bitgetCoins               = "coins"
	bitgetSymbols             = "symbols"
	bitgetMarket              = "market/"
	bitgetVIPFeeRate          = "vip-fee-rate"
	bitgetTickers             = "tickers"
	bitgetMergeDepth          = "merge-depth"
	bitgetOrderbook           = "orderbook"
	bitgetCandles             = "candles"
	bitgetHistoryCandles      = "history-candles"
	bitgetFillsHistory        = "fills-history"
	bitgetTicker              = "ticker"
	bitgetHistoryIndexCandles = "history-index-candles"
	bitgetHistoryMarkCandles  = "history-mark-candles"
	bitgetOpenInterest        = "open-interest"
	bitgetFundingTime         = "funding-time"
	bitgetSymbolPrice         = "symbol-price"
	bitgetHistoryFundRate     = "history-fund-rate"
	bitgetCurrentFundRate     = "current-fund-rate"
	bitgetContracts           = "contracts"
	bitgetQueryPositionLever  = "query-position-lever"
	bitgetCoinInfos           = "coinInfos"
	bitgetHourInterest        = "hour-interest"

	// Mixed endpoints
	bitgetSpot   = "spot/"
	bitgetFills  = "fills"
	bitgetMix    = "mix/"
	bitgetMargin = "margin/"
	bitgetEarn   = "earn/"
	bitgetLoan   = "loan"

	// Authenticated endpoints
	bitgetCommon                   = "common/"
	bitgetTradeRate                = "trade-rate"
	bitgetTax                      = "tax/"
	bitgetSpotRecord               = "spot-record"
	bitgetFutureRecord             = "future-record"
	bitgetMarginRecord             = "margin-record"
	bitgetP2PRecord                = "p2p-record"
	bitgetP2P                      = "p2p/"
	bitgetMerchantList             = "merchantList"
	bitgetMerchantInfo             = "merchantInfo"
	bitgetOrderList                = "orderList"
	bitgetAdvList                  = "advList"
	bitgetUser                     = "user/"
	bitgetCreate                   = "create-"
	bitgetVirtualSubaccount        = "virtual-subaccount"
	bitgetModify                   = "modify-"
	bitgetBatchCreateSubAccApi     = "batch-create-subaccount-and-apikey"
	bitgetList                     = "list"
	bitgetAPIKey                   = "apikey"
	bitgetConvert                  = "convert/"
	bitgetCurrencies               = "currencies"
	bitgetQuotedPrice              = "quoted-price"
	bitgetTrade                    = "trade"
	bitgetConvertRecord            = "convert-record"
	bitgetBGBConvert               = "bgb-convert"
	bitgetConvertCoinList          = "bgb-convert-coin-list"
	bitgetBGBConvertRecords        = "bgb-convert-records"
	bitgetPlaceOrder               = "/place-order"
	bitgetCancelOrder              = "/cancel-order"
	bitgetBatchOrders              = "/batch-orders"
	bitgetBatchCancel              = "/batch-cancel-order"
	bitgetCancelSymbolOrder        = "/cancel-symbol-order"
	bitgetOrderInfo                = "/orderInfo"
	bitgetUnfilledOrders           = "/unfilled-orders"
	bitgetHistoryOrders            = "/history-orders"
	bitgetPlacePlanOrder           = "/place-plan-order"
	bitgetModifyPlanOrder          = "/modify-plan-order"
	bitgetCancelPlanOrder          = "/cancel-plan-order"
	bitgetCurrentPlanOrder         = "/current-plan-order"
	bitgetPlanSubOrder             = "/plan-sub-order"
	bitgetPlanOrderHistory         = "/history-plan-order"
	bitgetBatchCancelPlanOrder     = "/batch-cancel-plan-order"
	bitgetAccount                  = "account"
	bitgetInfo                     = "/info"
	bitgetAssets                   = "/assets"
	bitgetSubaccountAssets         = "/subaccount-assets"
	bitgetWallet                   = "wallet/"
	bitgetModifyDepositAccount     = "modify-deposit-account"
	bitgetBills                    = "/bills"
	bitgetTransfer                 = "transfer"
	bitgetTransferCoinInfo         = "transfer-coin-info"
	bitgetSubaccountTransfer       = "subaccount-transfer"
	bitgetWithdrawal               = "withdrawal"
	bitgetSubaccountTransferRecord = "/sub-main-trans-record"
	bitgetTransferRecord           = "/transferRecords"
	bitgetDepositAddress           = "deposit-address"
	bitgetSubaccountDepositAddress = "subaccount-deposit-address"
	bitgetCancelWithdrawal         = "cancel-withdrawal"
	bitgetSubaccountDepositRecord  = "subaccount-deposit-records"
	bitgetWithdrawalRecord         = "withdrawal-records"
	bitgetDepositRecord            = "deposit-records"
	bitgetAccounts                 = "/accounts"
	bitgetSubaccountAssets2        = "/sub-account-assets"
	bitgetOpenCount                = "/open-count"
	bitgetSetLeverage              = "/set-leverage"
	bitgetSetMargin                = "/set-margin"
	bitgetSetMarginMode            = "/set-margin-mode"
	bitgetSetPositionMode          = "/set-position-mode"
	bitgetBill                     = "/bill"
	bitgetPosition                 = "position/"
	bitgetSinglePosition           = "single-position"
	bitgetAllPositions             = "all-position" // sic
	bitgetHistoryPosition          = "history-position"
	bitgetOrder                    = "order"
	bitgetClickBackhand            = "/click-backhand"
	bitgetBatchPlaceOrder          = "/batch-place-order"
	bitgetModifyOrder              = "/modify-order"
	bitgetBatchCancelOrders        = "/batch-cancel-orders"
	bitgetClosePositions           = "/close-positions"
	bitgetDetail                   = "/detail"
	bitgetOrdersPending            = "/orders-pending"
	bitgetFillHistory              = "/fill-history"
	bitgetOrdersHistory            = "/orders-history"
	bitgetCancelAllOrders          = "/cancel-all-orders"
	bitgetPlaceTPSLOrder           = "/place-tpsl-order"
	bitgetModifyTPSLOrder          = "/modify-tpsl-order"
	bitgetOrdersPlanPending        = "/orders-plan-pending"
	bitgetOrdersPlanHistory        = "/orders-plan-history"
	bitgetCrossed                  = "crossed"
	bitgetBorrowHistory            = "/borrow-history"
	bitgetRepayHistory             = "/repay-history"
	bitgetInterestHistory          = "/interest-history"
	bitgetLiquidationHistory       = "/liquidation-history"
	bitgetFinancialRecords         = "/financial-records"
	bitgetBorrow                   = "/borrow"
	bitgetRepay                    = "/repay"
	bitgetRiskRate                 = "/risk-rate"
	bitgetMaxBorrowableAmount      = "/max-borrowable-amount"
	bitgetMaxTransferOutAmount     = "/max-transfer-out-amount"
	bitgetInterestRateAndLimit     = "/interest-rate-and-limit"
	bitgetTierData                 = "/tier-data"
	bitgetFlashRepay               = "/flash-repay"
	bitgetQueryFlashRepayStatus    = "/query-flash-repay-status"
	bitgetBatchCancelOrder         = "/batch-cancel-order"
	bitgetOpenOrders               = "/open-orders"
	bitgetLiquidationOrder         = "/liquidation-order"
	bitgetIsolated                 = "isolated"
	bitgetSavings                  = "savings"
	bitgetProduct                  = "/product"
	bitgetRecords                  = "/records"
	bitgetSubscribeInfo            = "/subscribe-info"
	bitgetSubscribe                = "/subscribe"
	bitgetSubscribeResult          = "/subscribe-result"
	bitgetRedeem                   = "/redeem"
	bitgetRedeemResult             = "/redeem-result"
	bitgetSharkFin                 = "sharkfin"
	bitgetOngoingOrders            = "/ongoing-orders"
	bitgetRevisePledge             = "/revise-pledge"

	// Errors
	errUnknownEndpointLimit = "unknown endpoint limit %v"
)

var (
	errBusinessTypeEmpty             = errors.New("businessType cannot be empty")
	errPairEmpty                     = errors.New("currency pair cannot be empty")
	errCurrencyEmpty                 = errors.New("currency cannot be empty")
	errProductTypeEmpty              = errors.New("productType cannot be empty")
	errSubaccountEmpty               = errors.New("subaccounts cannot be empty")
	errNewStatusEmpty                = errors.New("newStatus cannot be empty")
	errNewPermsEmpty                 = errors.New("newPerms cannot be empty")
	errPassphraseEmpty               = errors.New("passphrase cannot be empty")
	errLabelEmpty                    = errors.New("label cannot be empty")
	errAPIKeyEmpty                   = errors.New("apiKey cannot be empty")
	errFromToMutex                   = errors.New("exactly one of fromAmount and toAmount must be set")
	errTraceIDEmpty                  = errors.New("traceID cannot be empty")
	errAmountEmpty                   = errors.New("amount cannot be empty")
	errPriceEmpty                    = errors.New("price cannot be empty")
	errTypeAssertTimestamp           = errors.New("unable to type assert timestamp")
	errTypeAssertOpenPrice           = errors.New("unable to type assert opening price")
	errTypeAssertHighPrice           = errors.New("unable to type assert high price")
	errTypeAssertLowPrice            = errors.New("unable to type assert low price")
	errTypeAssertClosePrice          = errors.New("unable to type assert close price")
	errTypeAssertBaseVolume          = errors.New("unable to type assert base volume")
	errTypeAssertQuoteVolume         = errors.New("unable to type assert quote volume")
	errTypeAssertUSDTVolume          = errors.New("unable to type assert USDT volume")
	errGranEmpty                     = errors.New("granularity cannot be empty")
	errEndTimeEmpty                  = errors.New("endTime cannot be empty")
	errSideEmpty                     = errors.New("side cannot be empty")
	errOrderTypeEmpty                = errors.New("orderType cannot be empty")
	errStrategyEmpty                 = errors.New("strategy cannot be empty")
	errLimitPriceEmpty               = errors.New("price cannot be empty for limit orders")
	errOrderClientEmpty              = errors.New("at least one of orderID and clientOrderID must not be empty")
	errOrderIDEmpty                  = errors.New("orderID cannot be empty")
	errOrdersEmpty                   = errors.New("orders cannot be empty")
	errTriggerPriceEmpty             = errors.New("triggerPrice cannot be empty")
	errTriggerTypeEmpty              = errors.New("triggerType cannot be empty")
	errNonsenseRequest               = errors.New("nonsense request expected error")
	errAccountTypeEmpty              = errors.New("accountType cannot be empty")
	errFromTypeEmpty                 = errors.New("fromType cannot be empty")
	errToTypeEmpty                   = errors.New("toType cannot be empty")
	errCurrencyAndPairEmpty          = errors.New("currency and pair cannot both be empty")
	errFromIDEmpty                   = errors.New("fromID cannot be empty")
	errToIDEmpty                     = errors.New("toID cannot be empty")
	errTransferTypeEmpty             = errors.New("transferType cannot be empty")
	errAddressEmpty                  = errors.New("address cannot be empty")
	errNoCandleData                  = errors.New("no candle data")
	errMarginCoinEmpty               = errors.New("marginCoin cannot be empty")
	errOpenAmountEmpty               = errors.New("openAmount cannot be empty")
	errOpenPriceEmpty                = errors.New("openPrice cannot be empty")
	errLeverageEmpty                 = errors.New("leverage cannot be empty")
	errMarginModeEmpty               = errors.New("marginMode cannot be empty")
	errPositionModeEmpty             = errors.New("positionMode cannot be empty")
	errNewClientOrderIDEmpty         = errors.New("newClientOrderID cannot be empty")
	errPlanTypeEmpty                 = errors.New("planType cannot be empty")
	errPlanOrderIDEmpty              = errors.New("planOrderID cannot be empty")
	errHoldSideEmpty                 = errors.New("holdSide cannot be empty")
	errExecutePriceEmpty             = errors.New("executePrice cannot be empty")
	errTakeProfitParamsInconsistency = errors.New("takeProfitTriggerPrice, takeProfitExecutePrice, and takeProfitTriggerType must either all be set or all be empty")
	errStopLossParamsInconsistency   = errors.New("stopLossTriggerPrice, stopLossExecutePrice, and stopLossTriggerType must either all be set or all be empty")
	errIDListEmpty                   = errors.New("idList cannot be empty")
	errLoanTypeEmpty                 = errors.New("loanType cannot be empty")
	errProductIDEmpty                = errors.New("productID cannot be empty")
	errPeriodTypeEmpty               = errors.New("periodType cannot be empty")
	errLoanCoinEmpty                 = errors.New("loanCoin cannot be empty")
	errCollateralCoinEmpty           = errors.New("collateralCoin cannot be empty")
	errTermEmpty                     = errors.New("term cannot be empty")
	errCollateralAmountEmpty         = errors.New("collateralAmount cannot be empty")
	errCollateralLoanMutex           = errors.New("exactly one of collateralAmount and loanAmount must be set")
	errReviseTypeEmpty               = errors.New("reviseType cannot be empty")
)

// QueryAnnouncement returns announcements from the exchange, filtered by type and time
func (bi *Bitget) QueryAnnouncements(ctx context.Context, annType string, startTime, endTime time.Time) (*AnnResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("annType", annType)
	params.Values.Set("language", "en_US")
	var resp *AnnResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, bitgetPublic+bitgetAnnouncements, params.Values,
		&resp)
}

// GetTime returns the server's time
func (bi *Bitget) GetTime(ctx context.Context) (*TimeResp, error) {
	var resp *TimeResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, bitgetPublic+bitgetTime, nil, &resp)
}

// GetTradeRate returns the fees the user would face for trading a given symbol
func (bi *Bitget) GetTradeRate(ctx context.Context, pair, businessType string) (*TradeRateResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if businessType == "" {
		return nil, errBusinessTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("businessType", businessType)
	var resp *TradeRateResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetCommon+bitgetTradeRate, vals, nil, &resp)
}

// GetSpotTransactionRecords returns the user's spot transaction records
func (bi *Bitget) GetSpotTransactionRecords(ctx context.Context, currency string, startTime, endTime time.Time, limit, pagination int64) (*SpotTrResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	params.Values.Set("coin", currency)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	var resp *SpotTrResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate1, http.MethodGet,
		bitgetTax+bitgetSpotRecord, params.Values, nil, &resp)
}

// GetFuturesTransactionRecords returns the user's futures transaction records
func (bi *Bitget) GetFuturesTransactionRecords(ctx context.Context, productType, currency string, startTime, endTime time.Time, limit, pagination int64) (*FutureTrResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	params.Values.Set("productType", productType)
	params.Values.Set("marginCoin", currency)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	var resp *FutureTrResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate1, http.MethodGet,
		bitgetTax+bitgetFutureRecord, params.Values, nil, &resp)
}

// GetMarginTransactionRecords returns the user's margin transaction records
func (bi *Bitget) GetMarginTransactionRecords(ctx context.Context, marginType, currency string, startTime, endTime time.Time, limit, pagination int64) (*MarginTrResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	params.Values.Set("marginType", marginType)
	params.Values.Set("coin", currency)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	var resp *MarginTrResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate1, http.MethodGet,
		bitgetTax+bitgetMarginRecord, params.Values, nil, &resp)
}

// GetP2PTransactionRecords returns the user's P2P transaction records
func (bi *Bitget) GetP2PTransactionRecords(ctx context.Context, currency string, startTime, endTime time.Time, limit, pagination int64) (*P2PTrResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	params.Values.Set("coin", currency)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	var resp *P2PTrResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate1, http.MethodGet,
		bitgetTax+bitgetP2PRecord, params.Values, nil, &resp)
}

// GetP2PMerchantList returns detailed information on a particular merchant
func (bi *Bitget) GetP2PMerchantList(ctx context.Context, online, merchantID string, limit, pagination int64) (*P2PMerListResp, error) {
	vals := url.Values{}
	vals.Set("online", online)
	vals.Set("merchantId", merchantID)
	if limit != 0 {
		vals.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		vals.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	var resp *P2PMerListResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetP2P+bitgetMerchantList, vals, nil, &resp)
}

// GetMerchantInfo returns detailed information on the user as a merchant
func (bi *Bitget) GetMerchantInfo(ctx context.Context) (*P2PMerInfoResp, error) {
	var resp *P2PMerInfoResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetP2P+bitgetMerchantInfo, nil, nil, &resp)
}

// GetMerchantP2POrders returns information on the user's P2P orders
func (bi *Bitget) GetMerchantP2POrders(ctx context.Context, startTime, endTime time.Time, limit, pagination, adNum, ordNum int64, status, side, cryptoCurrency, fiatCurrency string) (*P2POrdersResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("advNo", strconv.FormatInt(adNum, 10))
	params.Values.Set("orderNo", strconv.FormatInt(ordNum, 10))
	params.Values.Set("status", status)
	params.Values.Set("side", side)
	params.Values.Set("coin", cryptoCurrency)
	// params.Values.Set("language", "en-US")
	params.Values.Set("fiat", fiatCurrency)
	var resp *P2POrdersResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetP2P+bitgetOrderList, params.Values, nil, &resp)
}

// GetMerchantAdvertisementList returns information on a variety of merchant advertisements
func (bi *Bitget) GetMerchantAdvertisementList(ctx context.Context, startTime, endTime time.Time, limit, pagination, adNum, payMethodID int64, status, side, cryptoCurrency, fiatCurrency, orderBy, sourceType string) (*P2PAdListResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("advNo", strconv.FormatInt(adNum, 10))
	params.Values.Set("payMethodId", strconv.FormatInt(payMethodID, 10))
	params.Values.Set("status", status)
	params.Values.Set("side", side)
	params.Values.Set("coin", cryptoCurrency)
	// params.Values.Set("language", "en-US")
	params.Values.Set("fiat", fiatCurrency)
	params.Values.Set("orderBy", orderBy)
	params.Values.Set("sourceType", sourceType)
	var resp *P2PAdListResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetP2P+bitgetAdvList, params.Values, nil, &resp)
}

// CreateVirtualSubaccounts creates a batch of virtual subaccounts. These names must use English letters,
// no spaces, no numbers, and be exactly 8 characters long.
func (bi *Bitget) CreateVirtualSubaccounts(ctx context.Context, subaccounts []string) (*CrVirSubResp, error) {
	if len(subaccounts) == 0 {
		return nil, errSubaccountEmpty
	}
	path := bitgetUser + bitgetCreate + bitgetVirtualSubaccount
	req := map[string]interface{}{
		"subAccountList": subaccounts,
	}
	var resp *CrVirSubResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path,
		nil, req, &resp)
}

// ModifyVirtualSubaccount changes the permissions and/or status of a virtual subaccount
func (bi *Bitget) ModifyVirtualSubaccount(ctx context.Context, subaccountID, newStatus string, newPerms []string) (*SuccessBoolResp, error) {
	if subaccountID == "" {
		return nil, errSubaccountEmpty
	}
	if newStatus == "" {
		return nil, errNewStatusEmpty
	}
	if len(newPerms) == 0 {
		return nil, errNewPermsEmpty
	}
	path := bitgetUser + bitgetModify + bitgetVirtualSubaccount
	req := map[string]interface{}{
		"subAccountUid": subaccountID,
		"status":        newStatus,
		"permList":      newPerms,
	}
	var resp *SuccessBoolResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path,
		nil, req, &resp)
}

// CreateSubaccountAndAPIKey creates a subaccounts and an API key. Every account can have up to 20 sub-accounts,
// and every API key can have up to 10 API keys. The name of the sub-account must be exactly 8 English letters.
// The passphrase of the API key must be 8-32 letters and/or numbers. The label must be 20 or fewer characters.
// A maximum of 30 IPs can be a part of the whitelist.
func (bi *Bitget) CreateSubaccountAndAPIKey(ctx context.Context, subaccountName, passphrase, label string, whiteList, permList []string) (*CrSubAccAPIKeyResp, error) {
	if subaccountName == "" {
		return nil, errSubaccountEmpty
	}
	req := map[string]interface{}{
		"subAccountName": subaccountName,
		"passphrase":     passphrase,
		"label":          label,
		"ipList":         whiteList,
		"permList":       permList,
	}
	var resp *CrSubAccAPIKeyResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate1, http.MethodPost,
		bitgetUser+bitgetBatchCreateSubAccApi, nil, req, &resp)
}

// GetVirtualSubaccounts returns a list of the user's virtual sub-accounts
func (bi *Bitget) GetVirtualSubaccounts(ctx context.Context, limit, pagination int64, status string) (*GetVirSubResp, error) {
	vals := url.Values{}
	if limit != 0 {
		vals.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		vals.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	vals.Set("status", status)
	path := bitgetUser + bitgetVirtualSubaccount + "-" + bitgetList
	var resp *GetVirSubResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate2, http.MethodGet, path, vals,
		nil, &resp)
}

// CreateAPIKey creates an API key for the selected virtual sub-account
func (bi *Bitget) CreateAPIKey(ctx context.Context, subaccountID, passphrase, label string, whiteList, permList []string) (*AlterAPIKeyResp, error) {
	if subaccountID == "" {
		return nil, errSubaccountEmpty
	}
	if passphrase == "" {
		return nil, errPassphraseEmpty
	}
	if label == "" {
		return nil, errLabelEmpty
	}
	path := bitgetUser + bitgetCreate + bitgetVirtualSubaccount + "-" + bitgetAPIKey
	req := map[string]interface{}{
		"subAccountUid": subaccountID,
		"passphrase":    passphrase,
		"label":         label,
		"ipList":        whiteList,
		"permList":      permList,
	}
	var resp *AlterAPIKeyResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path,
		nil, req, &resp)
}

// ModifyAPIKey modifies the label, IP whitelist, and/or permissions of the API key associated with the selected
// virtual sub-account
func (bi *Bitget) ModifyAPIKey(ctx context.Context, subaccountID, passphrase, label, apiKey string, whiteList, permList []string) (*AlterAPIKeyResp, error) {
	if apiKey == "" {
		return nil, errAPIKeyEmpty
	}
	if passphrase == "" {
		return nil, errPassphraseEmpty
	}
	if label == "" {
		return nil, errLabelEmpty
	}
	if subaccountID == "" {
		return nil, errSubaccountEmpty
	}
	path := bitgetUser + bitgetModify + bitgetVirtualSubaccount + "-" + bitgetAPIKey
	req := make(map[string]interface{})
	req["subAccountUid"] = subaccountID
	req["passphrase"] = passphrase
	req["label"] = label
	req["subAccountApiKey"] = apiKey
	req["ipList"] = whiteList
	req["permList"] = permList
	var resp *AlterAPIKeyResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path,
		nil, req, &resp)
}

// GetAPIKeys lists the API keys associated with the selected virtual sub-account
func (bi *Bitget) GetAPIKeys(ctx context.Context, subaccountID string) (*GetAPIKeyResp, error) {
	if subaccountID == "" {
		return nil, errSubaccountEmpty
	}
	vals := url.Values{}
	vals.Set("subAccountUid", subaccountID)
	path := bitgetUser + bitgetVirtualSubaccount + "-" + bitgetAPIKey + "-" + bitgetList
	var resp *GetAPIKeyResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals,
		nil, &resp)
}

// GetConvertCoins returns a list of supported currencies, your balance in those currencies, and the maximum and
// minimum tradable amounts of those currencies
func (bi *Bitget) GetConvertCoins(ctx context.Context) (*ConvertCoinsResp, error) {
	var resp *ConvertCoinsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetConvert+bitgetCurrencies, nil, nil, &resp)
}

// GetQuotedPrice returns the price of a given amount of one currency in terms of another currency, and an
// ID for this quote, to be used in a subsequent conversion
func (bi *Bitget) GetQuotedPrice(ctx context.Context, fromCurrency, toCurrency string, fromAmount, toAmount float64) (*QuotedPriceResp, error) {
	if fromCurrency == "" || toCurrency == "" {
		return nil, errCurrencyEmpty
	}
	if (fromAmount == 0 && toAmount == 0) || (fromAmount != 0 && toAmount != 0) {
		return nil, errFromToMutex
	}
	vals := url.Values{}
	vals.Set("fromCoin", fromCurrency)
	vals.Set("toCoin", toCurrency)
	if fromAmount != 0 {
		vals.Set("fromCoinSize", strconv.FormatFloat(fromAmount, 'f', -1, 64))
	} else {
		vals.Set("toCoinSize", strconv.FormatFloat(toAmount, 'f', -1, 64))
	}
	var resp *QuotedPriceResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetConvert+bitgetQuotedPrice, vals, nil, &resp)
}

// CommitConversion commits a conversion previously quoted by GetQuotedPrice. This quote has to have been issued
// within the last 8 seconds.
func (bi *Bitget) CommitConversion(ctx context.Context, fromCurrency, toCurrency, traceID string, fromAmount, toAmount, price float64) (*CommitConvResp, error) {
	if fromCurrency == "" || toCurrency == "" {
		return nil, errCurrencyEmpty
	}
	if traceID == "" {
		return nil, errTraceIDEmpty
	}
	if fromAmount == 0 || toAmount == 0 {
		return nil, errAmountEmpty
	}
	if price == 0 {
		return nil, errPriceEmpty
	}
	req := map[string]interface{}{
		"fromCoin":     fromCurrency,
		"toCoin":       toCurrency,
		"traceId":      traceID,
		"fromCoinSize": strconv.FormatFloat(fromAmount, 'f', -1, 64),
		"toCoinSize":   strconv.FormatFloat(toAmount, 'f', -1, 64),
		"cnvtPrice":    strconv.FormatFloat(price, 'f', -1, 64),
	}
	var resp *CommitConvResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost,
		bitgetConvert+bitgetTrade, nil, req, &resp)
}

// GetConvertHistory returns a list of the user's previous conversions
func (bi *Bitget) GetConvertHistory(ctx context.Context, startTime, endTime time.Time, limit, pagination int64) (*ConvHistResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	var resp *ConvHistResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetConvert+bitgetConvertRecord, params.Values, nil, &resp)
}

// GetBGBConvertCoins returns a list of available currencies, with information on converting them to BGB
func (bi *Bitget) GetBGBConvertCoins(ctx context.Context) (*BGBConvertCoinsResp, error) {
	var resp *BGBConvertCoinsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetConvert+bitgetConvertCoinList, nil, nil, &resp)
}

// ConvertBGB converts all funds in the listed currencies to BGB
func (bi *Bitget) ConvertBGB(ctx context.Context, currencies []string) (*ConvertBGBResp, error) {
	if len(currencies) == 0 {
		return nil, errCurrencyEmpty
	}
	req := map[string]interface{}{
		"coinList": currencies,
	}
	var resp *ConvertBGBResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost,
		bitgetConvert+bitgetBGBConvert, nil, req, &resp)
}

// GetBGBConvertHistory returns a list of the user's previous BGB conversions
func (bi *Bitget) GetBGBConvertHistory(ctx context.Context, orderID, limit, pagination int64, startTime, endTime time.Time) (*BGBConvHistResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	var resp *BGBConvHistResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet,
		bitgetConvert+bitgetBGBConvertRecords, params.Values, nil, &resp)
}

// GetCoinInfo returns information on all supported spot currencies, or a single currency of the user's choice
func (bi *Bitget) GetCoinInfo(ctx context.Context, currency string) (*CoinInfoResp, error) {
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetSpot + bitgetPublic + bitgetCoins
	var resp *CoinInfoResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate3, path, vals, &resp)
}

// GetSymbolInfo returns information on all supported spot trading pairs, or a single pair of the user's choice
func (bi *Bitget) GetSymbolInfo(ctx context.Context, pair string) (*SymbolInfoResp, error) {
	vals := url.Values{}
	vals.Set("symbol", pair)
	path := bitgetSpot + bitgetPublic + bitgetSymbols
	var resp *SymbolInfoResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetSpotVIPFeeRate returns the different levels of VIP fee rates for spot trading
func (bi *Bitget) GetSpotVIPFeeRate(ctx context.Context) (*VIPFeeRateResp, error) {
	path := bitgetSpot + bitgetMarket + bitgetVIPFeeRate
	var resp *VIPFeeRateResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, nil, &resp)
}

// GetSpotTickerInformation returns the ticker information for all trading pairs, or a single pair of the user's
// choice
func (bi *Bitget) GetSpotTickerInformation(ctx context.Context, pair string) (*TickerResp, error) {
	vals := url.Values{}
	vals.Set("symbol", pair)
	path := bitgetSpot + bitgetMarket + bitgetTickers
	var resp *TickerResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetSpotMergeDepth returns part of the orderbook, with options to merge orders of similar price levels together,
// and to change how many results are returned. Limit's a string instead of the typical int64 because the API
// will accept a value of "max"
func (bi *Bitget) GetSpotMergeDepth(ctx context.Context, pair, precision, limit string) (*DepthResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("precision", precision)
	vals.Set("limit", limit)
	path := bitgetSpot + bitgetMarket + bitgetMergeDepth
	var resp *DepthResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetOrderbookDepth returns the orderbook for a given trading pair, with options to merge orders of similar price
// levels together, and to change how many results are returned.
func (bi *Bitget) GetOrderbookDepth(ctx context.Context, pair, step string, limit uint8) (*OrderbookResp, error) {
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("type", step)
	vals.Set("limit", strconv.FormatUint(uint64(limit), 10))
	path := bitgetSpot + bitgetMarket + bitgetOrderbook
	var resp *OrderbookResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetSpotCandlestickData returns candlestick data for a given trading pair
func (bi *Bitget) GetSpotCandlestickData(ctx context.Context, pair, granularity string, startTime, endTime time.Time, limit uint16, historic bool) (*CandleData, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if granularity == "" {
		return nil, errGranEmpty
	}
	var path string
	var params Params
	params.Values = make(url.Values)
	if historic {
		if endTime.IsZero() || endTime.Equal(time.Unix(0, 0)) {
			return nil, errEndTimeEmpty
		}
		path = bitgetSpot + bitgetMarket + bitgetHistoryCandles
		params.Values.Set("endTime", strconv.FormatInt(endTime.UnixMilli(), 10))
	} else {
		path = bitgetSpot + bitgetMarket + bitgetCandles
		err := params.prepareDateString(startTime, endTime, true, true)
		if err != nil {
			return nil, err
		}
	}
	return bi.candlestickHelper(ctx, pair, granularity, path, limit, params)
}

// GetRecentSpotFills returns the most recent trades for a given pair
func (bi *Bitget) GetRecentSpotFills(ctx context.Context, pair string, limit uint16) (*MarketFillsResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("limit", strconv.FormatInt(int64(limit), 10))
	path := bitgetSpot + bitgetMarket + bitgetFills
	var resp *MarketFillsResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, vals, &resp)
}

// GetSpotMarketTrades returns trades for a given pair within a particular time range, and/or before a certain ID
func (bi *Bitget) GetSpotMarketTrades(ctx context.Context, pair string, startTime, endTime time.Time, limit, pagination int64) (*MarketFillsResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetSpot + bitgetMarket + bitgetFillsHistory
	var resp *MarketFillsResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, params.Values, &resp)
}

// PlaceSpotOrder places a spot order on the exchange
func (bi *Bitget) PlaceSpotOrder(ctx context.Context, pair, side, orderType, strategy, clientOrderID string, price, amount float64, isCopyTradeLeader bool) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if side == "" {
		return nil, errSideEmpty
	}
	if orderType == "" {
		return nil, errOrderTypeEmpty
	}
	if strategy == "" {
		return nil, errStrategyEmpty
	}
	if orderType == "limit" && price == 0 {
		return nil, errLimitPriceEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"side":      side,
		"orderType": orderType,
		"force":     strategy,
		"price":     strconv.FormatFloat(price, 'f', -1, 64),
		"size":      strconv.FormatFloat(amount, 'f', -1, 64),
		"clientOid": clientOrderID,
	}
	path := bitgetSpot + bitgetTrade + bitgetPlaceOrder
	var resp *OrderResp
	// I suspect the two rate limits have to do with distinguishing ordinary traders, and traders who are also
	// copy trade leaders. Since this isn't detectable, it'll be handled in the relevant functions through a bool
	rLim := Rate10
	if isCopyTradeLeader {
		rLim = Rate1
	}
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, rLim, http.MethodPost, path, nil, req,
		&resp)
}

// CancelSpotOrderByID cancels an order on the exchange
func (bi *Bitget) CancelSpotOrderByID(ctx context.Context, pair, clientOrderID string, orderID int64) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if orderID == 0 && clientOrderID == "" {
		return nil, errOrderClientEmpty
	}
	req := map[string]interface{}{
		"symbol": pair,
	}
	if orderID != 0 {
		req["orderId"] = strconv.FormatInt(orderID, 10)
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	path := bitgetSpot + bitgetTrade + bitgetCancelOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// BatchPlaceSpotOrders places up to fifty orders on the exchange
func (bi *Bitget) BatchPlaceSpotOrders(ctx context.Context, pair string, orders []PlaceSpotOrderStruct, isCopyTradeLeader bool) (*BatchOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if len(orders) == 0 {
		return nil, errOrdersEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderList": orders,
	}
	path := bitgetSpot + bitgetTrade + bitgetBatchOrders
	var resp *BatchOrderResp
	rLim := Rate5
	if isCopyTradeLeader {
		rLim = Rate1
	}
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, rLim, http.MethodPost, path, nil, req,
		&resp)
}

// BatchCancelOrders cancels up to fifty orders on the exchange
func (bi *Bitget) BatchCancelOrders(ctx context.Context, pair string, orderIDs []OrderIDStruct) (*BatchOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if len(orderIDs) == 0 {
		return nil, errOrderIDEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderList": orderIDs,
	}
	path := bitgetSpot + bitgetTrade + bitgetBatchCancel
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// CancelOrderBySymbol cancels orders for a given symbol. Doesn't return information on failures/successes
func (bi *Bitget) CancelOrderBySymbol(ctx context.Context, pair string) (*SymbolResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	req := map[string]interface{}{
		"symbol": pair,
	}
	path := bitgetSpot + bitgetTrade + bitgetCancelSymbolOrder
	var resp *SymbolResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path, nil, req,
		&resp)
}

// GetSpotOrderDetails returns information on a single order
func (bi *Bitget) GetSpotOrderDetails(ctx context.Context, orderID int64, clientOrderID string) (*SpotOrderDetailResp, error) {
	if orderID == 0 && clientOrderID == "" {
		return nil, errOrderClientEmpty
	}
	vals := url.Values{}
	if orderID != 0 {
		vals.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if clientOrderID != "" {
		vals.Set("clientOid", clientOrderID)
	}
	path := bitgetSpot + bitgetTrade + bitgetOrderInfo
	return bi.spotOrderHelper(ctx, path, vals)
}

// GetUnfilledOrders returns information on the user's unfilled orders
func (bi *Bitget) GetUnfilledOrders(ctx context.Context, pair string, startTime, endTime time.Time, limit, pagination, orderID int64) (*UnfilledOrdersResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	path := bitgetSpot + bitgetTrade + bitgetUnfilledOrders
	var resp *UnfilledOrdersResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetHistoricalSpotOrders returns the user's spot order history
func (bi *Bitget) GetHistoricalSpotOrders(ctx context.Context, pair string, startTime, endTime time.Time, limit, pagination, orderID int64) (*SpotOrderDetailResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	path := bitgetSpot + bitgetTrade + bitgetHistoryOrders
	return bi.spotOrderHelper(ctx, path, params.Values)
}

// GetSpotFills returns information on the user's fulfilled orders in a certain pair
func (bi *Bitget) GetSpotFills(ctx context.Context, pair string, startTime, endTime time.Time, limit, pagination, orderID int64) (*SpotFillsResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	path := bitgetSpot + bitgetTrade + "/" + bitgetFills
	var resp *SpotFillsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// PlacePlanSpotOrder sets up an order to be placed after certain conditions are met
func (bi *Bitget) PlacePlanSpotOrder(ctx context.Context, pair, side, orderType, planType, triggerType, clientOrderID, strategy string, triggerPrice, executePrice, amount float64) (*OrderIDResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if side == "" {
		return nil, errSideEmpty
	}
	if triggerPrice == 0 {
		return nil, errTriggerPriceEmpty
	}
	if orderType == "" {
		return nil, errOrderTypeEmpty
	}
	if orderType == "limit" && executePrice == 0 {
		return nil, errLimitPriceEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	if triggerType == "" {
		return nil, errTriggerTypeEmpty
	}
	req := map[string]interface{}{
		"symbol":       pair,
		"side":         side,
		"triggerPrice": strconv.FormatFloat(triggerPrice, 'f', -1, 64),
		"orderType":    orderType,
		"executePrice": strconv.FormatFloat(executePrice, 'f', -1, 64),
		"planType":     planType,
		"size":         strconv.FormatFloat(amount, 'f', -1, 64),
		"triggerType":  triggerType,
		"force":        strategy,
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	path := bitgetSpot + bitgetTrade + bitgetPlacePlanOrder
	var resp *OrderIDResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodPost, path, nil, req,
		&resp)
}

// ModifyPlanSpotOrder alters the price, trigger price, amount, or order type of a plan order
func (bi *Bitget) ModifyPlanSpotOrder(ctx context.Context, orderID int64, clientOrderID, orderType string, triggerPrice, executePrice, amount float64) (*OrderIDResp, error) {
	if orderID == 0 && clientOrderID == "" {
		return nil, errOrderClientEmpty
	}
	if orderType == "" {
		return nil, errOrderTypeEmpty
	}
	if triggerPrice == 0 {
		return nil, errTriggerPriceEmpty
	}
	if orderType == "limit" && executePrice == 0 {
		return nil, errLimitPriceEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"orderType":    orderType,
		"triggerPrice": strconv.FormatFloat(triggerPrice, 'f', -1, 64),
		"executePrice": strconv.FormatFloat(executePrice, 'f', -1, 64),
		"size":         strconv.FormatFloat(amount, 'f', -1, 64),
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	if orderID != 0 {
		req["orderId"] = strconv.FormatInt(orderID, 10)
	}
	path := bitgetSpot + bitgetTrade + bitgetModifyPlanOrder
	var resp *OrderIDResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodPost, path, nil, req,
		&resp)
}

// CancelPlanSpotOrder cancels a plan order
func (bi *Bitget) CancelPlanSpotOrder(ctx context.Context, orderID int64, clientOrderID string) (*SuccessBoolResp, error) {
	if orderID == 0 && clientOrderID == "" {
		return nil, errOrderClientEmpty
	}
	req := make(map[string]interface{})
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	if orderID != 0 {
		req["orderId"] = strconv.FormatInt(orderID, 10)
	}
	path := bitgetSpot + bitgetTrade + bitgetCancelPlanOrder
	var resp *SuccessBoolResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodPost, path, nil, req,
		&resp)
}

// GetCurrentSpotPlanOrders returns the user's current plan orders
func (bi *Bitget) GetCurrentSpotPlanOrders(ctx context.Context, pair string, startTime, endTime time.Time, limit, pagination int64) (*PlanSpotOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetSpot + bitgetTrade + bitgetCurrentPlanOrder
	var resp *PlanSpotOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetSpotPlanSubOrder returns the sub-orders of a triggered plan order
func (bi *Bitget) GetSpotPlanSubOrder(ctx context.Context, orderID string) (*SubOrderResp, error) {
	if orderID == "" {
		return nil, errOrderIDEmpty
	}
	vals := url.Values{}
	vals.Set("planOrderId", orderID)
	path := bitgetSpot + bitgetTrade + bitgetPlanSubOrder
	var resp *SubOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodGet, path, vals, nil,
		&resp)
}

// GetSpotPlanOrderHistory returns the user's plan order history
func (bi *Bitget) GetSpotPlanOrderHistory(ctx context.Context, pair string, startTime, endTime time.Time, limit int32) (*PlanSpotOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(int64(limit), 10))
	}
	path := bitgetSpot + bitgetTrade + bitgetPlanOrderHistory
	var resp *PlanSpotOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodGet, path, params.Values,
		nil, &resp)
}

// BatchCancelSpotPlanOrders cancels all plan orders, with the option to restrict to only those for particular pairs
func (bi *Bitget) BatchCancelSpotPlanOrders(ctx context.Context, pairs []string) (*BatchOrderResp, error) {
	req := make(map[string]interface{})
	if len(pairs) > 0 {
		req["symbolList"] = pairs
	}
	path := bitgetSpot + bitgetTrade + bitgetBatchCancelPlanOrder
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetAccountInfo returns the user's account information
func (bi *Bitget) GetAccountInfo(ctx context.Context) (*AccountInfoResp, error) {
	path := bitgetSpot + bitgetAccount + bitgetInfo
	var resp *AccountInfoResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate1, http.MethodGet, path, nil, nil, &resp)
}

// GetAccountAssets returns information on the user's assets
func (bi *Bitget) GetAccountAssets(ctx context.Context, currency, assetType string) (*AccountAssetsResp, error) {
	vals := url.Values{}
	vals.Set("coin", currency)
	vals.Set("type", assetType)
	path := bitgetSpot + bitgetAccount + bitgetAssets
	var resp *AccountAssetsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil, &resp)
}

// GetSpotSubaccountAssets returns information on assets in the user's sub-accounts
func (bi *Bitget) GetSpotSubaccountAssets(ctx context.Context) (*SubaccountAssetsResp, error) {
	path := bitgetSpot + bitgetAccount + bitgetSubaccountAssets
	var resp *SubaccountAssetsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, nil, nil, &resp)
}

// ModifyDepositAccount changes which account is automatically used for deposits of a particular currency
func (bi *Bitget) ModifyDepositAccount(ctx context.Context, accountType, currency string) (*SuccessBoolResp2, error) {
	if accountType == "" {
		return nil, errAccountTypeEmpty
	}
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	req := map[string]interface{}{
		"coin":        currency,
		"accountType": accountType,
	}
	path := bitgetSpot + bitgetWallet + bitgetModifyDepositAccount
	var resp *SuccessBoolResp2
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetSpotAccountBills returns a section of the user's billing history
func (bi *Bitget) GetSpotAccountBills(ctx context.Context, currency, groupType, businessType string, startTime, endTime time.Time, limit, pagination int64) (*SpotAccBillResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	if currency != "" {
		params.Values.Set("coin", currency)
	}
	params.Values.Set("groupType", groupType)
	params.Values.Set("businessType", businessType)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetSpot + bitgetAccount + bitgetBills
	var resp *SpotAccBillResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// TransferAsset transfers a certain amount of a currency or pair between different productType accounts
func (bi *Bitget) TransferAsset(ctx context.Context, fromType, toType, currency, pair, clientOrderID string, amount float64) (*TransferResp, error) {
	if fromType == "" {
		return nil, errFromTypeEmpty
	}
	if toType == "" {
		return nil, errToTypeEmpty
	}
	if currency == "" && pair == "" {
		return nil, errCurrencyAndPairEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"fromType": fromType,
		"toType":   toType,
		"amount":   strconv.FormatFloat(amount, 'f', -1, 64),
		"coin":     currency,
		"symbol":   pair,
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	path := bitgetSpot + bitgetWallet + bitgetTransfer
	var resp *TransferResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetTransferableCoinList returns a list of coins that can be transferred between the provided accounts
func (bi *Bitget) GetTransferableCoinList(ctx context.Context, fromType, toType string) (*TransferableCoinsResp, error) {
	if fromType == "" {
		return nil, errFromTypeEmpty
	}
	if toType == "" {
		return nil, errToTypeEmpty
	}
	vals := url.Values{}
	vals.Set("fromType", fromType)
	vals.Set("toType", toType)
	path := bitgetSpot + bitgetWallet + bitgetTransferCoinInfo
	var resp *TransferableCoinsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// SubaccountTransfer transfers assets between sub-accounts
func (bi *Bitget) SubaccountTransfer(ctx context.Context, fromType, toType, currency, pair, clientOrderID, fromID, toID string, amount float64) (*TransferResp, error) {
	if fromType == "" {
		return nil, errFromTypeEmpty
	}
	if toType == "" {
		return nil, errToTypeEmpty
	}
	if currency == "" && pair == "" {
		return nil, errCurrencyAndPairEmpty
	}
	if fromID == "" {
		return nil, errFromIDEmpty
	}
	if toID == "" {
		return nil, errToIDEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"fromType": fromType,
		"toType":   toType,
		"amount":   strconv.FormatFloat(amount, 'f', -1, 64),
		"coin":     currency,
		"symbol":   pair,
		"fromId":   fromID,
		"toId":     toID,
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	path := bitgetSpot + bitgetWallet + bitgetSubaccountTransfer
	var resp *TransferResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// WithdrawFunds withdraws funds from the user's account
func (bi *Bitget) WithdrawFunds(ctx context.Context, currency, transferType, address, chain, innerAddressType, areaCode, tag, note, clientOrderID string, amount float64) (*OrderResp, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	if transferType == "" {
		return nil, errTransferTypeEmpty
	}
	if address == "" {
		return nil, errAddressEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"coin":         currency,
		"transferType": transferType,
		"address":      address,
		"chain":        chain,
		"innerToType":  innerAddressType,
		"areaCode":     areaCode,
		"tag":          tag,
		"size":         strconv.FormatFloat(amount, 'f', -1, 64),
		"remark":       note,
		"clientOid":    clientOrderID,
	}
	path := bitgetSpot + bitgetWallet + bitgetWithdrawal
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetSubaccountTransferRecord returns the user's sub-account transfer history
func (bi *Bitget) GetSubaccountTransferRecord(ctx context.Context, currency, subaccountID, clientOrderID string, startTime, endTime time.Time, limit, pagination int64) (*SubaccTfrRecResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("coin", currency)
	params.Values.Set("subUid", subaccountID)
	if clientOrderID != "" {
		params.Values.Set("clientOid", clientOrderID)
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetSpot + bitgetAccount + bitgetSubaccountTransferRecord
	var resp *SubaccTfrRecResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetTransferRecord returns the user's transfer history
func (bi *Bitget) GetTransferRecord(ctx context.Context, currency, fromType, clientOrderID string, startTime, endTime time.Time, limit, pagination int64) (*TransferRecResp, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	if fromType == "" {
		return nil, errFromTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("coin", currency)
	params.Values.Set("fromType", fromType)
	if clientOrderID != "" {
		params.Values.Set("clientOid", clientOrderID)
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetSpot + bitgetAccount + bitgetTransferRecord
	var resp *TransferRecResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetDepositAddressForCurrency returns the user's deposit address for a particular currency
func (bi *Bitget) GetDepositAddressForCurrency(ctx context.Context, currency, chain string) (*DepositAddressResp, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	vals.Set("chain", chain)
	path := bitgetSpot + bitgetWallet + bitgetDepositAddress
	var resp *DepositAddressResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetSubaccountDepositAddress returns the deposit address for a particular currency and sub-account
func (bi *Bitget) GetSubaccountDepositAddress(ctx context.Context, subaccountID, currency, chain string) (*DepositAddressResp, error) {
	if subaccountID == "" {
		return nil, errSubaccountEmpty
	}
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("subUid", subaccountID)
	vals.Set("coin", currency)
	vals.Set("chain", chain)
	path := bitgetSpot + bitgetWallet + bitgetSubaccountDepositAddress
	var resp *DepositAddressResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// CancelWithdrawal cancels a large withdrawal request that was placed in the last minute
func (bi *Bitget) CancelWithdrawal(ctx context.Context, orderID string) (*SuccessBool, error) {
	if orderID == "" {
		return nil, errOrderIDEmpty
	}
	req := map[string]interface{}{
		"orderId": orderID,
	}
	path := bitgetSpot + bitgetWallet + bitgetCancelWithdrawal
	var resp *SuccessBool
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetSubaccountDepositRecords returns the deposit history for a sub-account
func (bi *Bitget) GetSubaccountDepositRecords(ctx context.Context, subaccountID, currency string, orderID, pagination, limit int64, startTime, endTime time.Time) (*SubaccDepRecResp, error) {
	if subaccountID == "" {
		return nil, errSubaccountEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("subUid", subaccountID)
	params.Values.Set("coin", currency)
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetSpot + bitgetWallet + bitgetSubaccountDepositRecord
	var resp *SubaccDepRecResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetWithdrawalRecords returns the user's withdrawal history
func (bi *Bitget) GetWithdrawalRecords(ctx context.Context, currency, clientOrderID string, startTime, endTime time.Time, pagination, orderID, limit int64) (*WithdrawRecordsResp, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("coin", currency)
	params.Values.Set("clientOid", clientOrderID)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetSpot + bitgetWallet + bitgetWithdrawalRecord
	var resp *WithdrawRecordsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetDepositRecords returns the user's cryptocurrency deposit history
func (bi *Bitget) GetDepositRecords(ctx context.Context, crypto string, orderID, pagination, limit int64, startTime, endTime time.Time) (*CryptoDepRecResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	params.Values.Set("coin", crypto)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetSpot + bitgetWallet + bitgetDepositRecord
	var resp *CryptoDepRecResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetFuturesVIPFeeRate returns the different levels of VIP fee rates for futures trading
func (bi *Bitget) GetFuturesVIPFeeRate(ctx context.Context) (*VIPFeeRateResp, error) {
	path := bitgetMix + bitgetMarket + bitgetVIPFeeRate
	var resp *VIPFeeRateResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, nil, &resp)
}

// GetFuturesMergeDepth returns part of the orderbook, with options to merge orders of similar price levels together,
// and to change how many results are returned. Limit's a string instead of the typical int64 because the API
// will accept a value of "max"
func (bi *Bitget) GetFuturesMergeDepth(ctx context.Context, pair, productType, precision, limit string) (*DepthResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	vals.Set("precision", precision)
	vals.Set("limit", limit)
	path := bitgetMix + bitgetMarket + bitgetMergeDepth
	var resp *DepthResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetFuturesTicker returns the ticker information for a pair of the user's choice
func (bi *Bitget) GetFuturesTicker(ctx context.Context, pair, productType string) (*FutureTickerResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	path := bitgetMix + bitgetMarket + bitgetTicker
	var resp *FutureTickerResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetAllFuturesTickers returns the ticker information for all pairs
func (bi *Bitget) GetAllFuturesTickers(ctx context.Context, productType string) (*FutureTickerResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("productType", productType)
	path := bitgetMix + bitgetMarket + bitgetTickers
	var resp *FutureTickerResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetRecentFuturesFills returns the most recent trades for a given pair
func (bi *Bitget) GetRecentFuturesFills(ctx context.Context, pair, productType string, limit int64) (*MarketFillsResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	if limit != 0 {
		vals.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetMarket + bitgetFills
	var resp *MarketFillsResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetFuturesMarketTrades returns trades for a given pair within a particular time range, and/or before a certain ID
func (bi *Bitget) GetFuturesMarketTrades(ctx context.Context, pair, productType string, limit, pagination int64, startTime, endTime time.Time) (*MarketFillsResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("productType", productType)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMix + bitgetMarket + bitgetFillsHistory
	var resp *MarketFillsResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, params.Values, &resp)
}

// GetFuturesCandlestickData returns candlestick data for a given pair within a particular time range
func (bi *Bitget) GetFuturesCandlestickData(ctx context.Context, pair, productType, granularity string, startTime, endTime time.Time, limit uint16, mode CallMode) (*CandleData, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if granularity == "" {
		return nil, errGranEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("productType", productType)
	path := bitgetMix + bitgetMarket
	switch mode {
	case CallModeNormal:
		path += bitgetCandles
	case CallModeHistory:
		path += bitgetHistoryCandles
	case CallModeIndex:
		path += bitgetHistoryIndexCandles
	case CallModeMark:
		path += bitgetHistoryMarkCandles
	}
	return bi.candlestickHelper(ctx, pair, granularity, path, limit, params)
}

// GetOpenPositions returns the total positions of a particular pair
func (bi *Bitget) GetOpenPositions(ctx context.Context, pair, productType string) (*OpenPositionsResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	path := bitgetMix + bitgetMarket + bitgetOpenInterest
	var resp *OpenPositionsResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetNextFundingTime returns the settlement time and period of a particular contract
func (bi *Bitget) GetNextFundingTime(ctx context.Context, pair, productType string) (*FundingTimeResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	path := bitgetMix + bitgetMarket + bitgetFundingTime
	var resp *FundingTimeResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetFuturesPrices returns the current market, index, and mark prices for a given pair
func (bi *Bitget) GetFuturesPrices(ctx context.Context, pair, productType string) (*FuturesPriceResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	path := bitgetMix + bitgetMarket + bitgetSymbolPrice
	var resp *FuturesPriceResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetFundingHistorical returns the historical funding rates for a given pair
func (bi *Bitget) GetFundingHistorical(ctx context.Context, pair, productType string, limit, pagination int64) (*FundingHistoryResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	if limit != 0 {
		vals.Set("pageSize", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		vals.Set("pageNo", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMix + bitgetMarket + bitgetHistoryFundRate
	var resp *FundingHistoryResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetFundingCurrent returns the current funding rate for a given pair
func (bi *Bitget) GetFundingCurrent(ctx context.Context, pair, productType string) (*FundingCurrentResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	path := bitgetMix + bitgetMarket + bitgetCurrentFundRate
	var resp *FundingCurrentResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetContractConfig returns details for a given contract
func (bi *Bitget) GetContractConfig(ctx context.Context, pair, productType string) (*ContractConfigResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	path := bitgetMix + bitgetMarket + bitgetContracts
	var resp *ContractConfigResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, vals, &resp)
}

// GetOneFuturesAccount returns details for the account associated with a given pair, margin coin, and product type
func (bi *Bitget) GetOneFuturesAccount(ctx context.Context, pair, productType, marginCoin string) (*OneAccResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	vals.Set("marginCoin", marginCoin)
	path := bitgetMix + bitgetAccount + "/" + bitgetAccount
	var resp *OneAccResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetAllFuturesAccounts returns details for all accounts
func (bi *Bitget) GetAllFuturesAccounts(ctx context.Context, productType string) (*AllAccResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("productType", productType)
	path := bitgetMix + bitgetAccount + bitgetAccounts
	var resp *AllAccResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetFuturesSubaccountAssets returns details on the assets of all sub-accounts
func (bi *Bitget) GetFuturesSubaccountAssets(ctx context.Context, productType string) (*SubaccountFuturesResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("productType", productType)
	path := bitgetMix + bitgetAccount + bitgetSubaccountAssets2
	var resp *SubaccountFuturesResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetEstimatedOpenCount returns the estimated size of open orders for a given pair
func (bi *Bitget) GetEstimatedOpenCount(ctx context.Context, pair, productType, marginCoin string, openAmount, openPrice, leverage float64) (*EstOpenCountResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if openAmount == 0 {
		return nil, errOpenAmountEmpty
	}
	if openPrice == 0 {
		return nil, errOpenPriceEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	vals.Set("marginCoin", marginCoin)
	vals.Set("openAmount", strconv.FormatFloat(openAmount, 'f', -1, 64))
	vals.Set("openPrice", strconv.FormatFloat(openPrice, 'f', -1, 64))
	if leverage != 0 {
		vals.Set("leverage", strconv.FormatFloat(leverage, 'f', -1, 64))
	}
	path := bitgetMix + bitgetAccount + bitgetOpenCount
	var resp *EstOpenCountResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// ChangeLeverage changes the leverage for the given pair and product type
func (bi *Bitget) ChangeLeverage(ctx context.Context, pair, productType, marginCoin, holdSide string, leverage float64) (*LeverageResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if leverage == 0 {
		return nil, errLeverageEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"productType": productType,
		"marginCoin":  marginCoin,
		"holdSide":    holdSide,
		"leverage":    strconv.FormatFloat(leverage, 'f', -1, 64),
	}
	path := bitgetMix + bitgetAccount + bitgetSetLeverage
	var resp *LeverageResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path, nil, req,
		&resp)
}

// AdjustMargin adds or subtracts margin from a position
func (bi *Bitget) AdjustMargin(ctx context.Context, pair, productType, marginCoin, holdSide string, amount float64) error {
	if pair == "" {
		return errPairEmpty
	}
	if productType == "" {
		return errProductTypeEmpty
	}
	if marginCoin == "" {
		return errMarginCoinEmpty
	}
	if amount == 0 {
		return errAmountEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"productType": productType,
		"marginCoin":  marginCoin,
		"amount":      strconv.FormatFloat(amount, 'f', -1, 64),
	}
	if holdSide != "" {
		req["holdSide"] = holdSide
	}
	path := bitgetMix + bitgetAccount + bitgetSetMargin
	return bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path, nil, req, nil)
}

// ChangeMarginMode changes the margin mode for a given pair. Can only be done when there the user has no open
// positions or orders
func (bi *Bitget) ChangeMarginMode(ctx context.Context, pair, productType, marginCoin, marginMode string) (*LeverageResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if marginMode == "" {
		return nil, errMarginModeEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"productType": productType,
		"marginCoin":  marginCoin,
		"marginMode":  marginMode,
	}
	path := bitgetMix + bitgetAccount + bitgetSetMarginMode
	var resp *LeverageResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path, nil, req,
		&resp)
}

// ChangePositionMode changes the position mode for any pair. Having any positions or orders on any side of any pair
// may cause this to fail.
func (bi *Bitget) ChangePositionMode(ctx context.Context, productType, positionMode string) (*PosModeResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if positionMode == "" {
		return nil, errPositionModeEmpty
	}
	req := map[string]interface{}{
		"productType": productType,
		"posMode":     positionMode,
	}
	path := bitgetMix + bitgetAccount + bitgetSetPositionMode
	var resp *PosModeResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodPost, path, nil, req,
		&resp)
}

// GetFuturesAccountBills returns a section of the user's billing history
func (bi *Bitget) GetFuturesAccountBills(ctx context.Context, productType, pair, currency, businessType string, pagination, limit int64, startTime, endTime time.Time) (*FutureAccBillResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("productType", productType)
	params.Values.Set("symbol", pair)
	params.Values.Set("coin", currency)
	params.Values.Set("businessType", businessType)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetAccount + bitgetBill
	var resp *FutureAccBillResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetPositionTier returns the position configuration for a given pair
func (bi *Bitget) GetPositionTier(ctx context.Context, productType, pair string) (*PositionTierResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if pair == "" {
		return nil, errPairEmpty
	}
	vals := url.Values{}
	vals.Set("productType", productType)
	vals.Set("symbol", pair)
	path := bitgetMix + bitgetMarket + bitgetQueryPositionLever
	var resp *PositionTierResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, vals, &resp)
}

// GetSinglePosition returns position details for a given productType, pair, and marginCoin. The exchange recommends
// using the websocket feed  instead, as information from this endpoint may be delayed during settlement or market
// fluctuations
func (bi *Bitget) GetSinglePosition(ctx context.Context, productType, pair, marginCoin string) (*PositionResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if pair == "" {
		return nil, errPairEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	vals := url.Values{}
	vals.Set("productType", productType)
	vals.Set("symbol", pair)
	vals.Set("marginCoin", marginCoin)
	path := bitgetMix + bitgetPosition + bitgetSinglePosition
	var resp *PositionResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetAllPositions returns position details for a given productType and marginCoin. The exchange recommends using
// the websocket feed  instead, as information from this endpoint may be delayed during settlement or market
// fluctuations
func (bi *Bitget) GetAllPositions(ctx context.Context, productType, marginCoin string) (*PositionResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	vals := url.Values{}
	vals.Set("productType", productType)
	vals.Set("marginCoin", marginCoin)
	path := bitgetMix + bitgetPosition + bitgetAllPositions
	var resp *PositionResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodGet, path, vals, nil,
		&resp)
}

// GetHistoricalPositions returns historical position details, up to a maximum of three months ago
func (bi *Bitget) GetHistoricalPositions(ctx context.Context, pair, productType string, pagination, limit int64, startTime, endTime time.Time) (*HistPositionResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("productType", productType)
	fmt.Printf("pagination %v\n", pagination)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetPosition + bitgetHistoryPosition
	var resp *HistPositionResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodGet, path, params.Values,
		nil, &resp)
}

// PlaceFuturesOrder places a futures order on the exchange
func (bi *Bitget) PlaceFuturesOrder(ctx context.Context, pair, productType, marginMode, marginCoin, side, tradeSide, orderType, strategy, clientOID, stopSurplusPrice, stopLossPrice string, amount, price float64, reduceOnly, isCopyTradeLeader bool) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginMode == "" {
		return nil, errMarginModeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if side == "" {
		return nil, errSideEmpty
	}
	if orderType == "" {
		return nil, errOrderTypeEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	if orderType == "limit" && price == 0 {
		return nil, errLimitPriceEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"productType": productType,
		"marginMode":  marginMode,
		"marginCoin":  marginCoin,
		"side":        side,
		"tradeSide":   tradeSide,
		"orderType":   orderType,
		"force":       strategy,
		"size":        strconv.FormatFloat(amount, 'f', -1, 64),
		"price":       strconv.FormatFloat(price, 'f', -1, 64),
	}
	if clientOID != "" {
		req["clientOid"] = clientOID
	}
	if reduceOnly {
		req["reduceOnly"] = "YES"
	}
	req["presetStopSurplusPrice"] = stopSurplusPrice
	req["presetStopLossPrice"] = stopLossPrice
	path := bitgetMix + bitgetOrder + bitgetPlaceOrder
	rLim := Rate10
	if isCopyTradeLeader {
		rLim = Rate1
	}
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, rLim, http.MethodPost, path, nil, req,
		&resp)
}

// PlaceReversal attempts to close a position, in part or in whole, and opens a position of corresponding size
// on the opposite side. This operation may only be done in part under certain margin levels, market conditions,
// or other unspecified factors. If a reversal is attempted for an amount greater than the current outstanding position,
// that position will be closed, and a new position will be opened for the amount of the closed position; not the amount
// specified in the request. The side specified in the parameter should correspond to the side of the position you're
// attempting to close; if the original is open_long, use close_long; if the original is open_short, use close_short;
// if the original is sell_single, use buy_single. If the position is sell_single or buy_single, the amount parameter
// will be ignored, and the entire position will be closed, with a corresponding amount opened on the opposite side.
func (bi *Bitget) PlaceReversal(ctx context.Context, pair, marginCoin, productType, side, tradeSide, clientOID string, amount float64, isCopyTradeLeader bool) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"marginCoin":  marginCoin,
		"productType": productType,
		"side":        side,
		"tradeSide":   tradeSide,
		"size":        strconv.FormatFloat(amount, 'f', -1, 64),
		"orderType":   "market",
	}
	if clientOID != "" {
		req["clientOid"] = clientOID
	}
	path := bitgetMix + bitgetOrder + bitgetClickBackhand
	rLim := Rate10
	if isCopyTradeLeader {
		rLim = Rate1
	}
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, rLim, http.MethodPost, path, nil, req,
		&resp)
}

// BatchPlaceFuturesOrders places multiple orders at once. Can also be used to modify the take-profit and stop-loss
// of an open position.
func (bi *Bitget) BatchPlaceFuturesOrders(ctx context.Context, pair, productType, marginCoin, marginMode string, orders []PlaceFuturesOrderStruct, isCopyTradeLeader bool) (*BatchOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if marginMode == "" {
		return nil, errMarginModeEmpty
	}
	if len(orders) == 0 {
		return nil, errOrdersEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"productType": productType,
		"marginCoin":  marginCoin,
		"marginMode":  marginMode,
		"orderList":   orders,
	}
	path := bitgetMix + bitgetOrder + bitgetBatchPlaceOrder
	rLim := Rate10
	if isCopyTradeLeader {
		rLim = Rate1
	}
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, rLim, http.MethodPost, path, nil, req,
		&resp)
}

// ModifyFuturesOrder can change the size, price, take-profit, and stop-loss of an order. Size and price have to be
// modified at the same time, or the request will fail. If size and price are altered, the old order will be cancelled,
// and a new one will be created asynchronously. Due to the asynchronous creation of a new order, a new ClientOrderID
// must be supplied so it can be tracked.
func (bi *Bitget) ModifyFuturesOrder(ctx context.Context, orderID int64, clientOrderID, pair, productType, newClientOrderID string, newAmount, newPrice, newTakeProfit, newStopLoss float64) (*OrderResp, error) {
	if orderID == 0 && clientOrderID == "" {
		return nil, errOrderClientEmpty
	}
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if newClientOrderID == "" {
		return nil, errNewClientOrderIDEmpty
	}
	req := map[string]interface{}{
		"symbol":                    pair,
		"productType":               productType,
		"newClientOid":              newClientOrderID,
		"newPresetStopSurplusPrice": strconv.FormatFloat(newTakeProfit, 'f', -1, 64),
		"newPresetStopLossPrice":    strconv.FormatFloat(newStopLoss, 'f', -1, 64),
	}
	if orderID != 0 {
		req["orderId"] = orderID
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	if newAmount != 0 {
		req["newSize"] = strconv.FormatFloat(newAmount, 'f', -1, 64)
	}
	if newPrice != 0 {
		req["newPrice"] = strconv.FormatFloat(newPrice, 'f', -1, 64)
	}
	path := bitgetMix + bitgetOrder + bitgetModifyOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// CancelFuturesOrder cancels an order on the exchange
func (bi *Bitget) CancelFuturesOrder(ctx context.Context, pair, productType, marginCoin, clientOrderID string, orderID int64) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if clientOrderID == "" && orderID == 0 {
		return nil, errOrderClientEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"productType": productType,
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	if orderID != 0 {
		req["orderId"] = orderID
	}
	if marginCoin != "" {
		req["marginCoin"] = marginCoin
	}
	path := bitgetMix + bitgetOrder + bitgetCancelOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// BatchCancelFuturesOrders cancels multiple orders at once
func (bi *Bitget) BatchCancelFuturesOrders(ctx context.Context, orderIDs []OrderIDStruct, pair, productType, marginCoin string) (*BatchOrderResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"productType": productType,
		"orderList":   orderIDs,
	}
	if marginCoin != "" {
		req["marginCoin"] = marginCoin
	}
	path := bitgetMix + bitgetOrder + bitgetBatchCancelOrders
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// FlashClosePosition attempts to close a position at the best available price
func (bi *Bitget) FlashClosePosition(ctx context.Context, pair, holdSide, productType string) (*BatchOrderResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	req := map[string]interface{}{
		"symbol":      pair,
		"holdSide":    holdSide,
		"productType": productType,
	}
	path := bitgetMix + bitgetOrder + bitgetClosePositions
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate1, http.MethodPost, path, nil, req,
		&resp)
}

// GetFuturesOrderDetails returns details on a given order
func (bi *Bitget) GetFuturesOrderDetails(ctx context.Context, pair, productType, clientOrderID string, orderID int64) (*FuturesOrderDetailResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if clientOrderID == "" && orderID == 0 {
		return nil, errOrderClientEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	vals.Set("productType", productType)
	if clientOrderID != "" {
		vals.Set("clientOid", clientOrderID)
	}
	if orderID != 0 {
		vals.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	path := bitgetMix + bitgetOrder + bitgetDetail
	var resp *FuturesOrderDetailResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetFuturesFills returns fill details
func (bi *Bitget) GetFuturesFills(ctx context.Context, orderID, pagination, limit int64, pair, productType string, startTime, endTime time.Time) (*FuturesFillsResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("productType", productType)
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetOrder + "/" + bitgetFills
	var resp *FuturesFillsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetFuturesOrderFillHistory returns historical fill details
func (bi *Bitget) GetFuturesOrderFillHistory(ctx context.Context, pair, productType string, orderID, pagination, limit int64, startTime, endTime time.Time) (*FuturesFillsResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("productType", productType)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetOrder + bitgetFillHistory
	var resp *FuturesFillsResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetPendingFuturesOrders returns detailed information on pending futures orders
func (bi *Bitget) GetPendingFuturesOrders(ctx context.Context, orderID, pagination, limit int64, clientOrderID, pair, productType, status string, startTime, endTime time.Time) (*FuturesOrdResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("productType", productType)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	params.Values.Set("clientOid", clientOrderID)
	params.Values.Set("status", status)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetOrder + bitgetOrdersPending
	var resp *FuturesOrdResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetHistoricalFuturesOrders returns information on futures orders that are no longer pending
func (bi *Bitget) GetHistoricalFuturesOrders(ctx context.Context, orderID, pagination, limit int64, clientOrderID, pair, productType string, startTime, endTime time.Time) (*FuturesOrdResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("productType", productType)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	params.Values.Set("clientOid", clientOrderID)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetOrder + bitgetOrdersHistory
	var resp *FuturesOrdResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// CancelAllFuturesOrders cancels all pending orders
func (bi *Bitget) CancelAllFuturesOrders(ctx context.Context, pair, productType, marginCoin string, acceptableDelay time.Duration) (*BatchOrderResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	req := map[string]interface{}{
		"productType":   productType,
		"symbol":        pair,
		"requestTime":   time.Now().UnixMilli(),
		"receiveWindow": time.Unix(0, 0).Add(acceptableDelay).UnixMilli(),
	}
	if marginCoin != "" {
		req["marginCoin"] = marginCoin
	}
	path := bitgetMix + bitgetOrder + bitgetCancelAllOrders
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil,
		req, &resp)
}

// GetFuturesTriggerOrderByID returns information on a particular trigger order
func (bi *Bitget) GetFuturesTriggerOrderByID(ctx context.Context, planType, planOrderID, productType string) (*SubOrderResp, error) {
	if planType == "" {
		return nil, errPlanTypeEmpty
	}
	if planOrderID == "" {
		return nil, errPlanOrderIDEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	vals := url.Values{}
	vals.Set("planType", planType)
	vals.Set("planOrderId", planOrderID)
	vals.Set("productType", productType)
	path := bitgetMix + bitgetOrder + bitgetPlanSubOrder
	var resp *SubOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// PlaceTPSLFuturesOrder places a take-profit or stop-loss futures order
func (bi *Bitget) PlaceTPSLFuturesOrder(ctx context.Context, marginCoin, productType, pair, planType, triggerType, holdSide, rangeRate, clientOrderID string, triggerPrice, executePrice, amount float64) (*OrderResp, error) {
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if pair == "" {
		return nil, errPairEmpty
	}
	if planType == "" {
		return nil, errPlanTypeEmpty
	}
	if holdSide == "" {
		return nil, errHoldSideEmpty
	}
	if triggerPrice == 0 {
		return nil, errTriggerPriceEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"marginCoin":   marginCoin,
		"productType":  productType,
		"symbol":       pair,
		"planType":     planType,
		"triggerType":  triggerType,
		"holdSide":     holdSide,
		"rangeRate":    rangeRate,
		"clientOid":    clientOrderID,
		"triggerPrice": strconv.FormatFloat(triggerPrice, 'f', -1, 64),
		"executePrice": strconv.FormatFloat(executePrice, 'f', -1, 64),
		"size":         strconv.FormatFloat(amount, 'f', -1, 64),
	}
	path := bitgetMix + bitgetOrder + bitgetPlaceTPSLOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// PlaceTriggerFuturesOrder places a trigger futures order
func (bi *Bitget) PlaceTriggerFuturesOrder(ctx context.Context, planType, pair, productType, marginMode, marginCoin, triggerType, side, tradeSide, orderType, clientOrderID, takeProfitTriggerType, stopLossTriggerType string, amount, executePrice, callbackRatio, triggerPrice, takeProfitTriggerPrice, takeProfitExecutePrice, stopLossTriggerPrice, stopLossExecutePrice float64, reduceOnly bool) (*OrderResp, error) {
	if planType == "" {
		return nil, errPlanTypeEmpty
	}
	if pair == "" {
		return nil, errPairEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if marginMode == "" {
		return nil, errMarginModeEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if triggerType == "" {
		return nil, errTriggerTypeEmpty
	}
	if side == "" {
		return nil, errSideEmpty
	}
	if orderType == "" {
		return nil, errOrderTypeEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	if executePrice == 0 {
		return nil, errExecutePriceEmpty
	}
	if triggerPrice == 0 {
		return nil, errTriggerPriceEmpty
	}
	if (takeProfitTriggerType != "" || takeProfitTriggerPrice != 0 || takeProfitExecutePrice != 0) &&
		(takeProfitTriggerType == "" || takeProfitTriggerPrice == 0 || takeProfitExecutePrice == 0) {
		return nil, errTakeProfitParamsInconsistency
	}
	if (stopLossTriggerType != "" || stopLossTriggerPrice != 0 || stopLossExecutePrice != 0) &&
		(stopLossTriggerType == "" || stopLossTriggerPrice == 0 || stopLossExecutePrice == 0) {
		return nil, errStopLossParamsInconsistency
	}
	req := map[string]interface{}{
		"planType":      planType,
		"symbol":        pair,
		"productType":   productType,
		"marginMode":    marginMode,
		"marginCoin":    marginCoin,
		"triggerType":   triggerType,
		"side":          side,
		"tradeSide":     tradeSide,
		"orderType":     orderType,
		"size":          strconv.FormatFloat(amount, 'f', -1, 64),
		"price":         strconv.FormatFloat(executePrice, 'f', -1, 64),
		"triggerPrice":  strconv.FormatFloat(triggerPrice, 'f', -1, 64),
		"callbackRatio": strconv.FormatFloat(callbackRatio, 'f', -1, 64),
	}
	if reduceOnly {
		req["reduceOnly"] = "YES"
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	if takeProfitTriggerPrice != 0 || takeProfitExecutePrice != 0 || takeProfitTriggerType != "" {
		req["stopSurplusTriggerPrice"] = strconv.FormatFloat(takeProfitTriggerPrice, 'f', -1, 64)
		req["stopSurplusExecutePrice"] = strconv.FormatFloat(takeProfitExecutePrice, 'f', -1, 64)
		req["stopSurplusTriggerType"] = takeProfitTriggerType
	}
	if stopLossTriggerPrice != 0 || stopLossExecutePrice != 0 || stopLossTriggerType != "" {
		req["stopLossTriggerPrice"] = strconv.FormatFloat(stopLossTriggerPrice, 'f', -1, 64)
		req["stopLossExecutePrice"] = strconv.FormatFloat(stopLossExecutePrice, 'f', -1, 64)
		req["stopLossTriggerType"] = stopLossTriggerType
	}
	path := bitgetMix + bitgetOrder + bitgetPlacePlanOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// ModifyTPSLFuturesOrder modifies a take-profit or stop-loss futures order
func (bi *Bitget) ModifyTPSLFuturesOrder(ctx context.Context, orderID int64, clientOrderID, marginCoin, productType, pair, triggerType string, triggerPrice, executePrice, amount, rangeRate float64) (*OrderResp, error) {
	if orderID == 0 && clientOrderID == "" {
		return nil, errOrderClientEmpty
	}
	if marginCoin == "" {
		return nil, errMarginCoinEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	if pair == "" {
		return nil, errPairEmpty
	}
	if triggerPrice == 0 {
		return nil, errTriggerPriceEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"orderId":      orderID,
		"clientOid":    clientOrderID,
		"marginCoin":   marginCoin,
		"productType":  productType,
		"symbol":       pair,
		"triggerType":  triggerType,
		"triggerPrice": strconv.FormatFloat(triggerPrice, 'f', -1, 64),
		"executePrice": strconv.FormatFloat(executePrice, 'f', -1, 64),
		"size":         strconv.FormatFloat(amount, 'f', -1, 64),
	}
	if rangeRate != 0 {
		req["rangeRate"] = strconv.FormatFloat(rangeRate, 'f', -1, 64)
	}
	path := bitgetMix + bitgetOrder + bitgetModifyTPSLOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// CancelTriggerFuturesOrders cancels trigger futures orders
func (bi *Bitget) CancelTriggerFuturesOrders(ctx context.Context, orderIDList []OrderIDStruct, pair, productType, marginCoin, planType string) (*BatchOrderResp, error) {
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	req := map[string]interface{}{
		"productType": productType,
		"planType":    planType,
		"symbol":      pair,
		"orderIdList": orderIDList,
		"marginCoin":  marginCoin,
	}
	path := bitgetMix + bitgetOrder + bitgetCancelPlanOrder
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// ModifyTriggerFuturesOrder modifies a trigger futures order
func (bi *Bitget) ModifyTriggerFuturesOrder(ctx context.Context, orderID int64, clientOrderID, productType, triggerType, takeProfitTriggerType, stopLossTriggerType string, amount, executePrice, callbackRatio, triggerPrice, takeProfitTriggerPrice, takeProfitExecutePrice, stopLossTriggerPrice, stopLossExecutePrice float64) (*OrderResp, error) {
	if orderID == 0 && clientOrderID == "" {
		return nil, errOrderClientEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	req := map[string]interface{}{
		// See whether planType is accepted
		// See whether symbol is accepted
		"orderId":                   orderID,
		"clientOid":                 clientOrderID,
		"productType":               productType,
		"newTriggerType":            triggerType,
		"newSize":                   strconv.FormatFloat(amount, 'f', -1, 64),
		"newPrice":                  strconv.FormatFloat(executePrice, 'f', -1, 64),
		"newCallbackRatio":          strconv.FormatFloat(callbackRatio, 'f', -1, 64),
		"newTriggerPrice":           strconv.FormatFloat(triggerPrice, 'f', -1, 64),
		"newStopSurplusTriggerType": takeProfitTriggerType,
		"newStopLossTriggerType":    stopLossTriggerType,
	}
	if takeProfitTriggerPrice >= 0 {
		req["newStopSurplusTriggerPrice"] = strconv.FormatFloat(takeProfitTriggerPrice, 'f', -1, 64)
	}
	if takeProfitExecutePrice >= 0 {
		req["newStopSurplusExecutePrice"] = strconv.FormatFloat(takeProfitExecutePrice, 'f', -1, 64)
	}
	if stopLossTriggerPrice >= 0 {
		req["newStopLossTriggerPrice"] = strconv.FormatFloat(stopLossTriggerPrice, 'f', -1, 64)
	}
	if stopLossExecutePrice >= 0 {
		req["newStopLossExecutePrice"] = strconv.FormatFloat(stopLossExecutePrice, 'f', -1, 64)
	}
	path := bitgetMix + bitgetOrder + bitgetModifyPlanOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetPendingTriggerFuturesOrders returns information on pending trigger orders
func (bi *Bitget) GetPendingTriggerFuturesOrders(ctx context.Context, orderID, pagination, limit int64, clientOrderID, pair, planType, productType string, startTime, endTime time.Time) (*PlanFuturesOrdResp, error) {
	if planType == "" {
		return nil, errPlanTypeEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("planType", planType)
	params.Values.Set("productType", productType)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	params.Values.Set("clientOid", clientOrderID)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetOrder + bitgetOrdersPlanPending
	var resp *PlanFuturesOrdResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetHistoricalTriggerFuturesOrders returns information on historical trigger orders
func (bi *Bitget) GetHistoricalTriggerFuturesOrders(ctx context.Context, orderID, pagination, limit int64, clientOrderID, planType, planStatus, pair, productType string, startTime, endTime time.Time) (*HistTriggerFuturesOrdResp, error) {
	if planType == "" {
		return nil, errPlanTypeEmpty
	}
	if productType == "" {
		return nil, errProductTypeEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("planType", planType)
	params.Values.Set("productType", productType)
	params.Values.Set("planStatus", planStatus)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	params.Values.Set("clientOid", clientOrderID)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMix + bitgetOrder + bitgetOrdersPlanHistory
	var resp *HistTriggerFuturesOrdResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetSupportedCurrencies returns information on the currencies supported by the exchange
func (bi *Bitget) GetSupportedCurrencies(ctx context.Context) (*SupCurrencyResp, error) {
	path := bitgetMargin + bitgetCurrencies
	var resp *SupCurrencyResp
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, nil, &resp)
}

// GetCrossBorrowHistory returns the borrowing history for cross margin
func (bi *Bitget) GetCrossBorrowHistory(ctx context.Context, loanID, limit, pagination int64, currency string, startTime, endTime time.Time) (*BorrowHistCross, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if loanID != 0 {
		params.Values.Set("loanId", strconv.FormatInt(loanID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + bitgetBorrowHistory
	var resp *BorrowHistCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossRepayHistory returns the repayment history for cross margin
func (bi *Bitget) GetCrossRepayHistory(ctx context.Context, repayID, limit, pagination int64, currency string, startTime, endTime time.Time) (*RepayHistResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if repayID != 0 {
		params.Values.Set("repayId", strconv.FormatInt(repayID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + bitgetRepayHistory
	var resp *RepayHistResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossInterestHistory returns the interest history for cross margin
func (bi *Bitget) GetCrossInterestHistory(ctx context.Context, currency string, startTime, endTime time.Time, limit, pagination int64) (*InterHistCross, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + bitgetInterestHistory
	var resp *InterHistCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossLiquidationHistory returns the liquidation history for cross margin
func (bi *Bitget) GetCrossLiquidationHistory(ctx context.Context, startTime, endTime time.Time, limit, pagination int64) (*LiquidHistCross, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetCrossed + bitgetLiquidationHistory
	var resp *LiquidHistCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossFinancialHistory returns the financial history for cross margin
func (bi *Bitget) GetCrossFinancialHistory(ctx context.Context, marginType, currency string, startTime, endTime time.Time, limit, pagination int64) (*FinHistCross, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("marginType", marginType)
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + bitgetFinancialRecords
	var resp *FinHistCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossAccountAssets returns the account assets for cross margin
func (bi *Bitget) GetCrossAccountAssets(ctx context.Context, currency string) (*CrossAssetResp, error) {
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetAssets
	var resp *CrossAssetResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// CrossBorrow borrows funds for cross margin
func (bi *Bitget) CrossBorrow(ctx context.Context, currency, clientOrderID string, amount float64) (*BorrowCross, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"coin":         currency,
		"borrowAmount": strconv.FormatFloat(amount, 'f', -1, 64),
		"clientOid":    clientOrderID,
	}
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetBorrow
	var resp *BorrowCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// CrossRepay repays funds for cross margin
func (bi *Bitget) CrossRepay(ctx context.Context, currency string, amount float64) (*RepayCross, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"coin":        currency,
		"repayAmount": strconv.FormatFloat(amount, 'f', -1, 64),
	}
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetRepay
	var resp *RepayCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetCrossRiskRate returns the risk rate for cross margin
func (bi *Bitget) GetCrossRiskRate(ctx context.Context) (*RiskRateCross, error) {
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetRiskRate
	var resp *RiskRateCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, nil, nil,
		&resp)
}

// GetCrossMaxBorrowable returns the maximum amount that can be borrowed for cross margin
func (bi *Bitget) GetCrossMaxBorrowable(ctx context.Context, currency string) (*MaxBorrowCross, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetMaxBorrowableAmount
	var resp *MaxBorrowCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetCrossMaxTransferable returns the maximum amount that can be transferred out of cross margin
func (bi *Bitget) GetCrossMaxTransferable(ctx context.Context, currency string) (*MaxTransferCross, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetMaxTransferOutAmount
	var resp *MaxTransferCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetCrossInterestRateAndMaxBorrowable returns the interest rate and maximum borrowable amount for cross margin
func (bi *Bitget) GetCrossInterestRateAndMaxBorrowable(ctx context.Context, currency string) (*IntRateMaxBorrowCross, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + bitgetInterestRateAndLimit
	var resp *IntRateMaxBorrowCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetCrossTierConfiguration returns tier information for the user's VIP level
func (bi *Bitget) GetCrossTierConfiguration(ctx context.Context, currency string) (*TierConfigCross, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetMargin + bitgetCrossed + bitgetTierData
	var resp *TierConfigCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// CrossFlashRepay repays funds for cross margin, with the option to only repay for a particular currency
func (bi *Bitget) CrossFlashRepay(ctx context.Context, currency string) (*FlashRepayCross, error) {
	req := map[string]interface{}{
		"coin": currency,
	}
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetFlashRepay
	var resp *FlashRepayCross
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetCrossFlashRepayResult returns the result of the supplied flash repayments for cross margin
func (bi *Bitget) GetCrossFlashRepayResult(ctx context.Context, idList []int64) (*FlashRepayResult, error) {
	if len(idList) == 0 {
		return nil, errIDListEmpty
	}
	req := map[string]interface{}{
		"repayIdList": idList,
	}
	path := bitgetMargin + bitgetCrossed + "/" + bitgetAccount + bitgetQueryFlashRepayStatus
	var resp *FlashRepayResult
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// PlaceCrossOrder places an order using cross margin
func (bi *Bitget) PlaceCrossOrder(ctx context.Context, pair, orderType, loanType, strategy, clientOrderID, side string, price, baseAmount, quoteAmount float64) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if orderType == "" {
		return nil, errOrderTypeEmpty
	}
	if loanType == "" {
		return nil, errLoanTypeEmpty
	}
	if strategy == "" {
		return nil, errStrategyEmpty
	}
	if side == "" {
		return nil, errSideEmpty
	}
	if baseAmount == 0 && quoteAmount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderType": orderType,
		"loanType":  loanType,
		"force":     strategy,
		"clientOid": clientOrderID,
		"side":      side,
		"price":     strconv.FormatFloat(price, 'f', -1, 64),
		"baseSize":  strconv.FormatFloat(baseAmount, 'f', -1, 64),
		"quoteSize": strconv.FormatFloat(quoteAmount, 'f', -1, 64),
	}
	path := bitgetMargin + bitgetCrossed + bitgetPlaceOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// BatchPlaceCrossOrders places multiple orders using cross margin
func (bi *Bitget) BatchPlaceCrossOrders(ctx context.Context, pair string, orders []MarginOrderData) (*BatchOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if len(orders) == 0 {
		return nil, errOrdersEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderList": orders,
	}
	path := bitgetMargin + bitgetCrossed + bitgetBatchPlaceOrder
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// CancelCrossOrder cancels an order using cross margin
func (bi *Bitget) CancelCrossOrder(ctx context.Context, pair, clientOrderID string, orderID int64) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if clientOrderID == "" && orderID == 0 {
		return nil, errOrderClientEmpty
	}
	req := map[string]interface{}{
		"symbol": pair,
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	if orderID != 0 {
		req["orderId"] = orderID
	}
	path := bitgetMargin + bitgetCrossed + bitgetCancelOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// BatchCancelCrossOrders cancels multiple orders using cross margin
func (bi *Bitget) BatchCancelCrossOrders(ctx context.Context, pair string, orders []OrderIDStruct) (*BatchOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if len(orders) == 0 {
		return nil, errOrdersEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderList": orders,
	}
	path := bitgetMargin + bitgetCrossed + bitgetBatchCancelOrder
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetCrossOpenOrders returns the open orders for cross margin
func (bi *Bitget) GetCrossOpenOrders(ctx context.Context, pair, clientOrderID string, orderID, limit, pagination int64, startTime, endTime time.Time) (*MarginOpenOrds, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("clientOid", clientOrderID)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetCrossed + bitgetOpenOrders
	var resp *MarginOpenOrds
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossHistoricalOrders returns the historical orders for cross margin
func (bi *Bitget) GetCrossHistoricalOrders(ctx context.Context, pair, enterPointSource, clientOrderID string, orderID, limit, pagination int64, startTime, endTime time.Time) (*MarginHistOrds, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("enterPointSource", enterPointSource)
	params.Values.Set("clientOid", clientOrderID)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetCrossed + bitgetHistoryOrders
	var resp *MarginHistOrds
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossOrderFills returns the fills for cross margin orders
func (bi *Bitget) GetCrossOrderFills(ctx context.Context, pair string, orderID, pagination, limit int64, startTime, endTime time.Time) (*MarginOrderFills, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMargin + bitgetCrossed + "/" + bitgetFills
	var resp *MarginOrderFills
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetCrossLiquidationOrders returns the liquidation orders for cross margin
func (bi *Bitget) GetCrossLiquidationOrders(ctx context.Context, orderType, pair, fromCoin, toCoin string, startTime, endTime time.Time, limit, pagination int64) (*LiquidationResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("orderType", orderType)
	params.Values.Set("symbol", pair)
	params.Values.Set("fromCoin", fromCoin)
	params.Values.Set("toCoin", toCoin)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetCrossed + bitgetLiquidationOrder
	var resp *LiquidationResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedRepayHistory returns the repayment history for isolated margin
func (bi *Bitget) GetIsolatedRepayHistory(ctx context.Context, pair, currency string, repayID, limit, pagination int64, startTime, endTime time.Time) (*RepayHistResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if repayID != 0 {
		params.Values.Set("repayId", strconv.FormatInt(repayID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetIsolated + bitgetRepayHistory
	var resp *RepayHistResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedBorrowHistory returns the borrowing history for isolated margin
func (bi *Bitget) GetIsolatedBorrowHistory(ctx context.Context, pair, currency string, loanID, limit, pagination int64, startTime, endTime time.Time) (*BorrowHistIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if loanID != 0 {
		params.Values.Set("loanId", strconv.FormatInt(loanID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetIsolated + bitgetBorrowHistory
	var resp *BorrowHistIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedInterestHistory returns the interest history for isolated margin
func (bi *Bitget) GetIsolatedInterestHistory(ctx context.Context, pair, currency string, startTime, endTime time.Time, limit, pagination int64) (*InterHistIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetIsolated + bitgetInterestHistory
	var resp *InterHistIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedLiquidationHistory returns the liquidation history for isolated margin
func (bi *Bitget) GetIsolatedLiquidationHistory(ctx context.Context, pair string, startTime, endTime time.Time, limit, pagination int64) (*LiquidHistIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("symbol", pair)
	path := bitgetMargin + bitgetIsolated + bitgetLiquidationHistory
	var resp *LiquidHistIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedFinancialHistory returns the financial history for isolated margin
func (bi *Bitget) GetIsolatedFinancialHistory(ctx context.Context, pair, marginType, currency string, startTime, endTime time.Time, limit, pagination int64) (*FinHistIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("marginType", marginType)
	params.Values.Set("coin", currency)
	path := bitgetMargin + bitgetIsolated + bitgetFinancialRecords
	var resp *FinHistIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedAccountAssets returns the account assets for isolated margin
func (bi *Bitget) GetIsolatedAccountAssets(ctx context.Context, pair string) (*IsoAssetResp, error) {
	vals := url.Values{}
	vals.Set("symbol", pair)
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetAssets
	var resp *IsoAssetResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// IsolatedBorrow borrows funds for isolated margin
func (bi *Bitget) IsolatedBorrow(ctx context.Context, pair, currency, clientOrderID string, amount float64) (*BorrowIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"symbol":       pair,
		"coin":         currency,
		"borrowAmount": strconv.FormatFloat(amount, 'f', -1, 64),
		"clientOid":    clientOrderID,
	}
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetBorrow
	var resp *BorrowIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// IsolatedRepay repays funds for isolated margin
func (bi *Bitget) IsolatedRepay(ctx context.Context, amount float64, currency, pair, clientOrderID string) (*RepayIso, error) {
	if amount == 0 {
		return nil, errAmountEmpty
	}
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	if pair == "" {
		return nil, errPairEmpty
	}
	req := map[string]interface{}{
		"coin":        currency,
		"repayAmount": strconv.FormatFloat(amount, 'f', -1, 64),
		"symbol":      pair,
		"clientOid":   clientOrderID,
	}
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetRepay
	var resp *RepayIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetIsolatedRiskRate returns the risk rate for isolated margin
func (bi *Bitget) GetIsolatedRiskRate(ctx context.Context, pair string, pagination, limit int64) (*RiskRateIso, error) {
	vals := url.Values{}
	vals.Set("symbol", pair)
	if limit != 0 {
		vals.Set("pageSize", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		vals.Set("pageNum", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetRiskRate
	var resp *RiskRateIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetIsolatedInterestRateAndMaxBorrowable returns the interest rate and maximum borrowable amount for isolated margin
func (bi *Bitget) GetIsolatedInterestRateAndMaxBorrowable(ctx context.Context, pair string) (*IntRateMaxBorrowIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	path := bitgetMargin + bitgetIsolated + bitgetInterestRateAndLimit
	var resp *IntRateMaxBorrowIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetIsolatedTierConfiguration returns tier information for the user's VIP level
func (bi *Bitget) GetIsolatedTierConfiguration(ctx context.Context, pair string) (*TierConfigIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	path := bitgetMargin + bitgetIsolated + bitgetTierData
	var resp *TierConfigIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetIsolatedMaxBorrowable returns the maximum amount that can be borrowed for isolated margin
func (bi *Bitget) GetIsolatedMaxBorrowable(ctx context.Context, pair string) (*MaxBorrowIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetMaxBorrowableAmount
	var resp *MaxBorrowIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetIsolatedMaxTransferable returns the maximum amount that can be transferred out of isolated margin
func (bi *Bitget) GetIsolatedMaxTransferable(ctx context.Context, pair string) (*MaxTransferIso, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	vals := url.Values{}
	vals.Set("symbol", pair)
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetMaxTransferOutAmount
	var resp *MaxTransferIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// IsolatedFlashRepay repays funds for isolated margin, with the option to only repay for a set of up to 100 pairs
func (bi *Bitget) IsolatedFlashRepay(ctx context.Context, pairs []string) (*FlashRepayIso, error) {
	req := map[string]interface{}{
		"symbolList": pairs,
	}
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetFlashRepay
	var resp *FlashRepayIso
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetIsolatedFlashRepayResult returns the result of the supplied flash repayments for isolated margin
func (bi *Bitget) GetIsolatedFlashRepayResult(ctx context.Context, idList []int64) (*FlashRepayResult, error) {
	if len(idList) == 0 {
		return nil, errIDListEmpty
	}
	req := map[string]interface{}{
		"repayIdList": idList,
	}
	path := bitgetMargin + bitgetIsolated + "/" + bitgetAccount + bitgetQueryFlashRepayStatus
	var resp *FlashRepayResult
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// PlaceIsolatedOrder places an order using isolated margin
func (bi *Bitget) PlaceIsolatedOrder(ctx context.Context, pair, orderType, loanType, strategy, clientOrderID, side string, price, baseAmount, quoteAmount float64) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if orderType == "" {
		return nil, errOrderTypeEmpty
	}
	if loanType == "" {
		return nil, errLoanTypeEmpty
	}
	if strategy == "" {
		return nil, errStrategyEmpty
	}
	if side == "" {
		return nil, errSideEmpty
	}
	if baseAmount == 0 && quoteAmount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderType": orderType,
		"loanType":  loanType,
		"force":     strategy,
		"clientOid": clientOrderID,
		"side":      side,
		"price":     strconv.FormatFloat(price, 'f', -1, 64),
		"baseSize":  strconv.FormatFloat(baseAmount, 'f', -1, 64),
		"quoteSize": strconv.FormatFloat(quoteAmount, 'f', -1, 64),
	}
	path := bitgetMargin + bitgetIsolated + bitgetPlaceOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// BatchPlaceIsolatedOrders places multiple orders using isolated margin
func (bi *Bitget) BatchPlaceIsolatedOrders(ctx context.Context, pair string, orders []MarginOrderData) (*BatchOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if len(orders) == 0 {
		return nil, errOrdersEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderList": orders,
	}
	path := bitgetMargin + bitgetIsolated + bitgetBatchPlaceOrder
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// CancelIsolatedOrder cancels an order using isolated margin
func (bi *Bitget) CancelIsolatedOrder(ctx context.Context, pair, clientOrderID string, orderID int64) (*OrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if clientOrderID == "" && orderID == 0 {
		return nil, errOrderClientEmpty
	}
	req := map[string]interface{}{
		"symbol": pair,
	}
	if clientOrderID != "" {
		req["clientOid"] = clientOrderID
	}
	if orderID != 0 {
		req["orderId"] = orderID
	}
	path := bitgetMargin + bitgetIsolated + bitgetCancelOrder
	var resp *OrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// BatchCancelIsolatedOrders cancels multiple orders using isolated margin
func (bi *Bitget) BatchCancelIsolatedOrders(ctx context.Context, pair string, orders []OrderIDStruct) (*BatchOrderResp, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	if len(orders) == 0 {
		return nil, errOrdersEmpty
	}
	req := map[string]interface{}{
		"symbol":    pair,
		"orderList": orders,
	}
	path := bitgetMargin + bitgetIsolated + bitgetBatchCancelOrder
	var resp *BatchOrderResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetIsolatedOpenOrders returns the open orders for isolated margin
func (bi *Bitget) GetIsolatedOpenOrders(ctx context.Context, pair, clientOrderID string, orderID, limit, pagination int64, startTime, endTime time.Time) (*MarginOpenOrds, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("clientOid", clientOrderID)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetIsolated + bitgetOpenOrders
	var resp *MarginOpenOrds
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedHistoricalOrders returns the historical orders for isolated margin
func (bi *Bitget) GetIsolatedHistoricalOrders(ctx context.Context, pair, enterPointSource, clientOrderID string, orderID, limit, pagination int64, startTime, endTime time.Time) (*MarginHistOrds, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	params.Values.Set("enterPointSource", enterPointSource)
	params.Values.Set("clientOid", clientOrderID)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetIsolated + bitgetHistoryOrders
	var resp *MarginHistOrds
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedOrderFills returns the fills for isolated margin orders
func (bi *Bitget) GetIsolatedOrderFills(ctx context.Context, pair string, orderID, pagination, limit int64, startTime, endTime time.Time) (*MarginOrderFills, error) {
	if pair == "" {
		return nil, errPairEmpty
	}
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("symbol", pair)
	if orderID != 0 {
		params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetMargin + bitgetIsolated + "/" + bitgetFills
	var resp *MarginOrderFills
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetIsolatedLiquidationOrders returns the liquidation orders for isolated margin
func (bi *Bitget) GetIsolatedLiquidationOrders(ctx context.Context, orderType, pair, fromCoin, toCoin string, startTime, endTime time.Time, limit, pagination int64) (*LiquidationResp, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	params.Values.Set("orderType", orderType)
	params.Values.Set("symbol", pair)
	params.Values.Set("fromCoin", fromCoin)
	params.Values.Set("toCoin", toCoin)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetMargin + bitgetIsolated + bitgetLiquidationOrder
	var resp *LiquidationResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetSavingsProductList returns the list of savings products for a particular currency
func (bi *Bitget) GetSavingsProductList(ctx context.Context, currency, filter string) (*SavingsProductList, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	vals.Set("filter", filter)
	path := bitgetEarn + bitgetSavings + bitgetProduct
	var resp *SavingsProductList
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetSavingsBalance returns the savings balance and amount earned in BTC and USDT
func (bi *Bitget) GetSavingsBalance(ctx context.Context) (*SavingsBalance, error) {
	path := bitgetEarn + bitgetSavings + "/" + bitgetAccount
	var resp *SavingsBalance
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, nil, nil,
		&resp)
}

// GetSavingsAssets returns information on assets held over the last three months
func (bi *Bitget) GetSavingsAssets(ctx context.Context, periodType string, startTime, endTime time.Time, limit, pagination int64) (*SavingsAssets, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("periodType", periodType)
	path := bitgetEarn + bitgetSavings + bitgetAssets
	var resp *SavingsAssets
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetSavingsRecords returns information on transactions performed over the last three months
func (bi *Bitget) GetSavingsRecords(ctx context.Context, currency, periodType, orderType string, startTime, endTime time.Time, limit, pagination int64) (*SavingsRecords, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("coin", currency)
	params.Values.Set("periodType", periodType)
	params.Values.Set("orderType", orderType)
	path := bitgetEarn + bitgetSavings + bitgetRecords
	var resp *SavingsRecords
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetSavingsSubscriptionDetail returns detailed information on subscribing, for a single product
func (bi *Bitget) GetSavingsSubscriptionDetail(ctx context.Context, productID int64, periodType string) (*SavingsSubDetail, error) {
	if productID == 0 {
		return nil, errProductIDEmpty
	}
	if periodType == "" {
		return nil, errPeriodTypeEmpty
	}
	vals := url.Values{}
	vals.Set("productId", strconv.FormatInt(productID, 10))
	vals.Set("periodType", periodType)
	path := bitgetEarn + bitgetSavings + bitgetSubscribeInfo
	var resp *SavingsSubDetail
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// SubscribeSavings applies funds to a savings product
func (bi *Bitget) SubscribeSavings(ctx context.Context, productID int64, periodType string, amount float64) (*SaveResp, error) {
	if productID == 0 {
		return nil, errProductIDEmpty
	}
	if periodType == "" {
		return nil, errPeriodTypeEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"productId":  productID,
		"periodType": periodType,
		"amount":     strconv.FormatFloat(amount, 'f', -1, 64),
	}
	path := bitgetEarn + bitgetSavings + bitgetSubscribe
	var resp *SaveResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetSavingsSubscriptionResult returns the result of a subscription attempt
func (bi *Bitget) GetSavingsSubscriptionResult(ctx context.Context, orderID int64, periodType string) (*SaveResult, error) {
	if orderID == 0 {
		return nil, errOrderIDEmpty
	}
	if periodType == "" {
		return nil, errPeriodTypeEmpty
	}
	vals := url.Values{}
	vals.Set("orderId", strconv.FormatInt(orderID, 10))
	vals.Set("periodType", periodType)
	path := bitgetEarn + bitgetSavings + bitgetSubscribeResult
	var resp *SaveResult
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// RedeemSavings redeems funds from a savings product
func (bi *Bitget) RedeemSavings(ctx context.Context, productID, orderID int64, periodType string, amount float64) (*SaveResp, error) {
	if productID == 0 {
		return nil, errProductIDEmpty
	}
	if periodType == "" {
		return nil, errPeriodTypeEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"productId":  productID,
		"orderId":    orderID,
		"periodType": periodType,
		"amount":     strconv.FormatFloat(amount, 'f', -1, 64),
	}
	path := bitgetEarn + bitgetSavings + bitgetRedeem
	var resp *SaveResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetSavingsRedemptionResult returns the result of a redemption attempt
func (bi *Bitget) GetSavingsRedemptionResult(ctx context.Context, orderID int64, periodType string) (*SaveResult, error) {
	if orderID == 0 {
		return nil, errOrderIDEmpty
	}
	if periodType == "" {
		return nil, errPeriodTypeEmpty
	}
	vals := url.Values{}
	vals.Set("orderId", strconv.FormatInt(orderID, 10))
	vals.Set("periodType", periodType)
	path := bitgetEarn + bitgetSavings + bitgetRedeemResult
	var resp *SaveResult
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetEarnAccountAssets returns the assets in the earn account
func (bi *Bitget) GetEarnAccountAssets(ctx context.Context, currency string) (*EarnAssets, error) {
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetEarn + bitgetAccount + bitgetAssets
	var resp *EarnAssets
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetSharkFinProducts returns information on Shark Fin products
func (bi *Bitget) GetSharkFinProducts(ctx context.Context, currency string, limit, pagination int64) (*SharkFinProducts, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	if limit != 0 {
		vals.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		vals.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	path := bitgetEarn + bitgetSharkFin + bitgetProduct
	var resp *SharkFinProducts
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// GetSharkFinBalance returns the balance and amount earned in BTC and USDT for Shark Fin products
func (bi *Bitget) GetSharkFinBalance(ctx context.Context) (*SharkFinBalance, error) {
	path := bitgetEarn + bitgetSharkFin + "/" + bitgetAccount
	var resp *SharkFinBalance
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, nil, nil,
		&resp)
}

// GetSharkFinAssets returns information on assets held over the last three months for Shark Fin products
func (bi *Bitget) GetSharkFinAssets(ctx context.Context, status string, startTime, endTime time.Time, limit, pagination int64) (*SharkFinAssets, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("status", status)
	path := bitgetEarn + bitgetSharkFin + bitgetAssets
	var resp *SharkFinAssets
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetSharkFinRecords returns information on transactions performed over the last three months for Shark Fin products
func (bi *Bitget) GetSharkFinRecords(ctx context.Context, currency, transactionType string, startTime, endTime time.Time, limit, pagination int64) (*SharkFinRecords, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, true, true)
	if err != nil {
		return nil, err
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	params.Values.Set("coin", currency)
	params.Values.Set("type", transactionType)
	path := bitgetEarn + bitgetSharkFin + bitgetRecords
	var resp *SharkFinRecords
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// GetSharkFinSubscriptionDetail returns detailed information on subscribing, for a single product
func (bi *Bitget) GetSharkFinSubscriptionDetail(ctx context.Context, productID int64) (*SharkFinSubDetail, error) {
	if productID == 0 {
		return nil, errProductIDEmpty
	}
	vals := url.Values{}
	vals.Set("productId", strconv.FormatInt(productID, 10))
	path := bitgetEarn + bitgetSharkFin + bitgetSubscribeInfo
	var resp *SharkFinSubDetail
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// SubscribeSharkFin applies funds to a Shark Fin product
func (bi *Bitget) SubscribeSharkFin(ctx context.Context, productID int64, amount float64) (*SaveResp, error) {
	if productID == 0 {
		return nil, errProductIDEmpty
	}
	if amount == 0 {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"productId": productID,
		"amount":    strconv.FormatFloat(amount, 'f', -1, 64),
	}
	path := bitgetEarn + bitgetSharkFin + bitgetSubscribe
	var resp *SaveResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetSharkFinSubscriptionResult returns the result of a subscription attempt
func (bi *Bitget) GetSharkFinSubscriptionResult(ctx context.Context, orderID int64) (*SaveResult, error) {
	if orderID == 0 {
		return nil, errOrderIDEmpty
	}
	vals := url.Values{}
	vals.Set("orderId", strconv.FormatInt(orderID, 10))
	path := bitgetEarn + bitgetSharkFin + bitgetSubscribeResult
	var resp *SaveResult
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate5, http.MethodGet, path, vals, nil,
		&resp)
}

// GetLoanCurrencyList returns the list of currencies available for loan
func (bi *Bitget) GetLoanCurrencyList(ctx context.Context, currency string) (*LoanCurList, error) {
	if currency == "" {
		return nil, errCurrencyEmpty
	}
	vals := url.Values{}
	vals.Set("coin", currency)
	path := bitgetEarn + bitgetLoan + "/" + bitgetPublic + bitgetCoinInfos
	var resp *LoanCurList
	// return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
	// 	&resp)
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, vals, &resp)
}

// GetEstimatedInterestAndBorrowable returns the estimated interest and borrowable amount for a currency
func (bi *Bitget) GetEstimatedInterestAndBorrowable(ctx context.Context, loanCoin, collateralCoin, term string, collateralAmount float64) (*EstimateInterest, error) {
	if loanCoin == "" {
		return nil, errLoanCoinEmpty
	}
	if collateralCoin == "" {
		return nil, errCollateralCoinEmpty
	}
	if term == "" {
		return nil, errTermEmpty
	}
	if collateralAmount == 0 {
		return nil, errCollateralAmountEmpty
	}
	vals := url.Values{}
	vals.Set("loanCoin", loanCoin)
	vals.Set("pledgeCoin", collateralCoin)
	vals.Set("daily", term)
	vals.Set("pledgeAmount", strconv.FormatFloat(collateralAmount, 'f', -1, 64))
	path := bitgetEarn + bitgetLoan + "/" + bitgetPublic + bitgetHourInterest
	var resp *EstimateInterest
	return resp, bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate10, path, vals, &resp)
}

// BorrowFunds borrows funds for a currency, supplying a certain amount of currency as collateral
func (bi *Bitget) BorrowFunds(ctx context.Context, loanCoin, collateralCoin, term string, collateralAmount, loanAmount float64) (*BorrowResp, error) {
	if loanCoin == "" {
		return nil, errLoanCoinEmpty
	}
	if collateralCoin == "" {
		return nil, errCollateralCoinEmpty
	}
	if term == "" {
		return nil, errTermEmpty
	}
	if (collateralAmount == 0 && loanAmount == 0) || (collateralAmount != 0 && loanAmount != 0) {
		return nil, errCollateralLoanMutex
	}
	req := map[string]interface{}{
		"loanCoin":     loanCoin,
		"pledgeCoin":   collateralCoin,
		"daily":        term,
		"pledgeAmount": strconv.FormatFloat(collateralAmount, 'f', -1, 64),
		"loanAmount":   strconv.FormatFloat(loanAmount, 'f', -1, 64),
	}
	path := bitgetEarn + bitgetLoan + bitgetBorrow
	var resp *BorrowResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetOngoingLoans returns the ongoing loans, optionally filtered by currency
func (bi *Bitget) GetOngoingLoans(ctx context.Context, orderID int64, loanCoin, collateralCoin string) (*OngoingLoans, error) {
	vals := url.Values{}
	if orderID != 0 {
		vals.Set("orderId", strconv.FormatInt(orderID, 10))
	}
	vals.Set("loanCoin", loanCoin)
	vals.Set("pledgeCoin", collateralCoin)
	path := bitgetEarn + bitgetLoan + bitgetOngoingOrders
	var resp *OngoingLoans
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, vals, nil,
		&resp)
}

// RepayLoan repays a loan
func (bi *Bitget) RepayLoan(ctx context.Context, orderID int64, amount float64, repayUnlock, repayAll bool) (*RepayResp, error) {
	if orderID == 0 {
		return nil, errOrderIDEmpty
	}
	if amount == 0 && !repayAll {
		return nil, errAmountEmpty
	}
	req := map[string]interface{}{
		"orderId": orderID,
		"amount":  strconv.FormatFloat(amount, 'f', -1, 64),
	}
	if repayUnlock {
		req["repayUnlock"] = "yes"
	} else {
		req["repayUnlock"] = "no"
	}
	if repayAll {
		req["repayAll"] = "yes"
	} else {
		req["repayAll"] = "no"
	}
	path := bitgetEarn + bitgetLoan + bitgetRepay
	var resp *RepayResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// GetLoanRepayHistory returns the repayment records for a loan
func (bi *Bitget) GetLoanRepayHistory(ctx context.Context, orderID, pagination, limit int64, loanCoin, pledgeCoin string, startTime, endTime time.Time) (*RepayRecords, error) {
	var params Params
	params.Values = make(url.Values)
	err := params.prepareDateString(startTime, endTime, false, false)
	if err != nil {
		return nil, err
	}
	params.Values.Set("orderId", strconv.FormatInt(orderID, 10))
	params.Values.Set("loanCoin", loanCoin)
	params.Values.Set("pledgeCoin", pledgeCoin)
	if pagination != 0 {
		params.Values.Set("idLessThan", strconv.FormatInt(pagination, 10))
	}
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatInt(limit, 10))
	}
	path := bitgetEarn + bitgetLoan + bitgetRepayHistory
	var resp *RepayRecords
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodGet, path, params.Values,
		nil, &resp)
}

// ModifyPledgeRate modifies the amount of collateral pledged for a loan
func (bi *Bitget) ModifyPledgeRate(ctx context.Context, orderID int64, amoutn float64, pledgeCoin, reviseType string) (*ModPledgeResp, error) {
	if orderID == 0 {
		return nil, errOrderIDEmpty
	}
	if amoutn == 0 {
		return nil, errAmountEmpty
	}
	if pledgeCoin == "" {
		return nil, errCollateralCoinEmpty
	}
	if reviseType == "" {
		return nil, errReviseTypeEmpty
	}
	req := map[string]interface{}{
		"orderId":    orderID,
		"amount":     strconv.FormatFloat(amoutn, 'f', -1, 64),
		"pledgeCoin": pledgeCoin,
		"reviseType": reviseType,
	}
	path := bitgetEarn + bitgetLoan + bitgetRevisePledge
	var resp *ModPledgeResp
	return resp, bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate10, http.MethodPost, path, nil, req,
		&resp)
}

// SendAuthenticatedHTTPRequest sends an authenticated HTTP request
func (bi *Bitget) SendAuthenticatedHTTPRequest(ctx context.Context, ep exchange.URL, rateLim request.EndpointLimit, method, path string, queryParams url.Values, bodyParams map[string]interface{}, result interface{}) error {
	creds, err := bi.GetCredentials(ctx)
	if err != nil {
		return err
	}
	endpoint, err := bi.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	path = common.EncodeURLValues(path, queryParams)
	newRequest := func() (*request.Item, error) {
		payload := []byte("")
		if bodyParams != nil {
			payload, err = json.Marshal(bodyParams)
			if err != nil {
				return nil, err
			}
		}
		t := strconv.FormatInt(time.Now().UnixMilli(), 10)
		message := t + method + "/api/v2/" + path + string(payload)
		// The exchange also supports user-generated RSA keys, but we haven't implemented that yet
		var hmac []byte
		hmac, err = crypto.GetHMAC(crypto.HashSHA256, []byte(message), []byte(creds.Secret))
		if err != nil {
			return nil, err
		}
		headers := make(map[string]string)
		headers["ACCESS-KEY"] = creds.Key
		headers["ACCESS-SIGN"] = crypto.Base64Encode(hmac)
		headers["ACCESS-TIMESTAMP"] = t
		headers["ACCESS-PASSPHRASE"] = creds.ClientID
		headers["Content-Type"] = "application/json"
		headers["locale"] = "en-US"
		return &request.Item{
			Method:        method,
			Path:          endpoint + path,
			Headers:       headers,
			Body:          bytes.NewBuffer(payload),
			Result:        &result,
			Verbose:       bi.Verbose,
			HTTPDebugging: bi.HTTPDebugging,
			HTTPRecording: bi.HTTPRecording,
		}, nil
	}
	return bi.SendPayload(ctx, rateLim, newRequest, request.AuthenticatedRequest)
}

// SendHTTPRequest sends an unauthenticated HTTP request, with a few assumptions about the request;
// namely that it is a GET request with no body
func (bi *Bitget) SendHTTPRequest(ctx context.Context, ep exchange.URL, rateLim request.EndpointLimit, path string, queryParams url.Values, result interface{}) error {
	endpoint, err := bi.API.Endpoints.GetURL(ep)
	if err != nil {
		return err
	}
	newRequest := func() (*request.Item, error) {
		path = common.EncodeURLValues(path, queryParams)
		return &request.Item{
			Method:        "GET",
			Path:          endpoint + path,
			Result:        &result,
			Verbose:       bi.Verbose,
			HTTPDebugging: bi.HTTPDebugging,
			HTTPRecording: bi.HTTPRecording,
		}, nil
	}
	return bi.SendPayload(ctx, rateLim, newRequest, request.UnauthenticatedRequest)
}

func (p *Params) prepareDateString(startDate, endDate time.Time, ignoreUnsetStart, ignoreUnsetEnd bool) error {
	if startDate.After(endDate) && !(endDate.IsZero() || endDate.Equal(common.ZeroValueUnix)) {
		return common.ErrStartAfterEnd
	}
	if startDate.Equal(endDate) && !(startDate.IsZero() || startDate.Equal(common.ZeroValueUnix)) {
		return common.ErrStartEqualsEnd
	}
	if startDate.After(time.Now()) {
		return common.ErrStartAfterTimeNow
	}
	if startDate.IsZero() || startDate.Equal(common.ZeroValueUnix) {
		if !ignoreUnsetStart {
			return fmt.Errorf("start %w", common.ErrDateUnset)
		}
	} else {
		p.Values.Set("startTime", strconv.FormatInt(startDate.UnixMilli(), 10))
	}
	if endDate.IsZero() || endDate.Equal(common.ZeroValueUnix) {
		if !ignoreUnsetEnd {
			return fmt.Errorf("end %w", common.ErrDateUnset)
		}
	} else {
		p.Values.Set("endTime", strconv.FormatInt(endDate.UnixMilli(), 10))
	}
	return nil
}

// UnmarshalJSON unmarshals the JSON input into a UnixTimestamp type
func (t *UnixTimestamp) UnmarshalJSON(b []byte) error {
	var timestampStr string
	err := json.Unmarshal(b, &timestampStr)
	if err != nil {
		return err
	}
	if timestampStr == "" {
		*t = UnixTimestamp(time.Time{})
		return nil
	}
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return err
	}
	*t = UnixTimestamp(time.UnixMilli(timestamp).UTC())
	return nil
}

// String implements the stringer interface
func (t *UnixTimestamp) String() string {
	return t.Time().String()
}

// Time returns the time.Time representation of the UnixTimestamp
func (t *UnixTimestamp) Time() time.Time {
	return time.Time(*t)
}

// UnmarshalJSON unmarshals the JSON input into a YesNoBool type
func (y *YesNoBool) UnmarshalJSON(b []byte) error {
	var yn string
	err := json.Unmarshal(b, &yn)
	if err != nil {
		return err
	}
	switch yn {
	case "yes":
		*y = true
	case "no":
		*y = false
	}
	return nil
}

// MarshalJSON marshals the YesNoBool type into a JSON string
func (y YesNoBool) MarshalJSON() ([]byte, error) {
	if y {
		return json.Marshal("YES")
	}
	return json.Marshal("NO")
}

// UnmarshalJSON unmarshals the JSON input into a SuccessBool type
func (s *SuccessBool) UnmarshalJSON(b []byte) error {
	var success string
	err := json.Unmarshal(b, &success)
	if err != nil {
		return err
	}
	switch success {
	case "success":
		*s = true
	case "failure", "fail":
		*s = false
	}
	return nil
}

// UnmarshalJSON unmarshals the JSON input into an EmptyInt type
func (e *EmptyInt) UnmarshalJSON(b []byte) error {
	var num string
	err := json.Unmarshal(b, &num)
	if err != nil {
		return err
	}
	if num == "" {
		*e = 0
		return nil
	}
	i, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return err
	}
	*e = EmptyInt(i)
	return nil
}

// CandlestickHelper pulls out common candlestick functionality to avoid repetition
func (bi *Bitget) candlestickHelper(ctx context.Context, pair, granularity, path string, limit uint16, params Params) (*CandleData, error) {
	params.Values.Set("symbol", pair)
	params.Values.Set("granularity", granularity)
	if limit != 0 {
		params.Values.Set("limit", strconv.FormatUint(uint64(limit), 10))
	}
	var resp *CandleResponse
	err := bi.SendHTTPRequest(ctx, exchange.RestSpot, Rate20, path, params.Values, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, errNoCandleData
	}
	var spot bool
	var data CandleData
	if resp.Data[0][7] == nil {
		data.FuturesCandles = make([]OneFuturesCandle, len(resp.Data))
	} else {
		spot = true
		data.SpotCandles = make([]OneSpotCandle, len(resp.Data))
	}
	for i := range resp.Data {
		timeTemp, ok := resp.Data[i][0].(string)
		if !ok {
			return nil, errTypeAssertTimestamp
		}
		timeTemp = (timeTemp)[1 : len(timeTemp)-1]
		timeTemp2, err := strconv.ParseInt(timeTemp, 10, 64)
		if err != nil {
			return nil, err
		}
		openTemp, ok := resp.Data[i][1].(string)
		if !ok {
			return nil, errTypeAssertOpenPrice
		}
		highTemp, ok := resp.Data[i][2].(string)
		if !ok {
			return nil, errTypeAssertHighPrice
		}
		lowTemp, ok := resp.Data[i][3].(string)
		if !ok {
			return nil, errTypeAssertLowPrice
		}
		closeTemp, ok := resp.Data[i][4].(string)
		if !ok {
			return nil, errTypeAssertClosePrice
		}
		baseVolumeTemp := resp.Data[i][5].(string)
		if !ok {
			return nil, errTypeAssertBaseVolume
		}
		quoteVolumeTemp := resp.Data[i][6].(string)
		if !ok {
			return nil, errTypeAssertQuoteVolume
		}
		if spot {
			usdtVolumeTemp := resp.Data[i][7].(string)
			if !ok {
				return nil, errTypeAssertUSDTVolume
			}
			data.SpotCandles[i].Timestamp = time.Time(UnixTimestamp(time.UnixMilli(timeTemp2).UTC()))
			data.SpotCandles[i].Open, err = strconv.ParseFloat(openTemp, 64)
			if err != nil {
				return nil, err
			}
			data.SpotCandles[i].High, err = strconv.ParseFloat(highTemp, 64)
			if err != nil {
				return nil, err
			}
			data.SpotCandles[i].Low, err = strconv.ParseFloat(lowTemp, 64)
			if err != nil {
				return nil, err
			}
			data.SpotCandles[i].Close, err = strconv.ParseFloat(closeTemp, 64)
			if err != nil {
				return nil, err
			}
			data.SpotCandles[i].BaseVolume, err = strconv.ParseFloat(baseVolumeTemp, 64)
			if err != nil {
				return nil, err
			}
			data.SpotCandles[i].QuoteVolume, err = strconv.ParseFloat(quoteVolumeTemp, 64)
			if err != nil {
				return nil, err
			}
			data.SpotCandles[i].USDTVolume, err = strconv.ParseFloat(usdtVolumeTemp, 64)
			if err != nil {
				return nil, err
			}
		} else {
			data.FuturesCandles[i].Timestamp = time.Time(UnixTimestamp(time.UnixMilli(timeTemp2).UTC()))
			data.FuturesCandles[i].Entry, err = strconv.ParseFloat(openTemp, 64)
			if err != nil {
				return nil, err
			}
			data.FuturesCandles[i].High, err = strconv.ParseFloat(highTemp, 64)
			if err != nil {
				return nil, err
			}
			data.FuturesCandles[i].Low, err = strconv.ParseFloat(lowTemp, 64)
			if err != nil {
				return nil, err
			}
			data.FuturesCandles[i].Exit, err = strconv.ParseFloat(closeTemp, 64)
			if err != nil {
				return nil, err
			}
			data.FuturesCandles[i].BaseVolume, err = strconv.ParseFloat(baseVolumeTemp, 64)
			if err != nil {
				return nil, err
			}
			data.FuturesCandles[i].QuoteVolume, err = strconv.ParseFloat(quoteVolumeTemp, 64)
			if err != nil {
				return nil, err
			}
		}
	}
	return &data, nil
}

// spotOrderHelper is a helper function for unmarshalling spot order endpoints
func (bi *Bitget) spotOrderHelper(ctx context.Context, path string, vals url.Values) (*SpotOrderDetailResp, error) {
	var temp *OrderDetailTemp
	err := bi.SendAuthenticatedHTTPRequest(ctx, exchange.RestSpot, Rate20, http.MethodGet, path, vals, nil,
		&temp)
	if err != nil {
		return nil, err
	}
	resp := new(SpotOrderDetailResp)
	resp.Data = make([]SpotOrderDetailData, len(temp.Data))
	for i := range temp.Data {
		resp.Data[i].UserID = temp.Data[i].UserID
		resp.Data[i].Symbol = temp.Data[i].Symbol
		resp.Data[i].OrderID = temp.Data[i].OrderID
		resp.Data[i].ClientOrderID = temp.Data[i].ClientOrderID
		resp.Data[i].Price = temp.Data[i].Price
		resp.Data[i].Size = temp.Data[i].Size
		resp.Data[i].OrderType = temp.Data[i].OrderType
		resp.Data[i].Side = temp.Data[i].Side
		resp.Data[i].Status = temp.Data[i].Status
		resp.Data[i].PriceAverage = temp.Data[i].PriceAverage
		resp.Data[i].BaseVolume = temp.Data[i].BaseVolume
		resp.Data[i].QuoteVolume = temp.Data[i].QuoteVolume
		resp.Data[i].EnterPointSource = temp.Data[i].EnterPointSource
		resp.Data[i].CreationTime = temp.Data[i].CreationTime
		resp.Data[i].UpdateTime = temp.Data[i].UpdateTime
		resp.Data[i].OrderSource = temp.Data[i].OrderSource
		err = json.Unmarshal(temp.Data[i].FeeDetailTemp, &resp.Data[i].FeeDetail)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}
