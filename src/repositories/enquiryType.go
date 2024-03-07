package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"
)

func SaveToEnquiryType(tx models.EnquiryType) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}
