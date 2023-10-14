package database

import (
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

	GetPermissions(guildID string) (map[string]perms.PermsArray, error)
	SetPermissions(guildID, roleID string, perms perms.PermsArray) error

	// Votes

	GetVotes() (map[string]vote.Vote, error)
	AddUpdateVote(vote vote.Vote) error
	DeleteVote(voteID string) error

	// Data management

	FlushGuildData(guildID string) error
}
