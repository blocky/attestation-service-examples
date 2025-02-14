package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/blocky/as-demo/as"
)

type Args struct {
	Market string `json:"market"`
	CoinID string `json:"coin_id"`
}

type SecretArgs struct {
	CoinGeckoAPIKey string `json:"api_key"`
}

type Price struct {
	Market    string    `json:"market"`
	CoinID    string    `json:"coin_id"`
	Currency  string    `json:"currency"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type Output struct {
	IsErr bool
	Value Price
	Error string
}

func writeOutputToSharedMem(price Price, respErr error) uint64 {
	isErr := false
	errString := ""
	if respErr != nil {
		isErr = true
		errString = "Error executing function: " + respErr.Error()
	}

	output := Output{Value: price, Error: errString, IsErr: isErr}
	outputData, err := as.Marshal(output)
	if err != nil {
		// We panic on errors we cannot communicate back to function caller
		panic("Fatal error: could not marshal output data")
	}
	return as.ShareWithHost(outputData)
}

type CoinGeckoResponse struct {
	Tickers []struct {
		Base   string `json:"base"`
		Market struct {
			Name string `json:"name"`
		} `json:"market"`
		ConvertedLast struct {
			USD float64 `json:"usd"`
		} `json:"converted_last"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"tickers"`
}

func getPrice(market string, coinID string, apiKey string) (Price, error) {
	req := as.HostHTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/tickers", coinID),
		Headers: map[string][]string{
			"x-cg-demo-api-key": []string{apiKey},
		},
	}
	resp, err := as.HostFuncHTTPRequest(req)
	if err != nil {
		return Price{}, fmt.Errorf("making http request: %w", err)
	}

	if resp.StatusCode != 200 {
		return Price{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	coinGeckoResponse := CoinGeckoResponse{}
	err = json.Unmarshal(resp.Body, &coinGeckoResponse)
	if err != nil {
		return Price{}, fmt.Errorf(
			"unmarshaling  data: %w...%s", err,
			resp.Body,
		)
	}

	for _, ticker := range coinGeckoResponse.Tickers {
		if ticker.Market.Name == market {
			return Price{
				Market:    ticker.Market.Name,
				CoinID:    ticker.Base,
				Currency:  "USD",
				Price:     ticker.ConvertedLast.USD,
				Timestamp: ticker.Timestamp,
			}, nil
		}
	}

	return Price{}, fmt.Errorf("market %s not found", market)
}

//export myOracleFunc
func myOracleFunc(inputPtr, secretPtr uint64) uint64 {
	var input Args
	inputData := as.Bytes(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := errors.New("could not unmarshal input args: " + err.Error())
		return writeOutputToSharedMem(Price{}, outErr)
	}

	var secret SecretArgs
	secretData := as.Bytes(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := errors.New("could not unmarshal secret args: " + err.Error())
		return writeOutputToSharedMem(Price{}, outErr)
	}

	price, err := getPrice(input.Market, input.CoinID, secret.CoinGeckoAPIKey)
	if err != nil {
		outErr := errors.New("getting price: " + err.Error())
		return writeOutputToSharedMem(Price{}, outErr)
	}

	return writeOutputToSharedMem(price, nil)
}

func main() {}
