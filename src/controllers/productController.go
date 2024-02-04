package controllers

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/service"

	"github.com/gofiber/fiber/v2"
)

func GetProductsController(c *fiber.Ctx) error {
	res := service.GetProducts()
	return c.Status(http.StatusOK).JSON(res)
}

func GetProductCategoriesController(c *fiber.Ctx) error {
	res := service.GetProductCategories()
	return c.Status(http.StatusOK).JSON(res)
}

func GetProductsByCategoryID(c *fiber.Ctx) error {
	category_id := c.Params("category_id")
	res := service.GetProductByCategoryID(category_id)

	return c.Status(http.StatusOK).JSON(res)
}

func AddToCartController(c *fiber.Ctx) error {

	var request dto.CartItemIn
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res := service.AddToCart(request)

	return c.Status(http.StatusOK).JSON(res)
}

func GetProductsById(c *fiber.Ctx) error {
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if len(request.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No IDs provided"})
	}

	res, _ := service.GetProductsByIds(request.IDs)
	return c.Status(http.StatusOK).JSON(res)

}

func GetCartProductsByUserId(c *fiber.Ctx) error {
	user_id := c.Params("user_id")
	res := service.GetCartProductsByUserId(user_id)
	return c.Status(http.StatusOK).JSON(res)
}
