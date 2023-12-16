package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekurio/kikuri/pkg/discordutils"
	"github.com/zekurio/kikuri/pkg/perms"
)

var Ok = &Status{Code: 200}

type Status struct {
	Code    int
	Message string
}

type Error struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Context string `json:"context,omitempty"`
}

type AccessTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type ListResponse[T any] struct {
	N    int `json:"n"`
	Data []T `json:"data"`
}

func NewListResponse[T any](data []T) ListResponse[T] {
	return ListResponse[T]{len(data), data}
}

// User extends a discordgo.User to a
// response model.
type User struct {
	*discordgo.User

	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	BotOwner  bool      `json:"bot_owner"`
}

// FlatUser flattens a discordgo.User to a
// minimal response model.
type FlatUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	AvatarURL     string `json:"avatar_url"`
	Bot           bool   `json:"bot"`
}

func FlatUserFromDiscordUser(u *discordgo.User) (fu *FlatUser) {
	return &FlatUser{
		ID:            u.ID,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		AvatarURL:     u.AvatarURL(""),
		Bot:           u.Bot,
	}
}

// Member extends a discordgo.Member to
// a response model.
type Member struct {
	*discordgo.Member

	GuildName string    `json:"guild_name,omitempty"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	Privilege int       `json:"privilege"`
}

func MemberFromDiscordMember(m *discordgo.Member) *Member {
	if m == nil {
		return nil
	}

	created, _ := discordutils.GetDiscordSnowflakeCreationTime(m.User.ID)
	return &Member{
		Member:    m,
		AvatarURL: m.User.AvatarURL(""),
		CreatedAt: created,
	}
}

// Guild extends a discordgo.Guild to
// a response model.
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

	SelfMember *Member `json:"self_member"`
	IconURL    string  `json:"icon_url"`
}

func GuildFromDiscordGuild(g *discordgo.Guild, m *discordgo.Member, botOwnerID string) (ng *Guild, err error) {
	if g == nil {
		return
	}

	selfmm := MemberFromDiscordMember(m)

	if m != nil {
		switch {
		case discordutils.IsAdmin(g, m):
			selfmm.Privilege = 1
		case g.OwnerID == m.User.ID:
			selfmm.Privilege = 2
		case botOwnerID == m.User.ID:
			selfmm.Privilege = 3
		}
	}

	ng = &Guild{
		AfkChannelID:             g.AfkChannelID,
		Banner:                   g.Banner,
		Channels:                 g.Channels,
		Description:              g.Description,
		ID:                       g.ID,
		Icon:                     g.Icon,
		JoinedAt:                 g.JoinedAt,
		Large:                    g.Large,
		MemberCount:              g.MemberCount,
		MfaLevel:                 g.MfaLevel,
		Name:                     g.Name,
		OwnerID:                  g.OwnerID,
		PremiumSubscriptionCount: g.PremiumSubscriptionCount,
		PremiumTier:              g.PremiumTier,
		Region:                   g.Region,
		Roles:                    g.Roles,
		Splash:                   g.Splash,
		Unavailable:              g.Unavailable,
		VerificationLevel:        g.VerificationLevel,

		SelfMember: selfmm,
		IconURL:    g.IconURL(""),
	}

	return
}

// GuildReduced is a reduced response of
// models.Guild.
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

func GuildReducedFromGuild(g *discordgo.Guild) *GuildReduced {
	return &GuildReduced{
		ID:          g.ID,
		Name:        g.Name,
		Icon:        g.Icon,
		IconURL:     g.IconURL(""),
		Region:      g.Region,
		OwnerID:     g.OwnerID,
		JoinedAt:    g.JoinedAt,
		MemberCount: g.MemberCount,
	}
}

// PermissionsResponse wraps a
// permissions.PermissionsArra as response
// model.
type PermissionsResponse struct {
	Permissions perms.Array `json:"permissions"`
}

// GuildSettings is a response for the settings and
// preferences of a guild.
type GuildSettings struct {
	AutoRoles []string               `json:"auto_roles"`
	AutoVoice []string               `json:"auto_voice"`
	Perms     map[string]perms.Array `json:"perms"`
}

type SearchResult struct {
	Guilds  []*GuildReduced `json:"guilds"`
	Members []*Member       `json:"members"`
}

type APITokenResponse struct {
	Created    time.Time `json:"created"`
	Expires    time.Time `json:"expires"`
	LastAccess time.Time `json:"last_access"`
	Hits       int       `json:"hits"`
	Token      string    `json:"token,omitempty"`
}
