package service

import (
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserByDistId(dist_id string) fiber.Map {

	var user models.User
	var result *gorm.DB

	user, result = repositories.RecieveUserByID(dist_id, user)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}

	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}

	return fiber.Map{"data": user}
}
