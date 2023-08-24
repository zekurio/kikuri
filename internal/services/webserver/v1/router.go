package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
	"github.com/zekurio/daemon/internal/services/webserver/v1/controllers"
)

type Router struct {
	container di.Container
}

func (r *Router) SetContainer(container di.Container) {
	r.container = container
}

func (r *Router) Route(router fiber.Router) {
	new(controllers.PublicController).Setup(r.container, router.Group("/public"))
}
