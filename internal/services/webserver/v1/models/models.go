package models

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zekrotja/ken"
	"github.com/zekurio/daemon/internal/services/database"
	db_models "github.com/zekurio/daemon/internal/services/database/models"
	"github.com/zekurio/daemon/internal/services/permissions"
	"github.com/zekurio/daemon/pkg/discordutils"
	"github.com/zekurio/daemon/pkg/perms"
)

var Ok = &Status{200}

type Status struct {
	Code int `json:"code"`
}

type State struct {
	State bool `json:"state"`
}

type AccessTokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

type Error struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Context string `json:"context,omitempty"`
}

// ListResponse wraps a list response object
// with the list as Data and N as len(Data).
type ListResponse[T any] struct {
	N    int `json:"n"`
	Data []T `json:"data"`
}

func NewListResponse[T any](data []T) ListResponse[T] {
	return ListResponse[T]{len(data), data}
}

// User extends a discordgo.User as reponse
// model.
type User struct {
	*discordgo.User

	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	BotOwner  bool      `json:"bot_owner"`
}

// FlatUser shrinks the user object to the only
// necessary parts for the web interface.
type FlatUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	AvatarURL     string `json:"avatar_url"`
	Bot           bool   `json:"bot"`
}

// Member extends a discordgo.Member as
// response model.
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

// Guild extends a discordgo.Guild as
// response model.
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

// GuildReduced is a Guild model with fewer
// details than Guild model.
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

// PermissionsResponse wraps a
// permissions.PermissionsArra as response
// model.
type PermissionsResponse struct {
	Permissions perms.PermsArray `json:"permissions"`
}

// GuildSettings is the response model for
// guild settings and preferences.
type GuildSettings struct {
	Perms     map[string]perms.PermsArray `json:"perms"`
	AutoRoles []string                    `json:"autoroles"`
	AutoVoice []string                    `json:"autovoice"`
}

// PermissionsUpdate is the request model to
// update a permissions array.
type PermissionsUpdate struct {
	Perm     string   `json:"perm"`
	RoleIDs  []string `json:"role_ids"`
	Override bool     `json:"override"`
}

// ReasonRequest is a request model wrapping a
// Reason and Attachment URL.
type ReasonRequest struct {
	Reason         string     `json:"reason"`
	Timeout        *time.Time `json:"timeout"`
	Attachment     string     `json:"attachment"`
	AttachmentData string     `json:"attachment_data"`
}

// InviteSettingsRequest is the request model
// for setting the global invite setting.
type InviteSettingsRequest struct {
	GuildID    string `json:"guild_id"`
	Messsage   string `json:"message"`
	InviteCode string `json:"invite_code"`
}

// InviteSettingsResponse is the response model
// sent back when setting the global invite setting.
type InviteSettingsResponse struct {
	Guild     *Guild `json:"guild"`
	InviteURL string `json:"invite_url"`
	Message   string `json:"message"`
}

// Count is a simple response wrapper for a
// count number.
type Count struct {
	Count int `json:"count"`
}

type LandingPageResponse struct {
	LocalInvite        string `json:"localinvite"`
	PublicMainInvite   string `json:"publicmaininvite"`
	PublicCanaryInvite string `json:"publiccaranyinvite"`
}

// SystemInfo is the response model for a
// system info request.
type SystemInfo struct {
	Version    string    `json:"version"`
	CommitHash string    `json:"commit_hash"`
	BuildDate  time.Time `json:"build_date"`
	GoVersion  string    `json:"go_version"`

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

// APITokenResponse wraps the reponse model of
// an apit token request.
type APITokenResponse struct {
	Created    time.Time `json:"created"`
	Expires    time.Time `json:"expires"`
	LastAccess time.Time `json:"last_access"`
	Hits       int       `json:"hits"`
	Token      string    `json:"token,omitempty"`
}

// APITokenClaims extends the standard JWT claims
// by private claims used for api tokens.
type APITokenClaims struct {
	jwt.StandardClaims

	Salt string `json:"sp_salt,omitempty"`
}

// SessionTokenClaims extends the standard JWT
// claims by information used for session tokens.
//
// Currently, no additional information is
// extended but this wrapper is used tho to
// be able to add session information later.
type SessionTokenClaims struct {
	jwt.StandardClaims
}

// SlashCommandInfo wraps a slash command object
// containing all information of a slash command
// instance.
type SlashCommandInfo struct {
	Name            string                                `json:"name"`
	Description     string                                `json:"description"`
	Version         string                                `json:"version"`
	Options         []*discordgo.ApplicationCommandOption `json:"options"`
	Perms           string                                `json:"perms"`
	SubCommandPerms []permissions.SubCommandPerms         `json:"subperms"`
	DmCapable       bool                                  `json:"dm_capable"`
	Group           string                                `json:"group"`
}

type UsersettingsPrivacy struct {
	StarboardOptout bool `json:"starboard_optout"`
}

type PermissionsMap map[string]perms.PermsArray

type EnableStatus struct {
	Enabled bool `json:"enabled"`
}

type FlushGuildRequest struct {
	Validation string `json:"validation"`
	LeaveAfter bool   `json:"leave_after"`
}

type SearchResult struct {
	Guilds  []*GuildReduced `json:"guilds"`
	Members []*Member       `json:"members"`
}

type GuildAPISettingsRequest struct {
	db_models.GuildAPISettings
	NewToken   string `json:"token"`
	ResetToken bool   `json:"reset_token"`
}

// GuildFromGuild returns a Guild model from the passed
// discordgo.Guild g, discordgo.Member m and cmdHandler.
func GuildFromGuild(g *discordgo.Guild, m *discordgo.Member, db database.Database, botOwnerID string) (ng *Guild, err error) {
	if g == nil {
		return
	}

	selfmm := MemberFromMember(m)

	if m != nil {
		switch {
		case discordutils.IsAdmin(g, m):
			selfmm.Dominance = 1
		case g.OwnerID == m.User.ID:
			selfmm.Dominance = 2
		case botOwnerID == m.User.ID:
			selfmm.Dominance = 3
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

// GuildReducedFromGuild returns a GuildReduced from the passed
// discordgo.Guild g.
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

// MemberFromMember returns a Member from the passed
// discordgo.Member m.
func MemberFromMember(m *discordgo.Member) *Member {
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

func GetSlashCommandInfoFromCommand(cmd *ken.CommandInfo) (ci *SlashCommandInfo) {
	ci = new(SlashCommandInfo)

	ci.Name = cmd.ApplicationCommand.Name
	ci.Description = cmd.ApplicationCommand.Description
	ci.Options = cmd.ApplicationCommand.Options
	ci.Version = cmd.ApplicationCommand.Version
	ci.Perms = cmd.Implementations["Domain"][0].(string)
	ci.SubCommandPerms = cmd.Implementations["SubDomains"][0].([]permissions.SubCommandPerms)

	if v, ok := cmd.Implementations["IsDmCapable"]; ok && len(v) != 0 {
		ci.DmCapable = v[0].(bool)
	}

	domainSplit := strings.Split(ci.Perms, ".")
	ci.Group = strings.Join(domainSplit[1:len(domainSplit)-1], " ")
	ci.Group = strings.ToUpper(ci.Group)

	return
}

// FlatUserFromUser returns the reduced FlatUser object
// from the given user object.
func FlatUserFromUser(u *discordgo.User) (fu *FlatUser) {
	return &FlatUser{
		ID:            u.ID,
		Username:      u.Username,
		Discriminator: u.Discriminator,
		AvatarURL:     u.AvatarURL(""),
		Bot:           u.Bot,
	}
}
