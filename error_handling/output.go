package main

import (
	"encoding/json"
	"fmt"

	"github.com/blocky/as-demo/as"
)

func writeOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		as.Log(fmt.Sprintf("Error marshalling result: %v", err))
		return writeError(err)
	}
	return as.WriteToHost(data)
}

func writeError(err error) uint64 {
	data := Result{}.jsonMarshalWithError(err)
	return as.WriteToHost(data)
}
