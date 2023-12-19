package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/listeners"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/util/static"
)

func InitDiscord(ctn di.Container) (err error) {

	log.Info("Initializing bot session ...")

	session := ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	cfg := ctn.Get(static.DiConfig).(models.Config)

	session.Token = "Bot " + cfg.Discord.Token
	session.Identify.Intents = discordgo.MakeIntent(static.Intents)
	session.StateEnabled = false

	session.AddHandler(listeners.NewListenerReady(ctn).Handler)
	session.AddHandler(listeners.NewListenerGuildCreate(ctn).Handler)
	session.AddHandler(listeners.NewListenerVoiceStateUpdate(ctn).Handler)

	err = session.Open()
	if err != nil {
		log.Fatal("Failed connecting Discord bot session", err)
	}

	return
}
