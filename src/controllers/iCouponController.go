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
