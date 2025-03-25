package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/blocky/basm-go-sdk"
)

type Args struct {
	Sides uint64 `json:"sides"`
}

//export rollDie
func rollDie(inputPtr uint64, secretPtr uint64) uint64 {
	var input Args
	inputData := basm.ReadFromHost(inputPtr)
	err := json.Unmarshal(inputData, &input)
	switch {
	case err != nil:
		outErr := fmt.Errorf("could not unmarshal input args: %w", err)
		return WriteError(outErr)
	case input.Sides == 0:
		outErr := fmt.Errorf("die cannot have zero sides")
		return WriteError(outErr)
	}

	roll := rand.Intn(int(input.Sides)) + 1
	return WriteOutput(roll)
}

func main() {}
