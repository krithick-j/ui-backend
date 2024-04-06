package service

import (
	"time"
	"ui-back-end/src/dto"
	"ui-back-end/src/middleware"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func IsCheckqueAvailable(chequeDetailsIn dto.ChequeAvailableIn) (fiber.Map, int) {

	var chequeCounter, err = repositories.GetCheckoutFrequency(chequeDetailsIn.DistribId)
	if err.Error != nil {
		return fiber.Map{"error": err}, 500
	}

	const checkoutValue int = 4000

	if chequeCounter > 5 {
		return fiber.Map{"error": "Maximum checkout limit reached"}, 403
	}

	leftSumValue := chequeDetailsIn.LTc.BvPoint + chequeDetailsIn.LTc.LPoint + chequeDetailsIn.LTc.RPoint
	rightSumValue := chequeDetailsIn.RTc.BvPoint + chequeDetailsIn.RTc.LPoint + chequeDetailsIn.RTc.RPoint

	arr := []dto.Tc{chequeDetailsIn.LTc, chequeDetailsIn.RTc}

	if (leftSumValue >= checkoutValue) && (rightSumValue >= checkoutValue) {

		err := repositories.IncrementCheckoutFrequency(chequeDetailsIn.DistribId, chequeDetailsIn.Place, chequeCounter)

		if err.Error != nil {
			return fiber.Map{"error": err.Error}, 500
		}

		checkoutId := GenerateUniqueHexCode(10)

		for _, user := range arr {

			tx := models.BvTransaction{
				DisribId: user.DistribId,
				Place:    user.Place,
				OrderId:  checkoutId, //checkout id is saved in OrderID for now temporarily
				Date:     time.Now(),
				BvValue:  -checkoutValue,
			}
			repositories.SaveBvTransaction(tx)
		}

		return fiber.Map{"success": "checkout done"}, 200
	} else {
		return fiber.Map{"error": "cannot checkout"}, 400
	}
}

func TotalChequeValueByDistribId(TakeChequeIn dto.CheckoutIn) (dto.TakeChequeOut, dto.CheckoutFrequency, int) {

	CHECKOUT_VALUE := 4000
	var parentTCCheckoutFrequency, leftTCCheckoutFrequency, rightTCCheckoutFrequency int
	rank, _ := repositories.GetRankValueByDistribId(TakeChequeIn.DistribId)
	COUNT := 2 //Left and Right inside the tracking center

	types := []struct {
		Place             string
		CheckoutFrequency *int
	}{
		{"001", &parentTCCheckoutFrequency},
		{"002", &leftTCCheckoutFrequency},
		{"003", &rightTCCheckoutFrequency},
	}

	totalCheckoutFrequency := 0
	totalPoints := float32(0)
	placePointsObj := make([]dto.PlacePointsArr, len(types))

	for i, t := range types {
		var leftPoint, rightPoint int

		tcbv, _ := repositories.GetBVforTC(TakeChequeIn.DistribId, t.Place)
		for _, val := range tcbv {
			if val.Side == "left" {
				leftPoint = val.BValue
			}
			if val.Side == "right" {
				rightPoint = val.BValue
			}
		}
		*t.CheckoutFrequency = middleware.NCheckoutPossible(leftPoint, rightPoint, CHECKOUT_VALUE)

		checkoutFrequency := *t.CheckoutFrequency
		totalCheckoutFrequency += checkoutFrequency
		points := float32(checkoutFrequency*CHECKOUT_VALUE*COUNT) * rank

		placePointsObj[i] = dto.PlacePointsArr{
			Place: t.Place,
			Value: points,
		}
		totalPoints += points
	}

	pointsObj := dto.TakeChequeOut{
		TotalBalance:          totalPoints,
		TotalAvailableBalance: totalPoints,
		PlacePointsArr:        placePointsObj,
	}

	checkoutFrequencyObj := dto.CheckoutFrequency{
		TotalCheckoutFrequency:  totalCheckoutFrequency,
		ParentCheckoutFrequency: parentTCCheckoutFrequency,
		LeftCheckoutFrequency:   leftTCCheckoutFrequency,
		RightCheckoutFrequency:  rightTCCheckoutFrequency,
	}

	return pointsObj, checkoutFrequencyObj, 200
}

func TakeChequeByDistribIdAndPlace(TakeChequeIn dto.TakeChequeIn) (fiber.Map, int) {

	checkoutId := GenerateUniqueHexCode(10)

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

	if leftPoint >= 4000 && rightPoint >= 4000 {
		LeftInsideTcObj := models.BvTransaction{
			DisribId:     TakeChequeIn.DistribId,
			Place:        TakeChequeIn.Place,
			OrderId:      checkoutId,
			Date:         time.Now(),
			BvValue:      -4000,
			ActivateDate: time.Now().AddDate(0, 0, 7),
			Side:         "left",
		}
		RightInsidetCObj := models.BvTransaction{
			DisribId:     TakeChequeIn.DistribId,
			Place:        TakeChequeIn.Place,
			OrderId:      checkoutId,
			Date:         time.Now(),
			BvValue:      -4000,
			ActivateDate: time.Now().AddDate(0, 0, 7),
			Side:         "right",
		}
		repositories.SaveBvTransaction(LeftInsideTcObj)
		repositories.SaveBvTransaction(RightInsidetCObj)
		return fiber.Map{"data": "Cheque taken Successful"}, 200
	} else {
		return fiber.Map{"data": "cheque unsuccessfull!"}, 403
	}
}
