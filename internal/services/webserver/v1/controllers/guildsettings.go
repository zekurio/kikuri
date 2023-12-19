package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	sharedmodels "github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/database/dberr"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/services/webserver/v1/models"
	"github.com/zekurio/kikuri/internal/services/webserver/wsutil"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/arrayutils"
)

type GuildSettingsController struct {
	db      database.Database
	st      *dgrs.State
	session *discordgo.Session
	cfg     sharedmodels.Config
	pmw     *permissions.Permissions
}

func (c *GuildSettingsController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.st = ctn.Get(static.DiState).(*dgrs.State)
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = ctn.Get(static.DiConfig).(sharedmodels.Config)
	c.pmw = ctn.Get(static.DiPermissions).(*permissions.Permissions)

	router.Get("", c.getGuildSettings)
	router.Post("", c.postGuildSettings)
}

// @Summary Get Guild Settings
// @Description Returns the specified general guild settings.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.GuildSettings
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings [get]
func (c *GuildSettingsController) getGuildSettings(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	gs := new(models.GuildSettings)
	var err error

	if gs.AutoVoice, err = c.db.GetGuildAutoVoice(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	if gs.Perms, err = c.db.GetPermissions(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	return ctx.JSON(gs)
}

// @Summary Get Guild Settings
// @Description Returns the specified general guild settings.
// @Tags Guild Settings
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.GuildSettings true "Modified guild settings payload."
// @Success 200 {object} models.Status
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 403 {Object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings [post]
func (c *GuildSettingsController) postGuildSettings(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	var err error

	gs := new(models.GuildSettings)
	if err = ctx.BodyParser(gs); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if gs.AutoVoice != nil {
		if ok, _, err := c.pmw.HasPerms(guildID, uid, "ki.guild.config.autovoice"); err != nil {
			return err
		} else if !ok {
			return fiber.NewError(fiber.StatusForbidden, "missing permissions")
		}

		guildChannels, err := c.st.Channels(guildID)
		if err != nil {
			return err
		}

		guildChannelIDs := make([]string, len(guildChannels))
		for i, channel := range guildChannels {
			if channel.Type == discordgo.ChannelTypeGuildVoice {
				guildChannelIDs[i] = channel.ID
			}
		}

		if nc := arrayutils.Contained[string](gs.AutoVoice, guildChannelIDs); len(nc) != len(gs.AutoVoice) {
			return fiber.NewError(fiber.StatusBadRequest,
				"one or more channels are not part of this guild or are not voice channels")
		}

		if err = c.db.SetGuildAutoVoice(guildID, gs.AutoVoice); err != nil {
			return wsutil.IsErrInternalOrNotFound(err)
		}
	}

	return ctx.JSON(models.Ok)
}
