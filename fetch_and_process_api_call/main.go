package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type Result struct {
	Success bool
	Value   any
}

func (r Result) jsonMarshalWithError(err error) []byte {
	resultStr := fmt.Sprintf(`{ "Success": false, "Value": "%s" }`, err)
	data := []byte(resultStr)
	return data
}

func writeOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		as.Log(fmt.Sprintf("Error marshalling result: %s", err))
		return writeError(err)
	}
	return as.WriteToHost(data)
}

func writeError(err error) uint64 {
	data := Result{}.jsonMarshalWithError(err)
	return as.WriteToHost(data)
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

type Price struct {
	Market    string    `json:"market"`
	CoinID    string    `json:"coin_id"`
	Currency  string    `json:"currency"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
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

	if resp.StatusCode != http.StatusOK {
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
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return writeError(outErr)
	}

	var secret SecretArgs
	secretData := as.Bytes(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return writeErr(outErr.Error())
	}

	price, err := getPrice(input.Market, input.CoinID, secret.CoinGeckoAPIKey)
	if err != nil {
		outErr := fmt.Errorf("getting price: %w", err)
		return writeErr(outErr.Error())
	}

	return writePrice(price)
}

func main() {}
