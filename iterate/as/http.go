package as

import (
	"encoding/json"
	"errors"
	"runtime"
)

type HostHTTPRequestInput struct {
	Method  string                   `json:"method"`
	URL     string                   `json:"url"`
	Headers []HostHTTPRequestHeaders `json:"headers"`
	Body    []byte                   `json:"body"`
}

type HostHTTPRequestHeaders struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type HostHTTPRequestOutput struct {
	StatusCode int    `json:"status_code"`
	Body       []byte `json:"body"`
	Error      string `json:"error"`
}

func hostHTTPReqInToBytes(v HostHTTPRequestInput) ([]byte, error) {
	return json.Marshal(v)
}

func hostHTTPReqOutFromBytes(data []byte) (HostHTTPRequestOutput, error) {
	var v HostHTTPRequestOutput
	err := json.Unmarshal(data, &v)
	return v, err
}

//go:wasmimport env httpRequest
func _hostFuncHTTPRequest(ptr, size uint32) uint64

func HostFuncHTTPRequest(request HostHTTPRequestInput) (HostHTTPRequestOutput, error) {
	inputData, err := hostHTTPReqInToBytes(request)
	if err != nil {
		return HostHTTPRequestOutput{}, errors.New("marshaling request data: " + err.Error())
	}

	inOffset, inLen := bytesToPtr(inputData)
	outPtr := _hostFuncHTTPRequest(inOffset, inLen)
	runtime.KeepAlive(inputData)
	outputData := bytesFromPtr(outPtr)

	output, err := hostHTTPReqOutFromBytes(outputData)
	if err != nil {
		return HostHTTPRequestOutput{}, errors.New("unmarshaling result data: " + err.Error())
	}
	return output, nil
}
