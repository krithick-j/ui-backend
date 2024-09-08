package routes

import (
	"ui-back-end/src/controllers"

	"github.com/gofiber/fiber/v2"
)

func StaticRouter(router fiber.Router) {
	router.Get("/socialConn", controllers.SocialConnController)
}
