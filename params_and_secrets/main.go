package main

import (
	"encoding/json"
	"fmt"

	"github.com/blocky/basm-go-sdk"
)

type Parms struct {
	Data string `json:"data"`
}

type Secrets struct {
	Password string `json:"password"`
}

type Output struct {
	ExtractedData string `json:"extracted_data"`
}

//export extract
func extract(inputPtr uint64, secretPtr uint64) uint64 {
	var params Parms
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &params)
	if err != nil {
		outErr := fmt.Errorf("unmarshaling params: %w", err)
		return WriteError(outErr)
	}

	var secrets Secrets
	secretData := basm.ReadFromHost(secretPtr)
	err = json.Unmarshal(secretData, &secrets)
	if err != nil {
		outErr := fmt.Errorf("unmarshaling secrets: %w", err)
		return WriteError(outErr)
	}

	result, err := extractData(params.Data, secrets.Password)
	if err != nil {
		outErr := fmt.Errorf("extracting data: %w", err)
		return WriteError(outErr)
	}

	output := Output{ExtractedData: result}

	return WriteOutput(output)
}

func main() {}
