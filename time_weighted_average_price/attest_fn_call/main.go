package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blocky/as-demo/price"
	"github.com/blocky/basm-go-sdk"
	"github.com/blocky/basm-go-sdk/x/xbasm"
)

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

func extractPriceSamples(
	eAttest json.RawMessage,
	tAttest json.RawMessage,
	whitelist []basm.EnclaveMeasurement,
) (
	[]price.Price,
	error,
) {
	// bootstrap with empty samples if we don't have a transitive attestation
	if tAttest == nil {
		return []price.Price{}, nil
	}

	verifiedTA, err := basm.VerifyAttestation(
		basm.VerifyAttestationInput{
			EnclaveAttestedKey:       basm.EnclaveAttestation(eAttest),
			TransitiveAttestedClaims: basm.TransitiveAttestation(tAttest),
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

type ArgsIterate struct {
	CoinID     string                    `json:"coin_id"`
	NumSamples int                       `json:"num_samples"`
	EAttest    json.RawMessage           `json:"eAttest"`
	TAttest    json.RawMessage           `json:"tAttest"`
	Whitelist  []basm.EnclaveMeasurement `json:"whitelist"`
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

type ArgsTWAP struct {
	EAttest   json.RawMessage           `json:"eAttest"`
	TAttest   json.RawMessage           `json:"tAttest"`
	Whitelist []basm.EnclaveMeasurement `json:"whitelist"`
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

func main() {}
