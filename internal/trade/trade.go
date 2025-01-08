package trade

import (
	"crypto-bot/pkg/upbit"
	"fmt"
	"log"
	"math"
)

type Trader interface {
	Run()
}

type trader struct {
	client upbit.Client
}

func NewTrader(client upbit.Client) Trader {
	return &trader{
		client: client,
	}
}

func (t trader) Run() {
	market := "KRW-XRP" // Replace with the desired market (e.g., KRW-XRP)
	period := 20        // Lookback period for SMA and Bollinger Bands

	// Fetch historical prices from Upbit
	candles, err := t.client.ListMinuteCandles(upbit.ListCandlesRequest{
		Minutes: 5,
		Market:  market,
		Count:   200,
	})
	if err != nil {
		log.Fatalf("Failed to fetch historical prices: %v", err)
	}
	prices := make([]float64, len(candles))
	for i, candle := range candles {
		prices[i] = candle.TradePrice
	}
	// Process data and calculate signals
	for i := period; i < len(prices); i++ {
		mean, stdDev := t.calculateMeanAndStdDev(prices[:i], period)
		signal := generateSignal(prices[i], mean, stdDev)

		// Print results
		fmt.Printf("Price: %.2f | SMA: %.2f | StdDev: %.2f | Signal: %d\n",
			prices[i], mean, stdDev, signal)
	}

	// Example of executing trades (not implemented here)
}

func (t trader) calculateMeanAndStdDev(prices []float64, period int) (float64, float64) {
	if len(prices) < period {
		return 0, 0
	}

	// Calculate mean
	sum := 0.0
	for _, price := range prices[len(prices)-period:] {
		sum += price
	}
	mean := sum / float64(period)

	// Calculate standard deviation
	variance := 0.0
	for _, price := range prices[len(prices)-period:] {
		diff := price - mean
		variance += diff * diff
	}
	stdDev := math.Sqrt(variance / float64(period))

	return mean, stdDev
}

// generateSignal generates buy/sell signals based on Bollinger Bands
func generateSignal(price, mean, stdDev float64) int {
	upperBand := mean + 2*stdDev
	lowerBand := mean - 2*stdDev

	if price > upperBand {
		return -1 // Sell
	} else if price < lowerBand {
		return 1 // Buy
	}
	return 0 // Hold
}
