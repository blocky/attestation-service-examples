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
	Output  any
}

func (r Result) jsonMarshalWithError(err string) []byte {
	resultStr := fmt.Sprintf(
		`{
					"Success": false,
					"Error": "%s",
					"Output": null
				}`,
		err,
	)
	data := []byte(resultStr)
	return data
}

func writeOutput(output any) uint64 {
	result := Result{
		Success: true,
		Error:   "",
		Output:  output,
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
	successMsg := "Output from successFunc"
	return writeOutput(successMsg)
}

//export errorFunc
func errorFunc(inputPtr, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	as.Log(fmt.Sprintf("Error in errorFunc: %s", err.Error()))
	// todo: add log to host to show output in server
	return writeError(err.Error())
}

func main() {}
