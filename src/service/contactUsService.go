package service

import (
	"path/filepath"
	"ui-back-end/src/dto"
	"ui-back-end/src/middleware"
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

func SubmitContactUsQuery(form dto.ContactQueryFileForm, textfield dto.ContactUsIn) (fiber.Map, models.ContactQueryFile, int) {
	var contactQueryFileObj models.ContactQueryFile
	var AadhaarBackFileName, AadhaarFrontFileName, PanCardFileName, PassportSizeFileName string

	//mapping textField
	enquiryFields := models.EnquiryField{
		YourQuery:             textfield.YourQuery,
		BankAccountValidation: textfield.BankAccountValidation,
		RefundRequest:         textfield.RefundRequest,
	}

	//create unique name
	AadhaarBackFileName, err := middleware.GenerateUniqueFilename(textfield.DistribId, "aadhaarBack", filepath.Ext(form.AadhaarBack.Filename))
	if err != nil {
		return fiber.Map{"error": "Failed to generate filename"}, contactQueryFileObj, 500
	}
	AadhaarFrontFileName, err = middleware.GenerateUniqueFilename(textfield.DistribId, "aadhaarFront", filepath.Ext(form.AadhaarFront.Filename))
	if err != nil {
		return fiber.Map{"error": "Failed to generate filename"}, contactQueryFileObj, 500
	}
	PanCardFileName, err = middleware.GenerateUniqueFilename(textfield.DistribId, "panCard", filepath.Ext(form.PanCard.Filename))
	if err != nil {
		return fiber.Map{"error": "Failed to generate filename"}, contactQueryFileObj, 500
	}

	PassportSizeFileName, err = middleware.GenerateUniqueFilename(textfield.DistribId, "passportSize", filepath.Ext(form.PassportSize.Filename))
	if err != nil {
		return fiber.Map{"error": "Failed to generate filename"}, contactQueryFileObj, 500
	}

	//mapping filename to struct
	contactQueryFileObj = models.ContactQueryFile{
		AadhaarFront: AadhaarFrontFileName,
		AadhaarBack:  AadhaarBackFileName,
		PanCard:      PanCardFileName,
		PassportSize: PassportSizeFileName,
	}

	//save filepath to db
	contactUsObj := &models.ContactUs{
		ContactUsBasicDetails: textfield.ContactUsBasicDetails,
		EnquiryTypeID:         textfield.EnquiryTypeID,
		EnquiryField:          enquiryFields,
		ContactQueryFile:      contactQueryFileObj,
	}

	repositories.SaveContactUsQuery(contactUsObj)
	return fiber.Map{"success": "File Uploaded successfully"}, contactQueryFileObj, 200

}

func GetAllContactQueries() (fiber.Map, int) {

	var result *gorm.DB

	contactQueries, result := repositories.GetAllContactQueries()

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, 404
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, 500
	}

	return fiber.Map{"data": contactQueries}, 200
}

func SwitchContactQueryStatus(id string) (fiber.Map, int) {
	result := repositories.SwitchContactQueryStatusById(id)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, 404
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, 500
	}
	return fiber.Map{"data": "Contact Query Status Changed"}, 200

}
