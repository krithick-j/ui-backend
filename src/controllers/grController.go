package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func GetAllGrVisual(c *fiber.Ctx) error {
	filePath := "data/GR-Visual.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}

func GetAllGrVisualByDate(c *fiber.Ctx) error {
	filePath := "data/GR-VisualByDate.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
