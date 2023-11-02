package listeners

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/util/static"
)

type ListenerVoiceStateUpdate struct {
	db              database.Database
	st              *dgrs.State
	pmw             *permissions.Permissions
	voiceStateCache map[string]*discordgo.VoiceState
	autovcCache     map[string]string
}

func NewListenerVoiceStateUpdate(ctn di.Container) *ListenerVoiceStateUpdate {
	return &ListenerVoiceStateUpdate{
		db:              ctn.Get(static.DiDatabase).(database.Database),
		st:              ctn.Get(static.DiState).(*dgrs.State),
		pmw:             ctn.Get(static.DiPermissions).(*permissions.Permissions),
		voiceStateCache: map[string]*discordgo.VoiceState{},
		autovcCache:     map[string]string{},
	}
}

func (l *ListenerVoiceStateUpdate) Handler(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {

	allowed, _, err := l.pmw.HasPerms(s, e.GuildID, e.UserID, "ku.chat.autochannel")
	if err != nil || !allowed {
		return
	}

	vsOld := l.voiceStateCache[e.UserID]
	vsNew := e.VoiceState

	l.voiceStateCache[e.UserID] = vsNew

	avIDs, err := l.db.GetGuildAutoVoice(e.GuildID)
	if err != nil {
		return
	}
	avString := strings.Join(avIDs, ";")

	if vsOld == nil || (vsOld != nil && vsOld.ChannelID == "") {

		if !strings.Contains(avString, vsNew.ChannelID) {
			return
		}

		if err := l.createAutoVC(s, e.UserID, e.GuildID, vsNew.ChannelID); err != nil {
			return
		}

	} else if vsOld != nil && vsNew.ChannelID != "" && vsOld.ChannelID != vsNew.ChannelID {

		if vsNew.ChannelID == l.autovcCache[e.UserID] {

		} else if strings.Contains(avString, vsNew.ChannelID) && l.autovcCache[e.UserID] == "" {
			if l.autovcCache[e.UserID] == "" {
				if err := l.createAutoVC(s, e.UserID, e.GuildID, vsNew.ChannelID); err != nil {
					return
				}
			} else {
				if err := l.deleteAutoVC(s, e.UserID); err != nil {
					return
				}
			}
		} else if l.autovcCache[e.UserID] != "" {
			if err := l.deleteAutoVC(s, e.UserID); err != nil {
				return
			}
		}

	} else if vsOld != nil && vsNew.ChannelID == "" {
		if l.autovcCache[e.UserID] != "" {
			if err := l.deleteAutoVC(s, e.UserID); err != nil {
				return
			}
		}

	}
}

func (l *ListenerVoiceStateUpdate) createAutoVC(s *discordgo.Session, userID, guildID, parentChannelID string) error {
	parentCh, err := l.st.Channel(parentChannelID)
	if err != nil {
		return err
	}

	member, err := l.st.Member(guildID, userID)
	if err != nil {
		return err
	}

	var chName string
	if member.Nick != "" {
		chName = member.Nick + "'s " + parentCh.Name
	} else {
		chName = member.User.Username + "'s " + parentCh.Name
	}
	ch, err := s.GuildChannelCreate(guildID, chName, discordgo.ChannelTypeGuildVoice)
	if err != nil {
		return err
	}
	ch, err = s.ChannelEditComplex(ch.ID, &discordgo.ChannelEdit{
		ParentID: parentCh.ParentID,
		Position: parentCh.Position,
	})
	if err != nil {
		return err
	}
	l.autovcCache[userID] = ch.ID
	if err := s.GuildMemberMove(guildID, userID, &ch.ID); err != nil {
		return err
	}
	return nil
}

func (l *ListenerVoiceStateUpdate) deleteAutoVC(s *discordgo.Session, userID string) error {
	chID := l.autovcCache[userID]
	_, err := s.ChannelDelete(chID)
	if err != nil {
		return err
	}
	delete(l.autovcCache, userID)
	return nil
}
