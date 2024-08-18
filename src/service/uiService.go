package service

import (
	"fmt"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BvCounterByStartDate(payload dto.BvCounterIn, tx *gorm.DB) (fiber.Map, int) {

	startDate, err := time.Parse("02-01-2006", payload.StartDate)
	if err != nil {
		configs.Log.Errorln("Error parsing date:", err.Error())
		return fiber.Map{"error": "Error parsing date ", "err": err.Error()}, fiber.StatusInternalServerError
	}
	endDate := startDate.AddDate(0, 0, 7)

	lRCommissionBv, err := repositories.GetBVforTCOneRowByDate(payload.DistribId, payload.Place, startDate, endDate, tx)
	if err != nil {
		configs.Log.Errorln("Error on calling GetLRCommissionBvByDate repositories fn from BvCounterByStartDate service fn", err.Error())
		return fiber.Map{"error": "Error on calling GetLRCommissionBvByDate repositories fn from BvCounterByStartDate service fn", "err": err.Error()}, fiber.StatusInternalServerError
	}

	fmt.Println("lrcommission", lRCommissionBv)
	commissionBvObj := dto.LeftAndRight{
		Left:  float64(lRCommissionBv.LValue),
		Right: float64(lRCommissionBv.RValue),
	}

	res, status := RankBvSumByDate(payload.DistribId, payload.Place, startDate, endDate, tx)
	if status != fiber.StatusOK {
		return res, status
	}

	leftRankBv := res["leftRankBv"].(float64)
	rightRankBv := res["rightRankBv"].(float64)

	rankBvObj := dto.LeftAndRight{
		Left:  leftRankBv,
		Right: rightRankBv,
	}

	BvCounterOut := dto.BvCounterOut{
		RankBv:       rankBvObj,
		CommissionBv: commissionBvObj,
	}

	return fiber.Map{"data": BvCounterOut}, fiber.StatusOK
}

// returns left and right sum rank
func RankBvSumByDate(distrib_id string, place string, fromDate time.Time, toDate time.Time, tx *gorm.DB) (fiber.Map, int) {
	var leftRankBv = 0.0
	var rightRankBv = 0.0

	tc, err := repositories.GetTrackingCenter(distrib_id, place, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetTrackingCenter", "RankBvSumByDate", fiber.StatusInternalServerError, tx)
	}
	tcbv, err := repositories.GetBVforTCByDate(distrib_id, place, fromDate, toDate, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetBVforTCByDate", "RankBvSumByDate", fiber.StatusInternalServerError, tx)
	}
	for _, val := range tcbv {
		if val.Side == "left" {
			leftRankBv += float64(val.BValue)
		}
		if val.Side == "right" {
			rightRankBv += float64(val.BValue)
		}
	}

	if tc.LeftDistribID != "" {
		FindRecursiveTCForRankBv(&leftRankBv, &rightRankBv, tc.LeftDistribID, tc.LeftPlace, "left", tx)
	}
	if tc.RightDistribID != "" {
		FindRecursiveTCForRankBv(&leftRankBv, &rightRankBv, tc.RightDistribID, tc.RightPlace, "right", tx)
	}

	return fiber.Map{"leftRankBv": float64(leftRankBv), "rightRankBv": float64(rightRankBv)}, fiber.StatusOK
}
