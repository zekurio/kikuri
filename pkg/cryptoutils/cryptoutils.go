package cryptoutils

import (
	"crypto/rand"
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
