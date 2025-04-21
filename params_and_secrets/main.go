package main

import (
	"encoding/json"
	"fmt"

	"github.com/blocky/basm-go-sdk/basm"
)

type Params struct {
	Data string `json:"data"`
}

type Secrets struct {
	Password string `json:"password"`
}

type Output struct {
	DecryptedData string `json:"decrypted_data"`
}

//export extractInputs
func extractInputs(inputPtr uint64, secretPtr uint64) uint64 {
	var params Params
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

	result, err := decryptData(params.Data, secrets.Password)
	if err != nil {
		outErr := fmt.Errorf("decrypting data: %w", err)
		return WriteError(outErr)
	}

	output := Output{DecryptedData: result}

	return WriteOutput(output)
}

func main() {}
