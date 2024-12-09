package controllers

import (
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func TestController(c *fiber.Ctx) error {
	res, status := service.Test()
	return c.Status(status).JSON(res)
}
