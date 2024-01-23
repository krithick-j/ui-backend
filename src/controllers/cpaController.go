package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetSettlementCatalog(c *fiber.Ctx) error {
	filePath := "data/SettlementCatalog.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
