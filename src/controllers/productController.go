package controllers

import (
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
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
	distrib_id := c.Params("distrib_id")
	res := service.GetCartProductsByDistribId(distrib_id)
	return c.Status(http.StatusOK).JSON(res)
}

func DeleteCartProduct(c *fiber.Ctx) error {
	user_id := c.Query("user_id")
	product_id := c.Query("product_id")
	res, status := service.DeleteCartProduct(user_id, product_id)
	return c.Status(status).JSON(res)
}

func EditCartProducts(c *fiber.Ctx) error {
	user_id := c.Query("user_id")
	product_id := c.Query("product_id")

	var payload models.CartItem
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res, status := service.EditCartProducts(payload, user_id, product_id)
	return c.Status(status).JSON(res)
}

func CreateProduct(c *fiber.Ctx) error {

	var payload dto.ProductIn
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	res, status := service.CreateProduct(payload)
	return c.Status(status).JSON(res)
}

func GetOrderDetails(c *fiber.Ctx) error {
	distrib_id := c.Query("distrib_id")
	res, status := service.GetOrderDetails(distrib_id)

	return c.Status(status).JSON(res)
}
