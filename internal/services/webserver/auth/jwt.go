package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims

	UserID   string
	Username string
	Scopes   []string
}
