package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordutils"
)

type InviteController struct {
	session *discordgo.Session
	st      *dgrs.State
}

func (c *InviteController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.st = ctn.Get(static.DiState).(*dgrs.State)

	router.Get("", c.getInvite)
}

func (c *InviteController) getInvite(ctx *fiber.Ctx) error {
	self, err := c.st.SelfUser()
	if err != nil {
		return err
	}
	return ctx.Redirect(discordutils.GetInviteLink(self.ID))
}
