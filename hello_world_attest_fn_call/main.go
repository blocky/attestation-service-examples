package main

import (
	"fmt"

	"github.com/blocky/as-demo/as"
)

//export helloWorld
func helloWorld(inputPtr, secretPtr uint64) uint64 {
	msg := "Hello, World!"

	as.Log(fmt.Sprintf("Writing \"%s\" to host\n", msg))

	return as.WriteToHost([]byte(msg))
}

func main() {}
