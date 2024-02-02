package main

import (
	"fmt"
	"log"
	"ui-back-end/configs"
	"ui-back-end/src/controllers"
	"ui-back-end/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	fmt.Println("Initializing the DB")
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalln("\x1b[31mFailed to load environment variables!\x1b[0m \n", err.Error())
	}
	configs.DbConnect(&config)
}

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())

	// auth := app.Group("/auth")
	// auth.Route("/user", routes.UserAuthRouter)

	api := app.Group("/api")
	api.Route("/cpa", routes.CpaRouter)
	api.Route("/user", routes.UserRouter)
	api.Route("/product", routes.ProductRouter)
	api.Get("/allGrVisual", controllers.GetAllGrVisual)
	api.Get("/allGrVisualByDate", controllers.GetAllGrVisualByDate)
	api.Get("/carousalImages", controllers.CarousalImages)
	api.Get("/trackShipment", controllers.TrackShipment)
	api.Get("/orderAndPayment", controllers.OrderAndPayment)
	api.Get("/newReferrals", controllers.NewReferrals)
	api.Get("/allReferrals", controllers.AllReferrals)

	app.Listen(":8080")
}
