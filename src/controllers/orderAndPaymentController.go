package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func OrderAndPayment(c *fiber.Ctx) error {
	filePath := "data/OrderAndPayment.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
