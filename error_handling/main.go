package main

import (
	"errors"
	"fmt"

	"github.com/blocky/basm-go-sdk"
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

//export panicFunc
func panicFunc(inputPtr uint64, secretPtr uint64) uint64 {
	msg := "expected panic"
	basm.LogToHost(fmt.Sprintf("Logging \"%s\" to host\n", msg))
	panic(msg)
}

func main() {}
