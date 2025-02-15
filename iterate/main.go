package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blocky/as-demo/as"
)

type Args struct {
	TokenAddress string          `json:"token_address"`
	ChainID      string          `json:"chain_id"`
	EAttest      json.RawMessage `json:"eAttest"`
	TAttest      json.RawMessage `json:"tAttest"`
	Whitelist    json.RawMessage `json:"whitelist"`
}

type SecretArgs struct {
	CoinGeckoAPIKey string `json:"api_key"`
}

type Sample struct {
	Price     float64
	Timestamp time.Time
}

type Window struct {
	Average float64
	Samples []Sample
}

type Result struct {
	Success bool
	Value   Window
	Error   string
}

func extractSamples(eAttest, tAttest, whitelist json.RawMessage) ([]Sample, error) {
	// bootstrap with empty samples if we don't have a transitive attestation
	if tAttest == nil {
		return []Sample{}, nil
	}

	verifyOut, err := as.VerifyAttestation(
		as.HostVerifyAttestationInput{
			EnclaveAttestedKey:    eAttest,
			TransitiveAttestation: tAttest,
			AcceptableMeasures:    whitelist,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not verify previous attestation: %w", err)
	}

	var fixedRep [][]byte
	err = json.Unmarshal(verifyOut.RawClaims, &fixedRep)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal previous claims: %w", err)
	}

	prevResultData := fixedRep[3]
	var prevResult Result
	err = json.Unmarshal(prevResultData, &prevResult)
	switch {
	case err != nil:
		return nil, fmt.Errorf("could not unmarshal previous output: %w", err)
	case !prevResult.Success:
		return nil, fmt.Errorf("previous run was an error: %w", err)
	}

	return prevResult.Value.Samples, nil

}

func getNewSample(
	tokenAddress string,
	chainID string,
	apiKey string,
) (
	Sample,
	error,
) {
	coinID := "everclear"

	req := as.HostHTTPRequestInput{
		Method: "GET",
		URL:    fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/tickers", coinID),
		Headers: map[string][]string{
			"x-cg-demo-api-key": []string{apiKey},
		},
	}
	resp, err := as.HostFuncHTTPRequest(req)
	if err != nil {
		return Sample{}, fmt.Errorf("making http request: %w", err)
	}

	// I parsed out the data using the following jq query
	// cat out.json | jq '{ evercler2eth: .tickers[0].converted_last.eth, timestamp: .tickers[0].timestamp }'
	coinGeckoData := struct {
		Tickers []struct {
			ConvertedLast struct {
				Eth float64 `json:"eth"`
			} `json:"converted_last"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"tickers"`
	}{}
	err = json.Unmarshal(resp.Body, &coinGeckoData)
	if err != nil {
		return Sample{}, fmt.Errorf(
			"unmarshaling coin gecko data: %w...%s",
			err,
			resp.Body,
		)
	}

	price := coinGeckoData.Tickers[0].ConvertedLast.Eth
	timestamp := coinGeckoData.Tickers[0].Timestamp

	return Sample{Price: price, Timestamp: timestamp}, nil
}

func average(samples []Sample) (float64, error) {
	acc := float64(0)
	n := float64(len(samples))
	for _, s := range samples {
		acc += s.Price
	}

	return acc / n, nil
}

func advanceWindow(input Args, secret SecretArgs) (Window, error) {
	samples, err := extractSamples(input.EAttest, input.TAttest, input.Whitelist)
	if err != nil {
		return Window{}, fmt.Errorf("extracting samples: %w", err)
	}

	newSample, err := getNewSample(
		input.TokenAddress,
		input.ChainID,
		secret.CoinGeckoAPIKey,
	)
	if err != nil {
		return Window{}, fmt.Errorf("getting new sample %w: ", err)
	}

	nextSamples := append(samples, newSample)
	if len(nextSamples) > 5 {
		nextSamples = nextSamples[:5]
	}

	nextAvg, err := average(nextSamples)
	if err != nil {
		return Window{}, fmt.Errorf("computing average: %w", err)
	}

	next := Window{
		Average: nextAvg,
		Samples: nextSamples,
	}
	return next, nil
}

//export myTestFunc
func myTestFunc(inputPtr, secretPtr uint64) uint64 {
	var input Args
	inputData := as.Bytes(inputPtr)
	err := json.Unmarshal(inputData, &input)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return emitErr(outErr.Error())
	}

	var secret SecretArgs
	secretData := as.Bytes(secretPtr)
	err = json.Unmarshal(secretData, &secret)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal secret args: %w", err)
		return emitErr(outErr.Error())
	}

	nextWindow, err := advanceWindow(input, secret)
	if err != nil {
		outErr := fmt.Errorf("updating average price: %w", err)
		return emitErr(outErr.Error())
	}

	return emitWindow(nextWindow)
}

func main() {}

func emitErr(err string) uint64 {
	result := Result{
		Success: false,
		Error:   err,
	}
	return writeResultToSharedMem(result)
}

func emitWindow(window Window) uint64 {
	result := Result{
		Success: true,
		Value:   window,
	}
	return writeResultToSharedMem(result)
}

func writeResultToSharedMem(result Result) uint64 {
	outputData, err := as.Marshal(result)
	if err != nil {
		// We panic on errors we cannot communicate back to function caller
		panic("Fatal error: could not marshal output data")
	}
	return as.ShareWithHost(outputData)
}
