package hashutils

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"

	"github.com/bwmarrin/snowflake"
)

// HashSnowflake takes a snowflake as well as a salt value and returns a hash
func HashSnowflake(s string, salt []byte) (hash string, err error) {
	sid, err := snowflake.ParseString(s)
	if err != nil {
		return
	}

	idb := big.NewInt(sid.Int64() & int64(^uint(0)>>(64-48))).Bytes()
	comb := append(idb, salt...)
	hash = fmt.Sprintf("%x", sha256.Sum256(comb))

	return
}

// Sum returns a hash string from the given object v using
// the given hash function.
func Sum[TVal any](v TVal, hasher hash.Hash) (hash string, err error) {
	if _, err = hasher.Write([]byte(fmt.Sprintf("%v", v))); err != nil {
		return
	}

	hash = fmt.Sprintf("%x", hasher.Sum(nil))
	return
}

// SumMD5 returns a hash string from the given object v using
// the MD5 hash function.
func SumMD5[TVal any](v TVal) (hash string, err error) {
	return Sum(v, md5.New())
}

// SumSHA256 returns a hash string from the given object v using
// the SHA256 hash function.
func SumSHA256[TVal any](v TVal) (hash string, err error) {
	return Sum(v, sha256.New())
}

// Must is a helper function that wraps a call to a function returning
// (string, error) and panics if the error is non-nil.
func Must(hash string, err error) string {
	if err != nil {
		panic(err)
	}

	return hash
}
