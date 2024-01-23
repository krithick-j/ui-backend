package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CarousalImages(c *fiber.Ctx) error {
	filePath := "data/images.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
