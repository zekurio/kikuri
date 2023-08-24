package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

// Router is used to register routes for the webserver
// and is used to register endpoints
type Router interface {
	SetContainer(ctn di.Container)

	Route(router fiber.Router)
}
