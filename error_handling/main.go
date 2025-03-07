package main

import (
	"errors"

	"github.com/blocky/basm-go-sdk"
)

//export successFunc
func successFunc(inputPtr, secretPtr uint64) uint64 {
	type Output struct {
		Number int `json:"number"`
	}
	output := Output{Number: 42}
	return WriteOutput(output)
}

//export errorFunc
func errorFunc(inputPtr, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	return WriteError(err)
}

//export panicFunc
func panicFunc(inputPtr, secretPtr uint64) uint64 {
	msg := "expected panic"
	basm.LogToHost(msg)
	err := errors.New("expected error")
	return WriteError(err)
	// todo: clean up this example and move the log examples here
}

func main() {}
