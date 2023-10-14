package controllers

import (
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
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/arrayutils"
	"github.com/zekurio/daemon/pkg/perms"
)

type GuildsController struct {
	db      database.Database
	session *discordgo.Session
	cfg     config.Config
	pmw     *permissions.Permissions
	st      *dgrs.State
}

func (c *GuildsController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = ctn.Get(static.DiConfig).(config.Config)
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.pmw = ctn.Get(static.DiPermissions).(*permissions.Permissions)
	c.st = ctn.Get(static.DiState).(*dgrs.State)

	router.Get("", c.getGuilds)
	router.Get("/:guildid", c.getGuild)
	router.Post("/:guildid/permissions", c.postGuildPermissions)
}

// @Summary List Guilds
// @Description Returns a list of guilds the authenticated user has in common with shinpuru.
// @Tags Guilds
// @Accept json
// @Produce json
// @Success 200 {array} models.GuildReduced "Wrapped in models.ListResponse"
// @Failure 401 {object} models.Error
// @Router /guilds [get]
func (c *GuildsController) getGuilds(ctx *fiber.Ctx) (err error) {
	uid := ctx.Locals("uid").(string)

	guilds, err := c.st.Guilds()
	if err != nil {
		return err
	}

	userGuilds, err := c.st.UserGuilds(uid)
	if err != nil {
		return
	}

	guildRs := make([]*models.GuildReduced, len(userGuilds))
	i := 0
	for _, guild := range guilds {
		if arrayutils.ContainsAny(userGuilds, guild.ID) {
			guildRs[i] = models.GuildReducedFromGuild(guild)
			i++
		}
	}
	guildRs = guildRs[:i]

	return ctx.JSON(models.NewListResponse(guildRs))
}

// @Summary Get Guild
// @Description Returns a single guild object by it's ID.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.Guild
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id} [get]
func (c *GuildsController) getGuild(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	guildID := ctx.Params("guildid")

	memb, _ := c.st.Member(guildID, uid)
	if memb == nil {
		return fiber.ErrNotFound
	}

	guild, err := c.st.Guild(guildID, true)
	if err != nil {
		return err
	}

	gRes, err := models.GuildFromGuild(guild, memb, c.db, c.cfg.Discord.OwnerID)
	if err != nil {
		return err
	}

	return ctx.JSON(gRes)
}

// @Summary Apply Guild Permission Rule
// @Description Apply a new guild permission rule for a specified role.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Param payload body models.PermissionsUpdate true "The permission rule payload."
// @Success 200 {object} models.PermissionsMap
// @Failure 400 {object} models.Error
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id}/permissions [post]
func (c *GuildsController) postGuildPermissions(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	update := new(models.PermissionsUpdate)
	if err := ctx.BodyParser(update); err != nil {
		return fiber.ErrBadRequest
	}

	newPerms := update.Perm[1:]
	if !strings.HasPrefix(newPerms, "dm.guild") && !strings.HasPrefix(newPerms, "dm.etc") && !strings.HasPrefix(newPerms, "dm.chat") {
		return fiber.NewError(fiber.StatusBadRequest, "you can only give permissions over the 'dm.guild', 'dm.etc' and 'dm.chat' permission groups")
	}

	currPerms, err := c.db.GetPermissions(guildID)
	if err != nil && err == dberr.ErrNotFound {
		return fiber.ErrNotFound
	} else if err != nil {
		return err
	}

	for _, roleID := range update.RoleIDs {
		rperms, ok := currPerms[roleID]
		if !ok {
			rperms = perms.PermsArray{}
		}

		rperms, changed := rperms.Update(update.Perm, update.Override)

		if len(rperms) == 0 {
			delete(currPerms, roleID)
		} else {
			currPerms[roleID] = rperms
		}

		if changed {
			if err = c.db.SetPermissions(guildID, roleID, rperms); err != nil {
				return err
			}
		}
	}

	return ctx.JSON(currPerms)
}
