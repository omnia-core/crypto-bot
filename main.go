package main

import (
	"crypto-bot/pkg/config"
	"fmt"
)

func main() {
	cfg := config.ParseConfig()

	fmt.Println(cfg)

}
