package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/zekurio/kikuri/internal/embedded"

	"github.com/zekurio/kikuri/internal/models"

	"github.com/golang-jwt/jwt/v5"

	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/util/static"
)

var (
	jwtGenerationMethod = jwt.SigningMethodHS256
)

type AccessTokenHandlerImpl struct {
	sessionExpiration time.Duration
	sessionSecret     []byte
}

func NewAccessTokenHandlerImpl(container di.Container) *AccessTokenHandlerImpl {
	cfg := container.Get(static.DiConfig).(models.Config)

	return &AccessTokenHandlerImpl{
		sessionExpiration: time.Duration(cfg.Webserver.AccessToken.LifetimeSeconds) * time.Second,
		sessionSecret:     []byte(cfg.Webserver.AccessToken.Secret),
	}
}

func (ath *AccessTokenHandlerImpl) GetAccessToken(ident string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(ath.sessionExpiration)

	claims := jwt.RegisteredClaims{}
	claims.Issuer = fmt.Sprintf("kikuri v.%s", embedded.AppVersion)
	claims.Subject = ident
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.IssuedAt = jwt.NewNumericDate(now)

	token, err = jwt.NewWithClaims(jwtGenerationMethod, claims).
		SignedString(ath.sessionSecret)
	return
}

func (ath *AccessTokenHandlerImpl) ValidateAccessToken(token string) (ident string, err error) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(ath.sessionSecret), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		ident, ok = claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid claims")
		}
	} else {
		return "", errors.New("invalid token")
	}

	return ident, nil
}
