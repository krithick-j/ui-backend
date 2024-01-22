package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Post("/", func(c *fiber.Ctx) error {
		var user = User{Id: 1, Name: "Gadson", Current_badge: "Silver", Top_badge: "Gold"}
		resp, _ := json.Marshal(user)
		return c.Status(202).SendString(string(resp))
	})

	app.Listen(":8080")
}
