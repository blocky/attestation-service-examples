package main

import (
	"fmt"

	"github.com/blocky/basm-go-sdk/basm"
)

//export helloWorld
func helloWorld(inputPtr uint64, secretPtr uint64) uint64 {
	msg := "Hello, World!"

	basm.Log(fmt.Sprintf("Writing \"%s\" to host\n", msg))

	return basm.WriteToHost([]byte(msg))
}

func main() {}
