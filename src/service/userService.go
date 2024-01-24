package service

import (
	"github.com/gofiber/fiber/v2"
)

func GetUserByDistId(id string) fiber.Map {

	// var user models.User
	// var result *gorm.DB

	// user, result = repositories.RecieveUserByID(id, user)

	// if result.Error == gorm.ErrRecordNotFound {
	// 	return fiber.Map{"data": "Not Found"}
	// }

	// if result.Error != nil {
	// 	return fiber.Map{"error": result.Error}
	// }

	return fiber.Map{"data": id}
}
