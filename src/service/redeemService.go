package service

import (
	"ui-back-end/configs"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetEpBalance(distribId string) (fiber.Map, int) {
	var balance float64

	balance, result := repositories.GetEpBalance(distribId)
	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Infoln("Record not found")
		return fiber.Map{"data": "Not Found"}, fiber.StatusOK
	}
	if result.Error != nil {
		configs.Log.Errorln("Error on calling GetEpBalance service fn",result.Error.Error())
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": balance}, fiber.StatusOK
}
