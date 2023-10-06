package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/database/dberr"
	"github.com/zekurio/daemon/internal/util"
	"github.com/zekurio/daemon/internal/util/embedded"
	"github.com/zekurio/daemon/internal/util/static"
)

type DBRefreshTokenHandler struct {
	db      database.Database
	session *discordgo.Session
}

func NewDBRefreshTokenHandler(ctn di.Container) *DBRefreshTokenHandler {
	return &DBRefreshTokenHandler{
		db:      ctn.Get(static.DiDatabase).(database.Database),
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
	}
}

func (d *DBRefreshTokenHandler) GetRefreshToken(ident string) (token string, err error) {
	token, err = util.GetRandBase64Str(64)
	if err != nil {
		return
	}

	err = d.db.SetUserRefreshToken(ident, token, time.Now().Add(static.AuthSessionExpiration))
	return
}

func (rth *DBRefreshTokenHandler) ValidateAccessToken(token string) (ident string, err error) {
	ident, expires, err := rth.db.GetUserByRefreshToken(token)
	if err != nil {
		return
	}

	if time.Now().After(expires) {
		err = errors.New("expired")
	}

	u, _ := rth.session.User(ident)
	if u == nil {
		err = errors.New("invalid user")
		return
	}

	return
}

func (h *DBRefreshTokenHandler) RevokeToken(ident string) error {
	err := h.db.RevokeUserRefreshToken(ident)
	if err == dberr.ErrNotFound {
		err = nil
	}

	return err
}

type JWTAccessTokenHandler struct {
	sessionExpiration time.Duration
	sessionSecret     []byte
}

func NewJWTAccessTokenHandler(ctn di.Container) *JWTAccessTokenHandler {
	cfg := ctn.Get(static.DiConfig).(config.Config)
	return &JWTAccessTokenHandler{
		sessionExpiration: time.Duration(cfg.Webserver.AccessToken.LifetimeSeconds) * time.Second,
		sessionSecret:     []byte(cfg.Webserver.AccessToken.Secret),
	}
}

func (ath *JWTAccessTokenHandler) GetAccessToken(ident string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(ath.sessionExpiration)

	claims := jwt.RegisteredClaims{}
	claims.Subject = ident
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.Issuer = fmt.Sprintf("daemon v%s", embedded.AppVersion)
	claims.NotBefore = jwt.NewNumericDate(now)

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(ath.sessionSecret)

	return
}

func (ath *JWTAccessTokenHandler) ValidateAccessToken(token string) (ident string, err error) {
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

type apiTokenClaims struct {
	jwt.RegisteredClaims

	Salt string `json:"sp_salt,omitempty"`
}

func apiTokenClaimsFromMap(m jwt.MapClaims) apiTokenClaims {
	c := apiTokenClaims{
		RegisteredClaims: registeredClaimsFromMap(m),
	}

	c.Salt, _ = m["sp_salt"].(string)

	return c
}

func registeredClaimsFromMap(m jwt.MapClaims) jwt.RegisteredClaims {
	c := jwt.RegisteredClaims{}

	c.Subject, _ = m["sub"].(string)
	jwtExpDate, _ := m["exp"].(int64)
	c.ExpiresAt = jwt.NewNumericDate(time.Unix(jwtExpDate, 0))
	jwtIssuedDate, _ := m["iat"].(int64)
	c.IssuedAt = jwt.NewNumericDate(time.Unix(jwtIssuedDate, 0))
	c.Issuer, _ = m["iss"].(string)
	jwtNotBefore, _ := m["nbf"].(int64)
	c.NotBefore = jwt.NewNumericDate(time.Unix(jwtNotBefore, 0))

	return c
}

type DBAPITokenHandler struct {
	db      database.Database
	session *discordgo.Session
	secret  []byte
}

func NewDBAPITokenHandler(ctn di.Container) *DBAPITokenHandler {
	cfg := ctn.Get(static.DiConfig).(config.Config)

	return &DBAPITokenHandler{
		db:      ctn.Get(static.DiDatabase).(database.Database),
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
		secret:  []byte(cfg.Webserver.APITokenKey),
	}
}

func (apith *DBAPITokenHandler) GetAPIToken(ident string) (token string, expires time.Time, err error) {
	now := time.Now()
	expires = now.Add(static.AuthSessionExpiration)

	claims := apiTokenClaims{}
	claims.Subject = ident
	claims.ExpiresAt = jwt.NewNumericDate(expires)
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.Issuer = fmt.Sprintf("daemon v%s", embedded.AppVersion)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.Salt, err = util.GetRandBase64Str(64)
	if err != nil {
		return
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(apith.secret)

	return
}

func (apith *DBAPITokenHandler) ValidateAPIToken(token string) (ident string, err error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return apith.secret, nil
	})
	if jwtToken == nil && err != nil {
		return "", err
	}
	if !jwtToken.Valid || jwtToken.Claims.Valid() != nil {
		return "", nil
	}

	claimsMap, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}

	claims := apiTokenClaimsFromMap(claimsMap)

	tokenEntry, err := apith.db.GetAPIToken(claims.Subject)
	if err == dberr.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", err
	}

	if tokenEntry.Salt != claims.Salt {
		return "", err
	}

	tokenEntry.LastAccess = time.Now()
	tokenEntry.Hits++
	apith.db.SetAPIToken(tokenEntry)

	return claims.Subject, nil
}
