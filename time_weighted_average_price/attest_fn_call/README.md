# Attesting A Time-Weighted Average Price

This examples shows how to use the Blocky Attestation Service (Blocky AS) to
compute a time-weighted average price (TWAP) of a token. 

Before starting this example, make sure you are familiar with the
[Hello World - Attesting a Function Call](../../hello_world_attest_fn_call/README.md),
[Passing Input Parameters and Secrets](../../params_and_secrets/README.md),
[Error Handling - Attested Function Calls](../../error_handling_attest_fn_call/README.md), and
[Getting Coin Prices From CoinGecko](../../coin_prices_from_coingecko/README.md)
examples.

In this example, you'll learn how to:

- Verify a previously attested function call inside another function call
- Compute a time-weighted average price (TWAP) from a collection of attested
  price samples

## Setup

- Install the Blocky AS CLI by following the
  [setup instructions](https://blocky-docs.redocly.app/attestation-service/{{{AS_VERSION}}}/setup)
  in the Blocky AS documentation.
- Make sure you also have
  [Docker](https://www.docker.com/) and [jq](https://jqlang.org/) installed on
  your system.
- [Get a key for the CoinGecko API](https://docs.coingecko.com/reference/setting-up-your-api-key)
  and set it in `iteration-call.json.template` in the `api_key` field.

## Quick Start

To run this example, call:

```bash
make init
for i in {1..2}
do
  sleep 1
  make iteration
done
make twap
```

which shows you how to collect three price samples at 1 second intervals and
then compute a TWAP from the collected samples.

The output of the `make twap` should show the attested TWAP of wrapped ETH 
similar to:

```json
{
  "Success": true,
  "Error": "",
  "Value": 1896.93
}
```

## Walkthrough

### Step 1: Collect a price sample

Let's say you want to compute the time-weighted average price (TWAP) of a 
wrapped Ethereum (WETH) token. We can get the WETH price from the 
CoinGecko API by calling:

```bash
https://api.coingecko.com/api/v3/simple/price?ids=polygon-pos-bridged-weth-polygon-pos&vs_currencies=usd&include_last_updated_at=true&precision=full | jq .
```

which will return a JSON response like:

```json
{
  "polygon-pos-bridged-weth-polygon-pos": {
    "usd": 1929.1247885555736,
    "last_updated_at": 1742255916
  }
}
```

That's great, but to compute a TWAP, we need to collect multiple price samples. 
Let's do these things by defining a `getNewSamplePrice` function in 
[`main.go`](./main.go):

```go
type CoinGeckoResponse struct {
	Price struct {
		USD           float64 `json:"usd"`
		LastUpdatedAt int     `json:"last_updated_at"`
	} `json:"price"`
}

func getNewPriceSample(coinID string, apiKey string) (price.Price, error) {
	req := basm.HTTPRequestInput{
		Method: "GET",
		URL: fmt.Sprintf(
			"https://api.coingecko.com/api/v3/simple/price"+
				"?ids=%s"+
				"&vs_currencies=usd"+
				"&include_last_updated_at=true"+
				"&precision=full",
			coinID,
		),
		Headers: map[string][]string{
			"x-cg-demo-api-key": []string{apiKey},
		},
	}
	resp, err := basm.HTTPRequest(req)
	switch {
	case err != nil:
		return price.Price{}, fmt.Errorf("making http request: %w", err)
	case resp.StatusCode != http.StatusOK:
		return price.Price{}, fmt.Errorf(
			"http request failed with status code %d",
			resp.StatusCode,
		)
	}

	respBody := bytes.ReplaceAll(resp.Body, []byte(coinID), []byte("price"))

	var coinGeckoResponse CoinGeckoResponse
	err = json.Unmarshal(respBody, &coinGeckoResponse)
	if err != nil {
		return price.Price{}, fmt.Errorf(
			"unmarshaling CoinGecko data: %w...%s",
			err,
			resp.Body,
		)
	}

	timestamp := time.Unix(int64(coinGeckoResponse.Price.LastUpdatedAt), 0)

	return price.Price{
		Value:     coinGeckoResponse.Price.USD,
		Timestamp: timestamp,
	}, nil
}
```

where we fetch data from CoinGecko, parse the JSON response, record the current
time, set it in a `Price` struct defined in the `price` package in 
[`price/price.go`](./price/price.go), and return it. 
If the details of this flow are new to you, you may want to review the 
[Getting Coin Prices From CoinGecko](../../coin_prices_from_coingecko/README.md)
example, where we walk thought how to fetch and parse API data in more detail.
One thing to notice in this example is that we replace the value of the 
variable `coinID` in the response body with the string `"price"` to allow
for deterministic unmarshalling the response into the `CoinGeckoResponse` 
struct.


### Step 2: Attest the price samples

To compute a TWAP, we want to collect a number of price samples. The challenge
is that price APIs generally provide only spot prices, or historical data at
intervals that may be too large to compute a TWAP on a sufficiently granular
timescale. To solve this problem, we will follow an iterative process, where we
collect price samples at the desired interval, and then compute a TWAP from the
collected samples. The iterative step will read in a previously attested
collection of price samples, expand it with a new price sample, and then attest
the updated collection of price samples.

We define the iterative step in the `iteration` function, in
[`main.go`](./main.go):

```go
type ArgsIterate struct {
	CoinID     string                     `json:"coin_id"`
	NumSamples int                        `json:"num_samples"`
	EAttest    basm.EnclaveAttestation    `json:"eAttest"`
	TAttest    basm.TransitiveAttestation `json:"tAttest"`
	Whitelist  []basm.EnclaveMeasurement  `json:"whitelist"`
}

type SecretArgs struct {
	CoinGeckoAPIKey string `json:"api_key"`
}

//export iteration
func iteration(inputPtr uint64, secretPtr uint64) uint64 {
	var args ArgsIterate
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &args)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal args args: %w", err)
		return WriteError(outErr)
	}

	var secret SecretArgs
	secretData := basm.ReadFromHost(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return WriteError(outErr)
	}

	priceSamples, err := extractPriceSamples(args.EAttest, args.TAttest, args.Whitelist)
	if err != nil {
		outErr := fmt.Errorf("extracting priceSamples: %w", err)
		return WriteError(outErr)
	}

	newPriceSample, err := getNewPriceSample(args.CoinID, secret.CoinGeckoAPIKey)
	if err != nil {
		outErr := fmt.Errorf("getting new sample: %w", err)
		return WriteError(outErr)
	}

	nextPriceSamples := append(priceSamples, newPriceSample)
	if len(nextPriceSamples) > args.NumSamples {
		numToRemove := len(nextPriceSamples) - args.NumSamples
		nextPriceSamples = nextPriceSamples[numToRemove:]
	}

	return WriteOutput(nextPriceSamples)
}
```

where we: 

1. Extract `iteration` call arguments `args` from the `inputPtr` and `secret`
   arguments from the `secretPtr`.
2. Call `extractPriceSamples` to extract the previous collection of price
   samples `prevPriceSamples` from the transitively attested function call
   `args.TAttest`.
3. Collects a new price sample by calling `getNewPriceSample`.
4. Compute the new collection of price samples `nextPriceSamples` by appending
   the new price sample to the previous collection and then truncating the
   collection to `args.NumSamples` most recent samples.
5. Return `nextPriceSamples` to the Blocky AS server for attestation.

Let's dive a bit deeper into is the `extractPriceSamples` function:

```go
func extractPriceSamples(
	eAttest basm.EnclaveAttestation,
	tAttest basm.TransitiveAttestation,
	whitelist []basm.EnclaveMeasurement,
) (
	[]price.Price,
	error,
) {
	// bootstrap with empty samples if we don't have a transitive attestation
	if tAttest == "" {
		return []price.Price{}, nil
	}

	verifiedTA, err := basm.VerifyAttestation(
		basm.VerifyAttestationInput{
			EnclaveAttestedKey:       eAttest,
			TransitiveAttestedClaims: tAttest,
			AcceptableMeasures:       whitelist,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not verify previous attestation: %w", err)
	}

	verifiedClaims, err := xbasm.ParseFnCallClaims(verifiedTA.RawClaims)
	if err != nil {
		return nil, fmt.Errorf("could not parse claims: %w", err)
	}

	var prevResult Result
	err = json.Unmarshal(verifiedClaims.Output, &prevResult)
	switch {
	case err != nil:
		return nil, fmt.Errorf("could not unmarshal previous output: %w", err)
	case !prevResult.Success:
		return nil, fmt.Errorf("previous run was an error: %w", err)
	}

	prevPriceSamplesStr, err := json.Marshal(prevResult.Value)
	if err != nil {
		retErr := fmt.Errorf("could not marshal previous price samples: %w", err)
		return nil, retErr
	}

	var prevPriceSamples []price.Price
	err = json.Unmarshal(prevPriceSamplesStr, &prevPriceSamples)
	if err != nil {
		retErr := fmt.Errorf("could not unmarshal previous price samples: %w", err)
		return nil, retErr
	}

	return prevPriceSamples, nil
}
```

The `extractPriceSamples` function takes the enclave attested application public
key `eAttest`, the transitive attested function call `tAttest`, and a
`whitelist` of acceptable enclave measurements as parameters. It uses these to
call the`basm` 
[Blocky Attestation Service WASM Go SDK](https://github.com/blocky/basm-go-sdk/tree/{{{AS_VERSION}}})
`basm.VerifyAttestation` function to verify that `tAttest` has been signed by
the enclave attested application public key from `eAttest`, and checks that the
code measurement in `eAttest` is present in the `whitelist`. If you'd like to
learn more about the attestation verification process, please visit the
[Attestations in the Blocky Attestation Service](https://blocky-docs.redocly.app/attestation-service/{{{AS_VERSION}}}/concepts#attestations-in-the-blocky-attestation-service)
page in our documentation. The `extractPriceSamples` function proceeds to parse
out the verified transitive attestation claims `verifiedClaims` using the
experimental `xbasm` package of our SDK. (We use the `xbasm` package to stage
features we are considering for production use in the Blocky AS SDK.) From the
`verifiedClaims`, `extractPriceSamples` extracts the `Output` field and
unmarshalls it into `prevResult` of type a `Result` defined in 
[`output.go`](./output.go). In turn, `prevResult.Value` contains the previous
collection of price samples. Since the type of `Result.Value` is `any`, we first
marshal it to JSON, to then unmarshal it into a `[]price.Price` struct
`prevPriceSamples`. Finally, the `extractPriceSamples` function returns
`prevPriceSamples` to the caller, in our case the `iteration` function.

### Step 3: Collect price samples

To collect a sample, we define the call to the `iteration` function in
[`iteration-call.json.template`](./iteration-call.json.template):

```
{
    "code_file": "./tmp/x.wasm",
    "function": "iteration",
    "input": {
        "coin_id": "polygon-pos-bridged-weth-polygon-pos",
        "num_samples": 3,
        "eAttest": VAR_EATTEST,
        VAR_TATTEST
        "whitelist": [
            { "platform": "plain", "code": "plain" }
        ]
    },
    "secret": {
      "api_key": "CoinGeckoAPIKey"
    }
}
```

where we pass in the `coin_id` of the token we want to
price, the `num_samples` we want to collect, and the `whitelist` of acceptable
enclave measurements. The `VAR_EATTEST` and `VAR_TATTEST` placeholders will be
used to insert the enclave attested application public key and the
transitive attested function call, respectively in subsequent steps.

To collect the first price sample, we call:

```bash
make init
```

which will save the enclave attestation and transitive attestation output
resulting from invoking the`iteration` function in 
[`tmp/prev.json`](./tmp/prev.json).

To collect the next price sample, we call:

```bash
make iteration
```

If you inspect the `iteration` target in the [`Makefile`](./Makefile):

```makefile
prev: check
	$(eval prev_ea := $(shell jq '.enclave_attested_application_public_key.enclave_attestation' tmp/prev.json | sed 's/\//\\\//g' ))
	$(eval prev_ta := $(shell jq '.transitive_attested_function_call.transitive_attestation' tmp/prev.json ))

iteration: check prev build
	@sed \
		-e 's/VAR_TATTEST/"tAttest": ${prev_ta},/' \
		-e 's/VAR_EATTEST/${prev_ea}/' \
		iteration-call.json.template > tmp/iteration-call.json
	@cat tmp/iteration-call.json | bky-as attest-fn-call | jq . > tmp/prev.json
	@jq -r '.transitive_attested_function_call.claims.output | @base64d | fromjson' tmp/prev.json
```

you'll see that `iteration` target calls the `prev` target to load the `prev_ea`
and `prev_ta` from [`tmp/prev.json`](./tmp/prev.json). The `iteration` target
then replaces the `VAR_EATTEST` and `VAR_TATTEST` in 
[`iteration-call.json.template`](./iteration-call.json.template) with the
`prev_ea` and `prev_ta` values, respectively, and saves the result in 
[`tmp/iteration-call.json`](./tmp/iteration-call.json). With the enclave
attestation and transitive attestations as realized arguments in 
[`tmp/iteration-call.json`](./tmp/iteration-call.json), the `iteration` function
in [`main.go`](./main.go) will be able to parse out the previous samples in
`extractPriceSamples`. Finally, the `iteration` target calls the
`bky-as attest-fn-call` command to attest the function call in 
[`tmp/iteration-call.json`](./tmp/iteration-call.json) and saves the output in 
[`tmp/prev.json`](./tmp/prev.json).

### Step 4: Compute the TWAP

To compute the TWAP, we define the `twap` function in
[`main.go`](./main.go):

```go
type ArgsTWAP struct {
	EAttest   basm.EnclaveAttestation    `json:"eAttest"`
	TAttest   basm.TransitiveAttestation `json:"tAttest"`
	Whitelist []basm.EnclaveMeasurement  `json:"whitelist"`
}

//export twap
func twap(inputPtr uint64, secretPtr uint64) uint64 {
	var args ArgsTWAP
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &args)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal args args: %w", err)
		return WriteError(outErr)
	}

	priceSamples, err := extractPriceSamples(args.EAttest, args.TAttest, args.Whitelist)
	if err != nil {
		outErr := fmt.Errorf("extracting samples: %w", err)
		return WriteError(outErr)
	}

	twap, err := price.TWAP(priceSamples)
	if err != nil {
		outErr := fmt.Errorf("computing TWAP: %w", err)
		return WriteError(outErr)
	}

	return WriteOutput(twap)
}
```

which follows a similar pattern to the `iteration` function. It extracts
the `args` call arguments from the `inputPtr`, calls `extractPriceSamples` to
extract the collection of price samples from `args.TAttest`, and then calls the
`TWAP` function of the `price` package in [`price/price.go`](./price/price.go) to
compute the TWAP from the extracted price samples. Finally, the `twap`
function returns the computed `twap` price to the Blocky AS server for
attestation.

If we drill down into the `price.TWAP` function:

```go
type Price struct {
	Value     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

func TWAP(samples []Price) (float64, error) {
	if len(samples) == 0 {
		return 0, fmt.Errorf("no samples provided")
	}

	// Remove duplicate samples with the same timestamp
	uniqueSamples := make([]Price, 0, len(samples))
	seen := make(map[time.Time]struct{})
	for _, sample := range samples {
		if _, ok := seen[sample.Timestamp]; !ok {
			uniqueSamples = append(uniqueSamples, sample)
			seen[sample.Timestamp] = struct{}{}
		}
	}

	if len(uniqueSamples) == 1 {
		return uniqueSamples[0].Value, nil
	}

	// Sort samples from latest to earliest
	lessThan := func(i, j int) bool {
		return uniqueSamples[i].Timestamp.After(uniqueSamples[j].Timestamp)
	}
	sort.Slice(uniqueSamples, lessThan)

	var weightedSum, totalWeight float64

	// IMPORTANT: The value of the last sample is not included in the calculation
	// because it doesn't have a next sample to compare with. However, its
	// timestamp is used to calculate the weight of the previous sample.
	prev := uniqueSamples[0]
	for _, next := range uniqueSamples[1:] {
		timeDiff := prev.Timestamp.Sub(next.Timestamp).Microseconds()
		weight := float64(timeDiff)
		weightedSum += prev.Value * weight
		totalWeight += weight
		prev = next
	}

	if totalWeight == 0 {
		return 0, fmt.Errorf("total weight is zero, cannot compute TWAP")
	}

	return weightedSum / totalWeight, nil
}
```

we can see the TWAP calculation.

To obtain the TWAP, we define a call to the `twap` function in 
[`twap-call.json.template`](./twap-call.json.template):

```
{
    "code_file": "./tmp/x.wasm",
    "function": "twap",
    "input": {
        "eAttest": VAR_EATTEST,
        VAR_TATTEST
        "whitelist": [
            { "platform": "plain", "code": "plain" }
        ]
    }
}
```

where again the `VAR_EATTEST` and `VAR_TATTEST` placeholders will be
used to insert the enclave attested application public key and the
transitive attested function call with the collected price samples.

To invoke the `twap` function, we call:

```bash
make twap
```

which will save the output of `bky-as` running the `twap` function in 
[`tmp/twap.json`](./tmp/twap.json) and give us the parsed output similar to:

```json
{
  "Success": true,
  "Error": "",
  "Value": 1896.93
}

```

where `Value` is the computed TWAP of WETH.

## Next steps

Now that you have successfully run the example, you can start modifying it to
fit your own needs. For example, you can try computing a TWAPs for different
coins. You can also extend this example, by calling the `iteration` function at
the desired time interval, for example using a 
[cron](https://en.wikipedia.org/wiki/Cron) 
job, to control the granularity of the TWAP. Finally, you may explore the
[Time-Weighted Average Price On Chain](../on_chain)
example to learn how to bring the TWAP into a smart contract.
