package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"
)

func RecordICouponTx(tx models.ICouponTransaction) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}
