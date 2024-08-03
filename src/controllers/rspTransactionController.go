package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetTotalRspByDistribID(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetTotalRspByDistribId(distrib_id, tx)

	return c.Status(status).JSON(res)
}

func GetGroupRspByDistribID(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetGroupRspByDistribId(distrib_id, tx)

	if err := tx.Commit().Error; err != nil {
		configs.Log.Errorln("Error on committing transaction:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(status).JSON(res)
}
