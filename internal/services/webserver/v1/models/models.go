package models

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type Status struct {
	Code int `json:"code"`
}

type Member struct {
	*discordgo.Member

	GuildName  string    `json:"guild_name,omitempty"`
	AvatarURL  string    `json:"avatar_url"`
	CreatedAt  time.Time `json:"created_at"`
	Dominance  int       `json:"dominance"`
	Karma      int       `json:"karma"`
	KarmaTotal int       `json:"karma_total"`
	ChatMuted  bool      `json:"chat_muted"`
}

type Guild struct {
	ID                       string                      `json:"id"`
	Name                     string                      `json:"name"`
	Icon                     string                      `json:"icon"`
	Region                   string                      `json:"region"`
	AfkChannelID             string                      `json:"afk_channel_id"`
	OwnerID                  string                      `json:"owner_id"`
	JoinedAt                 time.Time                   `json:"joined_at"`
	Splash                   string                      `json:"splash"`
	MemberCount              int                         `json:"member_count"`
	VerificationLevel        discordgo.VerificationLevel `json:"verification_level"`
	Large                    bool                        `json:"large"`
	Unavailable              bool                        `json:"unavailable"`
	MfaLevel                 discordgo.MfaLevel          `json:"mfa_level"`
	Description              string                      `json:"description"`
	Banner                   string                      `json:"banner"`
	PremiumTier              discordgo.PremiumTier       `json:"premium_tier"`
	PremiumSubscriptionCount int                         `json:"premium_subscription_count"`

	Roles    []*discordgo.Role    `json:"roles"`
	Channels []*discordgo.Channel `json:"channels"`

	SelfMember         *Member   `json:"self_member"`
	IconURL            string    `json:"icon_url"`
	BackupsEnabled     bool      `json:"backups_enabled"`
	LatestBackupEntry  time.Time `json:"latest_backup_entry"`
	InviteBlockEnabled bool      `json:"invite_block_enabled"`
}

type GuildReduced struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Icon              string    `json:"icon"`
	IconURL           string    `json:"icon_url"`
	Region            string    `json:"region"`
	OwnerID           string    `json:"owner_id"`
	JoinedAt          time.Time `json:"joined_at"`
	MemberCount       int       `json:"member_count"`
	OnlineMemberCount int       `json:"online_member_count,omitempty"`
}

type SystemInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	GoVersion  string `json:"go_version"`

	Uptime    int64  `json:"uptime"`
	UptimeStr string `json:"uptime_str"`

	OS          string `json:"os"`
	Arch        string `json:"arch"`
	CPUs        int    `json:"cpus"`
	GoRoutines  int    `json:"go_routines"`
	StackUse    uint64 `json:"stack_use"`
	StackUseStr string `json:"stack_use_str"`
	HeapUse     uint64 `json:"heap_use"`
	HeapUseStr  string `json:"heap_use_str"`

	BotUserID string `json:"bot_user_id"`
	BotInvite string `json:"bot_invite"`

	Guilds int `json:"guilds"`
}

type ListResponse[T any] struct {
	N    int `json:"n"`
	Data []T `json:"data"`
}

func NewListResp[T any](data []T) ListResponse[T] {
	return ListResponse[T]{len(data), data}
}
