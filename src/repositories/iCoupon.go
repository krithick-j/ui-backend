package repositories

import (
	"time"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveICoupon(iCoupon models.ICoupon, tx *gorm.DB) error {

	err := tx.Create(&iCoupon).Error //insert into
	return err
}

func GetAllICouponsByDistribID(distribID string, tx *gorm.DB) ([]models.ICoupon, error) {
	var iCoupons []models.ICoupon
	err := tx.
		Find(&iCoupons, "distrib_id=? and active=true", distribID).
		Error
	return iCoupons, err

}

// remaining Balance of the Coupon
func GetICouponBalance(VID string, tx *gorm.DB) (float64, error) {
	var balance float64
	err :=
		tx.
			Table("i_coupon_transactions").
			Select("sum(value)").
			Where("v_id=?", VID).
			Find(&balance).
			Error
	return balance, err
}

// Retrieves all the values using Order ID
func GetICouponRowsByOrderID(reference string, vid string, tx *gorm.DB) (float64, error) {
	var balance float64
	err := tx.
		Table("i_coupon_transactions").
		Select("sum(value)").
		Where("reference=? AND v_id=?", reference, vid).
		Find(&balance).
		Error
	return balance, err
}

func GetICouponExpiryDate(VID string, tx *gorm.DB) (time.Time, error) {
	var expiresOn time.Time
	err :=
		tx.
			Table("i_coupons").
			Select("expires_on").
			Where("v_id=?", VID).
			Find(&expiresOn).
			Error
	return expiresOn, err
}

// Validate Coupon
func ValidateICoupon(VID string, Pin string, tx *gorm.DB) (models.ICoupon, error) {
	var iCoupon models.ICoupon
	err :=
		tx.
			Table("i_coupons").
			Where("v_id=? and pin=?", VID, Pin).
			First(&iCoupon).
			Error
	return iCoupon, err
}

// value means initial or Total value when coupon generated
func GetICouponValue(VID string, Pin string, tx *gorm.DB) (float64, error) {
	var couponValue float64
	err :=
		tx.
			Table("i_coupons").
			Select("value").
			Where("v_id=? and pin=?", VID, Pin).
			Take(&couponValue).
			Error
	return couponValue, err
}

// value means initial or Total value when coupon generated
func GetICouponValueByVID(VID string, tx *gorm.DB) (float64, error) {
	var couponValue float64
	err :=
		tx.
			Table("i_coupons").
			Select("value").
			Where("v_id=?", VID).
			Take(&couponValue).
			Error

	return couponValue, err
}

// updating iCoupon Balance in Icoupons table
func SoftDeleteCoupon(coupon models.ICoupon, tx *gorm.DB) (models.ICoupon, error) {
	err :=
		tx.
			Where("v_id=?", coupon.VID).
			Delete(&coupon).
			Error
	return coupon, err
}

func CloseCoupon(VID string, tx *gorm.DB) error {
	err :=
		tx.
			Table("i_coupons").
			Where("v_id=?", VID).
			Update("active", false).
			Error //0 means coupon closed //1 means active
	return err
}

// Update Balance of the Coupon in ICoupons table
func UpdateBalanceInICoupons(VID string, balance float64, tx *gorm.DB) (float64, error) {
	err :=
		tx.
			Table("i_coupons").
			Where("v_id=?", VID).
			Update("balance", balance).
			Error
	return balance, err
}
