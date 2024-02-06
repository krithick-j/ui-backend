package repositories

import (
	"errors"
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

func GetNextItem(distrib_id string, center_code string, place string) (string, string) {
	next_item := struct {
		LeftDistribID  string
		LeftPlace      string
		RightDistribID string
		RightPlace     string
	}{}

	configs.DB.Table("tracking_centers").Where("distrib_id=? AND center_code=?", distrib_id, center_code).First(&next_item)
	if place == "left" {
		return next_item.LeftDistribID, next_item.LeftPlace
	} else {
		return next_item.RightDistribID, next_item.RightPlace
	}
}

func GetTrackingCenter(distrib_id, center_id string) *models.TrackingCenter {
	tc := new(models.TrackingCenter)
	configs.DB.First(tc, "distrib_id = ? AND center_code = ?", distrib_id, center_id)
	return tc
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

func UpdateTC(tx *gorm.DB, distrib_id string, tracking_center string, center_code string, side string) error {
	var updatecols models.TrackingCenter
	if side == "left" {
		updatecols = models.TrackingCenter{LeftDistribID: distrib_id, LeftPlace: "001"}
	} else {
		updatecols = models.TrackingCenter{RightDistribID: distrib_id, RightPlace: "001"}
	}
	res := tx.Table("tracking_centers").Where("distrib_id=? AND center_code =?", tracking_center, center_code).Updates(updatecols)
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
