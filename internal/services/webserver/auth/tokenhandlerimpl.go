package auth

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/util/static"
)

type RefreshTokenHandlerImpl struct {
	db      database.Database
	session *discordgo.Session
}

func NewRefreshTokenHandlerImpl(container di.Container) *RefreshTokenHandlerImpl {
	return &RefreshTokenHandlerImpl{
		db:      container.Get(static.DiDatabase).(database.Database),
		session: container.Get(static.DiDiscordSession).(*discordgo.Session),
	}
}

func (rth *RefreshTokenHandlerImpl) GetRefreshToken(ident string) (token string, err error) {
	return "", nil
}

func (rth *RefreshTokenHandlerImpl) ValidateRefreshToken(token string) (ident string, err error) {
	return "", nil
}

func (rth *RefreshTokenHandlerImpl) RevokeToken(ident string) error {
	return nil
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
	return "", time.Time{}, nil
}

func (ath *AccessTokenHandlerImpl) ValidateAccessToken(token string) (ident string, err error) {
	return "", nil
}
