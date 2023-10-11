package listeners

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/ken"

	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/scheduler"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/internal/util/vote"
	"github.com/zekurio/daemon/pkg/discordutils"
)

type ListenerReady struct {
	db    database.Database
	ken   ken.IKen
	sched scheduler.Provider
}

func NewListenerReady(ctn di.Container) *ListenerReady {
	return &ListenerReady{
		db:    ctn.Get(static.DiDatabase).(database.Database),
		ken:   ctn.Get(static.DiCommandHandler).(ken.IKen),
		sched: ctn.Get(static.DiScheduler).(scheduler.Provider),
	}
}

func (l *ListenerReady) Handler(s *discordgo.Session, e *discordgo.Ready) {
	// vote
	votes, err := l.db.GetVotes()
	// first, add all votes to the running votes list back from the database
	if err != nil {
		log.Error("Failed getting votes from database")
		return
	} else {
		vote.RunningVotes = votes
		_, err = l.sched.Schedule("*/10 * * * * *", func() {
			now := time.Now()
			for _, v := range vote.RunningVotes {
				if (v.Expires != time.Time{}) && v.Expires.Before(now) {
					v.Close(l.ken, vote.StateExpired)
					if err = l.db.DeleteVote(v.ID); err != nil {
						log.Error("Failed deleting vote from database")
					}
				}
			}
		})
		if err != nil {
			log.Error("Failed scheduling vote expiration")
			return
		}
	}

	// second, add components back to the votes
	for _, v := range vote.RunningVotes {
		b := l.ken.Components().Add(v.MessageID, v.ChannelID)

		_, err = v.AddButtons(b)
		if err != nil {
			log.Error("Failed adding buttons to vote")
			return
		}
	}

	err = s.UpdateListeningStatus("slash commands [WIP]")
	if err != nil {
		return
	}
	log.Info("Signed in!",
		"Username", fmt.Sprintf("%s#%s", e.User.Username, e.User.Discriminator),
		"ID", e.User.ID)
	log.Infof("Invite link: %s", discordutils.GetInviteLink(s))

	l.sched.Start()
}
