package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/webserver/auth"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordutils"
)

type UsersController struct {
	session *discordgo.Session
	cfg     models.Config
	authMw  auth.Middleware
	st      *dgrs.State
}

func (c *UsersController) Setup(ctn di.Container, router fiber.Router) {
	c.session = ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	c.cfg = ctn.Get(static.DiConfig).(models.Config)
	c.authMw = ctn.Get(static.DiAuthMiddleware).(auth.Middleware)
	c.st = ctn.Get(static.DiState).(*dgrs.State)

	router.Get(":id", c.getUser)
}

// @Summary User
// @Description Returns a user by their id
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (c *UsersController) getUser(ctx *fiber.Ctx) error {
	uid := ctx.Params("id")

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
