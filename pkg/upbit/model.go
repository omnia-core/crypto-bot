package upbit

import (
	"time"

	eu "crypto-bot/pkg/errorutil"
)

type Account struct {
	Currency            string `json:"currency"`
	Balance             string `json:"balance"`
	Locked              string `json:"locked"`
	AvgBuyPrice         string `json:"avg_buy_price"`
	AvgBuyPriceModified bool   `json:"avg_buy_price_modified"`
	UnitCurrency        string `json:"unit_currency"`
}

type Accounts []Account

type MarketAllParams struct {
	IsDetails bool
}

type Market struct {
	Market      string      `json:"market"`
	KoreanName  string      `json:"korean_name"`
	EnglishName string      `json:"english_name"`
	MarketEvent MarketEvent `json:"market_event"`
}

type Markets []Market

type MarketEvent struct {
	Warning string `json:"warning"`
	Caution string `json:"caution"`
}

type ListCandlesParams struct {
	Minutes int
	Market  string
	Count   int
	To      *time.Time
}

func (r ListCandlesParams) valid() error {
	if r.Minutes == 0 {
		return eu.ErrInvalidMinutes
	}

	if r.Market == "" {
		return eu.ErrInvalidMarket
	}

	return nil
}

type Candle struct {
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

type Candles []Candle

func (r Candles) ToTradePricesSlice() []float64 {
	prices := make([]float64, len(r))
	for i, candle := range r {
		prices[i] = candle.TradePrice
	}
	return prices
}

type PlaceOrderParams struct {
	Market    string `json:"market"`
	Side      string `json:"side"`
	Volume    string `json:"volume"`
	Price     string `json:"price"`
	OrderType string `json:"ord_type"`
}

func (r PlaceOrderParams) valid() error {
	if r.Market == "" {
		return eu.ErrInvalidMarket
	}

	if r.Side == "" {
		return eu.ErrInvalidSide
	}

	if r.Volume == "" {
		return eu.ErrInvalidVolume
	}

	if r.Price == "" {
		return eu.ErrInvalidPrice
	}

	if r.OrderType == "" {
		return eu.ErrInvalidOrderType
	}

	return nil
}

type Order struct {
	UUID            string `json:"uuid"`
	Side            string `json:"side"`
	OrderType       string `json:"ord_type"`
	Price           string `json:"price"`
	State           string `json:"state"`
	Market          string `json:"market"`
	CreatedAt       string `json:"created_at"`
	Volume          string `json:"volume"`
	RemainingVolume string `json:"remaining_volume"`
	ReservedFee     string `json:"reserved_fee"`
	RemainingFee    string `json:"remaining_fee"`
	PaidFee         string `json:"paid_fee"`
	Locked          string `json:"locked"`
	ExecutedVolume  string `json:"executed_volume"`
	TradesCount     string `json:"trades_count"`
	TimeInForce     string `json:"time_in_force"`
	Identifier      string `json:"identifier"`
}
