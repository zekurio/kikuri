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

func (r *Router) Route(router fiber.Router) {
	authMw := r.ctn.Get(static.DiAuthMiddleware).(auth.Middleware)
	// TODO build routes
	new(controllers.AuthController).Setup(r.ctn, router.Group("/auth"))
	// TODO build middlewares
	router.Use(authMw.Handle)

	new(controllers.GuildSettingsController).Setup(r.ctn, router.Group("/guilds/:guildid/settings"))
	new(controllers.GuildsController).Setup(r.ctn, router.Group("/guilds"))
	new(controllers.UsersController).Setup(r.ctn, router.Group("/users"))
}
