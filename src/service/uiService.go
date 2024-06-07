package service

import (
	"fmt"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func BvCounterByStartDate(payload dto.BvCounterIn) (fiber.Map, int) {

	startDate, err := time.Parse("02-01-2006", payload.StartDate)
	if err != nil {
		configs.Log.Errorln("Error parsing date:", err.Error())
		return fiber.Map{"error": "Error parsing date ", "err": err.Error()}, fiber.StatusInternalServerError
	}
	endDate := startDate.AddDate(0, 0, 7)
	fmt.Println("start date", startDate)
	fmt.Println("end date", endDate)
	//logic to be written here
	lRCommissionBv, res := repositories.GetLRCommissionBvByDate(payload.DistribId, payload.Place, startDate, endDate)
	if res.Error != nil {
		configs.Log.Errorln("Error on calling GetLRCommissionBvByDate repositories fn from BvCounterByStartDate service fn", res.Error.Error())
		return fiber.Map{"error": "Error on calling GetLRCommissionBvByDate repositories fn from BvCounterByStartDate service fn", "err": err.Error()}, fiber.StatusInternalServerError
	}

	rankBvObj := dto.LeftAndRight{
		Left:  12, //actual value to be changed
		Right: 34,
	}

	commissionBvObj := dto.LeftAndRight{
		Left:  float64(lRCommissionBv.LValue),
		Right: float64(lRCommissionBv.RValue),
	}

	BvCounterOut := dto.BvCounterOut{
		RankBv:       rankBvObj,
		CommissionBv: commissionBvObj,
	}
	return fiber.Map{"data": BvCounterOut}, fiber.StatusOK
}
