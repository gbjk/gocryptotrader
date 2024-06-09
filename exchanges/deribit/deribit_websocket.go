package deribit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/common/crypto"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/nonce"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream"
	"github.com/thrasher-corp/gocryptotrader/exchanges/subscription"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
	"github.com/thrasher-corp/gocryptotrader/exchanges/trade"
	"github.com/thrasher-corp/gocryptotrader/log"
)

var deribitWebsocketAddress = "wss://www.deribit.com/ws" + deribitAPIVersion

const (
	rpcVersion    = "2.0"
	rateLimit     = 20
	errAuthFailed = 1002

	// public websocket channels
	tickerChannel                          = "ticker.{{.Pair}}.{{.Interval}}"
	orderbookChannel                       = "book.{{.Pair}}.{{.Interval}}" // book.{instrument_name}.{interval}
	chartTradesChannel                     = "chart.trades.%s.%s"           // chart.trades.{instrument_name}.{resolution} (Candles)
	tradesChannel                          = "trades.%s.%s.%s"              // trades.{kind}.{currency}.{interval}
	announcementsChannel                   = "announcements"
	priceIndexChannel                      = "deribit_price_index"
	priceRankingChannel                    = "deribit_price_ranking"
	priceStatisticsChannel                 = "deribit_price_statistics"
	volatilityIndexChannel                 = "deribit_volatility_index"
	estimatedExpirationPriceChannel        = "estimated_expiration_price"
	incrementalTickerChannel               = "incremental_ticker"
	instrumentStateChannel                 = "instrument.state"
	markPriceOptionsChannel                = "markprice.options"
	perpetualChannel                       = "perpetual."
	platformStateChannel                   = "platform_state"
	platformStatePublicMethodsStateChannel = "platform_state.public_methods_state"
	quoteChannel                           = "quote"
	requestForQuoteChannel                 = "rfq"

	// private websocket channels
	userAccessLogChannel          = "user.access_log"
	userChangesInstrumentsChannel = "user.changes."
	userChangesCurrencyChannel    = "user.changes"
	userLockChannel               = "user.lock"
	userMMPTriggerChannel         = "user.mmp_trigger"
	userOrdersChannel             = "user.orders.any.any.%s" // user.orders.{kind}.{currency}.{interval}
	userTradesChannel             = "user.trades.any.any.%s" // user.trades.{kind}.{currency}.{interval}
	userPortfolioChannel          = "user.portfolio"
)

var subscriptionNames = map[string]string{
	subscription.TickerChannel:    tickerChannel,
	subscription.OrderbookChannel: orderbookChannel,
	subscription.CandlesChannel:   chartTradesChannel,
	subscription.AllTradesChannel: tradesChannel,
	subscription.MyTradesChannel:  userTradesChannel,
}

var (
	indexENUMS = []string{"ada_usd", "algo_usd", "avax_usd", "bch_usd", "bnb_usd", "btc_usd", "doge_usd", "dot_usd", "eth_usd", "link_usd", "ltc_usd", "luna_usd", "matic_usd", "near_usd", "shib_usd", "sol_usd", "trx_usd", "uni_usd", "usdc_usd", "xrp_usd", "ada_usdc", "bch_usdc", "algo_usdc", "avax_usdc", "btc_usdc", "doge_usdc", "dot_usdc", "bch_usdc", "bnb_usdc", "eth_usdc", "link_usdc", "ltc_usdc", "luna_usdc", "matic_usdc", "near_usdc", "shib_usdc", "sol_usdc", "trx_usdc", "uni_usdc", "xrp_usdc", "btcdvol_usdc", "ethdvol_usdc"}

	pingMessage = WsSubscriptionInput{
		ID:             2,
		JSONRPCVersion: rpcVersion,
		Method:         "public/test",
		Params:         map[string][]string{},
	}
	setHeartBeatMessage = wsInput{
		ID:             1,
		JSONRPCVersion: rpcVersion,
		Method:         "public/set_heartbeat",
		Params: map[string]interface{}{
			"interval": 15,
		},
	}
)

// WsConnect starts a new connection with the websocket API
func (d *Deribit) WsConnect() error {
	if !d.Websocket.IsEnabled() || !d.IsEnabled() {
		return stream.ErrWebsocketNotEnabled
	}
	var dialer websocket.Dialer
	err := d.Websocket.Conn.Dial(&dialer, http.Header{})
	if err != nil {
		return err
	}
	d.Websocket.Wg.Add(1)
	go d.wsReadData()
	if d.Websocket.CanUseAuthenticatedEndpoints() {
		err = d.wsLogin(context.TODO())
		if err != nil {
			log.Errorf(log.ExchangeSys, "%v - authentication failed: %v\n", d.Name, err)
			d.Websocket.SetCanUseAuthenticatedEndpoints(false)
		}
	}
	return d.Websocket.Conn.SendJSONMessage(setHeartBeatMessage)
}

func (d *Deribit) wsLogin(ctx context.Context) error {
	if !d.IsWebsocketAuthenticationSupported() {
		return fmt.Errorf("%v AuthenticatedWebsocketAPISupport not enabled", d.Name)
	}
	creds, err := d.GetCredentials(ctx)
	if err != nil {
		return err
	}
	d.Websocket.SetCanUseAuthenticatedEndpoints(true)
	n := d.Requester.GetNonce(nonce.UnixNano).String()
	strTS := strconv.FormatInt(time.Now().UnixMilli(), 10)
	str2Sign := strTS + "\n" + n + "\n"
	hmac, err := crypto.GetHMAC(crypto.HashSHA256,
		[]byte(str2Sign),
		[]byte(creds.Secret))
	if err != nil {
		return err
	}

	request := wsInput{
		JSONRPCVersion: rpcVersion,
		Method:         "public/auth",
		ID:             d.Websocket.Conn.GenerateMessageID(false),
		Params: map[string]interface{}{
			"grant_type": "client_signature",
			"client_id":  creds.Key,
			"timestamp":  strTS,
			"nonce":      n,
			"signature":  crypto.HexEncodeToString(hmac),
		},
	}
	resp, err := d.Websocket.Conn.SendMessageReturnResponse(request.ID, request)
	if err != nil {
		d.Websocket.SetCanUseAuthenticatedEndpoints(false)
		return err
	}
	var response wsLoginResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return fmt.Errorf("%v %v", d.Name, err)
	}
	if response.Error != nil && (response.Error.Code > 0 || response.Error.Message != "") {
		return fmt.Errorf("%v Error:%v Message:%v", d.Name, response.Error.Code, response.Error.Message)
	}
	return nil
}

// wsReadData receives and passes on websocket messages for processing
func (d *Deribit) wsReadData() {
	defer d.Websocket.Wg.Done()

	for {
		resp := d.Websocket.Conn.ReadMessage()
		if resp.Raw == nil {
			return
		}

		err := d.wsHandleData(resp.Raw)
		if err != nil {
			d.Websocket.DataHandler <- err
		}
	}
}

func (d *Deribit) wsHandleData(respRaw []byte) error {
	var response WsResponse
	err := json.Unmarshal(respRaw, &response)
	if err != nil {
		return fmt.Errorf("%s - err %s could not parse websocket data: %s", d.Name, err, respRaw)
	}
	if response.Method == "heartbeat" {
		return d.Websocket.Conn.SendJSONMessage(pingMessage)
	}
	if response.ID > 2 {
		if !d.Websocket.Match.IncomingWithData(response.ID, respRaw) {
			return fmt.Errorf("can't send ws incoming data to Matched channel with RequestID: %d", response.ID)
		}
		return nil
	} else if response.ID > 0 {
		return nil
	}
	channels := strings.Split(response.Params.Channel, ".")
	switch channels[0] {
	case "announcements":
		announcement := &Announcement{}
		response.Params.Data = announcement
		err = json.Unmarshal(respRaw, &response)
		if err != nil {
			return err
		}
		d.Websocket.DataHandler <- announcement
	case "book":
		return d.processOrderbook(respRaw, channels)
	case "chart":
		return d.processCandleChart(respRaw, channels)
	case "deribit_price_index":
		indexPrice := &wsIndexPrice{}
		return d.processData(respRaw, indexPrice)
	case "deribit_price_ranking":
		priceRankings := &wsRankingPrices{}
		return d.processData(respRaw, priceRankings)
	case "deribit_price_statistics":
		priceStatistics := &wsPriceStatistics{}
		return d.processData(respRaw, priceStatistics)
	case "deribit_volatility_index":
		volatilityIndex := &wsVolatilityIndex{}
		return d.processData(respRaw, volatilityIndex)
	case "estimated_expiration_price":
		estimatedExpirationPrice := &wsEstimatedExpirationPrice{}
		return d.processData(respRaw, estimatedExpirationPrice)
	case "incremental_ticker":
		return d.processIncrementalTicker(respRaw, channels)
	case "instrument":
		instrumentState := &wsInstrumentState{}
		return d.processData(respRaw, instrumentState)
	case "markprice":
		markPriceOptions := []wsMarkPriceOptions{}
		return d.processData(respRaw, markPriceOptions)
	case "perpetual":
		perpetualInterest := &wsPerpetualInterest{}
		return d.processData(respRaw, perpetualInterest)
	case platformStateChannel:
		platformState := &wsPlatformState{}
		return d.processData(respRaw, platformState)
	case "quote": // Quote ticker information.
		return d.processQuoteTicker(respRaw, channels)
	case "rfq":
		rfq := &wsRequestForQuote{}
		return d.processData(respRaw, rfq)
	case "ticker":
		return d.processInstrumentTicker(respRaw, channels)
	case "trades":
		return d.processTrades(respRaw, channels)
	case "user":
		switch channels[1] {
		case "access_log":
			accessLog := &wsAccessLog{}
			return d.processData(respRaw, accessLog)
		case "changes":
			return d.processUserOrderChanges(respRaw, channels)
		case "lock":
			userLock := &WsUserLock{}
			return d.processData(respRaw, userLock)
		case "mmp_trigger":
			data := &WsMMPTrigger{
				Currency: channels[2],
			}
			return d.processData(respRaw, data)
		case "orders":
			return d.processUserOrders(respRaw, channels)
		case "portfolio":
			portfolio := &wsUserPortfolio{}
			return d.processData(respRaw, portfolio)
		case "trades":
			return d.processTrades(respRaw, channels)
		default:
			d.Websocket.DataHandler <- stream.UnhandledMessageWarning{
				Message: d.Name + stream.UnhandledMessage + string(respRaw),
			}
			return nil
		}
	case "public/test", "public/set_heartbeat":
	default:
		switch result := response.Result.(type) {
		case string:
			if result == "ok" {
				return nil
			}
		default:
			d.Websocket.DataHandler <- stream.UnhandledMessageWarning{
				Message: d.Name + stream.UnhandledMessage + string(respRaw),
			}
			return nil
		}
	}
	return nil
}

func (d *Deribit) processUserOrders(respRaw []byte, channels []string) error {
	if len(channels) != 4 && len(channels) != 5 {
		return fmt.Errorf("%w, expected format 'user.orders.{instrument_name}.raw, user.orders.{instrument_name}.{interval}, user.orders.{kind}.{currency}.raw, or user.orders.{kind}.{currency}.{interval}', but found %s", errMalformedData, strings.Join(channels, "."))
	}
	var response WsResponse
	orderData := []WsOrder{}
	response.Params.Data = orderData
	err := json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	orderDetails := make([]order.Detail, len(orderData))
	for x := range orderData {
		cp, a, err := d.getAssetPairByInstrument(orderData[x].InstrumentName)
		if err != nil {
			return err
		}
		oType, err := order.StringToOrderType(orderData[x].OrderType)
		if err != nil {
			return err
		}
		side, err := order.StringToOrderSide(orderData[x].Direction)
		if err != nil {
			return err
		}
		status, err := order.StringToOrderStatus(orderData[x].OrderState)
		if err != nil {
			return err
		}
		orderDetails[x] = order.Detail{
			Price:           orderData[x].Price,
			Amount:          orderData[x].Amount,
			ExecutedAmount:  orderData[x].FilledAmount,
			RemainingAmount: orderData[x].Amount - orderData[x].FilledAmount,
			Exchange:        d.Name,
			OrderID:         orderData[x].OrderID,
			Type:            oType,
			Side:            side,
			Status:          status,
			AssetType:       a,
			Date:            orderData[x].CreationTimestamp.Time(),
			LastUpdated:     orderData[x].LastUpdateTimestamp.Time(),
			Pair:            cp,
		}
	}
	d.Websocket.DataHandler <- orderDetails
	return nil
}

func (d *Deribit) processUserOrderChanges(respRaw []byte, channels []string) error {
	if len(channels) < 4 || len(channels) > 5 {
		return fmt.Errorf("%w, expected format 'trades.{instrument_name}.{interval} or trades.{kind}.{currency}.{interval}', but found %s", errMalformedData, strings.Join(channels, "."))
	}
	var response WsResponse
	changeData := &wsChanges{}
	response.Params.Data = changeData
	err := json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	td := make([]trade.Data, len(changeData.Trades))
	for x := range changeData.Trades {
		var side order.Side
		side, err = order.StringToOrderSide(changeData.Trades[x].Direction)
		if err != nil {
			return err
		}
		var cp currency.Pair
		var a asset.Item
		cp, a, err = d.getAssetPairByInstrument(changeData.Trades[x].InstrumentName)
		if err != nil {
			return err
		}

		td[x] = trade.Data{
			CurrencyPair: cp,
			Exchange:     d.Name,
			Timestamp:    changeData.Trades[x].Timestamp.Time(),
			Price:        changeData.Trades[x].Price,
			Amount:       changeData.Trades[x].Amount,
			Side:         side,
			TID:          changeData.Trades[x].TradeID,
			AssetType:    a,
		}
	}
	err = trade.AddTradesToBuffer(d.Name, td...)
	if err != nil {
		return err
	}
	orders := make([]order.Detail, len(changeData.Orders))
	for x := range orders {
		oType, err := order.StringToOrderType(changeData.Orders[x].OrderType)
		if err != nil {
			return err
		}
		side, err := order.StringToOrderSide(changeData.Orders[x].Direction)
		if err != nil {
			return err
		}
		status, err := order.StringToOrderStatus(changeData.Orders[x].OrderState)
		if err != nil {
			return err
		}
		cp, a, err := d.getAssetPairByInstrument(changeData.Orders[x].InstrumentName)
		if err != nil {
			return err
		}
		orders[x] = order.Detail{
			Price:           changeData.Orders[x].Price,
			Amount:          changeData.Orders[x].Amount,
			ExecutedAmount:  changeData.Orders[x].FilledAmount,
			RemainingAmount: changeData.Orders[x].Amount - changeData.Orders[x].FilledAmount,
			Exchange:        d.Name,
			OrderID:         changeData.Orders[x].OrderID,
			Type:            oType,
			Side:            side,
			Status:          status,
			AssetType:       a,
			Date:            changeData.Orders[x].CreationTimestamp.Time(),
			LastUpdated:     changeData.Orders[x].LastUpdateTimestamp.Time(),
			Pair:            cp,
		}
	}
	d.Websocket.DataHandler <- orders
	d.Websocket.DataHandler <- changeData.Positions
	return nil
}

func (d *Deribit) processQuoteTicker(respRaw []byte, channels []string) error {
	cp, a, err := d.getAssetPairByInstrument(channels[1])
	if err != nil {
		return err
	}
	var response WsResponse
	quoteTicker := &wsQuoteTickerInformation{}
	response.Params.Data = quoteTicker
	err = json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	d.Websocket.DataHandler <- &ticker.Price{
		ExchangeName: d.Name,
		Pair:         cp,
		AssetType:    a,
		LastUpdated:  quoteTicker.Timestamp.Time(),
		Bid:          quoteTicker.BestBidPrice,
		Ask:          quoteTicker.BestAskPrice,
		BidSize:      quoteTicker.BestBidAmount,
		AskSize:      quoteTicker.BestAskAmount,
	}
	return nil
}

func (d *Deribit) processTrades(respRaw []byte, channels []string) error {
	if len(channels) < 3 || len(channels) > 5 {
		return fmt.Errorf("%w, expected format 'trades.{instrument_name}.{interval} or trades.{kind}.{currency}.{interval}', but found %s", errMalformedData, strings.Join(channels, "."))
	}
	var response WsResponse
	var tradeList []wsTrade
	response.Params.Data = &tradeList
	err := json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	if len(tradeList) == 0 {
		return fmt.Errorf("%v, empty list of trades found", common.ErrNoResponse)
	}
	tradeDatas := make([]trade.Data, len(tradeList))
	for x := range tradeDatas {
		var cp currency.Pair
		var a asset.Item
		cp, a, err = d.getAssetPairByInstrument(tradeList[x].InstrumentName)
		if err != nil {
			return err
		}
		side, err := order.StringToOrderSide(tradeList[x].Direction)
		if err != nil {
			return err
		}
		tradeDatas[x] = trade.Data{
			CurrencyPair: cp,
			Exchange:     d.Name,
			Timestamp:    tradeList[x].Timestamp.Time(),
			Price:        tradeList[x].Price,
			Amount:       tradeList[x].Amount,
			Side:         side,
			TID:          tradeList[x].TradeID,
			AssetType:    a,
		}
	}
	return trade.AddTradesToBuffer(d.Name, tradeDatas...)
}

func (d *Deribit) processIncrementalTicker(respRaw []byte, channels []string) error {
	if len(channels) != 2 {
		return fmt.Errorf("%w, expected format 'incremental_ticker.{instrument_name}', but found %s", errMalformedData, strings.Join(channels, "."))
	}
	cp, a, err := d.getAssetPairByInstrument(channels[1])
	if err != nil {
		return err
	}
	var response WsResponse
	incrementalTicker := &WsIncrementalTicker{}
	response.Params.Data = incrementalTicker
	err = json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	d.Websocket.DataHandler <- &ticker.Price{
		ExchangeName: d.Name,
		Pair:         cp,
		AssetType:    a,
		LastUpdated:  incrementalTicker.Timestamp.Time(),
		BidSize:      incrementalTicker.BestBidAmount,
		AskSize:      incrementalTicker.BestAskAmount,
		High:         incrementalTicker.MaxPrice,
		Low:          incrementalTicker.MinPrice,
		Volume:       incrementalTicker.Stats.Volume,
		QuoteVolume:  incrementalTicker.Stats.VolumeUsd,
		Ask:          incrementalTicker.ImpliedAsk,
		Bid:          incrementalTicker.ImpliedBid,
	}
	return nil
}

func (d *Deribit) processInstrumentTicker(respRaw []byte, channels []string) error {
	if len(channels) != 3 {
		return fmt.Errorf("%w, expected format 'ticker.{instrument_name}.{interval}', but found %s", errMalformedData, strings.Join(channels, "."))
	}
	return d.processTicker(respRaw, channels)
}

func (d *Deribit) processTicker(respRaw []byte, channels []string) error {
	cp, a, err := d.getAssetPairByInstrument(channels[1])
	if err != nil {
		return err
	}
	var response WsResponse
	tickerPriceResponse := &wsTicker{}
	response.Params.Data = tickerPriceResponse
	err = json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	tickerPrice := &ticker.Price{
		ExchangeName: d.Name,
		Pair:         cp,
		AssetType:    a,
		LastUpdated:  tickerPriceResponse.Timestamp.Time(),
		Bid:          tickerPriceResponse.BestBidPrice,
		Ask:          tickerPriceResponse.BestAskPrice,
		BidSize:      tickerPriceResponse.BestBidAmount,
		AskSize:      tickerPriceResponse.BestAskAmount,
		Last:         tickerPriceResponse.LastPrice,
		High:         tickerPriceResponse.Stats.High,
		Low:          tickerPriceResponse.Stats.Low,
		Volume:       tickerPriceResponse.Stats.Volume,
	}
	if a != asset.Futures {
		tickerPrice.Low = tickerPriceResponse.MinPrice
		tickerPrice.High = tickerPriceResponse.MaxPrice
		tickerPrice.Last = tickerPriceResponse.MarkPrice
		tickerPrice.Ask = tickerPriceResponse.ImpliedAsk
		tickerPrice.Bid = tickerPriceResponse.ImpliedBid
	}
	d.Websocket.DataHandler <- tickerPrice
	return nil
}

func (d *Deribit) processData(respRaw []byte, result interface{}) error {
	var response WsResponse
	response.Params.Data = result
	err := json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	d.Websocket.DataHandler <- result
	return nil
}

func (d *Deribit) processCandleChart(respRaw []byte, channels []string) error {
	if len(channels) != 4 {
		return fmt.Errorf("%w, expected format 'chart.trades.{instrument_name}.{resolution}', but found %s", errMalformedData, strings.Join(channels, "."))
	}
	cp, a, err := d.getAssetPairByInstrument(channels[2])
	if err != nil {
		return err
	}
	var response WsResponse
	candleData := &wsCandlestickData{}
	response.Params.Data = candleData
	err = json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	d.Websocket.DataHandler <- stream.KlineData{
		Timestamp:  time.UnixMilli(candleData.Tick),
		Pair:       cp,
		AssetType:  a,
		Exchange:   d.Name,
		OpenPrice:  candleData.Open,
		HighPrice:  candleData.High,
		LowPrice:   candleData.Low,
		ClosePrice: candleData.Close,
		Volume:     candleData.Volume,
	}
	return nil
}

func (d *Deribit) processOrderbook(respRaw []byte, channels []string) error {
	var response WsResponse
	orderbookData := &wsOrderbook{}
	response.Params.Data = orderbookData
	err := json.Unmarshal(respRaw, &response)
	if err != nil {
		return err
	}
	if len(channels) == 3 {
		cp, a, err := d.getAssetPairByInstrument(orderbookData.InstrumentName)
		if err != nil {
			return err
		}
		asks := make(orderbook.Tranches, 0, len(orderbookData.Asks))
		for x := range orderbookData.Asks {
			if len(orderbookData.Asks[x]) != 3 {
				return errMalformedData
			}
			price, okay := orderbookData.Asks[x][1].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid orderbook price", errMalformedData)
			}
			amount, okay := orderbookData.Asks[x][2].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid amount", errMalformedData)
			}
			asks = append(asks, orderbook.Tranche{
				Price:  price,
				Amount: amount,
			})
		}
		bids := make(orderbook.Tranches, 0, len(orderbookData.Bids))
		for x := range orderbookData.Bids {
			if len(orderbookData.Bids[x]) != 3 {
				return errMalformedData
			}
			price, okay := orderbookData.Bids[x][1].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid orderbook price", errMalformedData)
			} else if price == 0.0 {
				continue
			}
			amount, okay := orderbookData.Bids[x][2].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid amount", errMalformedData)
			}
			bids = append(bids, orderbook.Tranche{
				Price:  price,
				Amount: amount,
			})
		}
		if len(asks) == 0 && len(bids) == 0 {
			return nil
		}
		if orderbookData.Type == "snapshot" {
			return d.Websocket.Orderbook.LoadSnapshot(&orderbook.Base{
				Exchange:        d.Name,
				VerifyOrderbook: d.CanVerifyOrderbook,
				LastUpdated:     orderbookData.Timestamp.Time(),
				Pair:            cp,
				Asks:            asks,
				Bids:            bids,
				Asset:           a,
				LastUpdateID:    orderbookData.ChangeID,
			})
		} else if orderbookData.Type == "change" {
			return d.Websocket.Orderbook.Update(&orderbook.Update{
				Asks:       asks,
				Bids:       bids,
				Pair:       cp,
				Asset:      a,
				UpdateID:   orderbookData.ChangeID,
				UpdateTime: orderbookData.Timestamp.Time(),
			})
		}
	} else if len(channels) == 5 {
		cp, a, err := d.getAssetPairByInstrument(orderbookData.InstrumentName)
		if err != nil {
			return err
		}
		asks := make(orderbook.Tranches, 0, len(orderbookData.Asks))
		for x := range orderbookData.Asks {
			if len(orderbookData.Asks[x]) != 2 {
				return errMalformedData
			}
			price, okay := orderbookData.Asks[x][0].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid orderbook price", errMalformedData)
			} else if price == 0 {
				continue
			}
			amount, okay := orderbookData.Asks[x][1].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid amount", errMalformedData)
			}
			asks = append(asks, orderbook.Tranche{
				Price:  price,
				Amount: amount,
			})
		}
		bids := make([]orderbook.Tranche, 0, len(orderbookData.Bids))
		for x := range orderbookData.Bids {
			if len(orderbookData.Bids[x]) != 2 {
				return errMalformedData
			}
			price, okay := orderbookData.Bids[x][0].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid orderbook price", errMalformedData)
			} else if price == 0 {
				continue
			}
			amount, okay := orderbookData.Bids[x][1].(float64)
			if !okay {
				return fmt.Errorf("%w, invalid amount", errMalformedData)
			}
			bids = append(bids, orderbook.Tranche{
				Price:  price,
				Amount: amount,
			})
		}
		if len(asks) == 0 && len(bids) == 0 {
			return nil
		}
		return d.Websocket.Orderbook.LoadSnapshot(&orderbook.Base{
			Asks:         asks,
			Bids:         bids,
			Pair:         cp,
			Asset:        a,
			Exchange:     d.Name,
			LastUpdateID: orderbookData.ChangeID,
			LastUpdated:  orderbookData.Timestamp.Time(),
		})
	}
	return nil
}

// generateSubscriptions returns a list of configured subscriptions
func (d *Deribit) generateSubscriptions() (subscription.List, error) {
	subs := subscription.List{}
	var errs error
	assetPairs, err := d.Features.Subscriptions.AssetPairs(d)
	if err != nil {
		return nil, err
	}
	for _, baseSub := range d.Features.Subscriptions {
		if !authed && baseSub.Authenticated {
			continue
		}

		switch baseSub.Channel {
		case announcementsChannel, userAccessLogChannel, platformStateChannel, userLockChannel, platformStatePublicMethodsStateChannel:
			if baseSub.Asset != asset.Empty {
				// TODO: non-asset based channels go here first
				continue
			}
		}

		if baseSub.Asset == asset.Empty {
			// TODO: non-asset based channels go here first
			continue
		}
		for a, pairs := range assetPairs {
			if baseSub.Asset != a && baseSub.Asset != asset.All {
				continue
			}
			// TODO: Think we might need to skip some pairs for futures, and some other assets

			// This switch will probably boil down to nothing
			for _, p := range pairs {
				if baseSub.Channel == perpetualChannel && !isPerp(p) {
					// Perpetual channel only applies to perpetual instruments, not other dated future contracts
					continue
				}
				s := baseSub.Clone()
				s.Pairs = currency.Pairs{p}
				s.Asset = a
				subs = append(subs, s)
			}
		}

		/*
			TODO: Currency handling
			currencies := b.Features.Subscriptions.Currencies()
			case requestForQuoteChannel, userMMPTriggerChannel, userPortfolioChannel:
				subscriptions = append(subscriptions, &subscription.Subscription{
					Channel: subscriptionChannels[x],
					Currencies: ...
				})
				isIndex
			case priceIndexChannel, priceRankingChannel, priceStatisticsChannel, volatilityIndexChannel, markPriceOptionsChannel, estimatedExpirationPriceChannel:
				for i := range indexENUMS {
					subscriptions = append(subscriptions, &subscription.Subscription{
						Channel: subscriptionChannels[x],
						Params: map[string]interface{}{
							"index_name": indexENUMS[i],
						},
					})
				}
		*/
	}
	return subs, errs
}

// isPerp returns if the pair quote is a PERPETUAL; e.g BTC-PERPETUAL or ALGO-USDC-PERPETUAL
func isPerp(p currency.Pair) bool {
	return strings.Contains(p.Quote.Upper().String(), "PERPETUAL")
}

func (d *Deribit) manageSub(op string, subs subscription.List) ([]WsSubscriptionInput, error) {
	msgs := make([]WsSubscriptionInput, len(subs))

	// TODO: Interval => resolution
	//"resolution": "1D",
	// TODO: Levels => depth
	// "depth": "10",

	for _, s := range subs {
		if len(s.Pairs) > 1 {
			return nil, subscription.ErrNotSinglePair
		}
		sub := WsSubscriptionInput{
			JSONRPCVersion: rpcVersion,
			ID:             d.Websocket.Conn.GenerateMessageID(false),
			Method:         "public/" + op,
			Params:         map[string][]string{},
		}
		if s.Authenticated {
			sub.Method = "private/" + op
		}

		/*
			var instrumentID string
			if len(s.Pairs) == 1 {
				pairFormat, err := d.GetPairFormat(s.Asset, true)
				if err != nil {
					return nil, err
				}
				p := s.Pairs[0].Format(pairFormat)
				if s.Asset == asset.Futures {
					instrumentID = d.formatFuturesTradablePair(p)
				} else {
					instrumentID = p.String()
				}
			}

				switch s.Channel {
				case announcementsChannel,
					userAccessLogChannel,
					platformStateChannel,
					platformStatePublicMethodsStateChannel,
					userLockChannel:
					sub.Params["channels"] = []string{s.Channel}
				case orderbookChannel:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					intervalString, err := d.GetResolutionFromInterval(s.Interval)
					if err != nil {
						return nil, err
					}
					group, okay := s.Params["group"].(string)
					if !okay {
						sub.Params["channels"] = []string{orderbookChannel + "." + instrumentID + "." + intervalString}
						break
					}
					depth, okay := s.Params["depth"].(string)
					if !okay {
						sub.Params["channels"] = []string{orderbookChannel + "." + instrumentID + "." + intervalString}
						break
					}
					sub.Params["channels"] = []string{orderbookChannel + "." + instrumentID + "." + group + "." + depth + "." + intervalString}
				case chartTradesChannel:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					resolution, okay := s.Params["resolution"].(string)
					if !okay {
						resolution = "1D"
					}
					sub.Params["channels"] = []string{chartTradesChannel + "." + d.formatFuturesTradablePair(s.Pairs[0]) + "." + resolution}
				case priceIndexChannel,
					priceRankingChannel,
					priceStatisticsChannel,
					volatilityIndexChannel,
					markPriceOptionsChannel,
					estimatedExpirationPriceChannel:
					indexName, okay := s.Params["index_name"].(string)
					if !okay {
						return nil, errUnsupportedIndexName
					}
					sub.Params["channels"] = []string{s.Channel + "." + indexName}
				case instrumentStateChannel:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					kind := d.GetAssetKind(s.Asset)
					currencyCode := getValidatedCurrencyCode(s.Pairs[0])
					sub.Params["channels"] = []string{"instrument.state." + kind + "." + currencyCode}
				case rawUsersOrdersKindCurrencyChannel:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					kind := d.GetAssetKind(s.Asset)
					currencyCode := getValidatedCurrencyCode(s.Pairs[0])
					sub.Params["channels"] = []string{"user.orders." + kind + "." + currencyCode + ".raw"}
				case quoteChannel,
					incrementalTickerChannel:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					sub.Params["channels"] = []string{s.Channel + "." + instrumentID}
				case rawUserOrdersChannel:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					sub.Params["channels"] = []string{"user.orders." + instrumentID + ".raw"}
				case requestForQuoteChannel,
					userMMPTriggerChannel,
					userPortfolioChannel:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					currencyCode := getValidatedCurrencyCode(s.Pairs[0])
					sub.Params["channels"] = []string{s.Channel + "." + currencyCode}
				case tradesChannel,
					userChangesInstrumentsChannel,
					userOrdersWithIntervalChannel,
					tickerChannel,
					perpetualChannel,
					userTradesChannelByInstrument:
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					if s.Interval.Duration() == 0 {
						sub.Params["channels"] = []string{s.Channel + instrumentID}
						continue
					}
					intervalString, err := d.GetResolutionFromInterval(s.Interval)
					if err != nil {
						return nil, err
					}
					sub.Params["channels"] = []string{s.Channel + instrumentID + "." + intervalString}
				case userChangesCurrencyChannel,
					tradesWithKindChannel,
					rawUsersOrdersWithKindCurrencyAndIntervalChannel,
					userTradesByKindCurrencyAndIntervalChannel:
					kind := d.GetAssetKind(s.Asset)
					if len(s.Pairs) != 1 {
						return nil, currency.ErrCurrencyPairEmpty
					}
					currencyCode := getValidatedCurrencyCode(s.Pairs[0])
					if s.Interval.Duration() == 0 {
						sub.Params["channels"] = []string{s.Channel + "." + kind + "." + currencyCode}
						continue
					}
					intervalString, err := d.GetResolutionFromInterval(s.Interval)
					if err != nil {
						return nil, err
					}
					sub.Params["channels"] = []string{s.Channel + "." + kind + "." + currencyCode + "." + intervalString}
				default:
					return nil, errUnsupportedChannel
				}
		*/
		msgs = append(msgs, sub)
	}
	return msgs, nil
}

// Subscribe sends a websocket message to receive data from the channel
func (d *Deribit) Subscribe(subs subscription.List) error {
	return d.handleSubscription("subscribe", subs)
}

// Unsubscribe sends a websocket message to stop receiving data from the channel
func (d *Deribit) Unsubscribe(subs subscription.List) error {
	return d.handleSubscription("unsubscribe", subs)
}

func (d *Deribit) handleSubscription(op string, subs subscription.List) error {
	for x := range subs {
		payloads, err := d.manageSub(op, subs)
		if err != nil {
			return err
		}
		data, err := d.Websocket.Conn.SendMessageReturnResponse(payloads[x].ID, payloads[x])
		if err != nil {
			return err
		}
		var response wsSubscriptionResponse
		err = json.Unmarshal(data, &response)
		if err != nil {
			return fmt.Errorf("%v %v", d.Name, err)
		}
		if payloads[x].ID == response.ID && len(response.Result) == 0 {
			log.Errorf(log.ExchangeSys, "subscription to channel %s was not successful", payloads[x].Params["channels"][0])
		}
	}
	return nil
}

func getValidatedCurrencyCode(pair currency.Pair) string {
	currencyCode := pair.Base.Upper().String()
	switch currencyCode {
	case currencyBTC, currencyETH,
		currencySOL, currencyUSDT,
		currencyUSDC, currencyEURR:
		return currencyCode
	default:
		switch {
		case strings.Contains(pair.String(), currencyUSDC):
			return currencyUSDC
		case strings.Contains(pair.String(), currencyUSDT):
			return currencyUSDT
		}
		return "any"
	}
}
