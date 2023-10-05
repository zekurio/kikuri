package inits

import (
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/webserver/auth"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordoauth"
)

func InitDiscordOAuth(ctn di.Container) *discordoauth.DiscordOAuth {
	cfg := ctn.Get(static.DiConfig).(config.Config)
	oauthHandler := ctn.Get(static.DiOAuthHandler).(auth.RequestHandler)

	doa, err := discordoauth.New(
		cfg.Discord.ClientID,
		cfg.Discord.ClientSecret,
		cfg.Webserver.PublicAddr+static.EndpointAuthCB,
		oauthHandler.LoginFailedHandler,
		oauthHandler.LoginSuccessHandler,
	)

	if err != nil {
		log.Fatal("discord oauth init failed", err)
	}

	return doa
}
