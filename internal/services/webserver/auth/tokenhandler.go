package auth

import "time"

// RefreshTokenHandler is an interface for handling refresh tokens.
type RefreshTokenHandler interface {
	GetRefreshToken(ident string) (token string, err error)

	ValidateRefreshToken(token string) (ident string, err error)

	RevokeToken(ident string) error
}

// AccessTokenHandler is an interface for handling access tokens.
type AccessTokenHandler interface {
	GetAccessToken(ident string) (token string, expires time.Time, err error)

	ValidateAccessToken(token string) (ident string, err error)
}
