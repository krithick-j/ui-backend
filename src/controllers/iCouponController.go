package controllers

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func CreateICoupon(c *fiber.Ctx) error {

	var iCouponIn dto.ICouponIn

	adminName := c.Params("admin_name")

	if err := c.BodyParser(&iCouponIn); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	res, status := service.AddICoupon(iCouponIn, adminName)
	return c.Status(status).JSON(res)
}

// func GetICouponsByDistribId(c *fiber.Ctx) error {

// 	distribID := c.Params("distrib_id")
// 	// res, status := service.GetAllICouponsByDistribID(distribID, distribID)
// 	return c.Status(status).JSON(res)
// }

func ValidateICoupon(c *fiber.Ctx) error {
	var payload dto.ValidateICouponIn

	distribID := c.Params("distrib_id")

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	res, status := service.ValidateICoupon(payload, distribID)
	return c.Status(status).JSON(res)
}
