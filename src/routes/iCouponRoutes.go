package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func ICouponRouter(router fiber.Router) {

	router.Post("/:admin_name", controllers.CreateICoupon)
	router.Post("/validate/:distrib_id", controllers.ValidateICoupon)
	// router.Get("/:distrib_id", controllers.GetICouponsByDistribId)
}
