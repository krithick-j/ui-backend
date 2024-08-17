package service

import (
	"ui-back-end/configs"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetEpBalance(distribId string, tx *gorm.DB) (fiber.Map, int) {
	var balance float64

	balance, err := repositories.GetEpBalance(distribId, tx)
	if err == gorm.ErrRecordNotFound {
		configs.Log.Infoln("Record not found")
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetEpBalance", "GetEpBalance", fiber.StatusInternalServerError, tx)
	}
	return fiber.Map{"data": balance}, fiber.StatusOK
}
