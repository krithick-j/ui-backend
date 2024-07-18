package service

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func Test(distribId string) (fiber.Map, int) {

	// get bv sum both active and not active
	bvSum, err := repositories.GetBvSumByDistribId(distribId, "product")
	if err != nil {
		configs.Log.Errorln("Error on calling GetBvSumByDistribId repositories fn from AddTc fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//get tc active and not active length
	tcArr, err := repositories.GetAllTrackingCenters(distribId)
	if err != nil {
		configs.Log.Errorln("Error on calling GetAllTrackingCenters repositories fn from AddTc fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//bvSum must be greater than tracking center length to add tracking center. Writing inversely
	reduceNumber := func(num int) int {
		str := fmt.Sprintf("%d", num)

		if num < 1000 {
			return 0
		} else if num > 1000 && num < 10000 {
			return int(str[0] - '0') // int('7' - '0')  '7' is 55 in ASCII, '0' is 48, so 55 - 48 = 7
		} else if num > 10000 && num < 100000 {
			return int(str[0]-'0')*10 + int(str[1]-'0') // int('7' - '0')  '7' is 55 in ASCII, '0' is 48, so 55 - 48 = 7
		} else {
			return 0
		}
	}

	bvSumFinal := reduceNumber(int(bvSum))

	// Get the next tracking center number
	tcLen := len(tcArr)

	if tcLen >= bvSumFinal {
		return fiber.Map{"data": "Tracking Center cannot be created"}, fiber.StatusForbidden
	} else if tcLen < bvSumFinal {

	}
	configs.Log.Infof("Test service completed")
	return fiber.Map{"data": bvSumFinal}, 200
}

func TestGetICouponArrayByOrderId(orderId string, distribId string) (fiber.Map, int) {

	ICouponArray, err := GetICouponArrayByOrderId(orderId, distribId)
	if err != nil {
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	configs.Log.Infof("Test service completed")
	return fiber.Map{"data": ICouponArray}, fiber.StatusOK
}
