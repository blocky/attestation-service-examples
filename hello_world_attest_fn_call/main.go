package main

import (
	"fmt"

	"github.com/blocky/basm-go-sdk"
)

//export helloWorld
func helloWorld(inputPtr, secretPtr uint64) uint64 {
	msg := "Hello, World!"

	basm.Log(fmt.Sprintf("Writing \"%s\" to host\n", msg))

	return basm.WriteToHost([]byte(msg))
}

func main() {}
