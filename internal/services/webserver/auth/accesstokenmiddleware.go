package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/util/static"
)

var (
	errInvalidAccessToken = fiber.NewError(fiber.StatusUnauthorized, "invalid access token")
)

// AccessTokenMiddleware implements Middleware
// for Access and API tokens passed via an
// Authorization header
type AccessTokenMiddleware struct {
	ath   AccessTokenHandler
	apith APITokenHandler
}

// NewAccessTokenMiddleware initializes a new instance
// of AccessTokenMiddleware.
func NewAccessTokenMiddleware(ctn di.Container) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{
		ath:   ctn.Get(static.DiAuthAccessTokenHandler).(AccessTokenHandler),
		apith: ctn.Get(static.DiAuthAPITokenHandler).(APITokenHandler),
	}
}

func (m *AccessTokenMiddleware) Handle(ctx *fiber.Ctx) (err error) {
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
		ident, err := m.ath.ValidateAccessToken(split[1])
		if err != nil || ident == "" {
			return errInvalidAccessToken
		}
		return next(ctx, ident)

	case "bearer":
		ident, err := m.apith.ValidateAPIToken(split[1])
		if err != nil || ident == "" {
			return fiber.ErrUnauthorized
		}
		return next(ctx, ident)

	default:
		return fiber.ErrUnauthorized
	}
}

func next(ctx *fiber.Ctx, ident string) error {
	ctx.Locals("uid", ident)
	return ctx.Next()
}
