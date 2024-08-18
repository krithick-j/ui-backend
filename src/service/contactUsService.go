package service

import (
	"path/filepath"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddEnquiryType(request dto.EnquiryTypeIn, tx *gorm.DB) (fiber.Map, int) {

	enquiryObj := &models.EnquiryType{
		Name:      request.Name,
		AdminName: request.AdminName,
	}

	err := repositories.SaveToEnquiryType(enquiryObj, tx)
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Errorln("Record not found", err.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling SaveToEnquiryType from AddEnquiryType service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	return fiber.Map{"success": "Enquiry type added successfully"}, fiber.StatusOK
}

func GetAllEnquiryType(tx *gorm.DB) (fiber.Map, int) {

	enquiryTypes, err := repositories.GetAllEnquiryType(tx)

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Errorln("Record not found", err.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetAllEnquiryType from GetAllEnquiryType service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": enquiryTypes}, fiber.StatusOK
}

func SubmitContactUsQuery(form dto.ContactQueryFileForm, textfield dto.ContactUsIn, tx *gorm.DB) (fiber.Map, models.ContactQueryFile, int) {
	var contactQueryFileObj models.ContactQueryFile
	var AadhaarBackFileName, AadhaarFrontFileName, PanCardFileName, PassportSizeFileName string

	//mapping textField
	enquiryFields := models.EnquiryField{
		YourQuery:             textfield.YourQuery,
		BankAccountValidation: textfield.BankAccountValidation,
		RefundRequest:         textfield.RefundRequest,
	}

	//create unique name
	AadhaarBackFileName, err := utils.GenerateUniqueFilename(textfield.DistribId, "aadhaarBack", filepath.Ext(form.AadhaarBack.Filename))
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("AadhaarBackFileName: Error on calling GenerateUniqueFilename from SubmitContactUsQuery service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, contactQueryFileObj, fiber.StatusInternalServerError
	}
	AadhaarFrontFileName, err = utils.GenerateUniqueFilename(textfield.DistribId, "aadhaarFront", filepath.Ext(form.AadhaarFront.Filename))
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("AadhaarFrontFileName: Error on calling GenerateUniqueFilename from SubmitContactUsQuery service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, contactQueryFileObj, fiber.StatusInternalServerError
	}
	PanCardFileName, err = utils.GenerateUniqueFilename(textfield.DistribId, "panCard", filepath.Ext(form.PanCard.Filename))
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("PanCardFileName: Error on calling GenerateUniqueFilename from SubmitContactUsQuery service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, contactQueryFileObj, fiber.StatusInternalServerError
	}
	PassportSizeFileName, err = utils.GenerateUniqueFilename(textfield.DistribId, "passportSize", filepath.Ext(form.PassportSize.Filename))
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("PassportSizeFileName: Error on calling GetAllContactQueries from SubmitContactUsQuery service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, contactQueryFileObj, fiber.StatusInternalServerError
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

	err = repositories.SaveContactUsQuery(contactUsObj, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling SaveContactUsQuery from GetAllContactQueries service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, contactQueryFileObj, fiber.StatusInternalServerError
	}

	return fiber.Map{"success": "File Uploaded successfully"}, contactQueryFileObj, fiber.StatusOK
}

func GetAllContactQueries(tx *gorm.DB) (fiber.Map, int) {

	contactQueries, err := repositories.GetAllContactQueries(tx)

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Errorln("Record not found", err.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetAllContactQueries from GetAllContactQueries service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	return fiber.Map{"data": contactQueries}, fiber.StatusOK
}

func SwitchContactQueryStatus(id string, tx *gorm.DB) (fiber.Map, int) {
	status, err := repositories.GetContactQueryStatusById(id, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetContactQueryStatusById from SwitchContactQueryStatus service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	err = repositories.SwitchContactQueryStatusById(status, id, tx)
	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		configs.Log.Errorln("Record not found", err.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling SwitchContactQueryStatusById from SwitchContactQueryStatus service fn ", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": "Contact Query Status Changed"}, fiber.StatusOK
}
