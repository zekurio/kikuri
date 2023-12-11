package listeners

import (
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/sarulabs/di/v2"

	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/database/dberr"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/arrayutils"
)

type ListenerMemberAdd struct {
	db database.Database
}

func NewListenerMemberAdd(ctn di.Container) *ListenerMemberAdd {
	return &ListenerMemberAdd{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (g *ListenerMemberAdd) Handler(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	autoroleIDs, err := g.db.GetGuildAutoRoles(e.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		log.With("err", err).Error("Failed getting auto role settings")
		return
	}

	invalidAutoRoleIDs := make([]string, 0)
	for _, rid := range autoroleIDs {
		err = s.GuildMemberRoleAdd(e.GuildID, e.User.ID, rid)
		if apiErr, ok := err.(*discordgo.RESTError); ok && apiErr.Message.Code == discordgo.ErrCodeUnknownRole {
			invalidAutoRoleIDs = append(invalidAutoRoleIDs, rid)
		} else if err != nil {
			log.With("err", err).Error("Failed setting autorole for member")
		}
	}

	if len(invalidAutoRoleIDs) > 0 {
		newAutoRoleIDs := make([]string, 0, len(autoroleIDs)-len(invalidAutoRoleIDs))
		for _, rid := range autoroleIDs {
			if !arrayutils.Contains(invalidAutoRoleIDs, rid) {
				newAutoRoleIDs = append(newAutoRoleIDs, rid)
			}
		}
		err = g.db.SetGuildAutoRoles(e.GuildID, newAutoRoleIDs)
		if err != nil {
			log.With("err", err).Error("Failed updating auto role settings")
		}
	}
}
