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

func getNewSample(tokenAddress string, chainID string) (Sample, error) {
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
		return Sample{}, fmt.Errorf("making http request: %w", err)
	}

	steerData := struct {
		Price float64 `json:"price"`
	}{}

	err = json.Unmarshal(resp.Body, &steerData)
	if err != nil {
		return Sample{}, fmt.Errorf(
			"unmarshaling coin gecko data: %w...%s",
			err,
			resp.Body,
		)
	}

	return Sample{
		Price:     steerData.Price,
		Timestamp: time.Now(),
	}, nil
}

func average(samples []Sample) (float64, error) {
	acc := float64(0)
	n := float64(len(samples))
	for _, s := range samples {
		acc += s.Price
	}

	return acc / n, nil
}

func advanceWindow(input Args) (Window, error) {
	samples, err := extractSamples(input.EAttest, input.TAttest, input.Whitelist)
	if err != nil {
		return Window{}, fmt.Errorf("extracting samples: %w", err)
	}

	newSample, err := getNewSample(input.TokenAddress, input.ChainID)
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

	nextWindow, err := advanceWindow(input)
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
