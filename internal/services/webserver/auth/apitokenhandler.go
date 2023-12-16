package auth

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/embedded"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/randutils"
)

type APITokenHandlerImpl struct {
	db      database.Database
	session *discordgo.Session
	secret  []byte
}

func NewAPITokenHandlerImpl(ctn di.Container) *APITokenHandlerImpl {
	cfg := ctn.Get(static.DiConfig).(models.Config)
	secret := []byte(cfg.Webserver.APITokenKey)

	return &APITokenHandlerImpl{
		db:      ctn.Get(static.DiDatabase).(database.Database),
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
		secret:  secret,
	}
}

func (ath *APITokenHandlerImpl) GetAPIToken(ident string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(static.ApiTokenExpiration)

	salt, err := randutils.GetRandBase64Str(16)
	if err != nil {
		return
	}

	claims := jwt.RegisteredClaims{}
	claims.Issuer = fmt.Sprintf("kikuri v.%s", embedded.AppVersion)
	claims.Subject = ident
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.IssuedAt = jwt.NewNumericDate(now)

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString(ath.secret)
	if err != nil {
		return
	}

	tokenEntry := models.APITokenEntry{
		Salt:    salt,
		Created: now,
		Expires: expires,
		UserID:  ident,
	}

	err = ath.db.SetAPIToken(tokenEntry)

	return
}

func (ath *APITokenHandlerImpl) ValidateAPIToken(token string) (ident string, err error) {
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			err = fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			return nil, err
		}
		return ath.secret, nil
	})

	if err != nil {
		return "", err
	}

	// validate the token from the database
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		ident, ok = claims["sub"].(string)
		if !ok {
			return "", fmt.Errorf("invalid claims")
		}

		tokenEntry, err := ath.db.GetAPIToken(ident)
		if err != nil {
			return "", err
		}

		if tokenEntry.Salt != claims["salt"].(string) {
			return "", fmt.Errorf("invalid token")
		}

		if tokenEntry.Expires.Before(time.Now()) {
			return "", fmt.Errorf("token expired")
		}

		return ident, nil
	}

	return "", fmt.Errorf("invalid token")
}

func (ath *APITokenHandlerImpl) RevokeToken(ident string) error {
	return ath.db.DeleteAPIToken(ident)
}
