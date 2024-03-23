package service

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddEnquiryType(request dto.EnquiryTypeIn) (fiber.Map, int) {

	enquiryObj := &models.EnquiryType{
		Name:      request.Name,
		AdminName: request.AdminName,
	}

	repositories.SaveToEnquiryType(enquiryObj)

	for _, field := range request.Fields {

		enquiryFieldObj := &models.Field{
			EnquiryTypeID: enquiryObj.ID,
			Name:          field.Name,
			Type:          field.Type,
			Required:      field.Required,
		}
		repositories.SaveToEnquiryField(enquiryFieldObj)
	}

	return fiber.Map{"success": "Enquiry type added successfully"}, 200
}

func GetAllEnquiryType() fiber.Map {

	var result *gorm.DB

	enquiryTypes, result := repositories.GetAllEnquiryType()

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": enquiryTypes}
}

func SubmitContactUsQuery(request dto.ContactUsIn) (fiber.Map, int) {

	// enquiryObj := &models.ContactUs{
	// 	Name: request.Name,
	// 	DistribId: request.DistribId,
	// 	Country: request.Country,
	// 	ContactNumber: request.ContactNumber,
	// 	EmailAddress: request.EmailAddress,
	// 	EnquiryTypeID: request.EnquiryTypeID,
	// 	AadhaarFront: request.AadhaarFront,
	// 	AadhaarBack: request.AadhaarBack,
	// 	PanCard: request.PanCard,
	// 	PassportSize: request.PassportSize,

	// }

	// repositories.SaveContactUsQuery(enquiryObj)

	return fiber.Map{"success": "Enquiry type added successfully"}, 200
}
