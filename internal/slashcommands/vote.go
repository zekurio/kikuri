package slashcommands

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"

	"github.com/zekurio/kikuri/internal/middlewares"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/internal/util/vote"
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
	return "ku.chat.vote"
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
	db := ctx.Get(static.DiDatabase).(database.Database)

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

	newVote := vote.Vote{
		ID:          ctx.GetEvent().ID,
		CreatorID:   ctx.User().ID,
		GuildID:     ctx.GetEvent().GuildID,
		ChannelID:   ctx.GetEvent().ChannelID,
		Description: body,
		Choices:     split,
		ImageURL:    imgLink,
		Expires:     expires,
		Buttons:     map[string]vote.OptionButton{},
		CurrentVote: map[string]vote.CurrentVote{},
	}

	emb, err := newVote.AsEmbed(ctx.GetSession())
	if err != nil {
		return err
	}

	fum := ctx.FollowUpEmbed(emb).Send()
	err = fum.Error
	if err != nil {
		return err
	}

	b := fum.AddComponents()

	newVote.MessageID = fum.Message.ID
	err = db.AddUpdateVote(newVote)
	if err != nil {
		return err
	}

	_, err = newVote.AddButtons(b)
	if err != nil {
		return err
	}

	vote.RunningVotes[newVote.ID] = newVote

	return
}

func (c *Vote) list(ctx ken.SubCommandContext) (err error) {

	emb := &discordgo.MessageEmbed{
		Description: "Your open votes on this guild:",
		Color:       static.ColorDefault,
		Fields:      make([]*discordgo.MessageEmbedField, 0),
	}

	for _, v := range vote.RunningVotes {
		if v.GuildID == ctx.GetEvent().GuildID {
			emb.Fields = append(emb.Fields, v.AsField())
		}
	}

	if len(emb.Fields) == 0 {
		emb.Description = "You don't have any open votes on this guild."
	}
	err = ctx.FollowUpEmbed(emb).Send().Error
	return err
}

func (c *Vote) expire(ctx ken.SubCommandContext) (err error) {
	db := ctx.Get(static.DiDatabase).(database.Database)

	expireDuration, err := timeutils.ParseDuration(ctx.Options().GetByName("timeout").StringValue())
	if err != nil {
		return ctx.FollowUpError(
			"Invalid duration format. Please take a look "+
				"[here](https://golang.org/pkg/time/#ParseDuration) how to format duration parameter.", "").
			Send().Error
	}

	id := ctx.Options().Get(0).StringValue()

	if id == "all" {
		return c.expireAllVotes(ctx, db, expireDuration)
	}

	return c.expireSingleVote(ctx, db, id, expireDuration)
}

func (c *Vote) expireAllVotes(ctx ken.SubCommandContext, db database.Database, expireDuration time.Duration) (err error) {
	for _, v := range vote.RunningVotes {
		if v.GuildID == ctx.GetEvent().GuildID {
			err := c.expireVote(ctx, db, &v, expireDuration)
			if err != nil {
				return err
			}
		}
	}
	return ctx.FollowUpError("No vote found.", "").Send().Error
}

func (c *Vote) expireSingleVote(ctx ken.SubCommandContext, db database.Database, id string, expireDuration time.Duration) (err error) {
	var currentVote *vote.Vote
	for _, v := range vote.RunningVotes {
		if v.GuildID == ctx.GetEvent().GuildID && v.ID == id {
			currentVote = &v
			break
		}
	}

	return c.expireVote(ctx, db, currentVote, expireDuration)
}

func (c *Vote) expireVote(ctx ken.SubCommandContext, db database.Database, ivote *vote.Vote, expireDuration time.Duration) (err error) {
	if err := ivote.SetExpire(ctx.GetSession(), expireDuration); err != nil {
		return err
	}
	if err := db.AddUpdateVote(*ivote); err != nil {
		return err
	}
	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Vote will expire <t:%d:R>", ivote.Expires.Unix()),
	}).Send().Error
}

func (c *Vote) close(ctx ken.SubCommandContext) (err error) {

	ken := ctx.GetKen()
	id := ctx.Options().GetByName("id").StringValue()

	if strings.ToLower(id) == "all" {
		var i int
		if err != nil {
			return err
		}

		for _, v := range vote.RunningVotes {
			if v.GuildID == ctx.GetEvent().GuildID {
				if err := v.Close(ken); err != nil {
					return err
				}

				i++
			}

			return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
				Description: fmt.Sprintf("Closed %d votes.", i),
			}).Send().Error
		}

		return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Closed %d votes.", i),
		}).Send().Error
	}

	var currentVote *vote.Vote
	for _, v := range vote.RunningVotes {
		if v.GuildID == ctx.GetEvent().GuildID && v.ID == id {
			currentVote = &v
			break
		}
	}

	if err := currentVote.Close(ken); err != nil {
		return err
	}

	return ctx.FollowUpEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("Closed vote %s.", currentVote.ID),
	}).Send().Error
}
