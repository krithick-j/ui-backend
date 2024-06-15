package service

import (
	"crypto/sha256"
	"fmt"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/middleware"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Get Total Cheque Value By Distrib ID
func TotalChequeValueByDistribId(TakeChequeIn dto.CheckoutIn) (fiber.Map, int) {

	CHEQUE_DRAW_VALUE := configs.GlobalConfig.ChequeDrawValue //Cheque draw value is the constant 4000
	placePointsObj := []dto.PlacePointsArr{}
	COUNT := 2 //Left and Right inside the tracking center
	totalAvailableBvPoints := 0.0
	totalPoints := 0.0
	totalbvPoints := 0.0

	//Get Rank value
	rank, res := repositories.GetRankValueByDistribId(TakeChequeIn.DistribId)
	if res.Error != nil {
		return fiber.Map{"data": res.Error.Error()}, fiber.StatusInternalServerError
	}
	//Active Points start.....
	directCommissionActiveValue, res := repositories.GetDirectCommissionActiveValueByDistribId(TakeChequeIn.DistribId)
	if res.Error != nil {
		configs.Log.Errorln("Error on calling GetDirectCommissionValueByDistribId from TotalChequeValueByDistribId", res.Error.Error())
		return fiber.Map{"data": res.Error.Error()}, fiber.StatusInternalServerError
	}

	//Get Tracking Center with is_active = 1
	trackingCentersActiveValue, status := GetTrackingCentersByDistribId(TakeChequeIn.DistribId, true, "", "")
	if status != fiber.StatusOK {
		configs.Log.Errorln("Error on calling GetTrackingCentersByDistribId from TotalChequeValueByDistribId", trackingCentersActiveValue["error"])
		return fiber.Map{"data": "something went wrong in getting tracking centers", "status": status}, status
	}
	configs.Log.Infoln("Get tracking centers completed")

	//use type assertion to map the variables
	trackingCentersActiveArray, ok := trackingCentersActiveValue["data"].(dto.TrackingCenterOut)
	if !ok {
		return fiber.Map{"data": "something went wrong in type assertion tracking centers", "status": ok}, fiber.StatusInternalServerError
	}

	//Formula and addition for total Available points
	for _, tc := range trackingCentersActiveArray.Tc {

		ChequeFrequency := middleware.NCheckoutPossible(tc.LPoint, tc.RPoint, CHEQUE_DRAW_VALUE)
		fmt.Println("Cheque frequency: ", ChequeFrequency, "Check draw value: ", CHEQUE_DRAW_VALUE, "count: ", COUNT, "rank: ", rank)
		bvAvailblePoints := float64(ChequeFrequency*CHEQUE_DRAW_VALUE*COUNT) * rank
		fmt.Println("Bv Points", bvAvailblePoints, "direct commission", directCommissionActiveValue)
		if bvAvailblePoints != 0 {
			placePointsObj = append(placePointsObj, dto.PlacePointsArr{
				Place: tc.Place,
				Value: bvAvailblePoints,
			})
		}
		totalAvailableBvPoints += bvAvailblePoints
	}
	//Active points ends

	//Active + Not Active points starts.....

	//Get Tracking Center with is_active = 0
	//I already have the is_active total points in placePointsObj. I just need to add it with is_active = 0
	trackingCentersAllValue, status := GetTrackingCentersByDistribId(TakeChequeIn.DistribId, false, "", "")
	if status != fiber.StatusOK {
		configs.Log.Errorln("Error on calling GetTrackingCentersByDistribId from TotalChequeValueByDistribId", trackingCentersActiveValue["error"])
		return fiber.Map{"data": "something went wrong in getting tracking centers", "status": status}, status
	}
	configs.Log.Infoln("Get tracking centers completed")

	//use type assertion to map the variables
	trackingCentersAllArray, ok := trackingCentersAllValue["data"].(dto.TrackingCenterOut)
	if !ok {
		return fiber.Map{"data": "something went wrong in type assertion tracking centers", "status": ok}, fiber.StatusInternalServerError
	}

	//Formula and addition
	for _, tc := range trackingCentersAllArray.Tc {

		ChequeFrequency := middleware.NCheckoutPossible(tc.LPoint, tc.RPoint, CHEQUE_DRAW_VALUE)
		fmt.Println("Cheque frequency: ", ChequeFrequency, "Check draw value: ", CHEQUE_DRAW_VALUE, "count: ", COUNT, "rank: ", rank)
		bvAllPoints := float64(ChequeFrequency*CHEQUE_DRAW_VALUE*COUNT) * rank
		fmt.Println("Bv Points", bvAllPoints, "direct commission", directCommissionActiveValue)
		//Adding is_active = 0 values
		totalbvPoints += bvAllPoints
	}

	//Adding is_active = 1 values
	for _, tc := range placePointsObj {
		totalbvPoints += tc.Value
	}

	directComissionAllValue, res := repositories.GetDirectCommissionAllValueByDistribId(TakeChequeIn.DistribId)
	if res.Error != nil {
		configs.Log.Errorln("Error on calling GetDirectCommissionValueByDistribId from TotalChequeValueByDistribId", res.Error.Error())
		return fiber.Map{"data": res.Error.Error()}, fiber.StatusInternalServerError
	}

	//Total points is the points without is_active filter
	totalPoints = totalbvPoints + directComissionAllValue

	//totalAvailablePoints is only is_active
	totalAvailablePoints := totalAvailableBvPoints + directCommissionActiveValue

	//Dto Out
	pointsObj := dto.TakeChequeOut{
		TotalBalance:                     totalPoints,          //time after 22 days is excluded
		TotalAvailableBalance:            totalAvailablePoints, //available balance is balance which he can take
		BvBalance:                        totalAvailableBvPoints,
		PlacePointsArr:                   placePointsObj,
		DirectCommissionBalance:          directComissionAllValue,
		DirectCommissionAvailableBalance: directCommissionActiveValue,
	}

	return fiber.Map{"points_obj": pointsObj, "code": fiber.StatusOK}, fiber.StatusOK
}

func TakeChequeByDistribIdAndPlace(TakeChequeIn dto.TakeChequeIn, tx *gorm.DB) (fiber.Map, int) {

	checkoutId := GenerateUniqueHexCode(10)

	var checkDrawValue int = configs.GlobalConfig.ChequeDrawValue

	var leftPoint, rightPoint int
	side := TakeChequeIn.Place

	tcbv, _ := repositories.GetBVforTC(TakeChequeIn.DistribId, side)
	for _, val := range tcbv {
		if val.Side == "left" {
			leftPoint = val.BValue
		}
		if val.Side == "right" {
			rightPoint = val.BValue
		}
	}

	if leftPoint < checkDrawValue || rightPoint < checkDrawValue {
		tx.Rollback()
		return fiber.Map{"data": "cheque unsuccessfull!"}, fiber.StatusForbidden
	}

	LeftInsideTcObj := models.BvTransaction{
		DistribId:    TakeChequeIn.DistribId,
		Place:        TakeChequeIn.Place,
		OrderId:      checkoutId,
		Date:         time.Now(),
		BvValue:      -float64(checkDrawValue),
		ActivateDate: time.Now().AddDate(0, 0, 22),
		Side:         "left",
		TransType:    "cheque",
	}

	RightInsidetCObj := models.BvTransaction{
		DistribId:    TakeChequeIn.DistribId,
		Place:        TakeChequeIn.Place,
		OrderId:      checkoutId,
		Date:         time.Now(),
		BvValue:      -float64(checkDrawValue),
		ActivateDate: time.Now().AddDate(0, 0, 22),
		Side:         "right",
		TransType:    "cheque",
	}

	res := repositories.SaveBvTransaction(LeftInsideTcObj)
	if res.Error != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}
	res = repositories.SaveBvTransaction(RightInsidetCObj)
	if res.Error != nil {
		tx.Rollback()
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": "Cheque taken Successful"}, fiber.StatusOK
}

func ChangeChequePin(payload dto.ChequePinIn) (fiber.Map, int) {
	currentPinHashFromPayload := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.CurrentPin)))

	currentPinHashFromDB, res := repositories.GetChequePinByDistribID(payload.DistribId)
	if res.Error != nil {
		configs.Log.Errorln("Error on calling GetChequePinByDistribID repositories fn from ChangeChequePin service fn")
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}

	if currentPinHashFromDB != currentPinHashFromPayload {
		configs.Log.Infoln("Invalid CPA pin")
		return fiber.Map{"error": "Invalid Current Pin"}, fiber.StatusBadRequest
	}

	res = repositories.ChangeCpaPin(payload.DistribId, fmt.Sprintf("%x", sha256.Sum256([]byte(payload.NewPin))))
	if res.Error != nil {
		configs.Log.Errorln("Error on calling ChangeCpaPin repositories fn from ChangeChequePin service fn")
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": "Cpa pin changed Successfully"}, fiber.StatusOK
}

func ChequeLogin(payload dto.ChequeLogin) (fiber.Map, int) {
	PinHashFromPayload := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.Pin)))

	currentPinHashFromDB, res := repositories.GetChequePinByDistribID(payload.DistribId)
	if res.Error != nil {
		configs.Log.Errorln("Error on calling GetChequePinByDistribID repositories fn from ChequeLogin service fn", res.Error.Error())
		return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError

	}

	if currentPinHashFromDB != PinHashFromPayload {
		configs.Log.Infoln("Invalid CPA pin")
		return fiber.Map{"error": "Invalid Current Pin"}, fiber.StatusBadRequest
	}

	return fiber.Map{"data": "Cpa pin login successfull"}, fiber.StatusOK
}
