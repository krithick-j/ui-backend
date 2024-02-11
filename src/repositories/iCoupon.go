package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
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

func GetAllICouponsByDistribID(distribID string, iCoupons []models.ICoupon) (*gorm.DB, []models.ICoupon) {
	result := configs.DB.Find(&iCoupons, "distrib_id=?", distribID)
	return result, iCoupons

}
func GetICoupon(VID string, Pin string, iCouponsOut dto.ValidateICouponOut, distribID string) (dto.ValidateICouponOut, *gorm.DB) {
	result := configs.DB.Table("i_coupons").Select("value").Where("distrib_id=? and v_id=? and pin=?", distribID, VID, Pin).Find(&iCouponsOut)
	return iCouponsOut, result
}
