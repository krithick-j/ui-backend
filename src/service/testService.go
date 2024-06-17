package service

import (
	"net/http"
	"ui-back-end/configs"

	"github.com/gofiber/fiber/v2"
)

func Test(orderId string) (fiber.Map, int) {

	invoicePdfPath, err := GenerateInvoice(orderId)
	if err != nil {
		return fiber.Map{"error": err.Error()}, fiber.StatusInternalServerError

	}
	configs.Log.Infof("Test service completed")
	return fiber.Map{"data": invoicePdfPath}, http.StatusOK
}
