package service

import (
	"ui-back-end/configs"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
)

func Test() (fiber.Map, int) {
	tx := configs.DB.Begin()
	invoicePdfPath, status, err := GenerateDistributorForm("IN-00001", tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GenerateInvoice", "Test", status, tx)
	}
	tx.Commit()
	return fiber.Map{"data": invoicePdfPath}, fiber.StatusOK
}
