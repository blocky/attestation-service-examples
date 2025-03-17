package main

import (
	"errors"
)

//export successFunc
func successFunc(inputPtr uint64, secretPtr uint64) uint64 {
	type Output struct {
		Number int `json:"number"`
	}
	output := Output{Number: 42}
	return WriteOutput(output)
}

//export errorFunc
func errorFunc(inputPtr uint64, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	return WriteError(err)
}

func main() {}
