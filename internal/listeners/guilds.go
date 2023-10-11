package listeners

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"

	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordutils"
)

type ListenerGuilds struct {
	cfg       config.Config
	st        *dgrs.State
	lockUntil *time.Time
}

func NewListenerGuildCreate(ctn di.Container) *ListenerGuilds {
	return &ListenerGuilds{
		cfg: ctn.Get(static.DiConfig).(config.Config),
		st:  ctn.Get(static.DiState).(*dgrs.State),
	}
}

func (l *ListenerGuilds) HandlerReady(s *discordgo.Session, e *discordgo.Ready) {
	now := time.Now().Add(10 * time.Second)
	l.lockUntil = &now
}

func (l *ListenerGuilds) Handler(s *discordgo.Session, e *discordgo.GuildCreate) {
	limit := l.cfg.Discord.GuildLimit
	if limit < 1 {
		return
	}

	if l.lockUntil == nil || time.Now().Before(*l.lockUntil) {
		return
	}

	time.Sleep(2 * time.Second)
	g, err := l.st.Guilds()
	if err != nil {
		log.Error("Failed fetching guilds from state", err)
		return
	}

	if len(g) <= limit {
		return
	}

	discordutils.SendEmbedMessageDM(s, e.OwnerID, &discordgo.MessageEmbed{
		Title:       "Guild Limit Reached",
		Description: fmt.Sprintf("The guild limit of %d has been reached. Daemon will leave this guild now.", limit),
	})

	if err = s.GuildLeave(e.Guild.ID); err != nil {
		log.Error("Failed leaving guild", err)
		return
	}
}
