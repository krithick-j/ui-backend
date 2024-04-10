package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func HistoryRouter(router fiber.Router) {

	router.Post("/ICoupon", controllers.GetICouponHistory)
	router.Post("/Bv", controllers.GetBvHistory)
}
