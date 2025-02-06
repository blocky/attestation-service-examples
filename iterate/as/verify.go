package as

import (
	"encoding/json"
	"errors"
	"runtime"
)

//go:wasmimport env verifyAttestation
func _hostFuncVerify(ptr, size uint32) uint64

type VerifyInput struct {
	EAttest   json.RawMessage `json:"eAttest"`
	TAttest   json.RawMessage `json:"tAttest"`
	Whitelist json.RawMessage `json:"whitelist"`
}

type VerifyOutput struct {
	RawClaims []byte `json:"raw_clamis"`
	Error     string `json:"error"`
}

func VerifyAttestation(eAttest, tAttest, whitelist json.RawMessage) ([]byte, error) {
	LogToHost("starting verification\n")
	in := VerifyInput{
		EAttest:   eAttest,
		TAttest:   tAttest,
		Whitelist: whitelist,
	}
	data, err := json.Marshal(in)
	if err != nil {
		return nil, errors.New("marshaling verify input: " + err.Error())
	}

	offset, length := bytesToPtr(data)
	outPtr := _hostFuncVerify(offset, length)
	runtime.KeepAlive(data)

	var out VerifyOutput
	outData := bytesFromPtr(outPtr)
	err = json.Unmarshal(outData, &out)
	switch {
	case err != nil:
		return nil, errors.New("unmarshaling verify output: " + err.Error())
	case out.Error != "":
		return nil, errors.New(out.Error)
	}

	return out.RawClaims, nil
}
