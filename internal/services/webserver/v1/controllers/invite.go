package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordutils"
)

type InviteController struct {
	session *discordgo.Session
}

func (c *InviteController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)

	router.Get("", c.getInvite)
}

func (c *InviteController) getInvite(ctx *fiber.Ctx) error {
	return ctx.Redirect(discordutils.GetInviteLink(c.session))
}
