package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetEpBalanceByDistribId(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetEpBalance(distrib_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}
