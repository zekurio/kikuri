package slashcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"

	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/colorutils"
	"github.com/zekurio/kikuri/pkg/discordutils"
	"github.com/zekurio/kikuri/pkg/embedbuilder"
)

type Guild struct {
	ken.EphemeralCommand
}

var (
	_ ken.SlashCommand         = (*Guild)(nil)
	_ permissions.CommandPerms = (*Guild)(nil)
)

func (c *Guild) Name() string {
	return "guild"
}

func (c *Guild) Description() string {
	return "Displays information about the current guild."
}

func (c *Guild) Version() string {
	return "1.0.0"
}

func (c *Guild) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Guild) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *Guild) Perm() string {
	return "ki.chat.guild"
}

func (c *Guild) SubPerms() []permissions.SubCommandPerms {
	return nil
}

func (c *Guild) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	s := ctx.GetSession()
	st := ctx.Get(static.DiState).(*dgrs.State)

	const maxGuildRoles = 16

	guild, err := st.Guild(ctx.GetEvent().GuildID)
	if err != nil {
		return
	}

	color, err := colorutils.GenerateColorFromImageURL(guild.IconURL("1024"))
	if err != nil {
		color = static.ColorDefault
	}

	guildChannels, err := s.GuildChannels(guild.ID)
	if err != nil {
		return
	}

	textChannels, voiceChannels, categoryChannels, threadChannels := 0, 0, 0, 0
	for _, c := range guildChannels {
		switch c.Type {
		case discordgo.ChannelTypeGuildCategory:
			categoryChannels++
		case discordgo.ChannelTypeGuildVoice:
			voiceChannels++
		case discordgo.ChannelTypeGuildPrivateThread:
		case discordgo.ChannelTypeGuildPublicThread:
			threadChannels++
		default:
			textChannels++
		}
	}

	channelInfo := fmt.Sprintf("Category Channels: `%d`\nText Channels: `%d`\nThreads: `%d`\nVoice Channels: `%d`",
		categoryChannels, textChannels, threadChannels, voiceChannels)

	lenRoles := len(guild.Roles) - 1
	if lenRoles > maxGuildRoles {
		lenRoles = maxGuildRoles + 1
	}

	roles := make([]string, lenRoles)
	i := 0
	for _, r := range guild.Roles {
		if r.ID == guild.ID {
			continue
		}
		if i == maxGuildRoles {
			roles[i] = "..."
			break
		}
		roles[i] = r.Mention()
		i++
	}

	createdTime, err := discordutils.GetDiscordSnowflakeCreationTime(guild.ID)
	if err != nil {
		return
	}

	emb := embedbuilder.New().
		SetTitle("About "+guild.Name).
		SetThumbnail(guild.IconURL("1024"), "", 100, 100).
		SetColor(color).
		AddField("Name", guild.Name).
		AddField("ID", fmt.Sprintf("```\n%s\n```", guild.ID)).
		AddField("Created", createdTime.Format(time.RFC1123)).
		AddField("Owner", fmt.Sprintf("<@%s>", guild.OwnerID)).
		AddField(fmt.Sprintf("Channels (%d)", len(guildChannels)), channelInfo).
		AddField("Member Count", fmt.Sprintf("Cached: %d", guild.MemberCount)).
		AddField(fmt.Sprintf("Roles (%d)", len(guild.Roles)-1), strings.Join(roles, ", ")).
		Build()

	return ctx.FollowUpEmbed(emb).Send().Error
}
