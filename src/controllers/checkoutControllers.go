package controllers

import (
	"net/http"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// func IsCheckqueAvailable(c *fiber.Ctx) error {

// 	var chequeAvailableIn dto.ChequeAvailableIn

// 	if err := c.BodyParser(&chequeAvailableIn); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
// 	}
// 	res, status := service.IsCheckqueAvailable(chequeAvailableIn)
// 	return c.Status(status).JSON(res)
// }

func TotalChequeValueByDistribId(c *fiber.Ctx) error {

	var CheckoutIn dto.CheckoutIn

	if err := c.BodyParser(&CheckoutIn); err != nil {
		configs.Log.Errorln("Error on parsing CheckoutIn from TotalChequeValueByDistribId controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.TotalChequeValueByDistribId(CheckoutIn)
	return c.Status(status).JSON(res)
}

func TakeChequeByDistribId(c *fiber.Ctx) error {

	var TakeChequeIn dto.TakeChequeIn

	if err := c.BodyParser(&TakeChequeIn); err != nil {
		configs.Log.Errorln("Error on calling TakeChequeIn from TakeChequeByDistribId controllers fn ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	tx := configs.DB.Begin()
	count, err := repositories.GetCheckoutFrequency(TakeChequeIn.DistribId)

	if err.Error == gorm.ErrRecordNotFound {
		repositories.CreateCheckoutFrequency(TakeChequeIn.DistribId, TakeChequeIn.Place)
	}

	if count%5 == 0 {
		tx.Rollback()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"data": "Maximum cheque limit reached!"})

	}

	res, status := service.TakeChequeByDistribIdAndPlace(TakeChequeIn, tx)

	if status == 200 {
		err := repositories.IncrementCheckoutFrequency(TakeChequeIn.DistribId, TakeChequeIn.Place, count)
		if err.Error != nil {
			return c.Status(http.StatusInternalServerError).JSON(err.Error.Error())
		}
		iCouponObj := dto.ICouponIn{
			DistribID: TakeChequeIn.DistribId,
			Coupons:   TakeChequeIn.Coupons,
		}

		service.AddICoupon(iCouponObj, iCouponObj.DistribID, tx) //admin is sent as empty string because user generating iCoupon
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
		}
		return c.Status(status).JSON(res)
	} else if status == 500 {
		return c.Status(fiber.StatusInternalServerError).JSON(res)

	} else {
		tx.Rollback()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"data": "Left and Right Points are insufficient"})
	}
}

func EditChequePin(c *fiber.Ctx) error {

	var payload dto.ChequePinIn

	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing payload from TotalChequeValueByDistribId controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.ChangeChequePin(payload)
	return c.Status(status).JSON(res)
}

func ChequeLogin(c *fiber.Ctx) error {

	var payload dto.ChequeLogin
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing payload from  controllers ChequeLogin fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.ChequeLogin(payload)
	return c.Status(status).JSON(res)
}
