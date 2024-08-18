package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetBvCounterByStartDate(c *fiber.Ctx) error {

	var payload dto.BvCounterIn

	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing payload from GetBvCounterByStartDate controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.BvCounterByStartDate(payload, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}
