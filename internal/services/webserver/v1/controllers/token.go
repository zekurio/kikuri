package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/services/database"
	"github.com/zekurio/kikuri/internal/services/database/dberr"
	"github.com/zekurio/kikuri/internal/services/webserver/auth"
	"github.com/zekurio/kikuri/internal/services/webserver/v1/models"
	"github.com/zekurio/kikuri/internal/util/static"
)

type TokenController struct {
	db    database.Database
	apith auth.APITokenHandler
}

func (c *TokenController) Setup(ctn di.Container, router fiber.Router) {
	c.db = ctn.Get(static.DiDatabase).(database.Database)
	c.apith = ctn.Get(static.DiAuthAPITokenHandler).(auth.APITokenHandler)

	router.Get("", c.getToken)
	router.Post("", c.postToken)
	router.Delete("", c.deleteToken)
}

// @Summary API Token Info
// @Description Returns API Token metadata, not the token itself.
// @Tags Tokens
// @Accept json
// @Produce json
// @Success 200 {object} models.APITokenResponse
// @Failure 401 {object} models.Error
// @Failure 404 {object} models.Error "Is returned when no token was generated before."
// @Router /token [get]
func (c *TokenController) getToken(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	token, err := c.db.GetAPIToken(uid)
	if dberr.IsErrNotFound(err) {
		return fiber.NewError(fiber.StatusNotFound, "no token found")
	} else if err != nil {
		return err
	}

	tokenResp := &models.APITokenResponse{
		Created:    token.Created,
		Expires:    token.Expires,
		Hits:       token.Hits,
		LastAccess: token.LastAccess,
	}

	return ctx.JSON(tokenResp)
}

// @Summary API Token Generation
// @Description Generates an API Token and returns it and its metadata.
// @Tags Tokens
// @Accept json
// @Produce json
// @Success 200 {object} models.APITokenResponse
// @Failure 401 {object} models.Error
// @Router /token [post]
func (c *TokenController) postToken(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	token, expires, err := c.apith.GetAPIToken(uid)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.APITokenResponse{
		Created: time.Now(),
		Expires: expires,
		Token:   token,
	})
}

// @Summary API Token Deletion
// @Description Deletes the users API token.
// @Tags Tokens
// @Accept json
// @Produce json
// @Success 200 {object} models.Status
// @Failure 401 {object} models.Error
// @Router /token [delete]
func (c *TokenController) deleteToken(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	err := c.apith.RevokeToken(uid)
	if err != nil {
		return err
	}

	return ctx.JSON(models.Ok)
}
