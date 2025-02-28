package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blocky/as-demo/as"
)

type Result struct {
	Success bool
	Value   any
}

func (r Result) jsonMarshalWithError(err string) []byte {
	resultStr := fmt.Sprintf(
		`{
					"Success": false,
					"Value": "%s"
				}`,
		err,
	)
	data := []byte(resultStr)
	return data
}

func writeOutput(output any) uint64 {
	result := Result{
		Success: true,
		Value:   output,
	}
	data, err := json.Marshal(result)
	if err != nil {
		as.Log(fmt.Sprintf("Error marshalling result: %s", err.Error()))
		return writeError(err.Error())
	}
	return as.WriteToHost(data)
}

func writeError(err string) uint64 {
	data := Result{}.jsonMarshalWithError(err)
	return as.WriteToHost(data)
}

//export successFunc
func successFunc(inputPtr, secretPtr uint64) uint64 {
	type Output struct {
		Number int `json:"number"`
	}
	output := Output{42}
	return writeOutput(output)
}

//export errorFunc
func errorFunc(inputPtr, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	return writeError(err.Error())
}

func main() {}
