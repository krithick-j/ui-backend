package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func SocialConnController(c *fiber.Ctx) error {

	tx := configs.DB.Begin()
	res, status := service.SocialConnService(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}
