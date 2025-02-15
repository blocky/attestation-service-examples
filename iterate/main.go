package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blocky/as-demo/as"
	"github.com/blocky/as-demo/price"
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

type Window struct {
	TWAP    float64
	Samples []price.Price
}

type Result struct {
	Success bool
	Value   Window
	Error   string
}

func extractPriceSamples(
	eAttest,
	tAttest,
	whitelist json.RawMessage,
) (
	[]price.Price,
	error,
) {
	// bootstrap with empty samples if we don't have a transitive attestation
	if tAttest == nil {
		return []price.Price{}, nil
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
			"unmarshaling coin gecko data: %w...%s",
			err,
			resp.Body,
		)
	}

	return price.Price{
		Value:     steerData.Price,
		Timestamp: time.Now(),
	}, nil
}

func advanceWindow(input Args) (Window, error) {
	samples, err := extractPriceSamples(input.EAttest, input.TAttest, input.Whitelist)
	if err != nil {
		return Window{}, fmt.Errorf("extracting samples: %w", err)
	}

	newPriceSample, err := getNewPriceSample(input.TokenAddress, input.ChainID)
	if err != nil {
		return Window{}, fmt.Errorf("getting new sample %w: ", err)
	}

	nextPriceSamples := append(samples, newPriceSample)
	if len(nextPriceSamples) > 5 {
		nextPriceSamples = nextPriceSamples[:5]
	}

	twap, err := price.TWAP(time.Now(), nextPriceSamples)
	if err != nil {
		return Window{}, fmt.Errorf("computing average: %w", err)
	}

	next := Window{
		TWAP:    twap,
		Samples: nextPriceSamples,
	}
	return next, nil
}

//export iteration
func iteration(inputPtr, secretPtr uint64) uint64 {
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
