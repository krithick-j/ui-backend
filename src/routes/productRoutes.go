package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRouter(router fiber.Router) {
	router.Get("/", controllers.GetProductsController)
	router.Get("/categories", controllers.GetProductCategoriesController)
	router.Get("/filterByCategory/:category_id", controllers.GetProductByCategoryID)
}
