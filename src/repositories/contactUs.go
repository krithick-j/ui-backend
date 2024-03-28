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

func GetAllContactQueries() ([]models.ContactUs, *gorm.DB) {
	var obj []models.ContactUs

	result := configs.DB.Find(&obj)
	return obj, result
}

func GetContactQueryStatusById(id string) (bool, *gorm.DB) {
	var status bool

	result := configs.DB.Model(models.ContactUs{}).Select("status").Find(&status, id)
	return status, result
}

func SwitchContactQueryStatusById(id string) *gorm.DB {
	status, _ := GetContactQueryStatusById(id)
	result := configs.DB.Model(models.ContactUs{}).Where("id=?", id).Update("status", !status)
	return result
}
