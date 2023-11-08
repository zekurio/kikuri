package slashcommands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"

	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordutils"
	"github.com/zekurio/kikuri/pkg/embedbuilder"
	"github.com/zekurio/kikuri/pkg/stringutils"
)

type Profile struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand         = (*Profile)(nil)
	_ permissions.CommandPerms = (*Profile)(nil)
)

func (c *Profile) Name() string {
	return "profile"
}

func (c *Profile) Description() string {
	return "Shows the profile of a user."
}

func (c *Profile) Version() string {
	return "1.1.0"
}

func (c *Profile) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Profile) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to be displayed.",
		},
	}
}

func (c *Profile) Perm() string {
	return "ku.chat.profile"
}

func (c *Profile) SubPerms() []permissions.SubCommandPerms {
	return nil
}

func (c *Profile) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	st := ctx.Get(static.DiState).(*dgrs.State)
	p := ctx.Get(static.DiPermissions).(*permissions.Permissions)

	var user *discordgo.User

	if resolved := ctx.GetEvent().ApplicationCommandData().Resolved; resolved != nil {
		for _, user = range ctx.GetEvent().ApplicationCommandData().Resolved.Users {
			break
		}
	}

	if user == nil {
		if userV, ok := ctx.Options().GetByNameOptional("user"); ok {
			user = userV.UserValue(ctx)
		}
	}

	member, err := st.Member(ctx.GetEvent().GuildID, user.ID)
	if err != nil {
		return
	}

	guild, err := st.Guild(ctx.GetEvent().GuildID)
	if err != nil {
		return
	}

	membRoleIDs := make(map[string]struct{})
	for _, rID := range member.Roles {
		membRoleIDs[rID] = struct{}{}
	}

	maxPos := len(guild.Roles)
	roleColor := static.ColorGrey
	for _, guildRole := range guild.Roles {
		if _, ok := membRoleIDs[guildRole.ID]; ok && guildRole.Position < maxPos && guildRole.Color != 0 {
			maxPos = guildRole.Position
			roleColor = guildRole.Color
		}
	}

	createdTime, err := discordutils.GetDiscordSnowflakeCreationTime(member.User.ID)
	if err != nil {
		return
	}

	perms, _, err := p.GetPerms(ctx.GetEvent().GuildID, member.User.ID)
	if err != nil {
		return
	}
	for i, perm := range perms {
		perms[i] = "`" + perm + "`"
	}

	roles := make([]string, len(member.Roles))
	for i, rID := range member.Roles {
		roles[i] = "<@&" + rID + ">"
	}

	emb := embedbuilder.New().
		SetTitle("Profile of "+member.User.Username).
		SetThumbnail(member.User.AvatarURL("256"), "", 100, 100).
		SetColor(roleColor).
		AddField("Nickname", member.Nick).
		AddField("ID", fmt.Sprintf("`%s`", member.User.ID)).
		AddField("Joined at", stringutils.EnsureNotEmpty(fmt.Sprintf("<t:%d:f> - <t:%d:R>", member.JoinedAt.Unix(), member.JoinedAt.Unix()),
			"*failed parsing timestamp*")).
		AddField("Created at", stringutils.EnsureNotEmpty(fmt.Sprintf("<t:%d:f> - <t:%d:R>", createdTime.Unix(), createdTime.Unix()),
			"*failed parsing timestamp*")).
		AddField("Bot Permissions", stringutils.EnsureNotEmpty(strings.Join(perms, "\n"), "*no perms set*")).
		AddField("Roles", stringutils.EnsureNotEmpty(strings.Join(roles, ", "), "*no roles set*"))

	if member.User.Bot {
		emb.SetDescription(":robot:  **This is a bot account**")
	}

	return ctx.FollowUpEmbed(emb.Build()).Send().Error
}
