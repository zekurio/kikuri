package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/zekurio/kikuri/internal/models"

	"github.com/golang-jwt/jwt/v4"

	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/util/embedded"
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
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return ath.sessionSecret, nil
	})
	if jwtToken == nil || err != nil || !jwtToken.Valid || jwtToken.Claims.Valid() != nil {
		return
	}

	claimsMap, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("invalid claims")
		return
	}

	ident, _ = claimsMap["sub"].(string)

	return
}
