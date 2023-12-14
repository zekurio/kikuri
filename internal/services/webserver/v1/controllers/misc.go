package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	sharedmodels "github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/webserver/auth"
	"github.com/zekurio/kikuri/internal/services/webserver/v1/models"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordutils"
)

type MiscController struct {
	st     *dgrs.State
	authMw auth.Middleware
	cfg    sharedmodels.Config
}

func (c *MiscController) Setup(ctn di.Container, router fiber.Router) {
	c.st = ctn.Get(static.DiState).(*dgrs.State)
	c.authMw = ctn.Get(static.DiAuthMiddleware).(auth.Middleware)
	c.cfg = ctn.Get(static.DiConfig).(sharedmodels.Config)

	router.Get("/me", c.authMw.Handle, c.getMe)
}

// @Summary Me
// @Description Returns the user object of the currently authenticated user.
// @Tags Misc
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /me [get]
func (c *MiscController) getMe(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	user, err := c.st.User(uid)
	if err != nil {
		return err
	}

	created, _ := discordutils.GetDiscordSnowflakeCreationTime(user.ID)

	res := &models.User{
		User:      user,
		AvatarURL: user.AvatarURL(""),
		CreatedAt: created,
		BotOwner:  uid == c.cfg.Discord.OwnerID,
	}

	return ctx.JSON(res)
}
