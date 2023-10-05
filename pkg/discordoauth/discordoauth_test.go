package discordoauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndSignWithPayload(t *testing.T) {
	d := DiscordOAuth{
		clientID:        "test_client_id",
		clientSecret:    "test_client_secret",
		redirectURI:     "http://localhost:8080/callback",
		stateSigningKey: []byte("test_state_signing_key"),
	}

	payload := map[string]string{
		"test_key": "test_value",
	}

	token, err := d.encodeAndSignWithPayload(payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := d.decodeAndVerifyStateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "test_client_id", claims.Issuer)
	assert.Equal(t, payload, claims.Payload)
}

func TestDecodeAndVerifyStateTokenInvalidToken(t *testing.T) {
	d := DiscordOAuth{
		clientID:        "test_client_id",
		clientSecret:    "test_client_secret",
		redirectURI:     "http://localhost:8080/callback",
		stateSigningKey: []byte("test_state_signing_key"),
	}

	_, err := d.decodeAndVerifyStateToken("invalid_token")
	assert.Error(t, err)
}

func TestDecodeAndVerifyStateTokenInvalidSigningKey(t *testing.T) {
	d1 := DiscordOAuth{
		clientID:        "test_client_id",
		clientSecret:    "test_client_secret",
		redirectURI:     "http://localhost:8080/callback",
		stateSigningKey: []byte("test_state_signing_key_1"),
	}

	d2 := DiscordOAuth{
		clientID:        "test_client_id",
		clientSecret:    "test_client_secret",
		redirectURI:     "http://localhost:8080/callback",
		stateSigningKey: []byte("test_state_signing_key_2"),
	}

	payload := map[string]string{
		"test_key": "test_value",
	}

	token, err := d1.encodeAndSignWithPayload(payload)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	_, err = d2.decodeAndVerifyStateToken(token)
	assert.Error(t, err)
}
