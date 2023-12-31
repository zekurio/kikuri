package inits

import (
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/webserver/auth"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordoauth"
)

func InitDiscordOAuth(ctn di.Container) *discordoauth.DiscordOAuth {
	cfg := ctn.Get(static.DiConfig).(models.Config)
	oauthHandler := ctn.Get(static.DiOAuthHandler).(auth.RequestHandler)

	doa, err := discordoauth.New(
		cfg.Discord.ClientID,
		cfg.Discord.ClientSecret,
		cfg.Webserver.PublicAddr+static.EndpointAuthCB,
		oauthHandler.LoginFailedHandler,
		oauthHandler.LoginSuccessHandler,
	)

	if err != nil {
		log.Fatal("Discord OAuth initialization failed")
	}

	return doa
}
