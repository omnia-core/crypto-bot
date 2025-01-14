package main

import (
	"crypto-bot/internal/trade"
	"crypto-bot/pkg/config"
	"crypto-bot/pkg/log"
	"crypto-bot/pkg/upbit"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

func main() {
	setupLogger()
	cfg := config.ParseConfig()

	fmt.Println(cfg)

	upbitClient := upbit.New(cfg)
	//balances, err := upbitClient.ListAccounts()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(balances)
	//
	//marketAll, err := upbitClient.GetMarketAll(upbit.MarketAllRequest{})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(marketAll)
	//
	//candles, err := upbitClient.ListMinuteCandles(upbit.ListCandlesRequest{
	//	Market:  "KRW-XRP",
	//	Minutes: 5,
	//	Count:   1,
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(candles)

	trader := trade.NewTrader(upbitClient)

	trader.Run()
}

func setupLogger() {
	log.Logger = lo.ToPtr(zerolog.New(os.Stderr))
}
