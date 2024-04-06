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

func GetAllICouponsByDistribID(distribID string, iCoupons []models.ICoupon) (*gorm.DB, []models.ICoupon) {
	result := configs.DB.Find(&iCoupons, "distrib_id=? and active=true", distribID)
	return result, iCoupons

}

// remaining Balance of the Coupon
func GetICouponBalance(VID string) (float64, *gorm.DB) {
	var balance float64
	result := configs.DB.Table("i_coupons").Select("balance").Where("v_id=?", VID).Take(&balance)
	return balance, result
}

// Validate Coupon
func ValidateICoupon(VID string, Pin string, iCoupon models.ICoupon) models.ICoupon {

	configs.DB.Table("i_coupons").Where("v_id=? and pin=?", VID, Pin).First(&iCoupon)
	return iCoupon
}

// value means initial or Total value when coupon generated
func GetICouponValue(VID string, Pin string) float64 {
	var couponValue float64
	result := configs.DB.Table("i_coupons").Select("value").Where("v_id=? and pin=?", VID, Pin).Take(&couponValue)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return couponValue
}

//updating iCoupon Balance in Icoupons table

func SoftDeleteCoupon(coupon models.ICoupon) (models.ICoupon, *gorm.DB) {
	result := configs.DB.Where("v_id=?", coupon.VID).Delete(&coupon)
	return coupon, result
}

func CloseCoupon(VID string) *gorm.DB {
	result := configs.DB.Table("i_coupons").Where("v_id=?", VID).Update("active", false) //0 means coupon closed //1 means active
	return result
}

// Update Balance of the Coupon in ICoupons table
func UpdateBalanceInICoupons(VID string, balance float64) (float64, *gorm.DB) {
	result := configs.DB.Table("i_coupons").Where("v_id=?", VID).Update("balance", balance)
	return balance, result
}
