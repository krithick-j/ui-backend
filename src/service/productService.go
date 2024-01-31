package service

import (
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetProducts() fiber.Map {
	var products []models.Product
	var result *gorm.DB
	products, result = repositories.GetAllProducts(products)
	println(result)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": products}
}

func GetProductCategories() fiber.Map {
	var productCategories []models.ProductCategory
	var result *gorm.DB
	productCategories, result = repositories.GetAllProductCategories(productCategories)
	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": productCategories}
}

func GetProductByCategoryID(category_id string) fiber.Map {

	var product []models.Product
	var result *gorm.DB

	product, result = repositories.GetAllProductByCategoryID(category_id, product)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}
	return fiber.Map{"data": product}
}
