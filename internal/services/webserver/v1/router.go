package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/webserver/auth"
	"github.com/zekurio/daemon/internal/services/webserver/v1/controllers"
	"github.com/zekurio/daemon/internal/util/static"
)

type Router struct {
	container di.Container
}

func (r *Router) SetContainer(container di.Container) {
	r.container = container
}

func (r *Router) Route(router fiber.Router) {
	authMw := r.container.Get(static.DiAuthMiddleware).(auth.Middleware)

	new(controllers.PublicController).Setup(r.container, router.Group("/public"))
	new(controllers.OthersController).Setup(r.container, router.Group("/others"))
	new(controllers.InviteController).Setup(r.container, router.Group("/invite"))
	new(controllers.AuthController).Setup(r.container, router.Group("/auth"))

	// Requires authentication token
	router.Use(authMw.Handle)

	new(controllers.GuildsController).Setup(r.container, router.Group("/guilds"))

}
