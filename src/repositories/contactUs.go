package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveToEnquiryType(tx *models.EnquiryType) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}

func SaveToEnquiryField(tx *models.Field) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}

func GetAllEnquiryType() (models.EnquiryType, *gorm.DB) {
	var obj models.EnquiryType

	result := configs.DB.Preload("Fields").Find(&obj)
	return obj, result
}
