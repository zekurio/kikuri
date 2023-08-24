package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

// Controller is used to register endpoints for
// a specific section of an API
type Controller interface {

	// Setup is called when the webserver is started
	// and is used to register endpoints
	Setup(container di.Container, router fiber.Router)
}
