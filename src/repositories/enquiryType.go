package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveToEnquiryType(tx models.EnquiryType) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}

func GetAllEnquiryType(name []string) ([]string, *gorm.DB) {
	result := configs.DB.Model(models.EnquiryType{}).Select("name").Find(&name) //fill the empty array variable
	return name, result                                                         //return the filled array variable
}
