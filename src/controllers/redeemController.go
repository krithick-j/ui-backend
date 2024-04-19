package controllers

import (
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetEpBalanceByDistribId(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")
	res, status := service.GetEpBalance(distrib_id)
	return c.Status(status).JSON(res)
}
