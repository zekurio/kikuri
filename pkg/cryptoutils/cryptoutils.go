package cryptoutils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
)

var ErrInvalidLen = errors.New("invalid length")

func GetRandByteArray(len int) (data []byte, err error) {
	if len <= 0 {
		return nil, ErrInvalidLen
	}

	data = make([]byte, len)
	_, err = rand.Read(data)

	return
}

func GetRandBase64Str(len int) (string, error) {
	if len <= 0 {
		return "", ErrInvalidLen
	}

	data, err := GetRandByteArray(len)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data)[:len], nil
}
