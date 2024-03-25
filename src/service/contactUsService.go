package service

import (
	"mime/multipart"
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

func SubmitContactUsQuery(form *multipart.Form) (fiber.Map, int) {

	// enquiryFields := models.EnquiryField{
	// 	YourQuery: form.Value["your_query"],
	// }

	//create unique name
	// AadhaarBack, err := middleware.GenerateUniqueFilename(request.DistribId, filepath.Ext(request.AadhaarBack.Filename))
	// if err != nil {
	// 	return fiber.Map{"error": "Failed to generate filename"}, 500
	// }
	// AadhaarFront, err := middleware.GenerateUniqueFilename(request.DistribId, filepath.Ext(request.AadhaarFront.Filename))
	// if err != nil {
	// 	return fiber.Map{"error": "Failed to generate filename"}, 500
	// }
	// PanCard, err := middleware.GenerateUniqueFilename(request.DistribId, filepath.Ext(request.PanCard.Filename))
	// if err != nil {
	// 	return fiber.Map{"error": "Failed to generate filename"}, 500
	// }
	// PassportSize, err := middleware.GenerateUniqueFilename(request.DistribId, filepath.Ext(request.PassportSize.Filename))
	// if err != nil {
	// 	return fiber.Map{"error": "Failed to generate filename"}, 500
	// }

	// //save file to server
	// //save filepath to db
	// contactUsObj := &models.ContactUs{
	// 	ContactUsBasicDetails: request.ContactUsBasicDetails,
	// 	EnquiryTypeID:         request.EnquiryTypeID,
	// 	EnquiryField:          enquiryFields,
	// 	ContactQueryFile:      request.ContactQueryFile,
	// }

	// repositories.SaveContactUsQuery(contactUsObj)
	return fiber.Map{"success": form}, 200

}
