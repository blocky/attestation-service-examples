package main

import (
	"errors"

	"github.com/blocky/as-demo/internal"
)

//export successFunc
func successFunc(inputPtr, secretPtr uint64) uint64 {
	type Output struct {
		Number int `json:"number"`
	}
	output := Output{42}
	return internal.WriteOutput(output)
}

//export errorFunc
func errorFunc(inputPtr, secretPtr uint64) uint64 {
	err := errors.New("expected error")
	return internal.WriteError(err)
}

func main() {}
