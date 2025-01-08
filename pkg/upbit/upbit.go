package upbit

import (
	"crypto-bot/pkg/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Client interface {
	ListAccounts() (AccountsResponse, error)
	GetMarketAll(request MarketAllRequest) (MarketsResponse, error)
	ListMinuteCandles(request ListCandlesRequest) (CandlesResponse, error)
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

func (c *client) ListAccounts() (AccountsResponse, error) {
	resp, err := c.sendRequest("GET", "/accounts", nil)
	if err != nil {
		return AccountsResponse{}, err
	}

	var accounts AccountsResponse
	err = json.Unmarshal(resp, &accounts)
	if err != nil {
		return AccountsResponse{}, err
	}

	return accounts, nil
}

// GetMarketAll 종목 코드 조회
func (c *client) GetMarketAll(request MarketAllRequest) (MarketsResponse, error) {
	resp, err := c.sendRequest("GET", "/market/all", nil)
	if err != nil {
		return MarketsResponse{}, err
	}

	var markets MarketsResponse
	err = json.Unmarshal(resp, &markets)
	if err != nil {
		return MarketsResponse{}, err
	}

	return markets, nil
}

// ListMinuteCandles 캔들 조회 - 분(Minute) 캔들 조회
func (c *client) ListMinuteCandles(request ListCandlesRequest) (CandlesResponse, error) {
	if err := request.valid(); err != nil {
		return CandlesResponse{}, err
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
		return CandlesResponse{}, err
	}

	var candles CandlesResponse
	err = json.Unmarshal(resp, &candles)
	if err != nil {
		return CandlesResponse{}, err
	}

	return candles, nil
}

func (c *client) sendRequest(method string, path string, request interface{}) ([]byte, error) {
	token, err := c.getToken()
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
		resp, err = c.client.R().SetBody(request).Post(path)
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

func (c *client) getToken() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"access_key": c.cfg.Upbit.AccessKey,
		"nonce":      id.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.cfg.Upbit.SecretKey))
}
