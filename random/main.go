package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/blocky/basm-go-sdk"
)

type Args struct {
	DieSides int `json:"die_sides"`
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
	case input.DieSides < 1:
		outErr := fmt.Errorf("die must have one or more sides")
		return WriteError(outErr)
	}

	roll := rand.Intn(input.DieSides) + 1
	return WriteOutput(roll)
}

func main() {}
