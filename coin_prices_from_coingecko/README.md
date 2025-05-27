# Getting Coin Prices From CoinGecko

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches data from the CoinGecko API and processes
it.

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../params_and_secrets/README.md),
and the
[Error Handling - Attested Function Calls](../error_handling_attest_fn_call/README.md)
examples.

In this example, you'll learn how to:

- Make an HTTP request to an external API in your function
- Parse a JSON response from an API

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/v0.1.0-beta.10/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the CoinGecko API](https://docs.coingecko.com/reference/setting-up-your-api-key)
  and set it in `fn-call.json` in the `api_key` field.

## Quick Start

To run this example, call:

```bash
make run
```

You will see the following output extracted from a Blocky AS response showing
you the price of Bitcoin in USD on the Binance market:

```json
{
  "Success": true,
  "Error": "",
  "Value": {
    "market": "Binance",
    "coin_id": "BTC",
    "currency": "USD",
    "price": 83799,
    "timestamp": "2025-03-04T03:20:35Z"
  }
}
```

> If you get a `401` error, you need to set your CoinGecko API key in
> [`fn-call.json`](./fn-call.json). 

## Walkthrough

Let's say you want to implement a simple price feed oracle that fetches the
price of Bitcoin in USD. We can get the price from the CoinGecko API using their
ticker API:

```bash
curl https://api.coingecko.com/api/v3/coins/bitcoin/tickers | jq .
```

If you run the above command, you will get a lot of information from multiple
markets. Let's say you just want to get the price reported by Binance. In this
example, you will write a Go function that fetches the ticker data, finds the
Binance entry, and parses out the price.

### Step 1: Create a parameterized oracle function

We'll implement the oracle as `priceFunc` in
[`main.go`](./main.go). As in previous examples, we will call this function
using the `bky-as` CLI by passing in the
[`fn-call.json`](./fn-call.json) file contents:

```json
{
  "code_file": "tmp/x.wasm",
  "function": "priceFunc",
  "input": {
    "market": "Binance",
    "coin_id": "bitcoin"
  },
  "secret": {
    "api_key": "{{.YOUR_COINGECKO_API_KEY}}"
  }
}
```

Notice the `input` section, which contains the parameters for `priceFunc`,
specifically the `market` field set to "Binance" and the `coin_id` field set to
"bitcoin". The `secret` section contains the `api_key` field, which you should
set to your CoinGecko API key. The `{{. }}` notation is a placeholder used
for automating testing of this example. You should remove this notation when
setting your api key so the final value is your key enclosed in double quotes.
Of course, you can change these values to get the price of other coins or
from other markets.

Next, we define the `priceFunc` function in [`main.go`](./main.go):

```go
type Args struct {
	Market string `json:"market"`
	CoinID string `json:"coin_id"`
}

type SecretArgs struct {
	CoinGeckoAPIKey string `json:"api_key"`
}

//export priceFunc
func priceFunc(inputPtr uint64, secretPtr uint64) uint64 {
	var input Args
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return WriteError(outErr)
	}

	var secret SecretArgs
	secretData := basm.ReadFromHost(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return WriteError(outErr)
	}

	price, err := getPriceFromCoinGecko(
		input.Market,
		input.CoinID,
		secret.CoinGeckoAPIKey,
	)
	if err != nil {
		outErr := fmt.Errorf("getting price: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(price)
}
```

First, we get the input parameters and secrets. Next, we call the `getPriceFromCoinGecko`
function to fetch the price of `input.CoinID` in the `input.Market` market using
the `secret.CoinGeckoAPIKey` API key. Finally, we return the `price` to user by
converting its data to fat pointer using the `WriteOutput` function and
returning the pointer from `priceFunc` to the Blocky AS server host runtime.

### Step 2: Make a request to the CoinGecko API

The `getPriceFromCoinGecko` function, in `priceFunc`, will make an HTTP request
to the CoinGecko API to fetch the price of a coin in a specific market.

Let's start by setting up a struct to parse the relevant fields from the
CoinGecko API response JSON:

```go
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
```

Next, we'll define the `getPriceFromCoinGecko` function to fetch and parse the data from the
CoinGecko API:

```go
type Price struct {
	Market    string    `json:"market"`
	CoinID    string    `json:"coin_id"`
	Currency  string    `json:"currency"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func getPriceFromCoinGecko(market string, coinID string, apiKey string) (Price, error) {
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/tickers", coinID),
		Headers: map[string][]string{
			"x-cg-demo-api-key": []string{apiKey},
		},
	}
	resp, err := basm.HTTPRequest(req)
	switch {
	case err != nil:
		return Price{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
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
```

The `getPriceFromCoinGecko` function takes the `market`, `coinID`, and `apiKey`
as arguments. First it constructs an HTTP request to the CoinGecko API using
`coinID` in the URL and the `apiKey` in the headers. It then sends the request
to the `basm.HTTPRequest` function, which makes the request through the
Blocky AS server networking stack. Next, it checks the response status code and
unmarshalls the JSON response into the `CoinGeckoResponse` struct. Finally, it
iterates through the tickers in the response to find the ticker for the
specified `market` and returns the price as a `Price` struct to the `priceFunc`
function. The `priceFunc` function returns a `Result` containing the `Price` to
the Blocky AS server to create an attestation over the function call and the
`Result` struct.

### Step 3: Run the oracle

To run `priceFunc`, you need call:

```bash
make run
```

You'll see output similar to the following:

```json
{
  "success": true,
  "value": {
    "market": "Binance",
    "coin_id": "BTC",
    "currency": "USD",
    "price": 104273,
    "timestamp": "2025-05-13T21:25:01Z"
  }
}
```

where `"success": true,` tells you that the function call was successful and
the `value` field gives you a JSON-serialized `Price` struct.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. For example, you can try passing in different parameters to
`priceFunc`, or changing out the API endpoint in `getPriceFromCoinGecko` to fetch data from
a different API, or even multiple APIs. You may also want to explore the
[Hello World - Bringing A Blocky AS Function Call Attestation On Chain](../hello_world_on_chain/README.md)
example to learn you can bring the `Price` struct into a smart contract.
