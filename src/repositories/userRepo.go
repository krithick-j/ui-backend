package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func GetUserByID(DistID string, user models.User) (models.User, *gorm.DB) {
	result := configs.DB.First(&user, "dist_id = ?", DistID)
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

func GetSide(distrib_id string, place string) string {
	var next_id string
	configs.DB.Table("users").Select(place).Where("dist_id=?", distrib_id).Row().Scan(&next_id)
	return next_id
}

func GetTrackingCenter(distrib_id, center_id string) *models.TrackingCenter {
	tc := new(models.TrackingCenter)
	configs.DB.First(tc, "distrib_id = ? AND center_code = ?", distrib_id, center_id)
	return tc
}

func CreateUser(user models.User) error {
	configs.DB.Create(&user)
	return nil
}

func CreateTCs(tcs []models.TrackingCenter) error {
	for _, tc := range tcs {
		res := configs.DB.Create(&tc)
		if res.Error != nil {
			fmt.Printf("Error %v\n", res.Error.Error())
		}
	}
	return nil
}

func UpdateTC(distrib_id string, tracking_center string, center_code string, side string) error {
	var updatecols models.TrackingCenter
	if side == "left" {
		updatecols = models.TrackingCenter{LeftDistribID: distrib_id, LeftPlace: "001"}
	} else {
		updatecols = models.TrackingCenter{RightDistribID: distrib_id, RightPlace: "001"}
	}
	res := configs.DB.Table("tracking_centers").Where("distrib_id=? AND center_code =?", tracking_center, center_code).Updates(updatecols)
	if res.Error != nil {
		fmt.Print(res.Error.Error())
	}
	return nil
}
