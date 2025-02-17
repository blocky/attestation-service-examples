package main

import (
	"encoding/json"
	"fmt"

	"github.com/blocky/as-demo/as"
	"github.com/blocky/as-demo/price"
)

type ArgsIterate struct {
	TokenAddress string          `json:"token_address"`
	ChainID      string          `json:"chain_id"`
	NumSamples   int             `json:"num_samples"`
	EAttest      json.RawMessage `json:"eAttest"`
	TAttest      json.RawMessage `json:"tAttest"`
	Whitelist    json.RawMessage `json:"whitelist"`
}

type ArgsTWAP struct {
	EAttest   json.RawMessage `json:"eAttest"`
	TAttest   json.RawMessage `json:"tAttest"`
	Whitelist json.RawMessage `json:"whitelist"`
}

type Result struct {
	Success bool
	Error   string
}

type PriceSamples []price.Price

type ResultPriceSamples struct {
	Result       Result
	PriceSamples PriceSamples
}

type ResultTWAP struct {
	Result Result
	TWAP   float64
}

func extractPriceSamples(
	eAttest,
	tAttest,
	whitelist json.RawMessage,
) (
	PriceSamples,
	error,
) {
	// bootstrap with empty samples if we don't have a transitive attestation
	if tAttest == nil {
		return PriceSamples{}, nil
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
	var prevResult ResultPriceSamples
	err = json.Unmarshal(prevResultData, &prevResult)
	switch {
	case err != nil:
		return nil, fmt.Errorf("could not unmarshal previous output: %w", err)
	case !prevResult.Result.Success:
		return nil, fmt.Errorf("previous run was an error: %w", err)
	}

	return prevResult.PriceSamples, nil

}

func getNewPriceSample(tokenAddress string, chainID string) (price.Price, error) {
	req := as.HostHTTPRequestInput{
		Method: "GET",
		URL: fmt.Sprintf(
			"https://app.steer.finance/api/token/price?tokenAddress=%s&chainId=%s",
			tokenAddress,
			chainID,
		),
	}
	resp, err := as.HostFuncHTTPRequest(req)
	if err != nil {
		return price.Price{}, fmt.Errorf("making http request: %w", err)
	}

	steerData := struct {
		Price float64 `json:"price"`
	}{}

	err = json.Unmarshal(resp.Body, &steerData)
	if err != nil {
		return price.Price{}, fmt.Errorf(
			"unmarshaling Steer data: %w...%s",
			err,
			resp.Body,
		)
	}

	now, err := as.TimeNow()
	if err != nil {
		return price.Price{}, fmt.Errorf("getting current time: %w", err)
	}

	return price.Price{
		Value:     steerData.Price,
		Timestamp: now,
	}, nil
}

//export iteration
func iteration(inputPtr, secretPtr uint64) uint64 {
	var args ArgsIterate
	inputData := as.Bytes(inputPtr)
	err := json.Unmarshal(inputData, &args)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal args args: %w", err)
		return writeErr(outErr.Error())
	}

	priceSamples, err := extractPriceSamples(args.EAttest, args.TAttest, args.Whitelist)
	if err != nil {
		outErr := fmt.Errorf("extracting priceSamples: %w", err)
		return writeErr(outErr.Error())
	}

	newPriceSample, err := getNewPriceSample(args.TokenAddress, args.ChainID)
	if err != nil {
		outErr := fmt.Errorf("getting new sample: %w", err)
		return writeErr(outErr.Error())
	}

	nextPriceSamples := append(priceSamples, newPriceSample)
	if len(nextPriceSamples) > args.NumSamples {
		nextPriceSamples = nextPriceSamples[1:]
	}

	return writePriceSamples(nextPriceSamples)
}

//export twap
func twap(inputPtr, secretPtr uint64) uint64 {
	var args ArgsTWAP
	inputData := as.Bytes(inputPtr)
	err := json.Unmarshal(inputData, &args)
	if err != nil {
		outErr := fmt.Errorf("could not unmarshal args args: %w", err)
		return writeErr(outErr.Error())
	}

	priceSamples, err := extractPriceSamples(args.EAttest, args.TAttest, args.Whitelist)
	if err != nil {
		outErr := fmt.Errorf("extracting samples: %w", err)
		return writeErr(outErr.Error())
	}

	twap, err := price.TWAP(priceSamples)
	if err != nil {
		outErr := fmt.Errorf("computing TWAP: %w", err)
		return writeErr(outErr.Error())
	}

	return writeTWAP(twap)
}

func main() {}

func writeErr(err string) uint64 {
	result := Result{
		Success: false,
		Error:   err,
	}
	return writeOutput(result)
}

func writePriceSamples(priceSamples PriceSamples) uint64 {
	result := ResultPriceSamples{
		Result: Result{
			Success: true,
			Error:   "",
		},
		PriceSamples: priceSamples,
	}
	return writeOutput(result)
}

func writeTWAP(twap float64) uint64 {
	result := ResultTWAP{
		Result: Result{
			Success: true,
			Error:   "",
		},
		TWAP: twap,
	}
	return writeOutput(result)
}

func writeOutput(output any) uint64 {
	outputData, err := as.Marshal(output)
	if err != nil {
		// We panic on errors we cannot communicate back to function caller
		panic("Fatal error: could not marshal output data")
	}
	return as.ShareWithHost(outputData)
}
