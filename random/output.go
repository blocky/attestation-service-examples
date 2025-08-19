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

func WriteOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		basm.Log(fmt.Sprintf("Error marshalling output Result: %v", err))
		return WriteError(err)
	}
	return basm.WriteToHost(data)
}

func WriteError(err error) uint64 {
	if err == nil {
		err = errors.New("WriteError called with nil error")
	}

	result := Result{
		Success: false,
		Error:   err.Error(),
		Value:   nil,
	}
	data, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		basm.Log(fmt.Sprintf("Error marshalling error Result: %v", marshalErr))
		data = []byte(`{ "success": false, "error": "failed to marshal result" , "value": null }`)
	}
	return basm.WriteToHost(data)
}
