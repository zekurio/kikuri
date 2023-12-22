package hashutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	uid1 = "352002717285089280"
	uid2 = "531861558834495498"
	uid3 = "852813245815324672"
)

func TestSnowflake(t *testing.T) {
	{
		salt := []byte("pepper")

		hash1, err := HashSnowflake(uid1, salt)
		assert.Nil(t, err)

		hash2, err := HashSnowflake(uid1, salt)
		assert.Nil(t, err)

		assert.Equal(t, hash1, hash2)
	}

	{
		salt := []byte("salt")

		hash1, err := HashSnowflake(uid1, salt)
		assert.Nil(t, err)

		hash2, err := HashSnowflake(uid2, salt)
		assert.Nil(t, err)

		assert.NotEqual(t, hash1, hash2)
	}

	{
		salt1 := []byte("pepper2")
		salt2 := []byte("salt2")

		hash1, err := HashSnowflake(uid1, salt1)
		assert.Nil(t, err)

		hash2, err := HashSnowflake(uid1, salt2)
		assert.Nil(t, err)

		assert.NotEqual(t, hash1, hash2)
	}
}

func TestSumMD5(t *testing.T) {
	hash, err := SumMD5("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "098f6bcd4621d373cade4e832627b4f6"
	if hash != expected {
		t.Errorf("Expected %s, got %s", expected, hash)
	}
}

func TestSumSHA256(t *testing.T) {
	hash, err := SumSHA256("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if hash != expected {
		t.Errorf("Expected %s, got %s", expected, hash)
	}
}
