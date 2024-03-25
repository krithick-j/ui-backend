package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

// Save object in enquiry_types table
func SaveToEnquiryType(tx *models.EnquiryType) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}

// get all objects from enquiry_type table
func GetAllEnquiryType() ([]models.EnquiryType, *gorm.DB) {
	var obj []models.EnquiryType

	result := configs.DB.Find(&obj)
	return obj, result
}

// insert object into contact_us table
func SaveContactUsQuery(tx *models.ContactUs) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}
