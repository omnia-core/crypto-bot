package upbit

import (
	"crypto-bot/pkg/config"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Client interface {
	ListAccounts() (Accounts, error)
	GetMarketAll(request MarketAllParams) (Markets, error)
	ListMinuteCandles(request ListCandlesParams) (Candles, error)

	PlaceOrder(request PlaceOrderParams) (Order, error)
}

type client struct {
	cfg    *config.Config
	client *resty.Client
}

func New(cfg *config.Config) Client {
	restyClient := resty.New().
		SetTimeout(60*time.Second).
		SetBaseURL(cfg.Upbit.BaseURL).
		SetHeader("Accept", "application/json")
	return &client{
		cfg:    cfg,
		client: restyClient,
	}
}

func (c *client) ListAccounts() (Accounts, error) {
	resp, err := c.sendRequest("GET", "/accounts", nil)
	if err != nil {
		return Accounts{}, err
	}

	var accounts Accounts
	err = json.Unmarshal(resp, &accounts)
	if err != nil {
		return Accounts{}, err
	}

	return accounts, nil
}

// GetMarketAll 종목 코드 조회
func (c *client) GetMarketAll(request MarketAllParams) (Markets, error) {
	resp, err := c.sendRequest("GET", "/market/all", nil)
	if err != nil {
		return Markets{}, err
	}

	var markets Markets
	err = json.Unmarshal(resp, &markets)
	if err != nil {
		return Markets{}, err
	}

	return markets, nil
}

// ListMinuteCandles 캔들 조회 - 분(Minute) 캔들 조회
func (c *client) ListMinuteCandles(request ListCandlesParams) (Candles, error) {
	if err := request.valid(); err != nil {
		return Candles{}, err
	}

	params := map[string]string{
		"market": request.Market,
		"count":  fmt.Sprintf("%d", request.Count),
	}
	if request.To != nil {
		params["to"] = request.To.Format(time.RFC3339)
	}

	resp, err := c.sendRequest("GET", fmt.Sprintf("/candles/minutes/%d", request.Minutes), params)
	if err != nil {
		return Candles{}, err
	}

	var candles Candles
	err = json.Unmarshal(resp, &candles)
	if err != nil {
		return Candles{}, err
	}

	return candles, nil
}

func (c *client) PlaceOrder(request PlaceOrderParams) (Order, error) {
	if err := request.valid(); err != nil {
		return Order{}, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		return Order{}, err
	}

	resp, err := c.sendRequest("POST", "/orders", body)
	if err != nil {
		return Order{}, err
	}

	var order Order
	err = json.Unmarshal(resp, &order)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (c *client) sendRequest(method string, path string, request interface{}) ([]byte, error) {
	token, err := c.getToken(request)
	if err != nil {
		return nil, err
	}

	c.client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))

	var resp *resty.Response
	switch method {
	case "GET":
		if request != nil {
			resp, err = c.client.R().SetQueryParams(request.(map[string]string)).Get(path)
		} else {
			resp, err = c.client.R().Get(path)
		}
	case "POST":
		resp, err = c.client.R().SetHeader("Content-Type", "application/json").SetBody(request).Post(path)
	default:
		return nil, fmt.Errorf("method %s not supported", method)
	}
	if err != nil {
		return nil, fmt.Errorf("upbit request failed: %w", err)
	}

	if resp.StatusCode() > 400 {
		return nil, fmt.Errorf("upbit error %s", resp.String())
	}

	return resp.Body(), nil
}

func (c *client) getToken(request interface{}) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"access_key": c.cfg.Upbit.AccessKey,
		"nonce":      id.String(),
	}

	if request != nil {
		queryHash := ""
		switch r := request.(type) {
		case map[string]string:
			queryHash = sha512Hash(createQueryString(r))
		case []byte:
			queryHash = sha512Hash(string(r))
		}
		claims["query_hash"] = queryHash
		claims["query_hash_alg"] = "SHA512"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.cfg.Upbit.SecretKey))
}

// createQueryString converts a map into a URL-encoded query string
func createQueryString(query map[string]string) string {
	values := url.Values{}
	for key, value := range query {
		values.Add(key, value)
	}
	return values.Encode()
}

// sha512Hash generates a SHA512 hash of a string
func sha512Hash(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
