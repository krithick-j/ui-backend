package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func AddTc(c *fiber.Ctx) error {
	var payload dto.AddTc
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing request from AddToCartController controller fn ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	res, status := service.AddTc(payload)
	return c.Status(status).JSON(res)
}

func AvailableTc(c *fiber.Ctx) error {

	distribId := c.Params("distrib_id")

	res, status := service.GetAvailableAddTc(distribId)
	return c.Status(status).JSON(res)
}
