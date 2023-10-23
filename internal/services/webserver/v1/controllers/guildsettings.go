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
	"github.com/zekurio/daemon/internal/util/static"
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
	//router.Post("", c.postGuildSettings)
}

func (c *GuildSettingsController) getGuildSettings(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	gs := new(models.GuildSettings)
	var err error

	if gs.AutoRoles, err = c.db.GetGuildAutoRoles(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	if gs.AutoVoice, err = c.db.GetGuildAutoVoice(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	if gs.Perms, err = c.db.GetPermissions(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	if gs.APIEnabled, err = c.db.GetGuildAPIEnabled(guildID); err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	return ctx.JSON(gs)
}
