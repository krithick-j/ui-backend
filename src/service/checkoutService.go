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

		err := repositories.IncrementCheckoutFrequency(chequeDetailsIn.DistribId, chequeCounter)

		if err.Error != nil {
			return fiber.Map{"error": err.Error}, 500
		}

		checkoutId := generateUniqueHexCode(10)

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

func TotalChequeValueByDistribId(TakeChequeIn dto.CheckoutIn) (fiber.Map, int) {

	var totalPoints float32 = 0.0
	CHECKOUT_VALUE := 4000
	rank, _ := repositories.GetRankValueByDistribId(TakeChequeIn.DistribId)
	COUNT := 2 //Left and Right inside the tracking center

	ruser := new(RecursiveUser)
	tc := repositories.GetTrackingCenter(TakeChequeIn.DistribId, "001")
	tcbv, _ := repositories.GetBVforTC(TakeChequeIn.DistribId, "001")

	//print ruser for better understanding
	ruser.Name = tc.Name
	ruser.TrackingCenter = tc.DistribID + " " + tc.Place
	ruser.IsActive = tc.IsActive
	for _, val := range tcbv {

		if val.Side == "left" {
			ruser.LeftPoint = val.BValue
		}
		if val.Side == "right" {
			ruser.RightPoint = val.BValue
		}
		if val.Side == "bv" {
			ruser.BV = val.BValue
		}
	}

	if tc.LeftDistribID == tc.DistribID {
		FindRecursiveTCOnlyDistribId(ruser, tc.LeftDistribID, tc.LeftPlace, "left")
	}
	if tc.RightDistribID == tc.DistribID {
		FindRecursiveTCOnlyDistribId(ruser, tc.RightDistribID, tc.RightPlace, "right")
	}

	parentTCCheckoutFrequency := middleware.NCheckoutPossible(ruser.LeftPoint, ruser.RightPoint, CHECKOUT_VALUE)
	leftTCCheckoutFrequency := middleware.NCheckoutPossible(ruser.Left.LeftPoint, ruser.Left.RightPoint, CHECKOUT_VALUE)
	rightTCCheckoutFrequency := middleware.NCheckoutPossible(ruser.Right.LeftPoint, ruser.Right.RightPoint, CHECKOUT_VALUE)

	totalCheckoutFrequency := parentTCCheckoutFrequency + leftTCCheckoutFrequency + rightTCCheckoutFrequency

	totalPoints = float32(totalCheckoutFrequency*CHECKOUT_VALUE*COUNT) * rank

	return fiber.Map{"data": totalPoints}, 200
}

func TakeChequeByDistribIdAndPlace(TakeChequeIn dto.TakeChequeIn) (fiber.Map, int) {

	checkoutId := generateUniqueHexCode(10)

	LeftObj := models.BvTransaction{
		DisribId: TakeChequeIn.DistribId,
		Place:    TakeChequeIn.Place,
		OrderId:  checkoutId,
		Date:     time.Now(),
		BvValue:  -4000,
		Side:     "left",
	}

	RightObj := models.BvTransaction{
		DisribId:     TakeChequeIn.DistribId,
		Place:        TakeChequeIn.Place,
		OrderId:      checkoutId,
		Date:         time.Now(),
		BvValue:      -4000,
		ActivateDate: time.Now().AddDate(0, 0, 7),
		Side:         "left",
	}

	repositories.SaveBvTransaction(LeftObj)
	repositories.SaveBvTransaction(RightObj)

	//generate coupon code and send mail
	

	return fiber.Map{"data": "Check Taken successfully"}, 200
}
