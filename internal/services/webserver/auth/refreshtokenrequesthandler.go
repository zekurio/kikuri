package auth

import (
	"github.com/zekurio/kikuri/internal/models"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/util/static"
	"github.com/zekurio/kikuri/pkg/debug"
	"github.com/zekurio/kikuri/pkg/discordoauth"
)

type RefreshTokenRequestHandler struct {
	session             *discordgo.Session
	accessTokenHandler  AccessTokenHandler
	refreshTokenHandler RefreshTokenHandler
}

func NewRefreshTokenRequestHandler(container di.Container) *RefreshTokenRequestHandler {
	return &RefreshTokenRequestHandler{
		session:             container.Get(static.DiDiscordSession).(*discordgo.Session),
		accessTokenHandler:  container.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
		refreshTokenHandler: container.Get(static.DiAuthRefreshTokenHandler).(RefreshTokenHandler),
	}
}

func (h *RefreshTokenRequestHandler) LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error {
	return fiber.NewError(status, msg)
}

func (h *RefreshTokenRequestHandler) BindRefreshToken(ctx *fiber.Ctx, uid string) error {
	user, _ := h.session.User(uid)
	if user == nil {
		return fiber.ErrUnauthorized
	}

	ctx.Locals("uid", uid)
	refreshToken, err := h.refreshTokenHandler.GetRefreshToken(uid)
	if err != nil {
		return err
	}

	expires := time.Now().Add(static.AuthSessionExpiration)
	ctx.Cookie(&fiber.Cookie{
		Name:     static.RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		Expires:  expires,
		HTTPOnly: true,
		Secure:   !debug.Enabled(),
	})

	return nil
}

func (h *RefreshTokenRequestHandler) LoginSuccessHandler(ctx *fiber.Ctx, res discordoauth.SuccessResponse) error {
	if err := h.BindRefreshToken(ctx, res.UserID); err != nil {
		return err
	}

	location := "/"
	if redirect, ok := res.State["redirect"].(string); ok {
		location += strings.TrimLeft(redirect, "/")
	}

	return ctx.Redirect(location, fiber.StatusTemporaryRedirect)
}

func (h *RefreshTokenRequestHandler) LogoutHandler(ctx *fiber.Ctx) error {
	if uid, ok := ctx.Locals("uid").(string); ok && uid != "" {
		if err := h.refreshTokenHandler.RevokeToken(uid); err != nil {
			return err
		}
	}

	ctx.ClearCookie(static.RefreshTokenCookieName)

	return ctx.JSON(models.Ok)
}
