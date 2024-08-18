package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func CheckoutRouter(router fiber.Router) {

	// router.Post("/available", controllers.IsCheckqueAvailable)
	router.Post("/chequeCount", controllers.TotalChequeValueByDistribId)
	router.Post("/checkFrequencyAmount", controllers.AddFrequencyAmount)
	router.Post("/cpaValues", controllers.GetValuesForCpa)
	router.Post("/takeCheque", controllers.TakeChequeByDistribId)
	router.Put("/changePin", controllers.EditChequePin)
	router.Post("/login", controllers.ChequeLogin)
}
