package service

import (
	"ui-back-end/configs"

	"github.com/gofiber/fiber/v2"
)

// func Test(payload dto.AddTc) (fiber.Map, int) {

// 	status, err := AddTc(payload)
// 	if err != nil {
// 		return fiber.Map{"error": err.Error()}, status
// 	}
// 	configs.Log.Infof("Test service completed")
// 	return fiber.Map{"data": "successfully completed"}, status
// }

func TestGetICouponArrayByOrderId(orderId string, distribId string) (fiber.Map, int) {

	ICouponArray, err := GetICouponArrayByOrderId(orderId, distribId)
	if err != nil {
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	configs.Log.Infof("Test service completed")
	return fiber.Map{"data": ICouponArray}, fiber.StatusOK
}
