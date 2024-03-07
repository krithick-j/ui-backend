package service

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func AddEnquiryType(request dto.EnquiryTypeIn) fiber.Map {
	for _, name := range request.Name {
		enquiryObj := models.EnquiryType{
			Name:      name,
			AdminName: request.AdminName,
		}
		repositories.SaveToEnquiryType(enquiryObj)
	}

	return fiber.Map{"success": "Enquiry type added successfully"}
}
