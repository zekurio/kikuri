package auth

import "time"

// Provides functionalities to manage refresh tokens.
type RefreshTokenHandler interface {

	// Generates a unique and secure token which can be used to recover the passed ident from it.
	GetRefreshToken(ident string) (token string, err error)

	// Validates a token. Returns an error if the token is invalid, otherwise returns the ident.
	ValidateRefreshToken(token string) (ident string, err error)

	// Marks the token linked to the passed ident as invalid.
	RevokeToken(ident string) error
}

// Provides functionalities to manage access tokens.
type AccessTokenHandler interface {

	// Generates a unique and secure token which can be used to recover the passed ident from it.
	// Also returns an expiration time after which the token will become invalid.
	GetAccessToken(ident string) (token string, expires time.Time, err error)

	// Validates a token. Returns an error if the token is invalid, otherwise returns the ident.
	ValidateAccessToken(token string) (ident string, err error)
}

// Provides functionalities to manage API tokens.
type APITokenHandler interface {

	// Generates a unique and secure token which can be used to recover the passed ident from it.
	// Also returns an expiration time after which the token will become invalid.
	GetAPIToken(ident string) (token string, expires time.Time, err error)

	// Validates a token. Returns an error if the token is invalid, otherwise returns the ident.
	ValidateAPIToken(token string) (ident string, err error)

	// Marks the token linked to the passed ident as invalid.
	RevokeToken(ident string) error
}
