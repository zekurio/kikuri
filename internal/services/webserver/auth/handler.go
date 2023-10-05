package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zekurio/daemon/pkg/discordoauth"
)

// RequestHandler provides fiber endpoints and handlers
// to authenticate users via an OAuth2 interface.
type RequestHandler interface {

	// LoginFailedHandler is called when an error occurred
	// during the authentication process.
	LoginFailedHandler(ctx *fiber.Ctx, status int, msg string) error

	// BindRefreshToken returns a refresh token for the user
	// and binds it to the user's session via a cookie.
	BindRefreshToken(ctx *fiber.Ctx, uid string) error

	// LoginSuccessHandler is called when the user has successfully
	// authenticated.
	LoginSuccessHandler(ctx *fiber.Ctx, res discordoauth.SuccessResult) error

	// LogoutHandler is called when the user wants to log out.
	LogoutHandler(ctx *fiber.Ctx) error
}
