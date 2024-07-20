package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func TcRouter(router fiber.Router) {
	router.Post("/addTc", controllers.AddTc)
	router.Get("/availableTc", controllers.AvailableTc)
}
