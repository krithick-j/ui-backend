package repositories

import (
	"crypto/sha256"
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
	configs.DB.Table("users").Select("max(dist_id)").Row().Scan(&lastNo)
	return lastNo
}

func GetSide(distrib_id string, place string) string {
	var next_id string
	configs.DB.Table("users").Select(place).Where("dist_id=?", distrib_id).Row().Scan(&next_id)
	return next_id
}

func CreateUser(user models.User) error {
	var side string
	if user.Place == "L" {
		side = "lside"
	} else {
		side = "rside"
	}
	tx := configs.DB.Begin()
	// Create User
	user.Pass = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Pass)))
	result := tx.Table("users").Create(&user)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	// Update Referrer Row
	result = tx.Table("users").Where("dist_id=?", user.RefDistID).Update(side, user.DistID)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}
