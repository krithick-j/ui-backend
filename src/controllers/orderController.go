package controllers

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func PlaceOrder(c *fiber.Ctx) error {

	var OrderIn dto.PlaceOrderIn

	distribId := c.Params("distrib_id")

	if err := c.BodyParser(&OrderIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.PlaceOrder(OrderIn, distribId)
	return c.Status(status).JSON(res)
}
