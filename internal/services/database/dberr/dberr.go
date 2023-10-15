package dberr

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
