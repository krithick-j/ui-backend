package controllers

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetICouponHistory(c *fiber.Ctx) error {
	var payload dto.ICouponHistoryIn
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.GetICouponHistory(payload)
	return c.Status(status).JSON(res)
}

func GetBvHistory(c *fiber.Ctx) error {
	var payload dto.BvHistoryIn
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.GetBvHistory(payload)
	return c.Status(status).JSON(res)
}
