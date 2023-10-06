package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

var (
	ErrInvalidLength = errors.New("invalid length")
)

// GetRandBase64Str returns a random base64 string with the given length
func GetRandBase64Str(len int) (string, error) {
	if len <= 0 {
		return "", ErrInvalidLength
	}

	data, err := GetRandByteArray(len)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data)[:len], nil
}

// GetRandByteArray returns a random byte array with the given length
func GetRandByteArray(len int) ([]byte, error) {
	if len <= 0 {
		return nil, ErrInvalidLength
	}

	data := make([]byte, len)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
