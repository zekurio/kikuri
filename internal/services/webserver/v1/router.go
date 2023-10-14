package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type Router struct {
	ctn di.Container
}

func (r *Router) SetContainer(ctn di.Container) {
	r.ctn = ctn
}

func (r *Router) Route(router fiber.Router) {
	// TODO build routes

	// TODO build middlewares
}
