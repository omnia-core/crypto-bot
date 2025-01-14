package upbit

import (
	"time"

	eu "crypto-bot/pkg/errorutil"
)

type AccountResponse struct {
	Currency            string `json:"currency"`
	Balance             string `json:"balance"`
	Locked              string `json:"locked"`
	AvgBuyPrice         string `json:"avg_buy_price"`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified"`
	UnitCurrency        string `json:"unit_currency"`
}

type AccountsResponse []AccountResponse

type MarketAllRequest struct {
	IsDetails bool
}

type MarketResponse struct {
	Market      string      `json:"market"`
	KoreanName  string      `json:"korean_name"`
	EnglishName string      `json:"english_name"`
	MarketEvent MarketEvent `json:"market_event"`
}

type MarketsResponse []MarketResponse

type MarketEvent struct {
	Warning string `json:"warning"`
	Caution string `json:"caution"`
}

type ListCandlesRequest struct {
	Minutes int
	Market  string
	Count   int
	To      *time.Time
}

func (r ListCandlesRequest) valid() error {
	if r.Minutes == 0 {
		return eu.ErrInvalidMinutes
	}

	if r.Market == "" {
		return eu.ErrInvalidMarket
	}

	return nil
}

type CandleResponse struct {
	Market               string  `json:"market"`
	CandleDateTimeUTC    string  `json:"candle_date_time_utc"`
	CandleDateTimeKST    string  `json:"candle_date_time_kst"`
	OpeningPrice         float64 `json:"opening_price"`
	HighPrice            float64 `json:"high_price"`
	LowPrice             float64 `json:"low_price"`
	TradePrice           float64 `json:"trade_price"`
	Timestamp            int64   `json:"timestamp"`
	CandleAccTradePrice  float64 `json:"candle_acc_trade_price"`
	CandleAccTradeVolume float64 `json:"candle_acc_trade_volume"`
	Unit                 int     `json:"unit"`
}

type CandlesResponse []CandleResponse

func (r CandlesResponse) ToTradePricesSlice() []float64 {
	prices := make([]float64, len(r))
	for i, candle := range r {
		prices[i] = candle.TradePrice
	}
	return prices
}
