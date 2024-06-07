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

	var (
		parentTCCheckoutFrequency, leftTCCheckoutFrequency, rightTCCheckoutFrequency int
		CHEQUE_DRAW_VALUE                                                            = configs.GlobalConfig.ChequeDrawValue
	)

	rank, res := repositories.GetRankValueByDistribId(TakeChequeIn.DistribId)
	if res.Error != nil {
		return fiber.Map{"data": res.Error.Error()}, fiber.StatusInternalServerError
	}

	directCommissionValue, res := repositories.GetDirectCommissionValueByDistribId(TakeChequeIn.DistribId)
	if res.Error != nil {
		return fiber.Map{"data": res.Error.Error()}, fiber.StatusInternalServerError
	}

	COUNT := 2 //Left and Right inside the tracking center

	types := []struct {
		Place             string
		CheckoutFrequency *int
	}{
		{"001", &parentTCCheckoutFrequency},
		{"002", &leftTCCheckoutFrequency},
		{"003", &rightTCCheckoutFrequency},
	}

	var totalCheckoutFrequency int
	var totalPoints float64
	var totalBvPoints float64
	placePointsObj := []dto.PlacePointsArr{}

	for _, t := range types {
		var leftPoint, rightPoint int

		tcbv, res := repositories.GetBVforTC(TakeChequeIn.DistribId, t.Place)
		if res.Error != nil {
			return fiber.Map{"error": res.Error.Error()}, fiber.StatusInternalServerError
		}

		for _, val := range tcbv {
			if val.Side == "left" {
				leftPoint = val.BValue
			}
			if val.Side == "right" {
				rightPoint = val.BValue
			}
		}
		*t.CheckoutFrequency = middleware.NCheckoutPossible(leftPoint, rightPoint, CHEQUE_DRAW_VALUE)

		checkoutFrequency := *t.CheckoutFrequency
		totalCheckoutFrequency += checkoutFrequency

		points := float64(checkoutFrequency*CHEQUE_DRAW_VALUE*COUNT)*rank + directCommissionValue
		bvPoints := float64(checkoutFrequency*CHEQUE_DRAW_VALUE*COUNT) * rank
		if points != 0 {
			placePointsObj = append(placePointsObj, dto.PlacePointsArr{
				Place: t.Place,
				Value: points,
			})
		}

		totalPoints += points
		totalBvPoints += bvPoints
	}

	pointsObj := dto.TakeChequeOut{
		TotalBalance:                     totalPoints,
		TotalAvailableBalance:            totalPoints,
		BvBalance:                        totalBvPoints,
		PlacePointsArr:                   placePointsObj,
		DirectCommissionBalance:          directCommissionValue,
		DirectCommissionAvailableBalance: directCommissionValue,
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

	if leftPoint >= checkDrawValue && rightPoint >= checkDrawValue {
		LeftInsideTcObj := models.BvTransaction{
			DistribId:    TakeChequeIn.DistribId,
			Place:        TakeChequeIn.Place,
			OrderId:      checkoutId,
			Date:         time.Now(),
			BvValue:      -float64(checkDrawValue),
			ActivateDate: time.Now().AddDate(0, 0, 7),
			Side:         "left",
			TransType:    "cheque",
		}
		RightInsidetCObj := models.BvTransaction{
			DistribId:    TakeChequeIn.DistribId,
			Place:        TakeChequeIn.Place,
			OrderId:      checkoutId,
			Date:         time.Now(),
			BvValue:      -float64(checkDrawValue),
			ActivateDate: time.Now().AddDate(0, 0, 7),
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
		return fiber.Map{"data": "Cheque taken Successful"}, 200
	} else {
		tx.Rollback()
		return fiber.Map{"data": "cheque unsuccessfull!"}, 403
	}
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
