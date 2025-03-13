package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blocky/as-demo/price"
	"github.com/blocky/basm-go-sdk"
	"github.com/blocky/basm-go-sdk/x/xbasm"
)

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
			EnclaveAttestedKey:       eAttest,
			TransitiveAttestedClaims: tAttest,
			AcceptableMeasures:       whitelist,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not verify previous attestation: %w", err)
	}

	claims, err := xbasm.ParseFnCallClaims(verifiedTA.RawClaims)
	if err != nil {
		return nil, fmt.Errorf("could not parse previous claims: %w", err)
	}

	var prevResult Result
	err = json.Unmarshal(claims.Output, &prevResult)
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

type ArgsIterate struct {
	TokenAddress string                    `json:"token_address"`
	ChainID      string                    `json:"chain_id"`
	NumSamples   int                       `json:"num_samples"`
	EAttest      json.RawMessage           `json:"eAttest"`
	TAttest      json.RawMessage           `json:"tAttest"`
	Whitelist    []basm.EnclaveMeasurement `json:"whitelist"`
}

//export iteration
func iteration(inputPtr, secretPtr uint64) uint64 {
	var args ArgsIterate
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &args)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal args args: %w", err)
		return WriteError(outErr)
	}

	priceSamples, err := extractPriceSamples(args.EAttest, args.TAttest, args.Whitelist)
	if err != nil {
		outErr := fmt.Errorf("extracting priceSamples: %w", err)
		return WriteError(outErr)
	}

	newPriceSample, err := getNewPriceSample(args.TokenAddress, args.ChainID)
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
func twap(inputPtr, secretPtr uint64) uint64 {
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

//export errorFunc
func errorFunc(inputPtr, secretPtr uint64) uint64 {
	err := errors.New("Expected error for testing")
	return WriteError(err)
}

func main() {}
