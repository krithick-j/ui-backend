package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetLoadRequest(c *fiber.Ctx) error {
	filePath := "data/cpa/LoadRequest.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}

func GetMyAccountSummary(c *fiber.Ctx) error {
	filePath := "data/cpa/SettlementCatalog.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
func GetSettlementCatalog(c *fiber.Ctx) error {
	filePath := "data/cpa/SettlementCatalog.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
func GetTransactionSummary(c *fiber.Ctx) error {
	filePath := "data/cpa/TransactionSummary.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
