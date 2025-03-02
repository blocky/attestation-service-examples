# Fetch and Process A Call to CoinGecko API

This example shows you how to use the Blocky Attestation Service (Blocky AS) to
attest a function call that fetches data from the CoinGecko API and processes 
it. 

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../hello_world_attest_fn_call/README.md)
and the
[Error Handling](../error-handling/README.md)
examples.

In this example, you'll learn how to:
- Pass in parameters and secrets to your function
- Make an HTTP request to an external API in your function
- Parse a JSON response from an API

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
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
  "Value": {
    "market": "Binance",
    "coin_id": "BTC",
    "currency": "USD",
    "price": 86156,
    "timestamp": "2025-03-02T00:33:34Z"
  },
  "Error": ""
}
```

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

### Step 1: Implement the oracle

We'll implement the oracle as `myOracleFunc` in 
[`main.go`](./main.go). As in previous examples, we will call this function
using the `bky-as` CLI by passing in the 
[`fn-call.json`](./fn-call.json) file contents:

```json
[
  {
    "code_file": "tmp/x.wasm",
    "function": "myOracleFunc",
    "input": {
      "market": "Binance",
      "coin_id": "bitcoin"
    },
    "secret": {
      "api_key": "CoinGecko API key"
    }
  }
]
```

Notice the `input` section, which contains the parameters for `myOracleFunc`, 
specifically the `market` field set to "Binance" and the `coin_id` field set to
"bitcoin". The `secret` section contains the `api_key` field, which you should
set to your CoinGecko API key. Of course, you can change these values to get
the price of other coins or from other markets.

Next, we define the `myOracleFunc` function:

















and parses the JSON response. In this example,
we'll fetch the price of Bitcoin in USD from the Binance market.

In this example, you can get the price of *Bitcoin* from CoinGecko by calling
their ticker API:

```bash
curl https://api.coingecko.com/api/v3/coins/bitcoin/tickers | jq .
```

That gives you a lot of information from multiple markets, but let's say you
just want to get the price reported by *Binance*. The goal is to write a Go
function that fetches the ticker data, finds the Binance entry, and parses out
the price.

## Implementing the oracle

Let's start by setting up a struct to parse the relevant fields from the
CoinGecko API reponse JSON:

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

Next, we'll define the `getPrice` function to fetch and parse the data from the
CoinGecko API:

```go
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
```

You will pass in:
- `market`: the market name to look for (e.g., "Binance")
- `coinID`: the coin ID to look for (e.g., "bitcoin")
- `apiKey`: the API key to use for authentication with the CoinGecko API

The `as` structs and functions will make the HTTP request from the WASM runtime
through the Blocky AS server networking stack. The rest of the function is
standard Go code.


Next you will define the `myOracleFunc` function for you to invoke through
`bky-as`:

```go
//export myOracleFunc
func myOracleFunc(inputPtr, secretPtr uint64) uint64 {
	var input Args
	inputData := as.Bytes(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return writeErr(outErr.Error())
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
```

The `myOracleFunc` accepts your input and secrets (encrypted by `bky-as` and
decrypted by the Blocky AS server in a TEE) and calls the `getPrice` function.
It then returns the results through shared memory to the Blocky AS server to
create an attestation to send back to the `bky-as` CLI.
Notice that the `myOracleFunc` function is marked with `//export myOracleFunc`,
which exports the function to a guest module of the WASM runtime so that
we can invoke it on a TEE. You'll also notice the `writeErr` and `writePrice`
functions that are used to return the results of the function call through
the Blocky AS server to the `bky-as` CLI.


## Running the oracle

To run the oracle, check out the
[https://github.com/blocky/attestation-service-examples](https://github.com/blocky/attestation-service-examples)
repository and navigate to the `fetch_and_process_api_call` directory.

```bash
git clone git@github.com:blocky/attestation-service-examples.git
cd attestation-service-examples/fetch_and_process_api_call
```

The first step is to build the WASM binary from `main.go`, which contains the
definitions of `myOracleFunc`, `getPrice`, and several other helper functions.
You can build the binary by running:

```bash
make build
```

which uses [TinyGo](https://tinygo.org/) to compile the Go code to the
`tmp/x.wasm` file.

To send `x.wasm` to the Blocky AS server you need to define a `fn-call.json`:

```json
[
  {
    "code_file": "tmp/x.wasm",
    "function": "myOracleFunc",
    "input": {
      "market": "Binance",
      "coin_id": "bitcoin"
    },
    "secret": {
      "api_key": "CoingGecko API key"
    }
  }
]
```

where:
- `code_file`: the path to the compiled WASM binary (in this case, `tmp/x.wasm`)
- `function`: the name of the function to call (in this case, `myOracleFunc`)
- `input`: the input parameters to pass to the function (in this case, the
  `market` and `coin_id`)
- `secret`: the secret parameters to pass to the function (in this case, the
  `api_key` for CoinGecko, which you can get from their
  [website](https://docs.coingecko.com/reference/setting-up-your-api-key))

To invoke the oracle call:

```bash
make run
```

which pipes the `fn-call.json` file to the `bky-as` CLI and parses the response
to give you the attested output written out by `myOracleFunc`:

```json
{
  "Success": true,
  "Value": {
    "market": "Binance",
    "coin_id": "BTC",
    "currency": "USD",
    "price": 97338,
    "timestamp": "2025-02-14T23:20:35Z"
  },
  "Error": ""
}
```

where
- `Success`: tells you where the function call was successful
- `Value`: the value returned by the function call
- `Error`: any error that occurred during the function call
  (if `Success` is `false`)


## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. Check out other examples in this repository, to learn what
else you can do with Blocky AS.

Of course can implement
more powerful oracles, for example to access multiple APIs and synthesize their
responses through more complex logic than just parsing a JSON response.
