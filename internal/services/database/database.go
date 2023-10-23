package database

import (
	"time"

	"github.com/zekurio/daemon/internal/util/vote"
	"github.com/zekurio/daemon/pkg/perms"
)

// Database is the interface for our database service
// which is then implemented by postgres
type Database interface {
	Close() error

	// Guild settings

	GetGuildAutoRoles(guildID string) (autoroles []string, err error)
	SetGuildAutoRoles(guildID string, roleIDs []string) error

	GetGuildAutoVoice(guildID string) (autovoices []string, err error)
	SetGuildAutoVoice(guildID string, channelIDs []string) error

	GetGuildAPIEnabled(guildID string) (enabled bool, err error)
	SetGuildAPIEnabled(guildID string, enabled bool) error

	// Permissions

	GetPermissions(guildID string) (permissions map[string]perms.PermsArray, err error)
	SetPermissions(guildID, roleID string, permissions perms.PermsArray) error

	// Votes

	GetVotes() (votes map[string]vote.Vote, err error)
	AddUpdateVote(vote vote.Vote) error
	DeleteVote(voteID string) error

	// Oauth2

	SetUserRefreshToken(ident, token string, expires time.Time) error
	GetUserByRefreshToken(token string) (ident string, expires time.Time, err error)
	RevokeUserRefreshToken(ident string) error

	// Data management

	FlushGuildData(guildID string) error
}
