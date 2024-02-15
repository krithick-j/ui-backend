package main

import (
	"fmt"
	"log"
	"ui-back-end/configs"
	"ui-back-end/src/controllers"
	"ui-back-end/src/routes"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/golang-jwt/jwt/v5"
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

	api := app.Group("/api")
	api.Post("/auth/login", controllers.Login)
	api.Post("/auth/register", controllers.UserRegistration)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	/* Hereafter all the endpoints will be secured */
	api.Route("/user", routes.UserRouter)
	api.Route("/cpa", routes.CpaRouter)
	api.Route("/iCoupon", routes.ICouponRouter)
	api.Route("/product", routes.ProductRouter)
	api.Route("/order", routes.OrdersRouter)
	api.Get("/allGrVisual", controllers.GetAllGrVisual)
	api.Get("/allGrVisualByDate", controllers.GetAllGrVisualByDate)
	api.Get("/carousalImages", controllers.CarousalImages)
	api.Get("/trackShipment", controllers.TrackShipment)
	api.Get("/newReferrals", controllers.NewReferrals)
	api.Get("/allReferrals", controllers.AllReferrals)

	app.Listen(":8080")
}
