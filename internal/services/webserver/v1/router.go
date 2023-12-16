package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/kikuri/internal/services/webserver/auth"
	"github.com/zekurio/kikuri/internal/services/webserver/v1/controllers"
	"github.com/zekurio/kikuri/internal/util/static"
)

type Router struct {
	ctn di.Container
}

func (r *Router) SetContainer(ctn di.Container) {
	r.ctn = ctn
}

// @Title Kikuri API
// @Description Kikuri API.
// @Version 1.0

// @Tag.Name Authorization
// @TagDescription Authorization endpoints.

// @Tag.Name Misc
// @TagDescription Miscellaneous endpoints.

// @Tag.Name Users
// @TagDescription User endpoints.

// @Tag.Name Guilds
// @TagDescription Guild endpoints.

// @Tag.Name Guild Settings
// @TagDescription Guild settings endpoints.

// @Tag.Name Guild Members
// @TagDescription Guild member endpoints.

// @Tag Name Search
// @TagDescription Search endpoints.

// @Tag Name Token
// @TagDescription API Token endpoints.

// @BasePath /api/v1

func (r *Router) Route(router fiber.Router) {
	authMw := r.ctn.Get(static.DiAuthMiddleware).(auth.Middleware)

	new(controllers.AuthController).Setup(r.ctn, router.Group("/auth"))
	new(controllers.MiscController).Setup(r.ctn, router)

	// Apply auth middleware to all routes below

	router.Use(authMw.Handle)

	new(controllers.UsersController).Setup(r.ctn, router.Group("/users"))
	new(controllers.GuildsController).Setup(r.ctn, router.Group("/guilds"))
	new(controllers.GuildSettingsController).Setup(r.ctn, router.Group("/guilds/:guildid/settings"))
	new(controllers.GuildMembersController).Setup(r.ctn, router.Group("/guilds/:guildid"))
	new(controllers.SearchController).Setup(r.ctn, router.Group("/search"))
	new(controllers.TokenController).Setup(r.ctn, router.Group("/token"))
}
