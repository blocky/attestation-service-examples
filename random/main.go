package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"

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

	randBytes := make([]byte, 8)
	n, err := crand.Read(randBytes)
	switch {
	case err != nil:
		outErr := fmt.Errorf("reading random bytes: %w", err)
		return WriteError(outErr)
	case n != 8:
		outErr := fmt.Errorf("reading random bytes: expected 8 bytes, got %d", n)
		return WriteError(outErr)
	}

	roll := (binary.BigEndian.Uint64(randBytes) % input.Sides) + 1
	return WriteOutput(roll)
}

func main() {}
