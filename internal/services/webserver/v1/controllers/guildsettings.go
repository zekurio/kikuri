package controllers

import (
	"fmt"
	"strings"

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

type GuildsSettingsController struct {
	db      database.Database
	session *discordgo.Session
	cfg     config.Config
	pmw     *permissions.Permissions
	st      *dgrs.State
}

func (c *GuildsSettingsController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = ctn.Get(static.DiConfig).(config.Config)
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.pmw = ctn.Get(static.DiPermissions).(*permissions.Permissions)
	c.st = ctn.Get(static.DiState).(*dgrs.State)

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
func (c *GuildsSettingsController) getGuildSettings(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	gs := new(models.GuildSettings)
	var err error

	if gs.Perms, err = c.db.GetPermissions(guildID); err != nil && err != dberr.ErrNotFound {
		return err
	}

	if gs.AutoRoles, err = c.db.GetGuildAutoRoles(guildID); err != nil && err != dberr.ErrNotFound {
		return err
	}

	if gs.AutoVoice, err = c.db.GetGuildAutoVoice(guildID); err != nil && err != dberr.ErrNotFound {
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
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/settings [post]
func (c *GuildsSettingsController) postGuildSettings(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	var err error

	gs := new(models.GuildSettings)
	if err = ctx.BodyParser(gs); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if gs.AutoRoles != nil {
		if ok, _, err := c.pmw.HasPerms(c.session, guildID, uid, "dm.guild.config.autorole"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrForbidden
		}

		if arrayutils.ContainsAny(gs.AutoRoles, "@everyone") {
			return fiber.NewError(fiber.StatusBadRequest,
				"@everyone can not be set as autorole")
		}

		guildRoles, err := c.st.Roles(guildID, true)
		if err != nil {
			return err
		}
		guildRoleIDs := make([]string, len(guildRoles))
		for i, role := range guildRoles {
			guildRoleIDs[i] = role.ID
		}

		if c := arrayutils.ContainsAny(gs.AutoRoles, guildRoleIDs...); !c {
			return fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("one or more of the specified roles does not exist on this guild (%s)",
					strings.Join(gs.AutoRoles, ", ")))
		}

		if err = c.db.SetGuildAutoRoles(guildID, gs.AutoRoles); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	if gs.AutoVoice != nil {
		if ok, _, err := c.pmw.HasPerms(c.session, guildID, uid, "dm.guild.config.autovoice"); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		} else if !ok {
			return fiber.ErrForbidden
		}

		guildChannels, err := c.st.Channels(guildID, true)
		if err != nil {
			return err
		}
		guildChannelIDs := make([]string, len(guildChannels))
		for i, channel := range guildChannels {
			guildChannelIDs[i] = channel.ID
		}

		if c := arrayutils.ContainsAny(gs.AutoVoice, guildChannelIDs...); !c {
			return fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("one or more of the specified channels does not exist on this guild (%s)",
					strings.Join(gs.AutoVoice, ", ")))
		}

		if err = c.db.SetGuildAutoVoice(guildID, gs.AutoVoice); err != nil {
			return wsutil.ErrInternalOrNotFound(err)
		}
	}

	return ctx.JSON(models.Ok)
}
