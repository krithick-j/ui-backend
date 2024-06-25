package controllers

import (
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func TestController(c *fiber.Ctx) error {
	orderId := c.Params("orderId") //string
	res, status, _ := service.GenerateInvoice(orderId)
	return c.Status(status).JSON(res)
}
