package controllers

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/models"
	"github.com/zekurio/kikuri/internal/services/database/dberr"
	"github.com/zekurio/kikuri/internal/services/webserver/auth"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/discordoauth"
)

type AuthController struct {
	dOauth *discordoauth.DiscordOAuth
	rth    auth.RefreshTokenHandler
	ath    auth.AccessTokenHandler
	authMw auth.Middleware
}

func (c *AuthController) Setup(ctn di.Container, router fiber.Router) {
	c.dOauth = ctn.Get(static.DiDiscordOAuth).(*discordoauth.DiscordOAuth)
	c.rth = ctn.Get(static.DiAuthRefreshTokenHandler).(auth.RefreshTokenHandler)
	c.ath = ctn.Get(static.DiAuthAccessTokenHandler).(auth.AccessTokenHandler)
	c.authMw = ctn.Get(static.DiAuthMiddleware).(auth.Middleware)

	router.Get("/login", c.getLogin)
	router.Get("/oauthcallback", c.dOauth.HandlerCallback)
	router.Post("/accesstoken", c.postAccessToken)
	router.Get("/check", c.authMw.Handle, c.getCheck)
	router.Post("/logout", c.authMw.Handle, c.postLogout)
}

func (c *AuthController) getLogin(ctx *fiber.Ctx) error {
	state := make(map[string]interface{})

	if redirect := ctx.Query("redirect"); redirect != "" {
		state["redirect"] = redirect
	}

	return c.dOauth.HandlerInitWithState(ctx, state)
}

func (c *AuthController) postAccessToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies(static.RefreshTokenCookieName)
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	ident, err := c.rth.ValidateRefreshToken(refreshToken)
	if err != nil && !dberr.IsErrNotFound(err) {
		log.Error("Failed validating refresh token", err)
	}
	if ident == "" {
		return fiber.ErrUnauthorized
	}

	token, expires, err := c.ath.GetAccessToken(ident)
	if err != nil {
		return err
	}

	return ctx.JSON(&models.AccessTokenResponse{
		Token:   token,
		Expires: expires,
	})
}

func (c *AuthController) getCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(models.Ok)
}

func (c *AuthController) postLogout(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	err := c.rth.RevokeToken(uid)
	if err != nil && !dberr.IsErrNotFound(err) {
		return err
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
