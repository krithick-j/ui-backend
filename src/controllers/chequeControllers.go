package controllers

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func IsChequeAvailable(c *fiber.Ctx) error {
	var chequeAvailableIn dto.Tc
	if err := c.BodyParser(&chequeAvailableIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.IsChequeAvailable(chequeAvailableIn)
	return c.Status(status).JSON(res)
}
