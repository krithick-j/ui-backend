package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func CheckoutRouter(router fiber.Router) {

	router.Post("/available", controllers.IsCheckqueAvailable)
}
