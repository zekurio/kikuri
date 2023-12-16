package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/util/static"
)

var (
	errInvalidAccessToken = fiber.NewError(fiber.StatusUnauthorized, "invalid access token")
)

type TokenMiddleware struct {
	acth  AccessTokenHandler
	apith APITokenHandler
}

func NewTokenMiddleware(ctn di.Container) *TokenMiddleware {
	return &TokenMiddleware{
		acth:  ctn.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
		apith: ctn.Get(static.DiAuthAPITokenHandler).(APITokenHandler),
	}
}

func (m *TokenMiddleware) Handle(ctx *fiber.Ctx) (err error) {
	var ident string

	authHeader := ctx.Get("authorization")
	if authHeader == "" {
		return errInvalidAccessToken
	}

	split := strings.Split(authHeader, " ")
	if len(split) < 2 {
		return errInvalidAccessToken
	}

	switch strings.ToLower(split[0]) {

	case "accesstoken":
		if ident, err = m.acth.ValidateAccessToken(split[1]); err != nil || ident == "" {
			return errInvalidAccessToken
		}

	case "bearer":
		if ident, err = m.apith.ValidateAPIToken(split[1]); err != nil || ident == "" {
			return errInvalidAccessToken
		}

	default:
		return fiber.ErrUnauthorized
	}

	return next(ctx, ident)
}

func next(ctx *fiber.Ctx, ident string) error {
	ctx.Locals("uid", ident)
	return ctx.Next()
}
