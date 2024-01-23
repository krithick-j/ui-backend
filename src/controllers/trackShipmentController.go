package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func TrackShipment(c *fiber.Ctx) error {
	filePath := "data/TrackShipment.json"
	return c.Status(http.StatusOK).SendFile(filePath)
}
