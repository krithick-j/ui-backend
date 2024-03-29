package service

import (
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func IsChequeAvailable(chequeDetailsIn dto.Tc) (fiber.Map, int) {
	var tx = configs.DB.Begin()
	chequeCounter, err := repositories.GetChequeFrequency(chequeDetailsIn.DistribId)
	if err.Error != nil {
		return fiber.Map{"error": err}, fiber.StatusInternalServerError // 500 internal server error
	}

	const checkoutValue int = 4000

	if chequeCounter > 4 {
		return fiber.Map{"error": "Maximum checkout limit reached"}, fiber.StatusForbidden // 403
 	}

	// leftSumValue := chequeDetailsIn.LTc.BvPoint + chequeDetailsIn.LTc.LPoint + chequeDetailsIn.LTc.RPoint
	// rightSumValue := chequeDetailsIn.RTc.BvPoint + chequeDetailsIn.RTc.LPoint + chequeDetailsIn.RTc.RPoint

	//arr := []dto.Tc{chequeDetailsIn.LTc, chequeDetailsIn.RTc}

	if chequeDetailsIn.LPoint < checkoutValue || chequeDetailsIn.RPoint < checkoutValue {
		return fiber.Map{"error": "Insufficient Bv values"}, fiber.StatusBadRequest // 400
	}
	err = repositories.IncrementCheckoutFrequency(chequeDetailsIn.DistribId, chequeCounter)

	if err.Error != nil {
		return fiber.Map{"error": err.Error}, fiber.StatusInternalServerError // 500 internal server error
	}

	checkoutID := generateUniqueHexCode(10)

	// for _, user := range arr {
	BVtx := models.BvTransaction{
		DisribId: chequeDetailsIn.DistribId,
		Place:    chequeDetailsIn.Place,
		OrderId:  checkoutID, //checkout id is saved in OrderID for now temporarily
		Date:     time.Now(),
		BvValue:  -checkoutValue,
	}
	repositories.UpdateLeftPointPlaceBv(chequeDetailsIn.DistribId,dto.PlaceBv{Place:chequeDetailsIn.Place , AddBv: -checkoutValue}, tx)
	repositories.UpdateLeftPointPlaceBv(chequeDetailsIn.DistribId,dto.PlaceBv{Place: chequeDetailsIn.Place, AddBv: -checkoutValue},tx)
	repositories.SaveBvTransaction(BVtx)
	// }
	return fiber.Map{"success": "checkout done"}, fiber.StatusOK // 200 success 
}
