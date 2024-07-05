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

	lRCommissionBv, res := repositories.GetBVforTCOneRowByDate(payload.DistribId, payload.Place, startDate, endDate)
	if res.Error != nil {
		configs.Log.Errorln("Error on calling GetLRCommissionBvByDate repositories fn from BvCounterByStartDate service fn", res.Error.Error())
		return fiber.Map{"error": "Error on calling GetLRCommissionBvByDate repositories fn from BvCounterByStartDate service fn", "err": res.Error.Error()}, fiber.StatusInternalServerError
	}

	fmt.Println("lrcommission", lRCommissionBv)
	commissionBvObj := dto.LeftAndRight{
		Left:  float64(lRCommissionBv.LValue),
		Right: float64(lRCommissionBv.RValue),
	}

	leftRankBv, rightRankBv := RankBvSumByDate(payload.DistribId, payload.Place, startDate, endDate)
	rankBvObj := dto.LeftAndRight{
		Left:  leftRankBv,
		Right: rightRankBv,
	}

	BvCounterOut := dto.BvCounterOut{
		RankBv:       rankBvObj,
		CommissionBv: commissionBvObj,
	}

	fmt.Println("bv Counter---_>", BvCounterOut)
	return fiber.Map{"data": BvCounterOut}, fiber.StatusOK
}

// returns left and right sum rank
func RankBvSumByDate(distrib_id string, place string, fromDate time.Time, toDate time.Time) (float64, float64) {
	var leftRankBv = 0.0
	var rightRankBv = 0.0

	tc := repositories.GetTrackingCenter(distrib_id, place)
	fmt.Println("from date-->", fromDate, "to date-->", toDate)
	tcbv, _ := repositories.GetBVforTCByDate(distrib_id, place, fromDate, toDate)
	for _, val := range tcbv {
		if val.Side == "left" {
			leftRankBv += float64(val.BValue)
		}
		if val.Side == "right" {
			rightRankBv += float64(val.BValue)
		}
	}

	if tc.LeftDistribID != "" {
		FindRecursiveTCForRankBv(&leftRankBv, &rightRankBv, tc.LeftDistribID, tc.LeftPlace, "left")
	}
	if tc.RightDistribID != "" {
		FindRecursiveTCForRankBv(&leftRankBv, &rightRankBv, tc.RightDistribID, tc.RightPlace, "right")
	}

	return float64(leftRankBv), float64(rightRankBv)
}
