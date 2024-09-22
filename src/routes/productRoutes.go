package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductRouter(router fiber.Router) {
	router.Get("/", controllers.GetProductsController)
	router.Post("/getProductsById", controllers.GetProductsById)
	router.Get("/getEpProducts/:category_id", controllers.GetEpProductsByCategoryId)
	router.Post("/admin/:admin_name", controllers.CreateProduct)
	router.Put("/admin", controllers.EditProduct)
	router.Post("/addToCart", controllers.AddToCartController)
	router.Get("/cartProducts/:distrib_id", controllers.GetCartProductsByUserId)
	router.Put("/cartProducts/", controllers.EditCartProducts)
	router.Delete("/cartProduct/", controllers.DeleteCartProduct)
	router.Delete("/allCartProducts/", controllers.DeleteAllCartProduct)
	router.Get("/categories", controllers.GetProductCategoriesController)
	router.Get("/filterByCategory", controllers.GetProductsByCategoryID)
	router.Get("/orderDetails/", controllers.GetOrderDetails)
}
