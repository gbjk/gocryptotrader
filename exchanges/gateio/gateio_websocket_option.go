package gateio

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/fill"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"github.com/thrasher-corp/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-corp/gocryptotrader/exchanges/stream"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ticker"
	"github.com/thrasher-corp/gocryptotrader/exchanges/trade"
)

const (
	// channels
	optionsContractTickersChannel        = "options.contract_tickers"
	optionsUnderlyingTickersChannel      = "options.ul_tickers"
	optionsTradesChannel                 = "options.trades"
	optionsUnderlyingTradesChannel       = "options.ul_trades"
	optionsUnderlyingPriceChannel        = "options.ul_price"
	optionsMarkPriceChannel              = "options.mark_price"
	optionsSettlementChannel             = "options.settlements"
	optionsContractsChannel              = "options.contracts"
	optionsContractCandlesticksChannel   = "options.contract_candlesticks"
	optionsUnderlyingCandlesticksChannel = "options.ul_candlesticks"
	optionsOrderbookChannel              = "options.order_book"
	optionsOrderbookTickerChannel        = "options.book_ticker"
	optionsOrderbookUpdateChannel        = "options.order_book_update"
	optionsOrdersChannel                 = "options.orders"
	optionsUserTradesChannel             = "options.usertrades"
	optionsLiquidatesChannel             = "options.liquidates"
	optionsUserSettlementChannel         = "options.user_settlements"
	optionsPositionCloseChannel          = "options.position_closes"
	optionsBalancesChannel               = "options.balances"
	optionsPositionsChannel              = "options.positions"
)

var defaultOptionsSubscriptions = []string{
	optionsContractTickersChannel,
	optionsUnderlyingTickersChannel,
	optionsTradesChannel,
	optionsUnderlyingTradesChannel,
	optionsContractCandlesticksChannel,
	optionsUnderlyingCandlesticksChannel,
	optionsOrderbookChannel,
	optionsOrderbookUpdateChannel,
}

var fetchedOptionsCurrencyPairSnapshotOrderbook = make(map[string]bool)

/*
	switch push.Channel {
	case optionsContractTickersChannel:
		return g.processOptionsContractTickers(push.Result)
	case optionsUnderlyingTickersChannel:
		return g.processOptionsUnderlyingTicker(push.Result)
	case optionsTradesChannel,
		optionsUnderlyingTradesChannel:
		return g.processOptionsTradesPushData(respRaw)
	case optionsUnderlyingPriceChannel:
		return g.processOptionsUnderlyingPricePushData(push.Result)
	case optionsMarkPriceChannel:
		return g.processOptionsMarkPrice(push.Result)
	case optionsSettlementChannel:
		return g.processOptionsSettlementPushData(push.Result)
	case optionsContractsChannel:
		return g.processOptionsContractPushData(push.Result)
	case optionsContractCandlesticksChannel,
		optionsUnderlyingCandlesticksChannel:
		return g.processOptionsCandlestickPushData(respRaw)
	case optionsOrderbookChannel:
		return g.processOptionsOrderbookSnapshotPushData(push.Event, push.Result, push.Time)
	case optionsOrderbookTickerChannel:
		return g.processOrderbookTickerPushData(respRaw)
	case optionsOrderbookUpdateChannel:
		return g.processFuturesAndOptionsOrderbookUpdate(push.Result, asset.Options)
	case optionsOrdersChannel:
		return g.processOptionsOrderPushData(respRaw)
	case optionsUserTradesChannel:
		return g.processOptionsUserTradesPushData(respRaw)
	case optionsLiquidatesChannel:
		return g.processOptionsLiquidatesPushData(respRaw)
	case optionsUserSettlementChannel:
		return g.processOptionsUsersPersonalSettlementsPushData(respRaw)
	case optionsPositionCloseChannel:
		return g.processPositionCloseData(respRaw)
	case optionsBalancesChannel:
		return g.processBalancePushData(respRaw, asset.Options)
	case optionsPositionsChannel:
		return g.processOptionsPositionPushData(respRaw)
	default:
		g.Websocket.DataHandler <- stream.UnhandledMessageWarning{
			Message: g.Name + stream.UnhandledMessage + string(respRaw),
		}
		return errors.New(stream.UnhandledMessage)
	}
*/

func (g *Gateio) processOptionsContractTickers(incoming []byte) error {
	var data OptionsTicker
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &ticker.Price{
		Pair:         data.Name,
		Last:         data.LastPrice.Float64(),
		Bid:          data.Bid1Price.Float64(),
		Ask:          data.Ask1Price.Float64(),
		AskSize:      data.Ask1Size,
		BidSize:      data.Bid1Size,
		ExchangeName: g.Name,
		AssetType:    asset.Options,
	}
	return nil
}

func (g *Gateio) processOptionsUnderlyingTicker(incoming []byte) error {
	var data WsOptionUnderlyingTicker
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &data
	return nil
}

func (g *Gateio) processOptionsTradesPushData(data []byte) error {
	saveTradeData := g.IsSaveTradeDataEnabled()
	if !saveTradeData &&
		!g.IsTradeFeedEnabled() {
		return nil
	}
	resp := struct {
		Time    int64             `json:"time"`
		Channel string            `json:"channel"`
		Event   string            `json:"event"`
		Result  []WsOptionsTrades `json:"result"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	trades := make([]trade.Data, len(resp.Result))
	for x := range resp.Result {
		trades[x] = trade.Data{
			Timestamp:    resp.Result[x].CreateTime.Time(),
			CurrencyPair: resp.Result[x].Contract,
			AssetType:    asset.Options,
			Exchange:     g.Name,
			Price:        resp.Result[x].Price,
			Amount:       resp.Result[x].Size,
			TID:          strconv.FormatInt(resp.Result[x].ID, 10),
		}
	}
	return g.Websocket.Trade.Update(saveTradeData, trades...)
}

func (g *Gateio) processOptionsUnderlyingPricePushData(incoming []byte) error {
	var data WsOptionsUnderlyingPrice
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &data
	return nil
}

func (g *Gateio) processOptionsMarkPrice(incoming []byte) error {
	var data WsOptionsMarkPrice
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &data
	return nil
}

func (g *Gateio) processOptionsSettlementPushData(incoming []byte) error {
	var data WsOptionsSettlement
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &data
	return nil
}

func (g *Gateio) processOptionsContractPushData(incoming []byte) error {
	var data WsOptionsContract
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &data
	return nil
}

func (g *Gateio) processOptionsCandlestickPushData(data []byte) error {
	resp := struct {
		Time    int64                          `json:"time"`
		Channel string                         `json:"channel"`
		Event   string                         `json:"event"`
		Result  []WsOptionsContractCandlestick `json:"result"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	klineDatas := make([]stream.KlineData, len(resp.Result))
	for x := range resp.Result {
		icp := strings.Split(resp.Result[x].NameOfSubscription, currency.UnderscoreDelimiter)
		if len(icp) < 3 {
			return errors.New("malformed options candlestick websocket push data")
		}
		currencyPair, err := currency.NewPairFromString(strings.Join(icp[1:], currency.UnderscoreDelimiter))
		if err != nil {
			return err
		}
		klineDatas[x] = stream.KlineData{
			Pair:       currencyPair,
			AssetType:  asset.Options,
			Exchange:   g.Name,
			StartTime:  resp.Result[x].Timestamp.Time(),
			Interval:   icp[0],
			OpenPrice:  resp.Result[x].OpenPrice.Float64(),
			ClosePrice: resp.Result[x].ClosePrice.Float64(),
			HighPrice:  resp.Result[x].HighestPrice.Float64(),
			LowPrice:   resp.Result[x].LowestPrice.Float64(),
			Volume:     resp.Result[x].Amount.Float64(),
		}
	}
	g.Websocket.DataHandler <- klineDatas
	return nil
}

func (g *Gateio) processOrderbookTickerPushData(incoming []byte) error {
	var data WsOptionsOrderbookTicker
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &data
	return nil
}

func (g *Gateio) processOptionsOrderbookSnapshotPushData(event string, incoming []byte, updatePushedAt time.Time) error {
	if event == "all" {
		var data WsOptionsOrderbookSnapshot
		err := json.Unmarshal(incoming, &data)
		if err != nil {
			return err
		}
		base := orderbook.Base{
			Asset:           asset.Options,
			Exchange:        g.Name,
			Pair:            data.Contract,
			LastUpdated:     data.Timestamp.Time(),
			UpdatePushedAt:  updatePushedAt,
			VerifyOrderbook: g.CanVerifyOrderbook,
		}
		base.Asks = make([]orderbook.Tranche, len(data.Asks))
		for x := range data.Asks {
			base.Asks[x].Amount = data.Asks[x].Size
			base.Asks[x].Price = data.Asks[x].Price.Float64()
		}
		base.Bids = make([]orderbook.Tranche, len(data.Bids))
		for x := range data.Bids {
			base.Bids[x].Amount = data.Bids[x].Size
			base.Bids[x].Price = data.Bids[x].Price.Float64()
		}
		return g.Websocket.Orderbook.LoadSnapshot(&base)
	}
	var data []WsFuturesOrderbookUpdateEvent
	err := json.Unmarshal(incoming, &data)
	if err != nil {
		return err
	}
	dataMap := map[string][2][]orderbook.Tranche{}
	for x := range data {
		ab, ok := dataMap[data[x].CurrencyPair]
		if !ok {
			ab = [2][]orderbook.Tranche{}
		}
		if data[x].Amount > 0 {
			ab[1] = append(ab[1], orderbook.Tranche{
				Price: data[x].Price.Float64(), Amount: data[x].Amount,
			})
		} else {
			ab[0] = append(ab[0], orderbook.Tranche{
				Price: data[x].Price.Float64(), Amount: -data[x].Amount,
			})
		}
		if !ok {
			dataMap[data[x].CurrencyPair] = ab
		}
	}
	if len(dataMap) == 0 {
		return errors.New("missing orderbook ask and bid data")
	}
	for key, ab := range dataMap {
		currencyPair, err := currency.NewPairFromString(key)
		if err != nil {
			return err
		}
		err = g.Websocket.Orderbook.LoadSnapshot(&orderbook.Base{
			Asks:            ab[0],
			Bids:            ab[1],
			Asset:           asset.Options,
			Exchange:        g.Name,
			Pair:            currencyPair,
			LastUpdated:     updatePushedAt,
			UpdatePushedAt:  updatePushedAt,
			VerifyOrderbook: g.CanVerifyOrderbook,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Gateio) processOptionsOrderPushData(data []byte) error {
	resp := struct {
		Time    int64            `json:"time"`
		Channel string           `json:"channel"`
		Event   string           `json:"event"`
		Result  []WsOptionsOrder `json:"result"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	orderDetails := make([]order.Detail, len(resp.Result))
	for x := range resp.Result {
		status, err := order.StringToOrderStatus(func() string {
			if resp.Result[x].Status == "finished" {
				return "cancelled"
			}
			return resp.Result[x].Status
		}())
		if err != nil {
			return err
		}
		orderDetails[x] = order.Detail{
			Amount:         resp.Result[x].Size,
			Exchange:       g.Name,
			OrderID:        strconv.FormatInt(resp.Result[x].ID, 10),
			Status:         status,
			Pair:           resp.Result[x].Contract,
			Date:           resp.Result[x].CreationTime.Time(),
			ExecutedAmount: resp.Result[x].Size - resp.Result[x].Left,
			Price:          resp.Result[x].Price,
			AssetType:      asset.Options,
			AccountID:      resp.Result[x].User,
		}
	}
	g.Websocket.DataHandler <- orderDetails
	return nil
}

func (g *Gateio) processOptionsUserTradesPushData(data []byte) error {
	if !g.IsFillsFeedEnabled() {
		return nil
	}
	resp := struct {
		Time    int64                `json:"time"`
		Channel string               `json:"channel"`
		Event   string               `json:"event"`
		Result  []WsOptionsUserTrade `json:"result"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	fills := make([]fill.Data, len(resp.Result))
	for x := range resp.Result {
		fills[x] = fill.Data{
			Timestamp:    resp.Result[x].CreateTime.Time(),
			Exchange:     g.Name,
			CurrencyPair: resp.Result[x].Contract,
			OrderID:      resp.Result[x].OrderID,
			TradeID:      resp.Result[x].ID,
			Price:        resp.Result[x].Price.Float64(),
			Amount:       resp.Result[x].Size,
		}
	}
	return g.Websocket.Fills.Update(fills...)
}

func (g *Gateio) processOptionsLiquidatesPushData(data []byte) error {
	resp := struct {
		Time    int64                 `json:"time"`
		Channel string                `json:"channel"`
		Event   string                `json:"event"`
		Result  []WsOptionsLiquidates `json:"result"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &resp
	return nil
}

func (g *Gateio) processOptionsUsersPersonalSettlementsPushData(data []byte) error {
	resp := struct {
		Time    int64                     `json:"time"`
		Channel string                    `json:"channel"`
		Event   string                    `json:"event"`
		Result  []WsOptionsUserSettlement `json:"result"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &resp
	return nil
}

func (g *Gateio) processOptionsPositionPushData(data []byte) error {
	resp := struct {
		Time    int64               `json:"time"`
		Channel string              `json:"channel"`
		Event   string              `json:"event"`
		Result  []WsOptionsPosition `json:"result"`
	}{}
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}
	g.Websocket.DataHandler <- &resp
	return nil
}
