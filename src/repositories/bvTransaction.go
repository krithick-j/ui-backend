package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"
)

func SaveBvTransaction(BVtx models.BvTransaction) {

	result := configs.DB.Create(&BVtx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}
