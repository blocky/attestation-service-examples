package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blocky/as-demo/as"
)

type Result struct {
	Success bool
	Error   string
	Value   any
}

func (r Result) jsonMarshalWithError(err error) []byte {
	if err == nil {
		err = errors.New("jsonMarshalWithError invoked with nil error")
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
		as.Log(fmt.Sprintf("Error marshalling Result: %v", err))
		return WriteError(err)
	}
	return as.WriteToHost(data)
}

func WriteError(err error) uint64 {
	data := Result{}.jsonMarshalWithError(err)
	return as.WriteToHost(data)
}
