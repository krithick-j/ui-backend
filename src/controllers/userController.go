package controllers

import (
	"net/http"
	"ui-back-end/src/models"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetUserByDistId(c *fiber.Ctx) error {
	id := c.Params("dist_id")

	res := service.GetUserByDistId(id)

	return c.Status(http.StatusOK).JSON(res)
}

func GetAllUsers(c *fiber.Ctx) error {
	res := service.GetUsers()
	return c.Status(http.StatusOK).JSON(res)
}

func GetUserTreeByDistId(c *fiber.Ctx) error {
	id := c.Params("dist_id")
	res := service.GetTreeUserByDistId(id)
	return c.Status(http.StatusOK).JSON(res)
}

func UserRegistration(c *fiber.Ctx) error {
	user := models.User{}
	c.BodyParser(&user)
	res := service.RegisterUser(user)
	return c.Status(http.StatusCreated).JSON(res)
}
