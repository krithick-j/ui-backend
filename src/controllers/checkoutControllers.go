package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
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
	tx := configs.DB.Begin()
	res, status := service.TotalChequeValueByDistribId(CheckoutIn, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func AddFrequencyAmount(c *fiber.Ctx) error {

	var payload dto.FrequencyForTc

	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing CheckoutIn from TotalChequeValueByDistribId controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.GetFrequencyAmount(payload, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GetChequeCountByDistribId(c *fiber.Ctx) error {

	var CheckoutIn dto.CheckoutIn

	if err := c.BodyParser(&CheckoutIn); err != nil {
		configs.Log.Errorln("Error on parsing CheckoutIn from GetChequeCountByDistribId controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.TotalChequeValueByDistribId(CheckoutIn, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GetValuesForCpa(c *fiber.Ctx) error {

	distribId := c.Params("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetValuesForCpa(distribId, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func TakeChequeByDistribId(c *fiber.Ctx) error {
	var TakeChequeIn dto.TakeChequeIn

	if err := c.BodyParser(&TakeChequeIn); err != nil {
		configs.Log.Errorln("Error on calling TakeChequeIn from TakeChequeByDistribId controllers fn ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	tx := configs.DB.Begin()
	res, status := service.TakeChequeByDistribIdAndPlace(TakeChequeIn, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func EditChequePin(c *fiber.Ctx) error {

	var payload dto.ChequePinIn

	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing payload from TotalChequeValueByDistribId controllers fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.ChangeChequePin(payload, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}

func ChequeLogin(c *fiber.Ctx) error {

	var payload dto.ChequeLogin
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing payload from  controllers ChequeLogin fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	tx := configs.DB.Begin()
	res, status := service.ChequeLogin(payload, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}
