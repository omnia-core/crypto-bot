package main

import (
	"crypto-bot/pkg/config"
	"crypto-bot/pkg/upbit"
	"fmt"
)

func main() {
	cfg := config.ParseConfig()

	fmt.Println(cfg)

	upbitClient := upbit.New(cfg)
	balances, err := upbitClient.GetBalances()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(balances))

	marketAll, err := upbitClient.GetMarketAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(marketAll))
}
