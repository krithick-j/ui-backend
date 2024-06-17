package controllers

import (
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func TestController(c *fiber.Ctx) error {
	orderId := c.Params("orderId") //string
	res, status := service.Test(orderId)
	return c.Status(status).JSON(res)
}
