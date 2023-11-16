package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/permissions"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/arrayutils"
)

type GuildsController struct {
	db      database.Database
	pmw     *permissions.Permissions
	session *discordgo.Session
	st      *dgrs.State
	cfg     models.Config
}

func (c *GuildsController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.pmw = ctn.Get(static.DiPermissions).(*permissions.Permissions)
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = ctn.Get(static.DiState).(*dgrs.State)
	c.cfg = ctn.Get(static.DiConfig).(models.Config)

	router.Get("", c.getGuilds)
	router.Get("/:guildid", c.getGuild)
}

// @Summary Get Guilds
// @Description Returns all guilds the bot and the user have in common.
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
// @Description Returns a single guild object by its ID.
// @Tags Guilds
// @Accept json
// @Produce json
// @Param id path string true "The ID of the guild."
// @Success 200 {object} models.Guild
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /guilds/{id} [get]
func (c *GuildsController) getGuild(ctx *fiber.Ctx) (err error) {
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

	gRes, err := models.GuildFromDiscordGuild(guild, memb, c.cfg.Discord.OwnerID)
	if err != nil {
		return err
	}

	return ctx.JSON(gRes)
}
