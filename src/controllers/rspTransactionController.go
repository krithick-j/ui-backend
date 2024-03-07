package controllers

import (
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetTotalRspByDistribID(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")
	res, status := service.GetTotalRspByDistribId(distrib_id)
	return c.Status(status).JSON(res)
}
