package trade

import (
	"crypto-bot/pkg/log"
	"crypto-bot/pkg/upbit"
	"fmt"
	"math"
	"time"
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

	logger := log.New()

	for {
		// Fetch historical prices
		candles, err := t.client.ListMinuteCandles(upbit.ListCandlesRequest{
			Market:  market,
			Minutes: 5,
			Count:   200,
		})
		if err != nil {
			logger.Errorf("Error fetching candles: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		prices := candles.ToTradePricesSlice()

		// Calculate mean and standard deviation
		mean, stdDev := t.calculateMeanAndStdDev(prices, period)

		// Check the latest price and generate a signal
		latestPrice := prices[0]
		signal := generateSignal(latestPrice, mean, stdDev)

		if signal == 1 {
			// Buy Signal
			fmt.Println("Buy signal detected!")
			//orderResp, err := client.placeOrder(market, "bid", "0.01", fmt.Sprintf("%.2f", latestPrice), "limit")
			//if err != nil {
			//	//log.Printf("Error placing buy order: %v", err)
			//} else {
			//	fmt.Printf("Buy order placed: %s\n", orderResp)
			//}
		} else if signal == -1 {
			// Sell Signal
			fmt.Println("Sell signal detected!")
			//orderResp, err := client.placeOrder(market, "ask", "0.01", fmt.Sprintf("%.2f", latestPrice), "limit")
			//if err != nil {
			//	//log.Printf("Error placing sell order: %v", err)
			//} else {
			//	fmt.Printf("Sell order placed: %s\n", orderResp)
			//}
		} else {
			fmt.Println("No trade signal. Holding...")
		}

		// Wait before fetching the next set of data
		//time.Sleep(60 * time.Second) // Wait 1 minute before next iteration
		time.Sleep(1 * time.Second) // Wait 1 second before next iteration
	}
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
