package controllers

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func IsCheckqueAvailable(c *fiber.Ctx) error {

	var chequeAvailableIn dto.ChequeAvailableIn

	if err := c.BodyParser(&chequeAvailableIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.IsCheckqueAvailable(chequeAvailableIn)
	return c.Status(status).JSON(res)
}

func TakeChequeByDistribId(c *fiber.Ctx) error {

	var TakeChequeIn dto.CheckoutIn

	if err := c.BodyParser(&TakeChequeIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.TotalChequeValueByDistribId(TakeChequeIn)
	return c.Status(status).JSON(res)
}