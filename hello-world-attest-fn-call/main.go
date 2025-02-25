package main

import (
	"github.com/blocky/as-demo/as"
)

//export helloWorld
func helloWorld(inputPtr, secretPtr uint64) uint64 {
	return as.WriteOutput("Hello, World!")
}

//export helloError
func helloError(inputPtr, secretPtr uint64) uint64 {
	as.Log("Returning an expected error")
	return as.WriteError("expected error")
}

func main() {}
