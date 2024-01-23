package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func NewReferrals(c *fiber.Ctx) error {
	filePath := "data/NewReferrals.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}

func AllReferrals(c *fiber.Ctx) error {
	filePath := "data/AllReferrals.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
