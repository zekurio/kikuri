package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type Controller interface {
	Setup(ctn di.Container, router fiber.Router)
}
