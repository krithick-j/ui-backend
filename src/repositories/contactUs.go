package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

// Save object in enquiry_types table
func SaveToEnquiryType(object *models.EnquiryType, tx *gorm.DB) error {

	err := tx.Create(&object).Error
	return err
}

// get all objects from enquiry_type table
func GetAllEnquiryType(tx *gorm.DB) ([]models.EnquiryType, error) {
	var obj []models.EnquiryType

	err := tx.Find(&obj).Error
	return obj, err
}

// insert object into contact_us table
func SaveContactUsQuery(obj *models.ContactUs, tx *gorm.DB) error {

	err := tx.Create(&obj).Error
	return err
}

func GetAllContactQueries(tx *gorm.DB) ([]models.ContactUs, error) {
	var obj []models.ContactUs

	err := tx.Find(&obj).Error
	return obj, err
}

func GetContactQueryStatusById(id string, tx *gorm.DB) (bool, error) {
	var status bool

	err := tx.Model(models.ContactUs{}).Select("status").Find(&status, id).Error
	return status, err
}

func SwitchContactQueryStatusById(status bool, id string, tx *gorm.DB) error {
	err := tx.Model(models.ContactUs{}).Where("id=?", id).Update("status", !status).Error
	return err
}
