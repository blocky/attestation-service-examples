# Time-Weighted Average Price

This examples shows how to use the Blocky Attestation Service (Blocky AS) to
compute a time-weighted average price (TWAP) of a token. 

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../../hello_world_attest_fn_call/README.md),
[Error Handling](../../error_handling/README.md), and
[Getting Coin Prices From CoinGecko](../../coin_prices_from_coingecko/README.md)
examples.

In this example, you'll learn how to:

- Verify a previously attested function call inside another function call

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.

## Quick Start

To run this example, collect the first price samples by calling:

```bash
make init
```

Then collect two additional samples by calling

```bash
make iteration
make iteration
```

Finally, compute the TWAP by calling:

```bash
make twap
```

Your output should show the attested TWAP of WETH like:

```json
{
  "Success": true,
  "Error": "",
  "Value": 1911.67
}
```

## Walkthrough

### Step 1: Collect a price sample

Let's say you want to compute the time-weighted average price (TWAP) of a 
a wrapped Ethereum (WETH) token. We can get the WETH price from the Steer 
Finance API by calling:

```bash
curl https://app.steer.finance/api/token/price?tokenAddress=0x7ceb23fd6bc0add59e62ac25578270cff1b9f619&chainId=137 | jq .
```

which will return a JSON response like:

```json
{
  "price": 1907.12
}
```

That's great, but to compute a TWAP, we need to collect multiple price samples,
and we need to record the time of each sample. Let's do that by defining a 
`getNewSamplePrice` function in [`main.go`](./main.go):

```go
type SteerData struct {
	Price float64 `json:"price"`
}

func getNewPriceSample(tokenAddress string, chainID string) (price.Price, error) {
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL: fmt.Sprintf(
			"https://app.steer.finance/api/token/price?tokenAddress=%s&chainId=%s",
			tokenAddress,
			chainID,
		),
	}
	resp, err := basm.HTTPRequest(req)
	if err != nil {
		return price.Price{}, fmt.Errorf("making http request: %w", err)
	}

	var steerData SteerData
	err = json.Unmarshal(resp.Body, &steerData)
	if err != nil {
		return price.Price{}, fmt.Errorf(
			"unmarshaling Steer data: %w...%s",
			err,
			resp.Body,
		)
	}

	now, err := TimeNow()
	if err != nil {
		return price.Price{}, fmt.Errorf("getting current time: %w", err)
	}

	return price.Price{
		Value:     steerData.Price,
		Timestamp: now,
	}, nil
}
```

where we fetch data from Steer, parse the JSON response, record the current
time, set it in a `Price` struct, and return it. If the details of this flow
are new to you, you may want to review the 
[Getting Coin Prices From CoinGecko](../../coin_prices_from_coingecko/README.md)
example, where we walk thought a similar flow in more detail.
One detail to notice is the `TimeNow` helper function, defined in
[`time.go`](./time.go), which fetches the current time from
[timeapi.io](https://timeapi.io/). In the future, we will support getting the
current time directly from the Blocky AS server running your function.

### Step 2: Attest the price samples

To compute a TWAP, we need to collect a number of samples. The challenge is that
price APIs generally provide only price samples, or historical data at intervals
that may be too large to compute a TWAP on a sufficiently granular timescale.
Our approach here will be to collect and attest one sample at a time at the 
desired interval.

