package binance

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/go-binance-api/config"
	"github.com/go-binance-api/logger"
	"strconv"
	"time"
)

type BinanceService struct {
	Client       *binance.Client
	FutureClient *futures.Client
	Positions    map[string]int64
	Balances     map[string]int64
}

func NewBinanceService() *BinanceService {
	futuresClient := binance.NewFuturesClient(config.Config.Binance.ApiKey, config.Config.Binance.SecretKey)
	client := binance.NewClient(config.Config.Binance.ApiKey, config.Config.Binance.SecretKey)
	return &BinanceService{
		Client:       client,
		FutureClient: futuresClient,
		Positions:    make(map[string]int64),
		Balances:     make(map[string]int64),
	}
}

func SubscribeFuturesOrderbook(symbol string, recv chan futures.WsDepthEvent) {
	wsDepthHandler := func(event *futures.WsDepthEvent) {
		recv <- *event
	}
	errHandler := func(err error) {
		panic("Binance SubscribeFuturesOrderbook error:" + err.Error())
	}
	doneC, _, err := futures.WsPartialDepthServeWithRate(symbol, 5, 100*time.Millisecond, wsDepthHandler, errHandler)
	if err != nil {
		panic(err.Error())
		return
	}
	<-doneC
}

func SubscribeSpotPrice(symbol string, recv chan *binance.WsPartialDepthEvent) {
	wsDepthHandler := func(event *binance.WsPartialDepthEvent) {
		recv <- event
	}
	errHandler := func(err error) {
		panic("Binance SubscribeSpotPrice error:" + err.Error())
	}
	doneC, _, err := binance.WsPartialDepthServe100Ms(symbol, "5", wsDepthHandler, errHandler)
	if err != nil {
		panic(err.Error())
		return
	}
	<-doneC
}

func SubscribeSpotKlinePrice(symbol string, recv chan *binance.WsKlineEvent) {
	wsDepthHandler := func(event *binance.WsKlineEvent) {
		recv <- event
	}
	errHandler := func(err error) {
		panic("Binance SubscribeSpotKlinePrice error:" + err.Error())
	}
	doneC, _, err := binance.WsKlineServe(symbol, "1s", wsDepthHandler, errHandler)
	if err != nil {
		panic(err.Error())
		return
	}
	<-doneC
}

func SubscribeFuturesKline(symbol string, recv chan futures.WsKlineEvent) {
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		recv <- *event
	}
	errHandler := func(err error) {
		panic("Binance SubscribeFuturesKline error:" + err.Error())
	}
	doneC, _, err := futures.WsKlineServe(symbol, "1m", wsKlineHandler, errHandler)
	if err != nil {
		panic(err.Error())
		return
	}
	<-doneC
}

func (b *BinanceService) SubscribeFutureUserData() {
	//b.FutureClient
	res, err := b.FutureClient.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func (b *BinanceService) GetPosition(symbol string) (float64, error) {
	// if !config.Config.Dev.Debug {
	res, err := b.FutureClient.NewGetPositionRiskService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, err
	}

	pos := &FuturePosition{}
	pos.EntryPrice = res[0].EntryPrice
	pos.IsAutoAddMargin = res[0].IsAutoAddMargin
	pos.IsolatedMargin = res[0].IsolatedMargin
	pos.IsolatedWallet = res[0].IsolatedWallet
	pos.Leverage = res[0].Leverage
	pos.LiquidationPrice = res[0].LiquidationPrice
	pos.MarginType = res[0].MarginType
	pos.MarkPrice = res[0].MarkPrice
	pos.MaxNotionalValue = res[0].MaxNotionalValue
	pos.Notional = res[0].Notional
	pos.PositionAmt = res[0].PositionAmt
	pos.Symbol = res[0].Symbol
	pos.PositionSide = res[0].PositionSide
	pos.UnRealizedProfit = res[0].UnRealizedProfit

	posAmt, _ := strconv.ParseFloat(pos.PositionAmt, 10)

	return posAmt, nil
}

func (b *BinanceService) Buy(symbol string, amount int) (*FutureOrder, error) {
	side := futures.SideTypeBuy
	return b.NewOrder(symbol, side, strconv.Itoa(amount))
}

func (b *BinanceService) Sell(symbol string, amount int) (*FutureOrder, error) {
	side := futures.SideTypeSell
	return b.NewOrder(symbol, side, strconv.Itoa(amount))
}

func (b *BinanceService) NewOrder(
	symbol string,
	side futures.SideType,
	quantity string,
) (*FutureOrder, error) {
	if config.NoSubmit {
		logger.Outlog("NewOrder Skipped :Not allowed")
		return nil, nil
	}
	order, err := b.FutureClient.
		NewCreateOrderService().
		PositionSide("BOTH").
		Symbol(symbol).Side(side).
		Quantity(quantity).
		ReduceOnly(false).
		Type(futures.OrderTypeMarket).
		Do(context.Background())

	if err != nil {
		logger.Outlog("Binance: Failed to create order" + err.Error())
		return nil, err
	}
	result := &FutureOrder{}
	result.ActivatePrice = order.ActivatePrice
	result.AvgPrice = order.AvgPrice
	result.ClientOrderID = order.ClientOrderID
	result.CumQuote = order.CumQuote
	result.OrderID = order.OrderID
	result.Price = order.Price
	result.Symbol = order.Symbol
	result.ExecutedQuantity = order.ExecutedQuantity

	return result, nil
}
