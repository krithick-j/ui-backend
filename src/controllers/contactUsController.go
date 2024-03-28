package controllers

import (
	"fmt"
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/middleware"
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
	//for parsing textfield in form data
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid text format", "err": err.Error()})
	}

	if form, err := c.MultipartForm(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form format", "err": err.Error()})
	} else {
		if fileFormObj, err := middleware.ParseForm(form.File); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot Parse input Form", "err": err.Error()})
		} else {
			res, formName, status := service.SubmitContactUsQuery(*fileFormObj, request)
			if status == 200 {

				if formName.AadhaarFront != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.AadhaarFront)
					c.SaveFile(&fileFormObj.AadhaarFront, path)
				}
				if formName.AadhaarBack != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.AadhaarBack)
					c.SaveFile(&fileFormObj.AadhaarBack, path)
				}
				if formName.PanCard != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.PanCard)
					c.SaveFile(&fileFormObj.PanCard, path)
				}
				if formName.PassportSize != "" {
					path := fmt.Sprintf("contactQueryUploads/%s", formName.PassportSize)
					c.SaveFile(&fileFormObj.PassportSize, path)
				}
			}
			return c.Status(status).JSON(res["success"])
		}
	}
}

func GetAllContactQueries(c *fiber.Ctx) error {

	res, status := service.GetAllContactQueries()

	return c.Status(status).JSON(res)
}

func SwitchContactQueryStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	res, status := service.SwitchContactQueryStatus(id)

	return c.Status(status).JSON(res)
}
