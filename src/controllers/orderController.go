package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func PlaceOrder(c *fiber.Ctx) error {

	var OrderIn dto.PlaceOrderIn

	if err := c.BodyParser(&OrderIn); err != nil {
		configs.Log.Errorln("Error on parsing OrderIn from PlaceOrder controller fn", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	tx := configs.DB.Begin()
	res, status := service.PlaceOrder(OrderIn, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}

	return c.Status(status).JSON(res)
}

func GetOrdersByDistribId(c *fiber.Ctx) error {

	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetOrdersByDistribId(distrib_id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)

}

func GetAllOrders(c *fiber.Ctx) error {
	tx := configs.DB.Begin()
	res, status := service.GetAllOrders(tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)

}
