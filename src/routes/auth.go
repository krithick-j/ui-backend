package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func NoAuthRouter(router fiber.Router) {
	router.Post("/", controllers.UserRegistration)
	router.Post("/login", controllers.Login)
}
