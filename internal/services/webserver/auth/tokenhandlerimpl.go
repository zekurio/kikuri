package auth

import (
	"errors"
	"github.com/zekurio/daemon/internal/services/database/dberr"
	"github.com/zekurio/daemon/pkg/cryptoutils"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/util/static"
)

type RefreshTokenHandlerImpl struct {
	db      database.Database
	st      *dgrs.State
	session *discordgo.Session
}

func NewRefreshTokenHandlerImpl(container di.Container) *RefreshTokenHandlerImpl {
	return &RefreshTokenHandlerImpl{
		db:      container.Get(static.DiDatabase).(database.Database),
		st:      container.Get(static.DiState).(*dgrs.State),
		session: container.Get(static.DiDiscordSession).(*discordgo.Session),
	}
}

// GetRefreshToken returns a refresh token for the given ident, and saves it to the database.
func (rth *RefreshTokenHandlerImpl) GetRefreshToken(ident string) (token string, err error) {
	token, err = cryptoutils.GetRandBase64Str(64)
	if err != nil {
		return
	}

	err = rth.db.SetUserRefreshToken(ident, token, time.Now().Add(static.AuthSessionExpiration))
	return
}

func (rth *RefreshTokenHandlerImpl) ValidateRefreshToken(token string) (ident string, err error) {
	ident, expires, err := rth.db.GetUserByRefreshToken(token)
	if err != nil {
		return
	}

	if expires.Before(time.Now()) {
		err = rth.RevokeToken(ident)
		return
	}

	user, err := rth.st.User(ident)
	if user == nil || err != nil {
		err = errors.New("user not found")
		return
	}

	return
}

func (rth *RefreshTokenHandlerImpl) RevokeToken(ident string) error {
	err := rth.db.RevokeUserRefreshToken(ident)
	if dberr.IsErrNotFound(err) {
		err = nil
	}
	return err
}

type AccessTokenHandlerImpl struct {
	sessionExpiration time.Duration
	sessionSecret     []byte
}

func NewAccessTokenHandlerImpl(container di.Container) *AccessTokenHandlerImpl {
	cfg := container.Get(static.DiConfig).(config.Config)

	return &AccessTokenHandlerImpl{
		sessionExpiration: time.Duration(cfg.Webserver.AccessToken.LifetimeSeconds) * time.Second,
		sessionSecret:     []byte(cfg.Webserver.AccessToken.Secret),
	}
}

func (ath *AccessTokenHandlerImpl) GetAccessToken(ident string) (token string, expires time.Time, err error) {
	return "", time.Time{}, nil // TODO implement
}

func (ath *AccessTokenHandlerImpl) ValidateAccessToken(token string) (ident string, err error) {
	return "", nil // TODO implement
}
