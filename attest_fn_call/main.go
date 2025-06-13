package main

import (
	"github.com/blocky/basm-go-sdk/basm"
)

//export helloWorld
func helloWorld(inputPtr uint64, secretPtr uint64) uint64 {
	msg := "Hello, World!"

	return basm.WriteToHost([]byte(msg))
}

func main() {}
