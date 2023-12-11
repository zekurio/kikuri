package listeners

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekurio/kikuri/internal/services/vote"

	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordutils"
)

type ListenerReady struct {
	db  database.Database
	st  *dgrs.State
	ken ken.IKen
	vh  vote.VotesProvider
}

func NewListenerReady(ctn di.Container) *ListenerReady {
	return &ListenerReady{
		db:  ctn.Get(static.DiDatabase).(database.Database),
		st:  ctn.Get(static.DiState).(*dgrs.State),
		ken: ctn.Get(static.DiCommandHandler).(ken.IKen),
		vh:  ctn.Get(static.DiVotes).(vote.VotesProvider),
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	err := s.UpdateListeningStatus("slash commands [WIP]")
	if err != nil {
		return
	}
	log.Info("Signed in!",
		"Username", fmt.Sprintf("%s#%s", e.User.Username, e.User.Discriminator),
		"ID", e.User.ID)

	self, err := l.st.SelfUser()
	if err != nil {
		return
	}

	// populate votes
	err = l.vh.Populate(l.ken)
	if err != nil {
		log.Error("Failed populating votes", err)
	}

	log.Infof("Invite link: %s", discordutils.GetInviteLink(self.ID))
}
