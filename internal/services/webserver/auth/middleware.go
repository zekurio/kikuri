package auth

import "github.com/gofiber/fiber/v2"

// Middleware provides an authorization middleware for the webserver
type Middleware interface {
	Handle(ctx *fiber.Ctx) error
}
