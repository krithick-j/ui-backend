package service

import (
	"github.com/gofiber/fiber/v2"
)

// func Test(distribId string) (fiber.Map, int) {

// 	configs.Log.Infof("Test service completed")
// 	return fiber.Map{"data": ""}, 200
// }

func Test(distribId string) (fiber.Map, int) {

	return fiber.Map{"data": "rspSum"}, fiber.StatusOK
}
