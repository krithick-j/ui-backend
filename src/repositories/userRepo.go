package repositories

import (
	"fmt"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func AuthUser(distrib_id, password string, tx *gorm.DB) (models.User, error) {
	user := models.User{}
	err :=
		tx.
			Where("distrib_id = ? AND pass =?", distrib_id, password).
			First(&user).
			Error
	return user, err
}

func GetUserByID(DistID string, tx *gorm.DB) (models.User, error) {
	user := models.User{}
	err := tx.First(&user, "distrib_id = ?", DistID).Error
	return user, err
}

func GetAllUsers(tx *gorm.DB) ([]models.User, error) {
	var user []models.User
	err := tx.Find(&user).Error
	return user, err
}

func GetLastId(tx *gorm.DB) (string, error) {
	var lastNo string
	err :=
		tx.
			Table("users").
			Select("max(distrib_id)").
			Row().
			Scan(&lastNo)
	return lastNo, err
}

// Get next Tracking Center based on side: returns next distribId and place
func GetNextItem(distrib_id string, place string, side string, tx *gorm.DB) (struct {
	LeftDistribID  string
	LeftPlace      string
	RightDistribID string
	RightPlace     string
}, error) {
	next_item := struct {
		LeftDistribID  string
		LeftPlace      string
		RightDistribID string
		RightPlace     string
	}{}

	//Needs to be improvised
	err :=
		tx.
			Table("tracking_centers").
			Where("distrib_id=? AND place=?", distrib_id, place).
			First(&next_item).
			Error
	return next_item, err

}

func CreateUser(user models.User, tx *gorm.DB) error {
	err := tx.Create(&user).Error
	return err
}

func CreateTCs(tc models.TrackingCenter, tx *gorm.DB) error {
	err :=
		tx.
			Create(&tc).
			Error
	return err
}

func UpdateTC(tx *gorm.DB, distrib_id string, place string, parent_distrib_id string, parent_ref_place string, side string) error {
	var updatecols models.TrackingCenter
	if side == "left" {
		updatecols = models.TrackingCenter{LeftDistribID: distrib_id, LeftPlace: place}
	} else if side == "right" {
		updatecols = models.TrackingCenter{RightDistribID: distrib_id, RightPlace: place}
	} else {
		return fmt.Errorf("ERROR IN UPDATE TC, SIDE VALUE IS %v", side)
	}
	res := tx.Table("tracking_centers").Where("distrib_id=? AND place =?", parent_distrib_id, parent_ref_place).Updates(updatecols)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func GetUserEmailByDistribID(distribID string, tx *gorm.DB) (string, error) {
	var email string
	err := tx.Table("users").Select("email_address").Where("distrib_id=?", distribID).Find(&email).Error
	return email, err
}

func EditUserByDistId(distrib_id string, userIn models.User, tx *gorm.DB) (models.User, error) {
	err :=
		tx.
			Table("users").
			Where("distrib_id=?", distrib_id).
			Updates(userIn).
			Error
	return userIn, err
}

func GetUserByRefDistribId(distrib_id string, tx *gorm.DB) ([]models.User, error) {
	var user []models.User
	err :=
		tx.
			Table("users").
			Where("ref_distrib_id", distrib_id).
			Find(&user).
			Error
	return user, err
}

// This function will return who referred you
func GetRefDistribIdByDistribId(distribId string, tx *gorm.DB) (string, error) {
	var refDistribId string
	err :=
		tx.
			Table("users").
			Select("ref_distrib_id").
			Where("distrib_id = ?", distribId).
			Take(&refDistribId).
			Error
	return refDistribId, err
}

// This function will return who I am referred
// Suppose IN-00001 referred two persons, then the two person will come as a list
func GetReferredUsersByDistribId(distribId string, tx *gorm.DB) ([]string, error) {
	var referredDistributors []string
	err :=
		tx.
			Table("users").
			Distinct("distrib_id").
			Select("distrib_id").
			Where("ref_distrib_id = ? AND distrib_id != ref_distrib_id", distribId).
			Find(&referredDistributors).
			Error
	return referredDistributors, err
}

func GetUserPassByDistribId(distribId string, tx *gorm.DB) (string, error) {
	var pass string
	err :=
		tx.
			Table("users").
			Where("distrib_id=?", distribId).
			Select("pass").
			Take(&pass).
			Error
	return pass, err
}

func UpdatePassword(distrib_id string, newHashPass string, tx *gorm.DB) error {
	err :=
		tx.
			Model("users").
			Where("distrib_id=?", distrib_id).
			Update("pass", newHashPass).
			Error
	return err
}

func SaveOTP(Type string, Value string, OTP string, tx *gorm.DB) error {
	obj :=
		models.OTPVerify{
			Type:  Type,
			Value: Value,
			OTP:   OTP,
		}
	err :=
		tx.
			Create(&obj).
			Error

	return err
}

func CheckAndUpdateOTP(Type string, Value string, OTP string, tx *gorm.DB) error {

	err :=
		tx.
			Model(&models.OTPVerify{}).
			Where("type=? AND value=? AND otp=?", Type, Value, OTP).
			Update("status", "verified").
			Error

	return err
}

func UpdateOTPVerified(Type string, Value string, tx *gorm.DB) error {
	err :=
		tx.
			Model(&models.OTPVerify{}).
			Where("type=? AND value=?", Type, Value).
			Update("status", "verified").Error
	return err
}

func UpdateKyc(user *models.User, tx *gorm.DB) error {
	err :=
		tx.
			Model(user).
			Where("distrib_id", user.DistribID).
			Updates(user).
			Error
	return err
}

func GetChequePinByDistribID(distribId string, tx *gorm.DB) (string, error) {
	var pin string
	err :=
		tx.
			Model(&models.User{}).
			Select("cpa_pin").
			Where("distrib_id", distribId).
			Find(&pin).
			Error
	return pin, err
}

func ChangeCpaPin(distribId string, newPin string, tx *gorm.DB) error {
	err :=
		tx.
			Model(&models.User{}).
			Where("distrib_id", distribId).
			Update("cpa_pin", newPin).
			Error
	return err
}

func GetCurrentRankValueByDistribId(distribId string, tx *gorm.DB) (float64, error) {
	var currentRank float64
	err :=
		tx.
			Model(&models.User{}).
			Select("current_rank").
			Where("distrib_id", distribId).
			Take(&currentRank).
			Error
	return currentRank, err
}

func GetCurrentRankArrByDistribId(referredDistribIds []string, tx *gorm.DB) ([]float64, error) {
	var referredRanks []float64
	err := tx.Model(&models.User{}).Where("distrib_id IN ?", referredDistribIds).Pluck("current_rank", &referredRanks).Error
	return referredRanks, err
}
