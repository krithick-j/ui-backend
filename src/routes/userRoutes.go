package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(router fiber.Router) {
	//User Crud
	router.Get("/", controllers.GetAllUsers)
	router.Get("/:dist_id", controllers.GetUserByDistId)
	router.Post("/", controllers.UserRegistration)

	//User Chain
	router.Get("/tree/:dist_id", controllers.GetUserTreeByDistId)

	//User Reports

}
