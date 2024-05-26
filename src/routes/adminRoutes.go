package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func AdminRouter(router fiber.Router) {
	router.Post("/product", controllers.CreateProduct)
	router.Post("/icoupon/:admin_name", controllers.CreateICoupon)
	router.Put("/product", controllers.EditProduct)
}
