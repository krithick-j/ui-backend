package service

import (
	"ui-back-end/configs"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
)

func Test() (fiber.Map, int) {
	tx := configs.DB.Begin()
	filename, status, err := GenerateInvoice("CEA9D2D712", "bv", tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GenerateInvoice", "Test", status, tx)
	}
	tx.Commit()
	return fiber.Map{"data": filename}, fiber.StatusOK
}
