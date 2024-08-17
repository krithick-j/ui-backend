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
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Get Total Cheque Value By Distrib ID
func TotalChequeValueByDistribId(TakeChequeIn dto.CheckoutIn, tx *gorm.DB) (fiber.Map, int) {

	CHEQUE_DRAW_VALUE := configs.GlobalConfig.ChequeDrawValue //Cheque draw value is the constant 4000
	placePointsObj := []dto.PlacePointsArr{}
	OneplacePointsObj := []dto.PlacePointsArr{}
	COUNT := 2 //Left and Right inside the tracking center
	totalAvailableBvPoints := 0.0
	totalPoints := 0.0
	totalbvPoints := 0.0

	//Get Rank value
	rank, err := repositories.GetRankValueByDistribId(TakeChequeIn.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.
			Errorln("Error on calling GetRankValueByDistribId from TotalChequeValueByDistribId", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//Active Points start.....
	directCommissionActiveValue, err := repositories.GetDirectCommissionActiveValueByDistribId(TakeChequeIn.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.
			Errorln("Error on calling GetDirectCommissionValueByDistribId from TotalChequeValueByDistribId", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	//Get Tracking Center with is_active = 1
	trackingCentersActiveValue, status := GetTrackingCentersByDistribId(TakeChequeIn.DistribId, true, "", "", tx)
	if status != fiber.StatusOK {
		return trackingCentersActiveValue, status
	}
	configs.Log.Infoln("Get tracking centers completed")

	//use type assertion to map the variables
	trackingCentersActiveArray, ok := trackingCentersActiveValue["data"].(dto.TrackingCenterOut)
	if !ok {
		tx.Rollback()
		return fiber.Map{"data": "something went wrong in type assertion tracking centers", "status": ok}, fiber.StatusInternalServerError
	}

	//Formula and addition for total Available points
	for _, tc := range trackingCentersActiveArray.Tc {

		ChequeFrequency := middleware.NCheckoutPossible(tc.LPoint, tc.RPoint, CHEQUE_DRAW_VALUE)
		fmt.Println("Cheque frequency: ", ChequeFrequency, "Check draw value: ", CHEQUE_DRAW_VALUE, "count: ", COUNT, "rank: ", rank)
		bvAvailblePoints := ChequeFrequency * float64(CHEQUE_DRAW_VALUE) * float64(COUNT) * rank
		if ChequeFrequency > 0 {
			OneFrequencyBvAvailablePoints := float64(CHEQUE_DRAW_VALUE) * float64(COUNT) * rank
			if OneFrequencyBvAvailablePoints != 0 {
				OneplacePointsObj = append(OneplacePointsObj, dto.PlacePointsArr{
					Place: tc.Place,
					Value: OneFrequencyBvAvailablePoints,
				})
			}
		}
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
	trackingCentersAllValue, status := GetTrackingCentersByDistribId(TakeChequeIn.DistribId, false, "", "", tx)
	if status != fiber.StatusOK {
		return trackingCentersActiveValue, status
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
		bvAllPoints := ChequeFrequency * float64(CHEQUE_DRAW_VALUE) * float64(COUNT) * rank
		fmt.Println("Bv Points", bvAllPoints, "direct commission", directCommissionActiveValue)
		//Adding is_active = 0 values
		totalbvPoints += bvAllPoints
	}

	//Adding is_active = 1 values
	for _, tc := range placePointsObj {
		totalbvPoints += tc.Value
	}

	directComissionAllValue, err := repositories.GetDirectCommissionAllValueByDistribId(TakeChequeIn.DistribId, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetDirectCommissionValueByDistribId", "TotalChequeValueByDistribId", fiber.StatusInternalServerError, tx)
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
		OneFrequencyPlacePointsArr:       OneplacePointsObj, //this is the array should be used in ICoupon generation page
		DirectCommissionBalance:          directComissionAllValue,
		DirectCommissionAvailableBalance: directCommissionActiveValue,
	}

	return fiber.Map{"points_obj": pointsObj, "code": fiber.StatusOK}, fiber.StatusOK
}

func TakeChequeByDistribIdAndPlace(TakeChequeIn dto.TakeChequeIn, tx *gorm.DB) (fiber.Map, int) {

	checkoutId := GenerateUniqueHexCode(10)
	CHEQUE_DRAW_VALUE := configs.GlobalConfig.ChequeDrawValue //Cheque draw value is the constant 4000
	COUNT := 2                                                //Left and Right inside the tracking center

	//Get Tracking Center with is_active = 1
	trackingCentersActiveValue, status := GetTrackingCentersByDistribId(TakeChequeIn.DistribId, true, "", "", tx)
	if status != fiber.StatusOK {
		return trackingCentersActiveValue, status
	}
	configs.Log.Infoln("Get tracking centers completed")

	for i := 1; i <= TakeChequeIn.ChequeCount; i++ {

		tcBv, err := repositories.GetBVforTCOneRow(TakeChequeIn.DistribId, TakeChequeIn.Place, tx)
		if err != nil {
			configs.Log.Errorln("Error on calling GetBVforTCOneRow from TakeChequeByDistribIdAndPlace", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}

		totalChequeCount := middleware.NCheckoutPossible(tcBv.LValue, tcBv.RValue, CHEQUE_DRAW_VALUE)
		if totalChequeCount < float64(TakeChequeIn.ChequeCount) {
			tx.Rollback()
			configs.Log.Errorln("Error from totalChequeCount validation. Cheque cannot be taken because of insufficient BV values")
			return fiber.Map{"data": "cheque unsuccessfull!"}, fiber.StatusForbidden
		}

		count, err := repositories.GetCheckoutFrequency(TakeChequeIn.DistribId, tx)
		if err != nil && err != gorm.ErrRecordNotFound {
			tx.Rollback()
			configs.Log.Errorln("Error on calling CreateCheckoutFrequency from TakeChequeByDistribIdAndPlace", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
		//This validation mostly won't be needed here because this is already done in get total cheque count frequency. Anyway it is used for extra validation
		if count%4 == 0 && err != gorm.ErrRecordNotFound && count != 0 {
			tx.Rollback()
			//take this cheque, this will not add in the cpa balance
		}

		if err == gorm.ErrRecordNotFound {
			if err = repositories.CreateCheckoutFrequency(TakeChequeIn.DistribId, TakeChequeIn.Place, tx); err != nil {
				tx.Rollback()
				configs.Log.Errorln("Error on calling CreateCheckoutFrequency from TakeChequeByDistribIdAndPlace", err.Error())
				return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
			}
		}

		if tcBv.LValue < CHEQUE_DRAW_VALUE || tcBv.RValue < CHEQUE_DRAW_VALUE {
			tx.Rollback()
			configs.Log.Errorln("Error from chequeDrawValue validation. Cheque cannot be taken because of insufficient BV values", err.Error())
			return fiber.Map{"data": "cheque unsuccessfull!"}, fiber.StatusForbidden
		}

		LeftInsideTcObj := models.BvTransaction{
			DistribId:    TakeChequeIn.DistribId,
			Place:        TakeChequeIn.Place,
			OrderId:      checkoutId,
			Date:         time.Now(),
			BvValue:      -float64(CHEQUE_DRAW_VALUE),
			ActivateDate: time.Now(),
			Side:         "left",
			TransType:    "cheque",
		}
		RightInsidetCObj := models.BvTransaction{
			DistribId:    TakeChequeIn.DistribId,
			Place:        TakeChequeIn.Place,
			OrderId:      checkoutId,
			Date:         time.Now(),
			BvValue:      -float64(CHEQUE_DRAW_VALUE),
			ActivateDate: time.Now(),
			Side:         "right",
			TransType:    "cheque",
		}

		err = repositories.SaveBvTransaction(LeftInsideTcObj, tx)
		if err != nil {
			tx.Rollback()
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}
		err = repositories.SaveBvTransaction(RightInsidetCObj, tx)
		if err != nil {
			tx.Rollback()
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}

		err = repositories.IncrementCheckoutFrequency(TakeChequeIn.DistribId, TakeChequeIn.Place, count, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling IncrementCheckoutFrequency from TakeChequeByDistribIdAndPlace service fn", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}

		err = handleCpaTransaction(TakeChequeIn, CHEQUE_DRAW_VALUE, COUNT, checkoutId, tx)
		if err != nil {
			tx.Rollback()
			configs.Log.Errorln("Error on calling handleCpaTransaction from TakeChequeByDistribIdAndPlace service fn", err.Error())
			return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
		}

	}
	return fiber.Map{"data": "Cheque taken Successful and added to CPA Balance"}, fiber.StatusOK
}

func handleCpaTransaction(TakeChequeIn dto.TakeChequeIn, CHEQUE_DRAW_VALUE int, COUNT int, checkoutId string, tx *gorm.DB) error {
	rank, err := repositories.GetRankValueByDistribId(TakeChequeIn.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetRankValueByDistribId from handleCpaTransaction service fn", err.Error())
		return err
	}

	CpaAmount := float64(TakeChequeIn.ChequeCount) * float64(CHEQUE_DRAW_VALUE) * float64(COUNT) * rank

	cpaTransactionObj := models.CpaTransaction{
		DistribId:    TakeChequeIn.DistribId,
		Reference:    checkoutId,
		ActivateDate: time.Now(),
		IsActive:     true,
		Amount:       CpaAmount,
	}
	err = repositories.SaveCpaTransaction(cpaTransactionObj, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.
			Errorln("Error on calling SaveCpaTransaction repositories fn from TakeChequeByDistribIdAndPlace service fn", err.Error())
		return err
	}
	return nil
}

func ChangeChequePin(payload dto.ChequePinIn, tx *gorm.DB) (fiber.Map, int) {
	currentPinHashFromPayload := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.CurrentPin)))

	currentPinHashFromDB, err := repositories.GetChequePinByDistribID(payload.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetChequePinByDistribID repositories fn from ChangeChequePin service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}

	if currentPinHashFromDB != currentPinHashFromPayload {
		tx.Rollback()
		configs.Log.Infoln("Invalid CPA pin")
		return fiber.Map{"error": "Invalid Current Pin"}, fiber.StatusBadRequest
	}

	err = repositories.ChangeCpaPin(payload.DistribId, fmt.Sprintf("%x", sha256.Sum256([]byte(payload.NewPin))), tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling ChangeCpaPin repositories fn from ChangeChequePin service fn")
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": "Cpa pin changed Successfully"}, fiber.StatusOK
}

func ChequeLogin(payload dto.ChequeLogin, tx *gorm.DB) (fiber.Map, int) {
	PinHashFromPayload := fmt.Sprintf("%x", sha256.Sum256([]byte(payload.Pin)))

	currentPinHashFromDB, err := repositories.GetChequePinByDistribID(payload.DistribId, tx)
	if err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on calling GetChequePinByDistribID repositories fn from ChequeLogin service fn", err.Error())
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError

	}

	if currentPinHashFromDB != PinHashFromPayload {
		tx.Rollback()
		configs.Log.Infoln("Invalid CPA pin")
		return fiber.Map{"error": "Invalid Current Pin"}, fiber.StatusBadRequest
	}

	return fiber.Map{"data": "Cpa pin login successfull"}, fiber.StatusOK
}
