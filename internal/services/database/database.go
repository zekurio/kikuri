package database

import (
	"time"

	"github.com/zekurio/daemon/internal/services/database/models"
	"github.com/zekurio/daemon/internal/util/vote"
	"github.com/zekurio/daemon/pkg/perms"
)

// Database is the interface for our database service
// which is then implemented by postgres
type Database interface {
	Close() error

	// Guild settings

	GetGuildAutoRoles(guildID string) ([]string, error)
	SetGuildAutoRoles(guildID string, roleIDs []string) error

	GetGuildAutoVoice(guildID string) ([]string, error)
	SetGuildAutoVoice(guildID string, channelIDs []string) error

	// Permissions

	GetPermissions(guildID string) (map[string]perms.Array, error)
	SetPermissions(guildID, roleID string, perms perms.Array) error

	// Votes

	GetVotes() (map[string]vote.Vote, error)
	AddUpdateVote(vote vote.Vote) error
	DeleteVote(voteID string) error

	// Guildapi

	GetGuildAPI(guildID string) (settings models.GuildAPISettings, err error)
	SetGuildAPI(guildID string, settings models.GuildAPISettings) error

	// User refresh tokens

	GetUserRefreshToken(userID string) (token string, err error)
	SetUserRefreshToken(userID, token string, expires time.Time) error
	RevokeUserRefreshToken(userID string) error

	GetUserByRefreshToken(token string) (userID string, expires time.Time, err error)

	// API tokens

	SetAPIToken(token models.APITokenEntry) error
	GetAPIToken(userID string) (models.APITokenEntry, error)
	DeleteAPIToken(userID string) error

	// Data management

	FlushGuildData(guildID string) error
}
