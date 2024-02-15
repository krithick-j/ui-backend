package controllers

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func PlaceOrder(c *fiber.Ctx) error {

	var OrderIn dto.PlaceOrderIn

	if err := c.BodyParser(&OrderIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	println("distrib id from place order controller", orderDetailsMap.DistribId)
	res, status := service.PlaceOrder(OrderIn, orderDetailsMap)
	return c.Status(status).JSON(res)
}
