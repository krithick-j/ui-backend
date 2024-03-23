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

	res, status := service.AddEnquiryType(request)

	return c.Status(status).JSON(res)
}

func GetAllEnquiryType(c *fiber.Ctx) error {

	res := service.GetAllEnquiryType()

	return c.Status(http.StatusOK).JSON(res)
}

func SubmitContactUsQuery(c *fiber.Ctx) error {
	var request dto.ContactUsIn
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res, status := service.SubmitContactUsQuery(request)

	return c.Status(status).JSON(res)
}
