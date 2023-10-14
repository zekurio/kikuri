package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTHandler struct {
	secret   []byte
	lifetime time.Duration
	issuer   string
}

func NewJWTHandler(secret string, issuer string, lifetime time.Duration) (*JWTHandler, error) {
	if secret == "" {
		return nil, errors.New("secret cannot be empty")
	}

	return &JWTHandler{
		secret:   []byte(secret),
		lifetime: lifetime,
		issuer:   issuer,
	}, nil
}

func (h *JWTHandler) GenerateToken(userID string, username string, scopes []string) (string, error) {
	claims := jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"scopes":   scopes,
		"exp":      time.Now().Add(h.lifetime).Unix(),
		"iss":      h.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}

func (h *JWTHandler) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}

		return h.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if !claims.VerifyIssuer(h.issuer, true) {
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}
