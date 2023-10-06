package auth

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/util"
	"github.com/zekurio/daemon/internal/util/static"
)

type DBRefreshTokenHandler struct {
	db      database.Database
	session *discordgo.Session
}

func New(ctn di.Container) *DBRefreshTokenHandler {
	return &DBRefreshTokenHandler{
		db:      ctn.Get(static.DiDatabase).(database.Database),
		session: ctn.Get(static.DiDiscordSession).(*discordgo.Session),
	}
}

func (h *DBRefreshTokenHandler) GetRefreshToken(ident string) (token string, err error) {
	token, err = util.GetRandBase64Str(64)
	if err != nil {
		return
	}

	err = h.db.SetUserRefreshToken(ident, token, time.Now().Add(static.AuthSessionExpiration))
	return
}
