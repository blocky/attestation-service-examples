package as

import (
	"encoding/json"
	"errors"
	"runtime"
)

type HostVerifyAttestationInput struct {
	EnclaveAttestedKey    json.RawMessage `json:"enclave_attested_app_public_key"`
	TransitiveAttestation json.RawMessage `json:"transitive_attestation"`
	AcceptableMeasures    json.RawMessage `json:"acceptable_measurements"`
}

type HostVerifyAttestationResult struct {
	IsOk  bool                        `json:"ok"`
	Value HostVerifyAttestationOutput `json:"value"`
	Error string                      `json:"error"`
}

type HostVerifyAttestationOutput struct {
	RawClaims []byte `json:"raw_claims"`
}

//go:wasmimport env verifyAttestation
func _hostFuncVerifyAttestation(offset, size uint32) uint64

func VerifyAttestation(
	input HostVerifyAttestationInput,
) (HostVerifyAttestationOutput, error) {
	inputData, err := Marshal(input)
	if err != nil {
		msg := "marshaling input data: " + err.Error()
		return HostVerifyAttestationOutput{}, errors.New(msg)
	}

	inOffset, inSize := bytesToOffsetSize(inputData)
	resultPtr := _hostFuncVerifyAttestation(inOffset, inSize)
	runtime.KeepAlive(inputData)
	resultData := bytesFromFatPtr(resultPtr)

	var result HostVerifyAttestationResult
	err = Unmarshal(resultData, &result)
	switch {
	case err != nil:
		msg := "unmarshaling result data: " + err.Error()
		return HostVerifyAttestationOutput{}, errors.New(msg)
	case !result.IsOk:
		msg := "host fn returned error: " + result.Error
		return HostVerifyAttestationOutput{}, errors.New(msg)
	}
	return result.Value, nil
}
