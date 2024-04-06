package controllers

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/repositories"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func IsCheckqueAvailable(c *fiber.Ctx) error {

	var chequeAvailableIn dto.ChequeAvailableIn

	if err := c.BodyParser(&chequeAvailableIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.IsCheckqueAvailable(chequeAvailableIn)
	return c.Status(status).JSON(res)
}

func TotalChequeValueByDistribId(c *fiber.Ctx) error {

	var CheckoutIn dto.CheckoutIn

	if err := c.BodyParser(&CheckoutIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.TotalChequeValueByDistribId(CheckoutIn)
	return c.Status(status).JSON(fiber.Map{"data": res})
}

func TakeChequeByDistribId(c *fiber.Ctx) error {

	var TakeChequeIn dto.TakeChequeIn

	if err := c.BodyParser(&TakeChequeIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	count, _ := repositories.GetCheckoutFrequency(TakeChequeIn.DistribId)

	if count < 5 {
		res, status := service.TakeChequeByDistribIdAndPlace(TakeChequeIn)

		if status == 200 {
			err := repositories.IncrementCheckoutFrequency(TakeChequeIn.DistribId, TakeChequeIn.Place, count)
			if err.Error != nil {
				return c.Status(http.StatusInternalServerError).JSON(err)
			}

			iCouponObj := dto.ICouponIn{
				DistribID: TakeChequeIn.DistribId,
				TxDetail:  TakeChequeIn.ICouponTxDetail,
				Coupons:   TakeChequeIn.Coupons,
			}

			service.AddICoupon(iCouponObj, "") //admin is sent as empty string because user generating iCoupon
			return c.Status(status).JSON(res)
		} else {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"data": "Left and Right Points are insufficient"})
		}

	} else {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"data": "Maximum checkout limit reached!"})
	}
}
