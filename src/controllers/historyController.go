package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetICouponHistory(c *fiber.Ctx) error {
	var payload dto.ICouponHistoryIn
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error parsing on payload from GetICouponHistory controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.GetICouponHistory(payload, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func GetBvHistory(c *fiber.Ctx) error {
	var payload dto.BvHistoryIn
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing GetBvHistory controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.GetBvHistory(payload, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}
