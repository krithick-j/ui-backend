package controllers

import (
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetPersonalRspByDistribID(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetPersonalRspByDistribId(distrib_id, tx)

	return c.Status(status).JSON(res)
}

func GetGroupRspByDistribID(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetGroupRspByDistribId(distrib_id, tx)

	if err := tx.Commit().Error; err != nil {
		configs.Log.Errorln("Error on committing transaction:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(status).JSON(res)
}

func GetGroupPerformanceByDistribId(c *fiber.Ctx) error {
	id := c.Params("distrib_id")
	tx := configs.DB.Begin()
	res, status := service.GetGroupPerformanceByDistribId(id, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		configs.Log.Errorln("Error on Committing tx from GetGroupPerformanceByDistribId controller", err.Error())
	}

	return c.Status(status).JSON(res)
}

func GetDirectBvByDistribID(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetDirectBvByDistribId(distrib_id, tx)

	if err := tx.Commit().Error; err != nil {
		configs.Log.Errorln("Error on committing transaction:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(status).JSON(res)
}

func GetTotalStepsByDistribId(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetTotalStepByDistribId(distrib_id, tx)

	if err := tx.Commit().Error; err != nil {
		configs.Log.Errorln("Error on committing transaction:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(status).JSON(res)
}

func GetRspValuesByDistribID(c *fiber.Ctx) error {
	distrib_id := c.Params("distrib_id")

	tx := configs.DB.Begin()
	res, status := service.GetRspValuesByDistribID(distrib_id, tx)

	if err := tx.Commit().Error; err != nil {
		configs.Log.Errorln("Error on committing transaction:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(status).JSON(res)
}

// Using the CPA value, ICoupon is generated
func SaveCpaICoupon(c *fiber.Ctx) error {
	var payload dto.TakeCpaAmount
	if err := c.BodyParser(&payload); err != nil {
		configs.Log.Errorln("Error on parsing request from AddToCartController controller fn ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	tx := configs.DB.Begin()
	res, status := service.SaveCpaICoupon(payload, tx)
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
	}
	return c.Status(status).JSON(res)
}
