package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func RecieveUserByID(DistID string, user models.User) (models.User, *gorm.DB) {
	result := configs.DB.First(&user, "dist_id = ?", DistID)
	return user, result
}
