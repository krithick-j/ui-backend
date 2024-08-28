package utils

import (
	"fmt"
	"ui-back-end/configs"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NotNilErrorMessage(err error, method string, serviceMethod string, status int, tx *gorm.DB) (fiber.Map, int) {
	tx.Rollback()
	configs.Log.Errorln(fmt.Sprintf("Error on calling %s from %s service fn: %s", method, serviceMethod, err.Error()))
	return fiber.Map{"error": err.Error(), "err": err}, status
}

func CommonErrorMessage(err error, Message string, status int, tx *gorm.DB) (fiber.Map, int) {
	tx.Rollback()
	configs.Log.Errorln(Message, err.Error())
	return fiber.Map{"error": err.Error(), "err": err}, status
}

func RecordNotFoundMessage(err error, tx *gorm.DB) (fiber.Map, int) {
	tx.Rollback() //To ignore this statement, give tx as nil
	configs.Log.Warnln("Record Not Found: ", err.Error())
	return fiber.Map{"error": err.Error(), "err": err}, fiber.StatusNotFound
}

func CommonMessage(Message string, status int, tx *gorm.DB) (fiber.Map, int) {
	tx.Rollback()
	configs.Log.Errorln(Message)
	return fiber.Map{"error": Message}, status
}

func SuccessMessage(data any, status int) (fiber.Map, int) {
	return fiber.Map{"data": data}, status
}
