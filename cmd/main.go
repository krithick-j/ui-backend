package main

import (
	"encoding/json"
	"ui-back-end/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type User struct {
	Id            int
	Name          string
	Current_badge string
	Top_badge     string
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())
	app.Post("/", func(c *fiber.Ctx) error {
		var user = User{Id: 1, Name: "Gadson", Current_badge: "Silver", Top_badge: "Gold"}
		resp, _ := json.Marshal(user)
		return c.Status(202).SendString(string(resp))
	})

	// auth := app.Group("api/auth")
	// auth.Route("/user", routes.UserAuthRouter)

	user := app.Group("api/user")
	user.Route("/gr", routes.GrRouter)
	user.Route("/cpa", routes.CpaRouter)
	app.Listen(":8080")
}
