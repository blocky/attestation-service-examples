package main

import (
	"time"
)

//export timeNow
func timeNow(_ uint64, _ uint64) uint64 {
	return WriteOutput(time.Now())
}

func main() {}
