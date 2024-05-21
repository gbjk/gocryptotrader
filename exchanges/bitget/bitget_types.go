package bitget

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/thrasher-corp/gocryptotrader/types"
)

// Params is used within functions to make the setting of parameters easier
type Params struct {
	url.Values
}

// UnixTimestamp is a type used to unmarshal unix millisecond timestamps returned from the exchange
type UnixTimestamp time.Time

// AnnResp holds information on announcements
type AnnResp struct {
	Data []struct {
		AnnID    string        `json:"annId"`
		AnnTitle string        `json:"annTitle"`
		AnnDesc  string        `json:"annDesc"`
		CTime    UnixTimestamp `json:"cTime"`
		Language string        `json:"language"`
		AnnURL   string        `json:"annUrl"`
	} `json:"data"`
}

// TimeResp holds information on the current server time
type TimeResp struct {
	Data struct {
		ServerTime UnixTimestamp `json:"serverTime"`
	} `json:"data"`
}

// TradeRateResp holds information on the current maker and taker fee rates
type TradeRateResp struct {
	Data struct {
		MakerFeeRate float64 `json:"makerFeeRate,string"`
		TakerFeeRate float64 `json:"takerFeeRate,string"`
	} `json:"data"`
}

// SpotTrResp holds information on spot transactions
type SpotTrResp struct {
	Data []struct {
		ID          int64         `json:"id,string"`
		Coin        string        `json:"coin"`
		SpotTaxType string        `json:"spotTaxType"`
		Amount      float64       `json:"amount,string"`
		Fee         float64       `json:"fee,string"`
		Balance     float64       `json:"balance,string"`
		Timestamp   UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// FutureTrResp holds information on futures transactions
type FutureTrResp struct {
	Data []struct {
		ID            int64         `json:"id,string"`
		Symbol        string        `json:"symbol"`
		MarginCoin    string        `json:"marginCoin"`
		FutureTaxType string        `json:"futureTaxType"`
		Amount        float64       `json:"amount,string"`
		Fee           float64       `json:"fee,string"`
		Timestamp     UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// MarginTrResp holds information on margin transactions
type MarginTrResp struct {
	Data []struct {
		ID            int64         `json:"id,string"`
		Coin          string        `json:"coin"`
		Symbol        string        `json:"symbol"`
		MarginTaxType string        `json:"marginTaxType"`
		Amount        float64       `json:"amount,string"`
		Fee           float64       `json:"fee,string"`
		Total         float64       `json:"total,string"`
		Timestamp     UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// P2PTrResp holds information on P2P transactions
type P2PTrResp struct {
	Data []struct {
		ID         int64         `json:"id,string"`
		Coin       string        `json:"coin"`
		P2PTaxType string        `json:"p2pTaxType"`
		Total      float64       `json:"total,string"`
		Timestamp  UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// P2PMerResp holds information on P2P merchant lists
type P2PMerListResp struct {
	Data struct {
		MerchantList []struct {
			RegisterTime        UnixTimestamp `json:"registerTime"`
			NickName            string        `json:"nickName"`
			IsOnline            string        `json:"isOnline"`
			MerchantID          int64         `json:"merchantId,string"`
			AvgPaymentTime      int64         `json:"avgPaymentTime,string"`
			AvgReleaseTime      int64         `json:"avgReleaseTime,string"`
			TotalTrades         int64         `json:"totalTrades,string"`
			TotalBuy            int64         `json:"totalBuy,string"`
			TotalSell           int64         `json:"totalSell,string"`
			TotalCompletionRate float64       `json:"totalCompletionRate,string"`
			Trades30D           int64         `json:"trades30d,string"`
			Sell30D             float64       `json:"sell30d,string"`
			Buy30D              float64       `json:"buy30d,string"`
			CompletionRate30D   float64       `json:"completionRate30d,string"`
		} `json:"merchantList"`
		MinMerchantID int64 `json:"minMerchantId,string"`
	} `json:"data"`
}

// YesNoBool is a type used to unmarshal strings that are either "yes" or "no" into bools
type YesNoBool bool

// P2PMerInfoResp holds information on P2P merchant information
type P2PMerInfoResp struct {
	Data struct {
		RegisterTime        UnixTimestamp `json:"registerTime"`
		NickName            string        `json:"nickName"`
		MerchantID          int64         `json:"merchantId,string"`
		AvgPaymentTime      int64         `json:"avgPaymentTime,string"`
		AvgReleaseTime      int64         `json:"avgReleaseTime,string"`
		TotalTrades         int64         `json:"totalTrades,string"`
		TotalBuy            int64         `json:"totalBuy,string"`
		TotalSell           int64         `json:"totalSell,string"`
		TotalCompletionRate float64       `json:"totalCompletionRate,string"`
		Trades30D           int64         `json:"trades30d,string"`
		Sell30D             float64       `json:"sell30d,string"`
		Buy30D              float64       `json:"buy30d,string"`
		CompletionRate30D   float64       `json:"completionRate30d,string"`
		KYCStatus           YesNoBool     `json:"kycStatus"`
		EmailBindStatus     YesNoBool     `json:"emailBindStatus"`
		MobileBindStatus    YesNoBool     `json:"mobileBindStatus"`
		Email               string        `json:"email"`
		Mobile              string        `json:"mobile"`
	} `json:"data"`
}

// P2POrdersResp holds information on P2P orders
type P2POrdersResp struct {
	Data struct {
		OrderList []struct {
			OrderID        int64         `json:"orderId,string"`
			OrderNum       int64         `json:"orderNo,string"`
			AdvNum         int64         `json:"advNo,string"`
			Side           string        `json:"side"`
			Count          float64       `json:"count,string"`
			FiatCurrency   string        `json:"fiat"`
			CryptoCurrency string        `json:"coin"`
			Price          float64       `json:"price,string"`
			WithdrawTime   UnixTimestamp `json:"withdrawTime"`
			RepresentTime  UnixTimestamp `json:"representTime"`
			ReleaseTime    UnixTimestamp `json:"releaseTime"`
			PaymentTime    UnixTimestamp `json:"paymentTime"`
			Amount         float64       `json:"amount,string"`
			Status         string        `json:"status"`
			BuyerRealName  string        `json:"buyerRealName"`
			SellerRealName string        `json:"sellerRealName"`
			CreationTime   UnixTimestamp `json:"ctime"`
			UpdateTime     UnixTimestamp `json:"utime"`
			PaymentInfo    struct {
				PayMethodName string `json:"paymethodName"`
				PayMethodID   string `json:"paymethodId"`
				PayMethodInfo []struct {
					Name     string    `json:"name"`
					Required YesNoBool `json:"required"`
					Type     string    `json:"type"`
					Value    string    `json:"value"`
				} `json:"paymethodInfo"`
			} `json:"paymentInfo"`
		} `json:"orderList"`
		MinOrderID int64 `json:"minOrderId,string"`
	} `json:"data"`
}

// P2PAdListResp holds information on P2P advertisements
type P2PAdListResp struct {
	Data struct {
		AdList []struct {
			AdID            int64         `json:"adId,string"`
			AdvNum          int64         `json:"advNo,string"`
			Side            string        `json:"side"`
			AdSize          float64       `json:"adSize,string"`
			Size            float64       `json:"size,string"`
			CryptoCurrency  string        `json:"coin"`
			Price           float64       `json:"price,string"`
			CryptoPrecision uint8         `json:"coinPrecision,string"`
			FiatCurrency    string        `json:"fiat"`
			FiatPrecision   uint8         `json:"fiatPrecision,string"`
			FiatSymbol      string        `json:"fiatSymbol"`
			Status          string        `json:"status"`
			Hide            YesNoBool     `json:"hide"`
			MaxTradeAmount  float64       `json:"maxTradeAmount,string"`
			MinTradeAmount  float64       `json:"minTradeAmount,string"`
			PayDuration     int64         `json:"payDuration,string"`
			TurnoverNum     int64         `json:"turnoverNum,string"`
			TurnoverRate    float64       `json:"turnoverRate,string"`
			Label           string        `json:"label"`
			CreationTime    UnixTimestamp `json:"ctime"`
			UpdateTime      UnixTimestamp `json:"utime"`
			UserLimitList   struct {
				MinCompleteNum     int64     `json:"minCompleteNum,string"`
				MaxCompleteNum     int64     `json:"maxCompleteNum,string"`
				PlaceOrderNum      int64     `json:"placeOrderNum,string"`
				AllowMerchantPlace YesNoBool `json:"allowMerchantPlace"`
				CompleteRate30D    float64   `json:"completeRate30d,string"`
				Country            string    `json:"country"`
			} `json:"userLimitList"`
			PaymentMethodList []struct {
				PaymentMethod string `json:"paymentMethod"`
				PaymentID     int64  `json:"paymentId,string"`
				PaymentInfo   []struct {
					Required bool   `json:"required"`
					Name     string `json:"name"`
					Type     string `json:"type"`
				} `json:"paymentInfo"`
			} `json:"paymentMethodList"`
			MerchantCertifiedList []struct {
				ImageURL string `json:"imageUrl"`
				Desc     string `json:"desc"`
			} `json:"merchantCertifiedList"`
		} `json:"advList"`
	} `json:"data"`
}

// CrVirSubResp contains information returned when creating virtual sub-accounts
type CrVirSubResp struct {
	Data struct {
		FailureList []struct {
			SubaccountName string `json:"subaAccountName"`
		} `json:"failureList"`
		SuccessList []struct {
			SubaccountUID  string        `json:"subAccountUid"`
			SubaccountName string        `json:"subaAccountName"`
			Status         string        `json:"status"`
			PermList       []string      `json:"permList"`
			Label          string        `json:"label"`
			CreationTime   UnixTimestamp `json:"cTime"`
			UpdateTime     UnixTimestamp `json:"uTime"`
		} `json:"successList"`
	} `json:"data"`
}

// SuccessBool is a type used to unmarshal strings that are either "success" or "failure" into bools
type SuccessBool bool

// SuccessBoolResp contains a success bool in the format returned by the exchange
type SuccessBoolResp struct {
	Data struct {
		Success SuccessBool `json:"result"`
	} `json:"data"`
}

// CrSubAccAPIKeyResp contains information returned when simultaneously creating a sub-account and
// an API key
type CrSubAccAPIKeyResp struct {
	Data []struct {
		SubaccountUID    string   `json:"subAccountUid"`
		SubaccountName   string   `json:"subAccountName"`
		Label            string   `json:"label"`
		SubaccountAPIKey string   `json:"subAccountApiKey"`
		SecretKey        string   `json:"secretKey"`
		PermList         []string `json:"permList"`
		IPList           []string `json:"ipList"`
	} `json:"data"`
}

// GetVirSubResp contains information on the user's virtual sub-accounts
type GetVirSubResp struct {
	Data struct {
		SubaccountList []struct {
			SubaccountUID  string        `json:"subAccountUid"`
			SubaccountName string        `json:"subAccountName"`
			Label          string        `json:"label"`
			Status         string        `json:"status"`
			PermList       []string      `json:"permList"`
			CreationTime   UnixTimestamp `json:"cTime"`
			UpdateTime     UnixTimestamp `json:"uTime"`
		} `json:"subAccountList"`
		EndID int64 `json:"endId,string"`
	} `json:"data"`
}

// AlterAPIKeyResp contains information returned when creating or modifying an API key
type AlterAPIKeyResp struct {
	Data struct {
		SubaccountUID    string   `json:"subAccountUid"`
		SubaccountApiKey string   `json:"subAccountApiKey"`
		SecretKey        string   `json:"secretKey"`
		PermList         []string `json:"permList"`
		Label            string   `json:"label"`
		IPList           []string `json:"ipList"`
	} `json:"data"`
}

// GetAPIKeyResp contains information on the user's API keys
type GetAPIKeyResp struct {
	Data []struct {
		SubaccountUID    string   `json:"subAccountUid"`
		SubaccountApiKey string   `json:"subAccountApiKey"`
		IPList           []string `json:"ipList"`
		PermList         []string `json:"permList"`
		Label            string   `json:"label"`
	} `json:"data"`
}

// ConvertCoinsResp contains information on the user's available currencies
type ConvertCoinsResp struct {
	Data []struct {
		Coin      string  `json:"coin"`
		Available float64 `json:"available,string"`
		MaxAmount float64 `json:"maxAmount,string"`
		MinAmount float64 `json:"minAmount,string"`
	} `json:"data"`
}

// QuotedPriceResp contains information on a queried conversion
type QuotedPriceResp struct {
	Data struct {
		FromCoin     string  `json:"fromCoin"`
		FromCoinSize float64 `json:"fromCoinSize,string"`
		ConvertPrice float64 `json:"cnvtPrice,string"`
		ToCoin       string  `json:"toCoin"`
		ToCoinSize   float64 `json:"toCoinSize,string"`
		TraceID      string  `json:"traceId"`
		Fee          float64 `json:"fee,string"`
	} `json:"data"`
}

// CommitConvResp contains information on a committed conversion
type CommitConvResp struct {
	Data struct {
		ToCoin       string        `json:"toCoin"`
		ToCoinSize   float64       `json:"toCoinSize,string"`
		ConvertPrice float64       `json:"cnvtPrice,string"`
		Timestamp    UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// ConvHistResp contains information on the user's conversion history
type ConvHistResp struct {
	Data struct {
		DataList []struct {
			ID           int64         `json:"id,string"`
			Timestamp    UnixTimestamp `json:"ts"`
			ConvertPrice float64       `json:"cnvtPrice,string"`
			Fee          float64       `json:"fee,string"`
			FromCoinSize float64       `json:"fromCoinSize,string"`
			FromCoin     string        `json:"fromCoin"`
			ToCoinSize   float64       `json:"toCoinSize,string"`
			ToCoin       string        `json:"toCoin"`
		} `json:"dataList"`
		EndID int64 `json:"endId,string"`
	} `json:"data"`
}

// BGBConvertCoinsResp contains information on the user's available currencies and conversions between those
// and BGB
type BGBConvertCoinsResp struct {
	Data struct {
		CoinList []struct {
			Coin         string  `json:"coin"`
			Available    float64 `json:"available,string"`
			BGBEstAmount float64 `json:"bgbEstAmount,string"`
			Precision    uint8   `json:"precision"`
			FeeDetail    []struct {
				FeeRate float64 `json:"feeRate,string"`
				Fee     float64 `json:"fee,string"`
			} `json:"feeDetail"`
			CurrentTime UnixTimestamp `json:"cTime"`
		} `json:"coinList"`
	} `json:"data"`
}

// ConvertBGBResp contains information on a series of conversions between BGB and other currencies
type ConvertBGBResp struct {
	Data struct {
		OrderList []struct {
			Coin    string `json:"coin"`
			OrderID int64  `json:"orderId,string"`
		} `json:"orderList"`
	} `json:"data"`
}

// BGBConvHistResp contains information on the user's conversion history between BGB and other currencies
type BGBConvHistResp struct {
	Data []struct {
		OrderID       int64   `json:"orderId,string"`
		FromCoin      string  `json:"fromCoin"`
		FromAmount    float64 `json:"fromAmount,string"`
		FromCoinPrice float64 `json:"fromCoinPrice,string"`
		ToCoin        string  `json:"toCoin"`
		ToAmount      float64 `json:"toAmount,string"`
		ToCoinPrice   float64 `json:"toCoinPrice,string"`
		FeeDetail     struct {
			FeeCoin string  `json:"feeCoin"`
			Fee     float64 `json:"fee,string"`
		} `json:"feeDetail"`
		Status       SuccessBool   `json:"status"`
		CreationTime UnixTimestamp `json:"cTime"`
	} `json:"data"`
}

// CoinInfoResp contains information on supported spot currencies
type CoinInfoResp struct {
	Data []struct {
		CoinID   uint32 `json:"coinId,string"`
		Coin     string `json:"coin"`
		Transfer bool   `json:"transfer,string"`
		Chains   []struct {
			Chain             string  `json:"chain"`
			NeedTag           bool    `json:"needTag,string"`
			Withdrawable      bool    `json:"withdrawable,string"`
			Rechargeable      bool    `json:"rechargeable,string"`
			WithdrawFee       float64 `json:"withdrawFee,string"`
			ExtraWithdrawFee  float64 `json:"extraWithdrawFee,string"`
			DepositConfirm    uint16  `json:"depositConfirm,string"`
			WithdrawConfirm   uint16  `json:"withdrawConfirm,string"`
			MinDepositAmount  float64 `json:"minDepositAmount,string"`
			MinWithdrawAmount float64 `json:"minWithdrawAmount,string"`
			BrowserURL        string  `json:"browserUrl"`
			ContractAddress   string  `json:"contractAddress"`
			WithdrawStep      uint8   `json:"withdrawStep,string"`
		} `json:"chains"`
	} `json:"data"`
}

// SymbolInfoResp contains information on supported spot trading pairs
type SymbolInfoResp struct {
	Data []struct {
		Symbol              string       `json:"symbol"`
		BaseCoin            string       `json:"baseCoin"`
		QuoteCoin           string       `json:"quoteCoin"`
		MinTradeAmount      float64      `json:"minTradeAmount,string"`
		MaxTradeAmount      float64      `json:"maxTradeAmount,string"`
		TakerFeeRate        float64      `json:"takerFeeRate,string"`
		MakerFeeRate        float64      `json:"makerFeeRate,string"`
		PricePrecision      uint8        `json:"pricePrecision,string"`
		QuantityPrecision   uint8        `json:"quantityPrecision,string"`
		QuotePrecision      uint8        `json:"quotePrecision,string"`
		MinTradeUSDT        float64      `json:"minTradeUSDT,string"`
		Status              string       `json:"status"`
		BuyLimitPriceRatio  types.Number `json:"buyLimitPriceRatio"`
		SellLimitPriceRatio types.Number `json:"sellLimitPriceRatio"`
	} `json:"data"`
}

// VIPFeeRateResp contains information on the different levels of VIP fee rates
type VIPFeeRateResp struct {
	Data []struct {
		Level        uint8   `json:"level,string"`
		DealAmount   float64 `json:"dealAmount,string"`
		AssetAmount  float64 `json:"assetAmount,string"`
		TakerFeeRate float64 `json:"takerFeeRate,string"`
		MakerFeeRate float64 `json:"makerFeeRate,string"`
		// 24-hour withdrawal limits
		BTCWithdrawAmount float64 `json:"btcWithdrawAmount,string"`
		USDWithdrawAmount float64 `json:"usdWithdrawAmount,string"`
	} `json:"data"`
}

// TickerResp contains information on tickers
type TickerResp struct {
	Data []struct {
		Symbol       string        `json:"symbol"`
		High24H      float64       `json:"high24h,string"`
		Open         float64       `json:"open,string"`
		LastPrice    float64       `json:"lastPr,string"`
		Low24H       float64       `json:"low24h,string"`
		QuoteVolume  float64       `json:"quoteVolume,string"`
		BaseVolume   float64       `json:"baseVolume,string"`
		USDTVolume   float64       `json:"usdtVolume,string"`
		BidPrice     float64       `json:"bidPr,string"`
		AskPrice     float64       `json:"askPr,string"`
		BidSize      float64       `json:"bidSz,string"`
		AskSize      float64       `json:"askSz,string"`
		OpenUTC      float64       `json:"openUTC,string"`
		Timestamp    UnixTimestamp `json:"ts"`
		ChangeUTC24H float64       `json:"changeUTC24h,string"`
		Change24H    float64       `json:"change24h,string"`
	} `json:"data"`
}

// DepthResp contains information on orderbook bids and asks, and any merging of orders done to them
type DepthResp struct {
	Data struct {
		Asks           [][2]float64  `json:"asks"`
		Bids           [][2]float64  `json:"bids"`
		Precision      string        `json:"precision"`
		Scale          float64       `json:"scale,string"`
		IsMaxPrecision YesNoBool     `json:"isMaxPrecision"`
		Timestamp      UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// OrderbookResp contains information on orderbook bids and asks
type OrderbookResp struct {
	Data struct {
		Asks      [][2]types.Number `json:"asks"`
		Bids      [][2]types.Number `json:"bids"`
		Timestamp UnixTimestamp     `json:"ts"`
	} `json:"data"`
}

// CandleResponse contains unsorted candle data
type CandleResponse struct {
	Data [][8]interface{} `json:"data"`
}

// OneSpotCandle contains a single candle
type OneSpotCandle struct {
	Timestamp   time.Time
	Open        float64
	High        float64
	Low         float64
	Close       float64
	BaseVolume  float64
	QuoteVolume float64
	USDTVolume  float64
}

// OneFuturesCandle contains a single candle
type OneFuturesCandle struct {
	Timestamp   time.Time
	Entry       float64
	High        float64
	Low         float64
	Exit        float64
	BaseVolume  float64
	QuoteVolume float64
}

// CandleData contains sorted candle data
type CandleData struct {
	SpotCandles    []OneSpotCandle
	FuturesCandles []OneFuturesCandle
}

// MarketFillsResp contains information on a batch of trades
type MarketFillsResp struct {
	Data []struct {
		Symbol    string        `json:"symbol"`
		TradeID   int64         `json:"tradeId,string"`
		Side      string        `json:"side"`
		Price     float64       `json:"price,string"`
		Size      float64       `json:"size,string"`
		Timestamp UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// OrderResp contains information on an order
type OrderResp struct {
	Data struct {
		OrderID       EmptyInt `json:"orderId"`
		ClientOrderID string   `json:"clientOid"`
	} `json:"data"`
}

// PlaceOrderStruct contains information on an order to be placed
type PlaceSpotOrderStruct struct {
	Side          string  `json:"side"`
	OrderType     string  `json:"orderType"`
	Strategy      string  `json:"force"`
	Price         float64 `json:"price,string"`
	Size          float64 `json:"size,string"`
	ClientOrderID string  `json:"clientOId"`
}

// EmptyInt is a type used to unmarshal empty string into 0, and numbers encoded as strings into int64
type EmptyInt int64

// BatchOrderResp contains information on the success or failure of a batch of orders to place or cancel
type BatchOrderResp struct {
	Data struct {
		SuccessList []struct {
			OrderID       EmptyInt `json:"orderId"`
			ClientOrderID string   `json:"clientOid"`
		} `json:"successList"`
		FailureList []struct {
			OrderID       EmptyInt `json:"orderId"`
			ClientOrderID string   `json:"clientOid"`
			ErrorCode     int64    `json:"errorCode,string"`
			ErrorMessage  string   `json:"errorMsg"`
		} `json:"failureList"`
	} `json:"data"`
}

// OrderIDStruct contains order IDs
type OrderIDStruct struct {
	OrderID       int64  `json:"orderId,string,omitempty"`
	ClientOrderID string `json:"clientOid,omitempty"`
}

// SymbolResp holds a single symbol
type SymbolResp struct {
	Data struct {
		Symbol string `json:"symbol"`
	} `json:"data"`
}

// OrderDetailTemp contains information on an order in a partially-unmarshalled state
type OrderDetailTemp struct {
	Data []struct {
		UserID           string          `json:"userId"` // Check whether this should be a different type
		Symbol           string          `json:"symbol"`
		OrderID          EmptyInt        `json:"orderId"`
		ClientOrderID    string          `json:"clientOid"`
		Price            float64         `json:"price,string"`
		Size             float64         `json:"size,string"`
		OrderType        string          `json:"orderType"`
		Side             string          `json:"side"`
		Status           string          `json:"status"`
		PriceAverage     float64         `json:"priceAvg,string"`
		BaseVolume       float64         `json:"baseVolume,string"`
		QuoteVolume      float64         `json:"quoteVolume,string"`
		EnterPointSource string          `json:"enterPointSource"`
		CreationTime     UnixTimestamp   `json:"cTime"`
		UpdateTime       UnixTimestamp   `json:"uTime"`
		OrderSource      string          `json:"orderSource"`
		FeeDetailTemp    json.RawMessage `json:"feeDetail"`
	}
}

// FeeDetail contains information on fees
type FeeDetail struct {
	AmountCoupons     float64 `json:"c"`
	AmountDeductedBGB float64 `json:"d"`
	AmountRemaining   float64 `json:"r"`
	AmountTotal       float64 `json:"t"`
	Deduction         bool    `json:"deduction"`
	FeeCoinCode       string  `json:"feeCoinCode"`
	TotalFee          float64 `json:"totalFee"`
	TotalDeductionFee float64 `json:"totalDeductionFee"`
}

// FeeDetailStore is a map of fee details for better unmarshalling
type FeeDetailStore map[string]FeeDetail

// SpotOrderDetailData contains information on an order for better unmarshalling
type SpotOrderDetailData struct {
	UserID           string // Check whether this should be a different type
	Symbol           string
	OrderID          EmptyInt
	ClientOrderID    string
	Price            float64
	Size             float64
	OrderType        string
	Side             string
	Status           string
	PriceAverage     float64
	BaseVolume       float64
	QuoteVolume      float64
	EnterPointSource string
	CreationTime     UnixTimestamp
	UpdateTime       UnixTimestamp
	OrderSource      string
	FeeDetail        FeeDetailStore
}

// SpotOrderDetailResp contains information on an order
type SpotOrderDetailResp struct {
	Data []SpotOrderDetailData
}

// UnfilledOrdersResp contains information on the user's unfilled orders
type UnfilledOrdersResp struct {
	Data []struct {
		UserID           string        `json:"userId"` // Check whether this should be a different type
		Symbol           string        `json:"symbol"`
		OrderID          EmptyInt      `json:"orderId"`
		ClientOrderID    string        `json:"clientOid"`
		PriceAverage     float64       `json:"priceAvg,string"`
		Size             float64       `json:"size,string"`
		OrderType        string        `json:"orderType"`
		Side             string        `json:"side"`
		Status           string        `json:"status"`
		BasePrice        float64       `json:"basePrice,string"`
		BaseVolume       float64       `json:"baseVolume,string"`
		QuoteVolume      float64       `json:"quoteVolume,string"`
		EnterPointSource string        `json:"enterPointSource"`
		OrderSource      string        `json:"orderSource"`
		CreationTime     UnixTimestamp `json:"cTime"`
		UpdateTime       UnixTimestamp `json:"uTime"`
	} `json:"data"`
}

// AbridgedFeeDetail contains some information on fees
type AbridgedFeeDetail struct {
	Deduction         YesNoBool    `json:"deduction"`
	FeeCoin           string       `json:"feeCoin"`
	TotalDeductionFee types.Number `json:"totalDeductionFee"`
	TotalFee          float64      `json:"totalFee,string"`
}

// SpotFillsResp contains information on the user's fulfilled orders
type SpotFillsResp struct {
	Data []struct {
		UserID       string            `json:"userId"` // Check whether this should be a different type
		Symbol       string            `json:"symbol"`
		OrderID      EmptyInt          `json:"orderId"`
		TradeID      int64             `json:"tradeId,string"`
		OrderType    string            `json:"orderType"`
		Side         string            `json:"side"`
		PriceAverage float64           `json:"priceAvg,string"`
		Size         float64           `json:"size,string"`
		Amount       float64           `json:"amount,string"`
		FeeDetail    AbridgedFeeDetail `json:"feeDetail"`
		TradeScope   string            `json:"tradeScope"`
		CreationTime UnixTimestamp     `json:"cTime"`
		UpdateTime   UnixTimestamp     `json:"uTime"`
	} `json:"data"`
}

// OrderIDResp contains order IDs in the format returned by the exchange
type OrderIDResp struct {
	Data OrderIDStruct `json:"data"`
}

// PlanSpotOrderResp contains information on plan orders
type PlanSpotOrderResp struct {
	Data struct {
		NextFlag   bool     `json:"nextFlag"`
		IDLessThan EmptyInt `json:"idLessThan"`
		OrderList  []struct {
			OrderID          int64         `json:"orderId,string"`
			ClientOrderID    string        `json:"clientOid"`
			Symbol           string        `json:"symbol"`
			TriggerPrice     float64       `json:"triggerPrice,string"`
			OrderType        string        `json:"orderType"`
			ExecutePrice     types.Number  `json:"executePrice"`
			PlanType         string        `json:"planType"`
			Size             float64       `json:"size,string"`
			Status           string        `json:"status"`
			Side             string        `json:"side"`
			TriggerType      string        `json:"triggerType"`
			EnterPointSource string        `json:"enterPointSource"`
			CreationTime     UnixTimestamp `json:"cTime"`
			UpdateTime       UnixTimestamp `json:"uTime"`
		} `json:"orderList"`
	} `json:"data"`
}

// SubOrderResp contains information on sub-orders
type SubOrderResp struct {
	Data struct {
		OrderID int64   `json:"orderId,string"`
		Price   float64 `json:"price,string"`
		Type    string  `json:"type"`
		Status  string  `json:"status"`
	} `json:"data"`
}

// AccountInfoResp contains information on the user's account
type AccountInfoResp struct {
	Data struct {
		UserID       int64         `json:"userId,string"`
		InviterID    int64         `json:"inviterId,string"`
		ChannelCode  string        `json:"channelCode"`
		Channel      string        `json:"channel"`
		IPs          string        `json:"ips"`
		Authorities  []string      `json:"authorities"`
		ParentID     int64         `json:"parentId"`
		TraderType   string        `json:"traderType"`
		RegisterTime UnixTimestamp `json:"regisTime"`
	} `json:"data"`
}

// AssetData contains information on the amount of an assset an account owns
type AssetData struct {
	Coin           string        `json:"coin"`
	Available      float64       `json:"available,string"`
	Frozen         float64       `json:"frozen,string"`
	Locked         float64       `json:"locked,string"`
	LimitAvailable float64       `json:"limitAvailable,string"`
	UpdateTime     UnixTimestamp `json:"uTime"`
}

// AccountAssetsResp contains information on the user's assets
type AccountAssetsResp struct {
	Data []AssetData `json:"data"`
}

// SubaccountAssetsResp contains information on assets in a user's sub-accounts
type SubaccountAssetsResp struct {
	Data []struct {
		UserID     int64       `json:"userId,string"`
		AssetsList []AssetData `json:"assetsList"`
	} `json:"data"`
}

// SuccessBoolResp2 contains a success bool in a secondary format returned by the exchange
type SuccessBoolResp2 struct {
	Success SuccessBool `json:"data"`
}

// SpotAccBillResp contains information on the user's billing history
type SpotAccBillResp struct {
	Data []struct {
		CreationTime UnixTimestamp `json:"cTime"`
		Coin         string        `json:"coin"`
		GroupType    string        `json:"groupType"`
		BusinessType string        `json:"businessType"`
		Size         float64       `json:"size,string"`
		Balance      float64       `json:"balance,string"`
		Fees         float64       `json:"fees,string"`
		BillID       int64         `json:"billId,string"`
	} `json:"data"`
}

// TransferResp contains information on an asset transfer
type TransferResp struct {
	Data struct {
		TransferID    int64  `json:"transferId,string"`
		ClientOrderID string `json:"clientOid"`
	} `json:"data"`
}

// TransferableCoinsResp contains a list of coins that can be transferred between the provided accounts
type TransferableCoinsResp struct {
	Data []string `json:"data"`
}

// SubaccTfrRecResp contains detailed information on asset transfers between sub-accounts
type SubaccTfrRecResp struct {
	Data []struct {
		Coin          string        `json:"coin"`
		Status        string        `json:"status"`
		ToType        string        `json:"toType"`
		FromType      string        `json:"fromType"`
		Size          float64       `json:"size,string"`
		Timestamp     UnixTimestamp `json:"ts"`
		ClientOrderID string        `json:"clientOid"`
		TransferID    int64         `json:"transferId,string"`
		FromUserID    int64         `json:"fromUserId,string"`
		ToUserID      int64         `json:"toUserId,string"`
	} `json:"data"`
}

// TransferRecResp contains detailed information on asset transfers
type TransferRecResp struct {
	Data []struct {
		Coin          string        `json:"coin"`
		Status        string        `json:"status"`
		ToType        string        `json:"toType"`
		ToSymbol      string        `json:"toSymbol"`
		FromType      string        `json:"fromType"`
		FromSymbol    string        `json:"fromSymbol"`
		Size          float64       `json:"size,string"`
		Timestamp     UnixTimestamp `json:"ts"`
		ClientOrderID string        `json:"clientOid"`
		TransferID    int64         `json:"transferId,string"`
	} `json:"data"`
}

// DepositAddressResp contains information on a deposit address
type DepositAddressResp struct {
	Data struct {
		Address string `json:"address"`
		Chain   string `json:"chain"`
		Coin    string `json:"coin"`
		Tag     string `json:"tag"`
		URL     string `json:"url"`
	} `json:"data"`
}

// SubaccDepRecResp contains detailed information on deposits to sub-accounts
type SubaccDepRecResp struct {
	Data []struct {
		OrderID      int64         `json:"orderId,string"`
		TradeID      int64         `json:"tradeId,string"`
		Coin         string        `json:"coin"`
		Size         float64       `json:"size,string"`
		Status       string        `json:"status"`
		FromAddress  string        `json:"fromAddress"`
		ToAddress    string        `json:"toAddress"`
		Chain        string        `json:"chain"`
		Destination  string        `json:"dest"`
		CreationTime UnixTimestamp `json:"cTime"`
		UpdateTime   UnixTimestamp `json:"uTime"`
	} `json:"data"`
}

// WithdrawRecordsResp contains detailed information on withdrawals
type WithdrawRecordsResp struct {
	Data []struct {
		OrderID       int64         `json:"orderId,string"`
		TradeID       int64         `json:"tradeId,string"`
		Coin          string        `json:"coin"`
		ClientOrderID string        `json:"clientOid"`
		OrderType     string        `json:"type"`
		Destination   string        `json:"dest"`
		Size          float64       `json:"size,string"`
		Fee           float64       `json:"fee,string"`
		Status        string        `json:"status"`
		FromAddress   string        `json:"fromAddress"`
		ToAddress     string        `json:"toAddress"`
		Chain         string        `json:"chain"`
		Confirm       uint32        `json:"confirm,string"`
		Tag           string        `json:"tag"`
		CreationTime  UnixTimestamp `json:"cTime"`
		UpdateTime    UnixTimestamp `json:"uTime"`
	} `json:"data"`
}

// CryptoDepRecResp contains detailed information on cryptocurrency deposits
type CryptoDepRecResp struct {
	Data []struct {
		OrderID      int64         `json:"orderId,string"`
		TradeID      int64         `json:"tradeId,string"`
		Coin         string        `json:"coin"`
		OrderType    string        `json:"type"`
		Size         float64       `json:"size,string"`
		Status       string        `json:"status"`
		FromAddress  string        `json:"fromAddress"`
		ToAddress    string        `json:"toAddress"`
		Chain        string        `json:"chain"`
		Destination  string        `json:"dest"`
		CreationTime UnixTimestamp `json:"cTime"`
		UpdateTime   UnixTimestamp `json:"uTime"`
	} `json:"data"`
}

// FutureTickerResp contains information on a futures ticker
type FutureTickerResp struct {
	Data []struct {
		Symbol            string        `json:"symbol"`
		LastPrice         float64       `json:"lastPr,string"`
		AskPrice          float64       `json:"askPr,string"`
		BidPrice          float64       `json:"bidPr,string"`
		BidSize           float64       `json:"bidSz,string"`
		AskSize           float64       `json:"askSz,string"`
		High24H           float64       `json:"high24h,string"`
		Low24H            float64       `json:"low24h,string"`
		Timestamp         UnixTimestamp `json:"ts"`
		Change24H         float64       `json:"change24h,string"`
		BaseVolume        float64       `json:"baseVolume,string"`
		QuoteVolume       float64       `json:"quoteVolume,string"`
		USDTVolume        float64       `json:"usdtVolume,string"`
		OpenUTC           float64       `json:"openUtc,string"`
		ChangeUTC24H      float64       `json:"changeUtc24h,string"`
		IndexPrice        float64       `json:"indexPrice,string"`
		FundingRate       float64       `json:"fundingRate,string"`
		HoldingAmount     float64       `json:"holdingAmount,string"`
		DeliveryStartTime UnixTimestamp `json:"deliveryStartTime"`
		DeliveryTime      UnixTimestamp `json:"deliveryTime"`
		DeliveryStatus    string        `json:"deliveryStatus"`
		Open24H           float64       `json:"open24h,string"`
	} `json:"data"`
}

// CallMode represents the call mode for the futures candlestick endpoints
type CallMode uint8

const (
	// CallModeNormal represents the normal call mode
	CallModeNormal CallMode = iota
	// CallModeHistory represents the history call mode
	CallModeHistory
	// CallModeIndex represents the historical index call mode
	CallModeIndex
	// CallModeMark represents the historical mark call mode
	CallModeMark
)

// OpenInterestResp contains information on open positions
type OpenPositionsResp struct {
	Data struct {
		OpenInterestList []struct {
			Symbol string  `json:"symbol"`
			Size   float64 `json:"size,string"`
		} `json:"openInterestList"`
		Timestamp UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// FundingTimeResp contains information on funding times
type FundingTimeResp struct {
	Data []struct {
		Symbol          string        `json:"symbol"`
		NextFundingTime UnixTimestamp `json:"nextFundingTime"`
		RatePeriod      uint16        `json:"ratePeriod,string"`
	} `json:"data"`
}

// FuturesPriceResp contains information on futures prices
type FuturesPriceResp struct {
	Data []struct {
		Symbol     string        `json:"symbol"`
		Price      float64       `json:"price,string"`
		IndexPrice float64       `json:"indexPrice,string"`
		MarkPrice  float64       `json:"markPrice,string"`
		Timestamp  UnixTimestamp `json:"ts"`
	} `json:"data"`
}

// FundingHistoryResp contains information on funding history
type FundingHistoryResp struct {
	Data []struct {
		Symbol      string        `json:"symbol"`
		FundingRate float64       `json:"fundingRate,string"`
		FundingTime UnixTimestamp `json:"fundingTime"`
	} `json:"data"`
}

// FundingCurrentResp contains information on current funding rates
type FundingCurrentResp struct {
	Data []struct {
		Symbol      string  `json:"symbol"`
		FundingRate float64 `json:"fundingRate,string"`
	} `json:"data"`
}

// ContractConfigResp contains information on contract details
type ContractConfigResp struct {
	Data []struct {
		Symbol                string        `json:"symbol"`
		BaseCoin              string        `json:"baseCoin"`
		QuoteCoin             string        `json:"quoteCoin"`
		BuyLimitPriceRatio    float64       `json:"buyLimitPriceRatio,string"`
		SellLimitPriceRatio   float64       `json:"sellLimitPriceRatio,string"`
		FeeRateUpRatio        float64       `json:"feeRateUpRatio,string"`
		MakerFeeRate          float64       `json:"makerFeeRate,string"`
		TakerFeeRate          float64       `json:"takerFeeRate,string"`
		OpenCostUpRatio       float64       `json:"openCostUpRatio,string"`
		SupportMarginCoins    []string      `json:"supportMarginCoins"`
		MinTradeNum           float64       `json:"minTradeNum,string"`
		PriceEndStep          float64       `json:"priceEndStep,string"`
		VolumePlace           float64       `json:"volumePlace,string"`
		PricePlace            float64       `json:"pricePlace,string"`
		SizeMultiplier        float64       `json:"sizeMultiplier,string"`
		SymbolType            string        `json:"symbolType"`
		MinTradeUSDT          float64       `json:"minTradeUSDT,string"`
		MaxSymbolOrderNum     int64         `json:"maxSymbolOrderNum,string"`
		MaxSymbolOpenOrderNum int64         `json:"maxSymbolOpenOrderNum,string"`
		MaxPositionNum        int64         `json:"maxPositionNum,string"`
		SymbolStatus          string        `json:"symbolStatus"`
		OffTime               int64         `json:"offTime,string"`
		LimitOpenTime         int64         `json:"limitOpenTime,string"`
		DeliveryTime          EmptyInt      `json:"deliveryTime"`
		DeliveryStartTime     EmptyInt      `json:"deliveryStartTime"`
		DeliveryPeriod        EmptyInt      `json:"deliveryPeriod"`
		LaunchTime            EmptyInt      `json:"launchTime"`
		FundInterval          uint16        `json:"fundInterval,string"`
		MinLever              float64       `json:"minLever,string"`
		MaxLever              float64       `json:"maxLever,string"`
		PosLimit              float64       `json:"posLimit,string"`
		MaintainTime          UnixTimestamp `json:"maintainTime"`
	} `json:"data"`
}

// OneAccResp contains information on a single account
type OneAccResp struct {
	Data struct {
		MarginCoin            string       `json:"marginCoin"`
		Locked                float64      `json:"locked,string"`
		Available             float64      `json:"available,string"`
		CrossedMaxAvailable   float64      `json:"crossedMaxAvailable,string"`
		IsolatedMaxAvailable  float64      `json:"isolatedMaxAvailable,string"`
		MaxTransferOut        float64      `json:"maxTransferOut,string"`
		AccountEquity         float64      `json:"accountEquity,string"`
		USDTEquity            float64      `json:"usdtEquity,string"`
		BTCEquity             float64      `json:"btcEquity,string"`
		CrossedRiskRate       float64      `json:"crossedRiskRate,string"`
		CrossedMarginleverage float64      `json:"crossedMarginleverage"`
		IsolatedLongLever     float64      `json:"isolatedLongLever"`
		IsolatedShortLever    float64      `json:"isolatedShortLever"`
		MarginMode            string       `json:"marginMode"`
		PositionMode          string       `json:"posMode"`
		UnrealizedPL          types.Number `json:"unrealizedPL"`
		Coupon                types.Number `json:"coupon,string"`
		CrossedUnrealizedPL   types.Number `json:"crossedUnrealizedPL"`
		IsolatedUnrealizedPL  types.Number `json:"isolatedUnrealizedPL"`
	} `json:"data"`
}

// FutureAccDetails contains information on a user's futures account
type FutureAccDetails struct {
	MarginCoin           string       `json:"marginCoin"`
	Locked               float64      `json:"locked,string"`
	Available            float64      `json:"available,string"`
	CrossedMaxAvailable  float64      `json:"crossedMaxAvailable,string"`
	IsolatedMaxAvailable float64      `json:"isolatedMaxAvailable,string"`
	MaxTransferOut       float64      `json:"maxTransferOut,string"`
	AccountEquity        float64      `json:"accountEquity,string"`
	USDTEquity           float64      `json:"usdtEquity,string"`
	BTCEquity            float64      `json:"btcEquity,string"`
	CrossedRiskRate      float64      `json:"crossedRiskRate,string"`
	UnrealizedPL         types.Number `json:"unrealizedPL"`
	Coupon               types.Number `json:"coupon"`
	CrossedUnrealizedPL  types.Number `json:"crossedUnrealizedPL"`
	IsolatedUnrealizedPL types.Number `json:"isolatedUnrealizedPL"`
}

// AllAccResp contains information on all accounts
type AllAccResp struct {
	Data []FutureAccDetails `json:"data"`
}

// SubaccountFuturesResp contains information on futures details of a user's sub-accounts
type SubaccountFuturesResp struct {
	Data []struct {
		UserID    int64              `json:"userId"`
		AssetList []FutureAccDetails `json:"assetList"`
	} `json:"data"`
}

// EstOpenCountResp contains information on the estimated size of open orders
type EstOpenCountResp struct {
	Data struct {
		Size float64 `json:"size,string"`
	} `json:"data"`
}

// LeverageResp contains information on the leverage of a position
type LeverageResp struct {
	Data struct {
		Symbol              string       `json:"symbol"`
		MarginCoin          string       `json:"marginCoin"`
		LongLeverage        float64      `json:"longLeverage,string"`
		ShortLeverage       float64      `json:"shortLeverage,string"`
		CrossMarginLeverage types.Number `json:"crossMarginLeverage"`
		MarginMode          string       `json:"marginMode"`
	} `json:"data"`
}

// PosModeResp contains information on the position mode
type PosModeResp struct {
	Data struct {
		PositionMode string `json:"posMode"`
	} `json:"data"`
}

// FutureAccBillResp contains information on futures billing history
type FutureAccBillResp struct {
	Data struct {
		Bills []struct {
			OrderID      int64         `json:"orderId,string"`
			Symbol       string        `json:"symbol"`
			Amount       float64       `json:"amount,string"`
			Fee          float64       `json:"fee,string"`
			FeeByCoupon  types.Number  `json:"feeByCoupon"`
			FeeCoin      string        `json:"feeCoin"`
			BusinessType string        `json:"businessType"`
			Coin         string        `json:"coin"`
			CreationTime UnixTimestamp `json:"cTime"`
		} `json:"bills"`
	} `json:"data"`
}

// PositionTierResp contains information on position configurations
type PositionTierResp struct {
	Data []struct {
		Symbol         string  `json:"symbol"`
		Level          uint8   `json:"level,string"`
		StartUnit      float64 `json:"startUnit,string"`
		EndUnit        float64 `json:"endUnit,string"`
		Leverage       float64 `json:"leverage,string"`
		KeepMarginRate float64 `json:"keepMarginRate,string"`
	} `json:"data"`
}

// PositionResp contains information on positions
type PositionResp struct {
	Data []struct {
		MarginCoin       string        `json:"marginCoin"`
		Symbol           string        `json:"symbol"`
		HoldSide         string        `json:"holdSide"`
		OpenDelegateSize float64       `json:"openDelegateSize,string"`
		MarginSize       float64       `json:"marginSize,string"`
		Available        float64       `json:"available,string"`
		Locked           float64       `json:"locked,string"`
		Total            float64       `json:"total,string"`
		Leverage         float64       `json:"leverage,string"`
		AchievedProfits  float64       `json:"achievedProfits,string"`
		OpenPriceAverage float64       `json:"openPriceAvg,string"`
		MarginMode       string        `json:"marginMode"`
		PositionMode     string        `json:"posMode"`
		UnrealizedPL     float64       `json:"unrealizedPL,string"`
		LiquidationPrice float64       `json:"liquidationPrice,string"`
		KeepMarginRate   float64       `json:"keepMarginRate,string"`
		MarkPrice        float64       `json:"markPrice,string"`
		MarginRatio      float64       `json:"marginRatio,string"`
		CreationTime     UnixTimestamp `json:"cTime"`
	} `json:"data"`
}

// HistPositionResp contains information on historical positions
type HistPositionResp struct {
	Data struct {
		List []struct {
			MarginCoin         string        `json:"marginCoin"`
			Symbol             string        `json:"symbol"`
			HoldSide           string        `json:"holdSide"`
			OpenAveragePrice   float64       `json:"openAvgPrice,string"`
			CloseAveragePrice  float64       `json:"closeAvgPrice,string"`
			MarginMode         string        `json:"marginMode"`
			OpenTotalPosition  float64       `json:"openTotalPos,string"`
			CloseTotalPosition float64       `json:"closeTotalPos,string"`
			PNL                float64       `json:"pnl,string"`
			NetProfit          float64       `json:"netProfit,string"`
			TotalFunding       float64       `json:"totalFunding,string"`
			OpenFee            float64       `json:"openFee,string"`
			CloseFee           float64       `json:"closeFee,string"`
			UpdateTime         UnixTimestamp `json:"uTime"`
			CreationTime       UnixTimestamp `json:"cTime"`
		} `json:"list"`
		EndID int64 `json:"endId,string"`
	} `json:"data"`
}

// PlaceFuturesOrderStruct contains information on an order to be placed
type PlaceFuturesOrderStruct struct {
	Size            float64   `json:"size,string"`
	Price           float64   `json:"price,string"`
	Side            string    `json:"side"`
	TradeSide       string    `json:"tradeSide"`
	OrderType       string    `json:"orderType"`
	Strategy        string    `json:"force"`
	ClientOID       string    `json:"clientOId"`
	ReduceOnly      YesNoBool `json:"reduceOnly"`
	TakeProfitValue float64   `json:"presetStopSurplusPrice,string,omitempty"`
	StopLossValue   float64   `json:"presetStopLossPrice,string,omitempty"`
}

// FuturesOrderDetailResp contains information on a futures order
type FuturesOrderDetailResp struct {
	Symbol                 string        `json:"symbol"`
	Size                   float64       `json:"size,string"`
	OrderID                EmptyInt      `json:"orderId"`
	ClientOrderID          string        `json:"clientOid"`
	BaseVolume             float64       `json:"baseVolume,string"`
	PriceAverage           float64       `json:"priceAvg,string"`
	Fee                    types.Number  `json:"fee"`
	Price                  float64       `json:"price,string"`
	State                  string        `json:"state"`
	Side                   string        `json:"side"`
	Force                  string        `json:"force"`
	TotalProfits           float64       `json:"totalProfits,string"`
	PositionSide           string        `json:"posSide"`
	MarginCoin             string        `json:"marginCoin"`
	PresetStopSurplusPrice float64       `json:"presetStopSurplusPrice,string"`
	PresetStopLossPrice    float64       `json:"presetStopLossPrice,string"`
	QuoteVolume            float64       `json:"quoteVolume,string"`
	OrderType              string        `json:"orderType"`
	Leverage               float64       `json:"leverage,string"`
	MarginMode             string        `json:"marginMode"`
	ReduceOnly             YesNoBool     `json:"reduceOnly"`
	EnterPointSource       string        `json:"enterPointSource"`
	TradeSide              string        `json:"tradeSide"`
	PositionMode           string        `json:"posMode"`
	OrderSource            string        `json:"orderSource"`
	CreationTime           UnixTimestamp `json:"cTime"`
	UpdateTime             UnixTimestamp `json:"uTime"`
}

// FuturesFillsResp contains information on fulfilled futures orders
type FuturesFillsResp struct {
	Data struct {
		FillList []struct {
			TradeID          int64               `json:"tradeId,string"`
			Symbol           string              `json:"symbol"`
			OrderID          int64               `json:"orderId,string"`
			Price            float64             `json:"price,string"`
			BaseVolume       float64             `json:"baseVolume,string"`
			FeeDetail        []AbridgedFeeDetail `json:"feeDetail"`
			Side             string              `json:"side"`
			QuoteVolume      float64             `json:"quoteVolume,string"`
			Profit           float64             `json:"profit,string"`
			EnterPointSource string              `json:"enterPointSource"`
			TradeSide        string              `json:"tradeSide"`
			PositionMode     string              `json:"posMode"`
			TradeScope       string              `json:"tradeScope"`
			CreationTime     UnixTimestamp       `json:"cTime"`
		} `json:"fillList"`
		EndID EmptyInt `json:"endId"`
	} `json:"data"`
}

// FuturesOrdResp contains information on futures orders
type FuturesOrdResp struct {
	Data struct {
		EntrustedList []struct {
			Symbol                 string        `json:"symbol"`
			Size                   float64       `json:"size,string"`
			OrderID                int64         `json:"orderId,string"`
			ClientOrderID          string        `json:"clientOid"`
			BaseVolume             float64       `json:"baseVolume,string"`
			Fee                    types.Number  `json:"fee"`
			Price                  types.Number  `json:"price"`
			PriceAverage           types.Number  `json:"priceAvg"`
			Status                 string        `json:"status"`
			Side                   string        `json:"side"`
			Force                  string        `json:"force"`
			TotalProfits           float64       `json:"totalProfits,string"`
			PositionSide           string        `json:"posSide"`
			MarginCoin             string        `json:"marginCoin"`
			QuoteVolume            float64       `json:"quoteVolume,string"`
			Leverage               float64       `json:"leverage,string"`
			MarginMode             string        `json:"marginMode"`
			EnterPointSource       string        `json:"enterPointSource"`
			TradeSide              string        `json:"tradeSide"`
			PositionMode           string        `json:"posMode"`
			OrderType              string        `json:"orderType"`
			OrderSource            string        `json:"orderSource"`
			CreationTime           UnixTimestamp `json:"cTime"`
			UpdateTime             UnixTimestamp `json:"uTime"`
			PresetStopSurplusPrice types.Number  `json:"presetStopSurplusPrice"`
			PresetStopLossPrice    types.Number  `json:"presetStopLossPrice"`
		} `json:"entrustedList"`
		EndID EmptyInt `json:"endId"`
	} `json:"data"`
}

// PlanFuturesOrdResp contains information on planned futures orders
type PlanFuturesOrdResp struct {
	Data struct {
		EntrustedList []struct {
			PlanType               string        `json:"planType"`
			Symbol                 string        `json:"symbol"`
			Size                   float64       `json:"size,string"`
			OrderID                int64         `json:"orderId,string"`
			ClientOrderID          string        `json:"clientOid"`
			Price                  float64       `json:"price,string"`
			CallbackRatio          types.Number  `json:"callbackRatio"`
			TriggerPrice           float64       `json:"triggerPrice,string"`
			TriggerType            string        `json:"triggerType"`
			PlanStatus             string        `json:"planStatus"`
			Side                   string        `json:"side"`
			PositionSide           string        `json:"posSide"`
			MarginCoin             string        `json:"marginCoin"`
			MarginMode             string        `json:"marginMode"`
			EnterPointSource       string        `json:"enterPointSource"`
			TradeSide              string        `json:"tradeSide"`
			PositionMode           string        `json:"posMode"`
			OrderType              string        `json:"orderType"`
			OrderSource            string        `json:"orderSource"`
			CreationTime           UnixTimestamp `json:"cTime"`
			UpdateTime             UnixTimestamp `json:"uTime"`
			PresetTakeProfitPrice  types.Number  `json:"presetStopSurplusPrice"`
			TakeprofitTriggerPrice types.Number  `json:"stopSurplusTriggerPrice"`
			TakeProfitTriggerType  string        `json:"stopSurplusTriggerType"`
			PresetStopLossPrice    types.Number  `json:"presetStopLossPrice"`
			StopLossTriggerPrice   types.Number  `json:"stopLossTriggerPrice"`
			StopLossTriggerType    string        `json:"stopLossTriggerType"`
		} `json:"entrustedList"`
		EndID EmptyInt `json:"endId"`
	} `json:"data"`
}

// HistTriggerFuturesOrdResp contains information on historical trigger futures orders
type HistTriggerFuturesOrdResp struct {
	Data struct {
		EntrustedList []struct {
			PlanType               string        `json:"planType"`
			Symbol                 string        `json:"symbol"`
			Size                   float64       `json:"size,string"`
			OrderID                int64         `json:"orderId,string"`
			ExecuteOrderID         int64         `json:"executeOrderId,string"`
			ClientOrderID          string        `json:"clientOid"`
			PlanStatus             string        `json:"planStatus"`
			Price                  float64       `json:"price,string"`
			PriceAverage           float64       `json:"priceAvg,string"`
			BaseVolume             float64       `json:"baseVolume,string"`
			CallbackRatio          types.Number  `json:"callbackRatio"`
			TriggerPrice           float64       `json:"triggerPrice,string"`
			TriggerType            string        `json:"triggerType"`
			Side                   string        `json:"side"`
			PositionSide           string        `json:"posSide"`
			MarginCoin             string        `json:"marginCoin"`
			MarginMode             string        `json:"marginMode"`
			EnterPointSource       string        `json:"enterPointSource"`
			TradeSide              string        `json:"tradeSide"`
			PositionMode           string        `json:"posMode"`
			OrderType              string        `json:"orderType"`
			CreationTime           UnixTimestamp `json:"cTime"`
			UpdateTime             UnixTimestamp `json:"uTime"`
			PresetTakeProfitPrice  types.Number  `json:"presetStopSurplusPrice"`
			TakeprofitTriggerPrice types.Number  `json:"stopSurplusTriggerPrice"`
			TakeProfitTriggerType  string        `json:"stopSurplusTriggerType"`
			PresetStopLossPrice    types.Number  `json:"presetStopLossPrice"`
			StopLossTriggerPrice   types.Number  `json:"stopLossTriggerPrice"`
			StopLossTriggerType    string        `json:"stopLossTriggerType"`
		} `json:"entrustedList"`
		EndID EmptyInt `json:"endId"`
	} `json:"data"`
}

// SupCurrencyResp contains information on supported currencies
type SupCurrencyResp struct {
	Data []struct {
		Symbol                    string  `json:"symbol"`
		BaseCoin                  string  `json:"baseCoin"`
		QuoteCoin                 string  `json:"quoteCoin"`
		MaxCrossedLeverage        float64 `json:"maxCrossedLeverage,string"`
		MaxIsolatedLeverage       float64 `json:"maxIsolatedLeverage,string"`
		WarningRiskRatio          float64 `json:"warningRiskRatio,string"`
		LiquidationRiskRatio      float64 `json:"liquidationRiskRatio,string"`
		MinTradeAmount            float64 `json:"minTradeAmount,string"`
		MaxTradeAmount            float64 `json:"maxTradeAmount,string"`
		TakerFeeRate              float64 `json:"takerFeeRate,string"`
		MakerFeeRate              float64 `json:"makerFeeRate,string"`
		PricePrecision            uint8   `json:"pricePrecision,string"`
		QuantityPrecision         uint8   `json:"quantityPrecision,string"`
		MinTradeUSDT              float64 `json:"minTradeUSDT,string"`
		IsBorrowable              bool    `json:"isBorrowable"`
		UserMinBorrow             float64 `json:"userMinBorrow,string"`
		Status                    string  `json:"status"`
		IsIsolatedBaseBorrowable  bool    `json:"isIsolatedBaseBorrowable"`
		IsIsolatedQuoteBorrowable bool    `json:"isIsolatedQuoteBorrowable"`
		IsCrossBorrowable         bool    `json:"isCrossBorrowable"`
	} `json:"data"`
}

// BorrowHistResp contains information on borrowing history
type BorrowHistResp struct {
	Data struct {
		ResultList []struct {
			LoanID       int64         `json:"loanId,string"`
			Coin         string        `json:"coin"`
			BorrowAmount float64       `json:"borrowAmount,string"`
			BorrowType   string        `json:"borrowType"`
			CreationTime UnixTimestamp `json:"cTime"`
			UpdateTime   UnixTimestamp `json:"uTime"`
		} `json:"resultList"`
		MaxID EmptyInt `json:"maxId"`
		MinID EmptyInt `json:"minId"`
	} `json:"data"`
}
