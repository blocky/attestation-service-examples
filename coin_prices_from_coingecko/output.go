package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blocky/basm-go-sdk/basm"
)

type Result struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Value   any    `json:"value"`
}

func (r Result) JSONMarshalWithError(err error) []byte {
	if err == nil {
		err = errors.New("JSONMarshalWithError invoked with nil error")
	}
	resultStr := fmt.Sprintf(
		`{ "Success": false, "Error": "%s" , "Value": null }`,
		err.Error(),
	)
	return []byte(resultStr)
}

func WriteOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		basm.Log(fmt.Sprintf("Error marshalling Result: %v", err))
		return WriteError(err)
	}
	return basm.WriteToHost(data)
}

func WriteError(err error) uint64 {
	data := Result{}.JSONMarshalWithError(err)
	return basm.WriteToHost(data)
}
