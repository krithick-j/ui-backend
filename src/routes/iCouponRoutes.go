package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func ICouponRouter(router fiber.Router) {

	router.Post("/create/:admin_name", controllers.CreateICoupon)            // /admin --> path parameter syntax ("/:variable_name")
	router.Post("/validate/:distrib_id", controllers.ValidateICoupon) //validate/IN-00015
	router.Get("/:distrib_id", controllers.GetICouponsByDistribId)
	// router.Get("/:distrib_id", controllers.GetICouponsByDistribId)
}

// http://localhost:8080/api/iCoupon/IN-00015
