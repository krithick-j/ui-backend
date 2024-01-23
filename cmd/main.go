package main

import (
	"ui-back-end/src/controllers"
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
	// app.Post("/", func(c *fiber.Ctx) error {
	// 	var user = User{Id: 1, Name: "Gadson", Current_badge: "Silver", Top_badge: "Gold"}
	// 	resp, _ := json.Marshal(user)
	// 	return c.Status(202).SendString(string(resp))
	// })

	// auth := app.Group("api/auth")
	// auth.Route("/user", routes.UserAuthRouter)

	api := app.Group("/api")
	api.Get("/allGrVisual", controllers.GetAllGrVisual)
	api.Get("/allGrVisualByDate", controllers.GetAllGrVisualByDate)
	api.Route("/cpa", routes.CpaRouter)
	api.Get("/carousalImages", controllers.CarousalImages)
	api.Get("/trackShipment", controllers.TrackShipment)
	api.Get("/orderAndPayment", controllers.OrderAndPayment)
	api.Get("/newReferrals", controllers.NewReferrals)
	api.Get("/allReferrals", controllers.AllReferrals)

	app.Listen(":8080")
}
