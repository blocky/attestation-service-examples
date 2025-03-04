package main

import (
	"errors"
)

//export successFunc
func successFunc(inputPtr, secretPtr uint64) uint64 {
	type Output struct {
		Number int `json:"number"`
	}
	output := Output{Number: 42}
	return writeOutput(output)
}

//export errorFunc
func errorFunc(inputPtr, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	return writeError(err)
}

func main() {}
