package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func OrdersRouter(router fiber.Router) {

	router.Post("/", controllers.PlaceOrder)
	router.Get("/", controllers.GetAllOrders)
	router.Get("/:distrib_id", controllers.GetOrdersByDistribId)
	router.Get("/trackShipment", controllers.TrackShipment)
}
