package controllers

import (
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func TestController(c *fiber.Ctx) error {
	distridId := c.Params("distribId") //string
	res, status := service.Test(distridId)
	return c.Status(status).JSON(res)
}
