package service

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddEnquiryType(request dto.EnquiryTypeIn) (fiber.Map, int) {
	for _, name := range request.Name {
		enquiryObj := models.EnquiryType{
			Name:      name,
			AdminName: request.AdminName,
		}
		repositories.SaveToEnquiryType(enquiryObj)
	}

	return fiber.Map{"success": "Enquiry type added successfully"}, 200
}

func GetAllEnquiryType() fiber.Map {

	var name []string //create empty variable
	var result *gorm.DB

	name, result = repositories.GetAllEnquiryType(name) //passing empty array variable

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": name}
}
