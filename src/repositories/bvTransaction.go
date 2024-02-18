package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetCurrentBvFromTc(distrib_id string, place string, bv int) (int, *gorm.DB) {
	result := configs.DB.Table("tracking_centers").Select("bv").Where("distrib_id=? AND place=?", distrib_id, place).Take(&bv)
	return bv, result
}

func RecordBvTx(tx models.BvTransaction) {

	result := configs.DB.Create(&tx)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}
}
