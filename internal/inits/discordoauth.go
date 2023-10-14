package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/pkg/discordoauth"
)

func InitDiscordOAuth(ctn di.Container) *discordoauth.DiscordOAuth {
	return &discordoauth.DiscordOAuth{}
}
