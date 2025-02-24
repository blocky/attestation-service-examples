package main

import (
	"github.com/blocky/as-demo/as"
)

//export helloWorld
func helloWorld(_, _ uint64) uint64 {
	return as.WriteOutput("Hello, World!")
}

//export helloError
func helloError(_, _ uint64) uint64 {
	as.Logf("Returning an expected error")
	return as.WriteError("expected error")
}

func main() {}
