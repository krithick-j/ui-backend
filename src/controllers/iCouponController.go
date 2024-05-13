package controllers

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func CreateICoupon(c *fiber.Ctx) error {

	var iCouponIn dto.ICouponIn

	adminName := c.Params("admin_name")
	fmt.Println(adminName)
	if err := c.BodyParser(&iCouponIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	fmt.Println(iCouponIn)
	tx := configs.DB.Begin()
	res, status := service.AddICoupon(iCouponIn, adminName, tx)
	// If no error occurred, commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())

	}
	return c.Status(status).JSON(res)
}

func GetICouponsByDistribId(c *fiber.Ctx) error {

	DistribId := c.Params("distrib_id") //searching distrib_id parameter exists in the URL
	res, status := service.GetAllICouponsByDistribId(DistribId)
	return c.Status(status).JSON(res)
}

func ValidateICoupon(c *fiber.Ctx) error {
	var payload dto.ValidateICouponIn

	distribID := c.Params("distrib_id")

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	res, status := service.ValidateICoupon(payload, distribID)
	return c.Status(status).JSON(res)
}
