package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"
)

func SaveICoupon(iCoupon models.ICoupon) error {

	result := configs.DB.Create(&iCoupon)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return nil
}
