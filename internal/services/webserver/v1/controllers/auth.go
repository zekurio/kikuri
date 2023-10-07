package controllers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekrotja/dgrs"
	"github.com/zekrotja/ken"
	"github.com/zekurio/daemon/internal/services/database/dberr"
	"github.com/zekurio/daemon/internal/services/webserver/auth"
	"github.com/zekurio/daemon/internal/services/webserver/v1/models"
	"github.com/zekurio/daemon/internal/util/static"
	"github.com/zekurio/daemon/pkg/discordoauth"
)

type AuthController struct {
	discordOAuth *discordoauth.DiscordOAuth
	rth          auth.RefreshTokenHandler
	ath          auth.AccessTokenHandler
	authMw       auth.Middleware
	oauthHandler auth.RequestHandler
	st           *dgrs.State
	session      *discordgo.Session
	cmdHandler   *ken.Ken
}

func (c *AuthController) Setup(container di.Container, router fiber.Router) {
	c.discordOAuth = container.Get(static.DiDiscordOAuth).(*discordoauth.DiscordOAuth)
	c.rth = container.Get(static.DiAuthRefreshTokenHandler).(auth.RefreshTokenHandler)
	c.ath = container.Get(static.DiAuthAccessTokenHandler).(auth.AccessTokenHandler)
	c.authMw = container.Get(static.DiAuthMiddleware).(auth.Middleware)
	c.st = container.Get(static.DiState).(*dgrs.State)
	c.session = container.Get(static.DiDiscordSession).(*discordgo.Session)
	c.oauthHandler = container.Get(static.DiOAuthHandler).(auth.RequestHandler)
	c.cmdHandler = container.Get(static.DiCommandHandler).(*ken.Ken)

	router.Get("/login", c.getLogin)
	router.Get("/oauthcallback", c.discordOAuth.HandlerCallback)
	router.Post("/accesstoken", c.postAccessToken)
	router.Get("/check", c.authMw.Handle, c.getCheck)
	router.Post("/logout", c.authMw.Handle, c.postLogout)
}

func (c *AuthController) getLogin(ctx *fiber.Ctx) error {
	state := make(map[string]string)

	if redirect := ctx.Query("redirect"); redirect != "" {
		state["redirect"] = redirect
	}

	return c.discordOAuth.HandlerInitWithState(ctx, state)
}

// @Summary Access Token Exchange
// @Description Exchanges the cookie-passed refresh token with a generated access token.
// @Tags Authorization
// @Accept json
// @Produce json
// @Success 200 {object} models.AccessTokenResponse
// @Failure 401 {object} models.Error
// @Router /auth/accesstoken [post]
func (c *AuthController) postAccessToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Cookies(static.RefreshTokenCookieName)
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	ident, err := c.rth.ValidateRefreshToken(refreshToken)
	if err != nil && err != dberr.ErrNotFound {
		log.With("err", err).Error("Failed validating refresh token")
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

// @Summary Authorization Check
// @Description Returns OK if the request is authorized.
// @Tags Authorization
// @Accept json
// @Produce json
// @Success 200 {object} models.Status
// @Failure 401 {object} models.Error
// @Router /auth/check [get]
func (c *AuthController) getCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(models.Ok)
}

// @Summary Logout
// @Description Reovkes the currently used access token and clears the refresh token.
// @Tags Authorization
// @Accept json
// @Produce json
// @Success 200 {object} models.Status
// @Router /auth/logout [post]
func (c *AuthController) postLogout(ctx *fiber.Ctx) error {
	uid := ctx.Locals("uid").(string)

	err := c.rth.RevokeToken(uid)
	if err != nil && err != dberr.ErrNotFound {
		return err
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
