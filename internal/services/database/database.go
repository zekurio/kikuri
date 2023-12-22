package database

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/zekurio/kikuri/internal/models"

	"github.com/zekurio/kikuri/pkg/perms"
)

// Database is the interface for our database service
// which is then implemented by postgres
type Database interface {
	Close() error

	// Guild settings

	GetGuildAutoVoice(guildID string) (autovoices []string, err error)
	SetGuildAutoVoice(guildID string, channelIDs []string) error

	// Permissions

	GetPermissions(guildID string) (permissions map[string]perms.Array, err error)
	SetPermissions(guildID, roleID string, permissions perms.Array) error

	// Votes

	GetVotes() (votes map[string]models.Vote, err error)
	AddUpdateVote(vote models.Vote) error
	DeleteVote(voteID string) error

	// Oauth2

	SetUserRefreshToken(ident, token string, expires time.Time) error
	GetUserByRefreshToken(token string) (ident string, expires time.Time, err error)
	RevokeUserRefreshToken(ident string) error

	// API tokens

	SetAPIToken(token models.APITokenEntry) error
	GetAPIToken(userID string) (models.APITokenEntry, error)
	DeleteAPIToken(userID string) error

	// User settings

	GetRedditKarma(userID, guildID string) (int, error)
	GetRedditKarmaSum(userID string) (int, error)
	GetRedditGuildEntries(guildID string, limit int) ([]models.GuildReddit, error)
	SetRedditKarma(userID, guildID string, val int) error
	UpdateRedditKarma(userID, guildID string, diff int) error

	SetRedditState(guildID string, state bool) error
	GetRedditState(guildID string) (bool, error)

	SetRedditEmotes(guildID, emotesInc, emotesDec string) error
	GetRedditEmotes(guildID string) (emotesInc, emotesDec string, err error)

	SetRedditTokens(guildID string, tokens int) error
	GetRedditTokens(guildID string) (int, error)

	SetRedditPenalty(guildID string, state bool) error
	GetRedditPenalty(guildID string) (bool, error)

	GetRedditBlockList(guildID string) ([]string, error)
	IsRedditBlockListed(guildID, userID string) (bool, error)
	AddRedditBlockList(guildID, userID string) error
	RemoveRedditBlockList(guildID, userID string) error

	GetRedditRules(guildID string) ([]models.RedditRule, error)
	CheckRedditRule(guildID, checksum string) (ok bool, err error)
	AddOrUpdateRedditRule(rule models.RedditRule) error
	RemoveRedditRule(guildID string, id snowflake.ID) error

	// Data management

	FlushGuildData(guildID string) error
}
