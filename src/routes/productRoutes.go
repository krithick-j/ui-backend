package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRouter(router fiber.Router) {
	router.Get("/", controllers.GetProductsController)
	router.Post("/getProductsById", controllers.GetProductsById)
	router.Post("/", controllers.CreateProduct)
	router.Post("/addToCart", controllers.AddToCartController)
	router.Get("/cartProducts/:distrib_id", controllers.GetCartProductsByUserId)
	router.Put("/cartProducts/", controllers.EditCartProducts)
	router.Delete("/cartProduct/", controllers.DeleteCartProduct)
	router.Delete("/allCartProducts/", controllers.DeleteAllCartProduct)
	router.Get("/categories", controllers.GetProductCategoriesController)
	router.Get("/filterByCategory/:category_id", controllers.GetProductsByCategoryID)
	router.Get("/orderDetails/", controllers.GetOrderDetails)
}
