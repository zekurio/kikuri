package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sarulabs/di/v2"
)

type GuildMembersController struct{}

func (c *GuildMembersController) Setup(container di.Container, router fiber.Router) {}
