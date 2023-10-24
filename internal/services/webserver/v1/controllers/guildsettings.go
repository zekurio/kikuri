package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/daemon/internal/services/config"
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/database/dberr"
	"github.com/zekurio/daemon/internal/services/permissions"
	"github.com/zekurio/daemon/internal/services/webserver/v1/models"
	"github.com/zekurio/daemon/internal/services/webserver/wsutil"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/arrayutils"
)

type GuildSettingsController struct {
	db      database.Database
	st      *dgrs.State
	session *discordgo.Session
	cfg     config.Config
	pmw     *permissions.Permissions
}

func (c *GuildSettingsController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.st = ctn.Get(static.DiState).(*dgrs.State)
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = ctn.Get(static.DiConfig).(config.Config)
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

	if gs.APIEnabled, err = c.db.GetGuildAPIEnabled(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	if !gs.APIEnabled {
		return ctx.JSON(models.GuildSettingsEmpty{APIEnabled: false})
	}

	if gs.AutoRoles, err = c.db.GetGuildAutoRoles(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

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

	// check if api is enabled
	gs.APIEnabled, err = c.db.GetGuildAPIEnabled(guildID)
	if err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	if !gs.APIEnabled {
		return fiber.NewError(fiber.StatusForbidden, "api is not enabled for this guild")
	}

	if gs.AutoRoles != nil {
		if ok, _, err := c.pmw.HasPerms(c.session, guildID, uid, "dm.guild.config.autorole"); err != nil {
			return err
		} else if !ok {
			return fiber.NewError(fiber.StatusForbidden, "missing permissions")
		}

		if arrayutils.ContainsAny[string](gs.AutoRoles, "@everyone") {
			return fiber.NewError(fiber.StatusBadRequest,
				"@everyone can not be set as an autorole")
		}

		guildRoles, err := c.st.Roles(guildID, true)
		if err != nil {
			return err
		}
		guildRoleIDs := make([]string, len(guildRoles))
		for i, role := range guildRoles {
			guildRoleIDs[i] = role.ID
		}

		if nc := arrayutils.Contained[string](gs.AutoRoles, guildRoleIDs); len(nc) != len(gs.AutoRoles) {
			return fiber.NewError(fiber.StatusBadRequest,
				"one or more roles are not part of this guild")
		}

		if err = c.db.SetGuildAutoRoles(guildID, gs.AutoRoles); err != nil {
			return wsutil.IsErrInternalOrNotFound(err)
		}
	}

	if gs.AutoVoice != nil {
		if ok, _, err := c.pmw.HasPerms(c.session, guildID, uid, "dm.guild.config.autovoice"); err != nil {
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
