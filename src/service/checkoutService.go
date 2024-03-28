package service

import (
	"time"
	"ui-back-end/src/dto"
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
