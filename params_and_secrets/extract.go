package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

func extractData(inputData, pwd string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(inputData)
	if err != nil {
		return "", err
	}

	pad := func(key string) []byte {
		p := make([]byte, 32)
		copy(p, key)
		return p
	}

	block, err := aes.NewCipher(pad(pwd))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	extractedData, err := aesGCM.Open(nil, data[:12], data[12:], nil)
	if err != nil {
		return "", errors.New("incorrect password")
	}

	return string(extractedData), nil
}
