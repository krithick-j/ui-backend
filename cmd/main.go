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
	configs.LoggerConfig()
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalln("\x1b[31mFailed to load environment variables!\x1b[0m \n", err.Error())
	}
	configs.DbConnect(config)
}

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Get("/media/:filename", controllers.GetMediaFile)
	api := app.Group("/api")
	api.Post("/auth/login", controllers.Login)
	api.Post("/auth/register", controllers.UserRegistration)
	api.Post("/auth/emailcode", controllers.SendEmailCode)
	api.Put("/auth/emailcode", controllers.VerifyEmailCode)
	api.Post("/auth/phonecode", controllers.SendPhoneCode)
	api.Put("/auth/phonecode", controllers.VerifyPhoneCode)
	api.Route("/user", routes.UserRouter)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	/* Hereafter all the endpoints will be secured */
	api.Route("/cpa", routes.CpaRouter)
	api.Route("/iCoupon", routes.ICouponRouter)
	api.Route("/product", routes.ProductRouter)
	api.Route("/order", routes.OrdersRouter)
	api.Route("/rsp", routes.RspRouter)
	api.Route("/history", routes.HistoryRouter)
	api.Route("/ui", routes.UiRouter)
	api.Route("/redeem", routes.RedeemRouter)
	api.Route("/cheque", routes.CheckoutRouter)
	api.Route("/contactCenter", routes.ContactCenter)
	api.Route("/admin", routes.AdminRouter)
	app.Listen(fmt.Sprintf(":%d", configs.GlobalConfig.AppPort))
}
