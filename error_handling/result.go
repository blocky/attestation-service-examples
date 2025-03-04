package main

import (
	"errors"
	"fmt"
)

type Result struct {
	Success bool
	Error   string
	Value   any
}

func (r Result) jsonMarshalWithError(err error) []byte {
	if err == nil {
		err = errors.New("jsonMarshalWithError invoked with nil error")
	}
	resultStr := fmt.Sprintf(
		`{ "Success": false, "Error": "%s" , "Value": null }`,
		err.Error(),
	)
	return []byte(resultStr)
}
