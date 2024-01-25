package controllers

import (
	"net/http"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetUserByDistId(c *fiber.Ctx) error {
	id := c.Params("dist_id")

	res := service.GetUserByDistId(id)

	return c.Status(http.StatusOK).JSON(res)
}
