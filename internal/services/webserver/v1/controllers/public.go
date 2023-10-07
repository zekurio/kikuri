package controllers

import (
	"github.com/zekurio/daemon/internal/services/database"
	"github.com/zekurio/daemon/internal/services/webserver/v1/models"
	"github.com/zekurio/daemon/internal/util/static"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type PublicController struct {
	session *discordgo.Session
	db      database.Database
}

func (c *PublicController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.db = ctn.Get(static.DiDatabase).(database.Database)

	router.Get("/guilds/:guildid", c.getGuild)
}

// @Summary Get Public Guild
// @Description Returns public guild information.
// @Tags Public
// @Accept json
// @Produce json
// @Param id path string true "Guild ID"
// @Success 200 {object} models.GuildReduced
// @Router /public/guilds/{id} [get]
func (c *PublicController) getGuild(ctx *fiber.Ctx) error {
	guildID := ctx.Params("guildid")

	guild, err := c.session.Guild(guildID)
	if err != nil {
		return err
	}

	// populate GuildReduced
	guildReduced := &models.GuildReduced{
		ID:          guild.ID,
		Name:        guild.Name,
		Icon:        guild.Icon,
		IconURL:     guild.IconURL(""),
		Region:      guild.Region,
		OwnerID:     guild.OwnerID,
		JoinedAt:    guild.JoinedAt,
		MemberCount: guild.MemberCount,
	}

	ctx.Set("Access-Control-Allow-Methods", "GET")
	ctx.Set("Access-Control-Allow-Headers", "*")

	return ctx.JSON(guildReduced)
}
