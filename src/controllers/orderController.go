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

	res, status := service.PlaceOrder(OrderIn)
	return c.Status(status).JSON(res)
}

func GetOrdersByDistribId(c *fiber.Ctx) error {

	distrib_id := c.Params("distrib_id")

	res, status := service.GetOrdersByDistribId(distrib_id)
	return c.Status(status).JSON(res)

}
