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

	// Permissions

	GetPermissions(guildID string) (map[string]perms.PermsArray, error)
	SetPermissions(guildID, roleID string, perms perms.PermsArray) error

	// Votes

	GetVotes() (map[string]vote.Vote, error)
	AddUpdateVote(vote vote.Vote) error
	DeleteVote(voteID string) error

	// Oauth2

	SetUserRefreshToken(ident, token string, expires time.Time) error
	GetUserByRefreshToken(token string) (ident string, expires time.Time, err error)
	RevokeUserRefreshToken(ident string) error

	// Data management

	FlushGuildData(guildID string) error
}
