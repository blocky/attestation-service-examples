package main

import (
	"encoding/json"
	"fmt"

	"github.com/blocky/basm-go-sdk/basm"
)

type Result struct {
	Success bool
	Error   string
	Value   any
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
	result := Result{
		Success: false,
		Error:   err.Error(),
		Value:   nil,
	}
	data, marshalErr := json.Marshal(result)
	if marshalErr != nil {
		basm.Log(fmt.Sprintf("Error marshalling Result: %v", marshalErr))
	}
	return basm.WriteToHost(data)
}
