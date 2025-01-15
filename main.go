package main

import (
	"crypto-bot/internal/trade"
	"crypto-bot/pkg/config"
	"crypto-bot/pkg/log"
	"crypto-bot/pkg/upbit"
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

func main() {
	setupLogger()
	cfg := config.ParseConfig()

	upbitClient := upbit.New(cfg)

	trader := trade.NewTrader(upbitClient)

	trader.Run()
}

func setupLogger() {
	log.Logger = lo.ToPtr(zerolog.New(os.Stderr))
}
