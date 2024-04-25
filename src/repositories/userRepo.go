package repositories

import (
	"errors"
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func AuthUser(distrib_id, password string) (models.User, error) {
	user := models.User{}
	res := configs.DB.Where("distrib_id = ? AND pass =?", distrib_id, password).First(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return models.User{}, res.Error
	}
	return user, nil
}

func GetUserByID(DistID string, user models.User) (models.User, *gorm.DB) {
	result := configs.DB.First(&user, "distrib_id = ?", DistID)
	return user, result
}

func GetAllUsers(user []models.User) ([]models.User, *gorm.DB) {
	result := configs.DB.Find(&user)
	return user, result
}

func GetLastId() string {
	var lastNo string
	configs.DB.Table("users").Select("max(distrib_id)").Row().Scan(&lastNo)
	return lastNo
}

func GetNextItem(distrib_id string, place string, side string) (string, string) {
	next_item := struct {
		LeftDistribID  string
		LeftPlace      string
		RightDistribID string
		RightPlace     string
	}{}

	configs.DB.Table("tracking_centers").Where("distrib_id=? AND place=?", distrib_id, place).First(&next_item)
	fmt.Println("place -------------->", place)
	if side == "left" {
		return next_item.LeftDistribID, next_item.LeftPlace
	} else if side == "right" {
		return next_item.RightDistribID, next_item.RightPlace
	} else {
		return fmt.Sprintf("SIDE NOT PROPERLY GIVEN %v", side), "error"
	}
}

func CreateUser(tx *gorm.DB, user models.User) error {
	res := tx.Create(&user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func CreateTCs(tx *gorm.DB, tcs []models.TrackingCenter) error {
	for _, tc := range tcs {
		res := tx.Create(&tc)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func UpdateTC(tx *gorm.DB, distrib_id string, parent_distrib_id string, place string, side string) error {
	var updatecols models.TrackingCenter
	if side == "left" {
		fmt.Print("hi inside if left")
		updatecols = models.TrackingCenter{LeftDistribID: distrib_id, LeftPlace: "001"}
	} else if side == "right" {
		updatecols = models.TrackingCenter{RightDistribID: distrib_id, RightPlace: "001"}
	} else {
		return fmt.Errorf("ERROR IN UPDATE TC, SIDE VALUE IS %v", side)
	}
	res := tx.Table("tracking_centers").Where("distrib_id=? AND place =?", parent_distrib_id, place).Updates(updatecols)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func GetUserEmailByDistribID(distribID string) (*gorm.DB, string) {
	var email string
	result := configs.DB.Table("users").Select("email_address").Where("distrib_id=?", distribID).Find(&email)
	return result, email
}

func EditUserByDistId(distrib_id string, userIn models.User, user models.User) (models.User, *gorm.DB) {
	result := configs.DB.Model(models.User{}).Where("distrib_id=?", distrib_id).Updates(userIn)
	return userIn, result
}

func GetUserByRefDistribId(distrib_id string, user []models.User) ([]models.User, *gorm.DB) {
	result := configs.DB.Where("ref_distrib_id", distrib_id).Find(&user)
	return user, result
}

func GetRefDistribIdByDistribId(distribId string) (string, *gorm.DB) {
	var refDistribId string
	result := configs.DB.Model(&models.User{}).Select("ref_distrib_id").Take(&refDistribId)
	return refDistribId, result
}

func GetUserPassByDistribId(distribId string) (string, *gorm.DB) {
	var pass string
	result := configs.DB.Model(&models.User{}).Where("distrib_id=?", distribId).Select("pass").Take(&pass)
	fmt.Println("User pass", pass)
	return pass, result
}

func UpdatePassword(distrib_id string, newHashPass string) *gorm.DB {
	result := configs.DB.Model(models.User{}).Where("distrib_id=?", distrib_id).Update("pass", newHashPass)
	return result
}

func SaveOTP(Type string, Value string, OTP string) error {
	fmt.Println("Entering here as well")
	err := configs.DB.Create(&models.OTPVerify{Type: Type, Value: Value, OTP: OTP}).Error
	return err
}
