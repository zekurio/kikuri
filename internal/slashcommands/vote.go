package slashcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/zekurio/kikuri/internal/models"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"

	"github.com/zekurio/kikuri/internal/middlewares"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/services/vote"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/timeutils"
)

type Vote struct{}

var (
	_ ken.SlashCommand            = (*Vote)(nil)
	_ permissions.CommandPerms    = (*Vote)(nil)
	_ middlewares.CommandCooldown = (*Vote)(nil)
)

func (c *Vote) Name() string {
	return "vote"
}

func (c *Vote) Description() string {
	return "Create a vote."
}

func (c *Vote) Version() string {
	return "1.1.0"
}

func (c *Vote) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Vote) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create a new vote.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "body",
					Description: "The vote body content.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "choices",
					Description: "The choices - split by `,`.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "imageurl",
					Description: "An optional image URL.",
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "timeout",
					Description: "Timeout of the vote (i.e. `1h`, `30m`, ...)",
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List currently running votes.",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "expire",
			Description: "Set the expiration of a running vote.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "The ID of the vote or `all` if you want to close all.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "timeout",
					Description: "Timeout of the vote (i.e. `1h`, `30m`, ...)",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "close",
			Description: "Close a running vote.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "The ID of the vote or `all` if you want to close all.",
					Required:    true,
				},
			},
		},
	}
}

func (c *Vote) Perm() string {
	return "ki.chat.vote"
}

func (c *Vote) SubPerms() []permissions.SubCommandPerms {
	return []permissions.SubCommandPerms{
		{
			Perm:        "close",
			Explicit:    true,
			Description: "Allows closing votes of other users.",
		},
	}
}

func (c *Vote) Cooldown() int {
	return 120
}

func (c *Vote) Run(ctx ken.Context) (err error) {
	if err = ctx.Defer(); err != nil {
		return
	}

	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "create", Run: c.create},
		ken.SubCommandHandler{Name: "list", Run: c.list},
		ken.SubCommandHandler{Name: "expire", Run: c.expire},
		ken.SubCommandHandler{Name: "close", Run: c.close},
	)

	return
}

func (c *Vote) create(ctx ken.SubCommandContext) (err error) {
	vh := ctx.Get(static.DiVotes).(vote.VotesProvider)
	body := ctx.Options().GetByName("body").StringValue()
	choices := ctx.Options().GetByName("choices").StringValue()
	split := strings.Split(choices, ",")
	if len(split) < 2 || len(split) > 10 {
		return ctx.FollowUpError(
			"Invalid arguments.", "").
			Send().Error
	}
	for i, e := range split {
		if len(e) < 1 {
			return ctx.FollowUpError(
				"Possibilities can not be empty.", "").
				Send().Error
		}
		split[i] = strings.Trim(e, " \t")
	}

	var imgLink string
	if imgLinkV, ok := ctx.Options().GetByNameOptional("imageurl"); ok {
		imgLink = imgLinkV.StringValue()
	}

	var expires time.Time
	if expiresV, ok := ctx.Options().GetByNameOptional("timeout"); ok {
		expiresDuration, err := timeutils.ParseDuration(expiresV.StringValue())
		if err != nil {
			return ctx.FollowUpError(
				"Invalid duration format. Please take a look "+
					"[here](https://golang.org/pkg/time/#ParseDuration) how to format duration parameter.", "").
				Send().Error
		}
		expires = time.Now().Add(expiresDuration)
	}

	newVote := models.Vote{
		ID:          ctx.GetEvent().ID,
		CreatorID:   ctx.User().ID,
		GuildID:     ctx.GetEvent().GuildID,
		ChannelID:   ctx.GetEvent().ChannelID,
		Description: body,
		Choices:     split,
		ImageURL:    imgLink,
		Expires:     expires,
		Buttons:     map[string]models.OptionButton{},
		CurrentVote: map[string]models.CurrentVote{},
	}

	if err = vh.Create(ctx, newVote); err != nil {
		return ctx.FollowUpError(
			"Failed creating vote.", "").
			Send().Error
	}

	return
}

func (c *Vote) list(ctx ken.SubCommandContext) (err error) {
	vh := ctx.Get(static.DiVotes).(vote.VotesProvider)

	emb := &discordgo.MessageEmbed{
		Description: "Your open votes on this guild:",
		Color:       static.ColorDefault,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	for _, currVote := range vh.GetAllFromGuild(ctx.GetEvent().GuildID) {
		emb.Fields = append(emb.Fields, currVote.AsField())
	}

	// if no votes are open, send an embed a different embed
	if len(emb.Fields) == 0 {
		emb.Description = "You have no open votes on this guild."
	}

	err = ctx.FollowUpEmbed(emb).Send().Error
	return err
}

func (c *Vote) expire(ctx ken.SubCommandContext) (err error) {
	expireDuration, err := timeutils.ParseDuration(ctx.Options().GetByName("timeout").StringValue())
	if err != nil {
		return ctx.FollowUpError(
			"Invalid duration format. Please take a look "+
				"[here](https://golang.org/pkg/time/#ParseDuration) how to format duration parameter.", "").
			Send().Error
	}

	id := ctx.Options().Get(0).StringValue()

	if id == "all" {
		return c.expireAllVotes(ctx, expireDuration)
	}

	return c.expireSingleVote(ctx, id, expireDuration)
}

func (c *Vote) expireAllVotes(ctx ken.SubCommandContext, expireDuration time.Duration) (err error) {
	// get all votes from database for the current guild
	vh := ctx.Get(static.DiVotes).(vote.VotesProvider)
	votes := vh.GetAllFromGuild(ctx.GetEvent().GuildID)

	// iterate over all votes
	for _, currVote := range votes {
		if err := c.expireVote(ctx, currVote, expireDuration); err != nil {
			return err
		}
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("All votes will expire <t:%d:R>", time.Now().Add(expireDuration).Unix()),
	}).Send().Error
}

func (c *Vote) expireSingleVote(ctx ken.SubCommandContext, id string, expireDuration time.Duration) (err error) {
	vh := ctx.Get(static.DiVotes).(vote.VotesProvider)

	// get currVote from map
	currVote, err := vh.Get(id)
	if err != nil {
		return ctx.FollowUpError("Vote not found.", "").Send().Error
	}

	err = c.expireVote(ctx, currVote, expireDuration)

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Vote %s will expire <t:%d:R>", id, time.Now().Add(expireDuration).Unix()),
	}).Send().Error
}

func (c *Vote) expireVote(ctx ken.SubCommandContext, ivote models.Vote, expireDuration time.Duration) (err error) {
	vh := ctx.Get(static.DiVotes).(vote.VotesProvider)

	if err := ivote.SetExpire(ctx.GetSession(), expireDuration); err != nil {
		return err
	}

	if err = vh.Update(ivote); err != nil {
		return err
	}

	return err
}

func (c *Vote) close(ctx ken.SubCommandContext) (err error) {
	vh := ctx.Get(static.DiVotes).(vote.VotesProvider)
	id := ctx.Options().GetByName("id").StringValue()

	if strings.ToLower(id) == "all" {
		votesClosed, err := vh.CloseAll(ctx.GetKen(), ctx.GetEvent().GuildID)
		if err != nil {
			return ctx.FollowUpError("Failed closing all votes.", "").Send().Error
		}

		return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Closed %d votes.", votesClosed),
		}).Send().Error
	}

	if err := vh.Close(ctx.GetKen(), id, models.StateClosed); err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Closed vote %s.", id),
	}).Send().Error
}
