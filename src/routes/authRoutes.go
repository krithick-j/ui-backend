package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(router fiber.Router) {
	router.Post("/login", controllers.Login)
	router.Post("/register", controllers.UserRegistration)
	router.Post("/emailcode", controllers.SendEmailCode)
	router.Put("/emailcode", controllers.VerifyEmailCode)
	router.Post("/phonecode", controllers.SendPhoneCode)
	router.Put("/phonecode", controllers.VerifyPhoneCode)
}
