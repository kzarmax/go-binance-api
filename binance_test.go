package binance_test

import (
	"fmt"
	binance2 "github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-binance-api/binance"
	"github.com/go-binance-api/config"
	"github.com/go-binance-api/logger"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"
)

func init() {
	config.BasePath = "../../../"
	config.Initialize()
	logfile := filepath.Join("logs", time.Now().UTC().Format("2006_01_02_00_00_00")+".log")
	logger.Setup(logfile)
}

func subscribePrice() {
	recv := make(chan futures.WsDepthEvent)
	go binance.SubscribeFuturesOrderbook("SOLBUSD", recv)
	for {
		priceData := <-recv
		bidPrice, _, err := priceData.Bids[0].Parse()
		if err != nil {
			spew.Dump("subscribePrice Error" + err.Error())
			continue
		}
		askPrice, _, err := priceData.Asks[0].Parse()
		spew.Dump(fmt.Sprintf("OrderBook Now %d => event: %d; bid: %.4f, ask: %.4f", time.Now().UnixMilli(), priceData.Time, bidPrice, askPrice))
	}
}

func subscribeKPrice() {
	recv := make(chan futures.WsKlineEvent)
	go binance.SubscribeFuturesKline("SOLBUSD", recv)
	for {
		priceData := <-recv
		openPrice, err := strconv.ParseFloat(priceData.Kline.Open, 64)
		if err != nil {
			spew.Dump("subscribeKPrice Error" + err.Error())
			continue
		}
		closePrice, err := strconv.ParseFloat(priceData.Kline.Close, 64)
		if err != nil {
			spew.Dump("subscribeKPrice Error" + err.Error())
			continue
		}
		spew.Dump(fmt.Sprintf("Kline Now %d => event: %d; open: %.4f, close: %.4f", time.Now().UnixMilli(), priceData.Time, openPrice, closePrice))
	}
}

func subscribeSpotPrice(symbol string) {
	recv := make(chan *binance2.WsPartialDepthEvent)
	go binance.SubscribeSpotPrice(symbol, recv)
	for {
		priceData := <-recv
		bidPrice, _, err := priceData.Bids[0].Parse()
		if err != nil {
			spew.Dump("subscribeSpotPrice Error" + err.Error())
			continue
		}
		askPrice, _, err := priceData.Asks[0].Parse()
		spew.Dump(fmt.Sprintf("%s Now %d => bid: %.5f, ask: %.5f", symbol, time.Now().UnixMilli(), bidPrice, askPrice))
	}
}

func subscribeSpotKlinePrice(symbol string) {
	recv := make(chan *binance2.WsKlineEvent)
	go binance.SubscribeSpotKlinePrice(symbol, recv)
	for {
		priceData := <-recv
		openPrice, err := strconv.ParseFloat(priceData.Kline.Open, 64)
		if err != nil {
			spew.Dump("subscribeKPrice Error" + err.Error())
			continue
		}
		closePrice, err := strconv.ParseFloat(priceData.Kline.Close, 64)
		if err != nil {
			spew.Dump("subscribeKPrice Error" + err.Error())
			continue
		}
		spew.Dump(fmt.Sprintf("%s Kline Now %d => event: %d; open: %.4f, close: %.4f", symbol, time.Now().UnixMilli(), priceData.Time, openPrice, closePrice))
	}
}

func getPosition() {
	service := binance.NewBinanceService()
	pos, err := service.GetPosition("SOLBUSD")
	if err != nil {
		print(err)
	} else {
		fmt.Printf("Pos: %.3f\n", pos)
	}
}

func buy() {
	service := binance.NewBinanceService()
	result, err := service.Buy("SOLBUSD", 1)
	if err != nil {
		print(err)
	} else {
		spew.Dump(result)
	}
}

func sell() {
	service := binance.NewBinanceService()
	result, err := service.Sell("SOLBUSD", 1)
	if err != nil {
		print(err)
	} else {
		spew.Dump(result)
	}
}

func TestBinance(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	//go subscribePrice()
	//go subscribeKPrice()
	go subscribeSpotPrice("USDCUSDT")
	go subscribeSpotKlinePrice("USDCUSDT")
	go subscribeSpotPrice("BUSDUSDT")
	go subscribeSpotKlinePrice("BUSDUSDT")
	wg.Wait()
	//getPosition()
	//buy()
	time.Sleep(3 * time.Second)

}
