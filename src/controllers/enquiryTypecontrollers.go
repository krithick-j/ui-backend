package controllers

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func AddEnquiryType(c *fiber.Ctx) error {
	var request dto.EnquiryTypeIn
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res := service.AddEnquiryType(request)

	return c.Status(http.StatusOK).JSON(res)
}
