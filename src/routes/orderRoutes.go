package routes

import (
	"github.com/gofiber/fiber/v2"
)

func OrdersRouter(router fiber.Router) {

	router.Post("/:distrib_id", controllers.PlaceOrder)
}
