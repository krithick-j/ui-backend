package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveICoupon(iCoupon models.ICoupon) error {

	result := configs.DB.Create(&iCoupon) //insert into
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return nil
}

func GetAllICouponsByDistribID(distribID string, iCoupon []models.ICoupon) (*gorm.DB, []models.ICoupon) {
	result := configs.DB.Find(&iCoupon, "distrib_id=?", distribID)
	return result, iCoupon
}

func GetICoupon(VID string, Pin string, iCoupon models.ICoupon, distribID string) (models.ICoupon, *gorm.DB) {
	result := configs.DB.First(&iCoupon, "v_id=? AND pin=? AND distrib_id=?", VID, Pin, distribID)
	return iCoupon, result
}
