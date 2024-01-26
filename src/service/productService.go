package service

import (
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
)

func GetProducts() fiber.Map {
	print("jo")
	var products []models.Product
	var result error
	products, result = repositories.GetAllProducts(products)
	println(result)
	// if result.Error == gorm.ErrRecordNotFound {
	// 	return fiber.Map{"data": "Not Found"}
	// }
	// if result.Error != nil {
	// 	return fiber.Map{"error": result.Error}
	// }
	return fiber.Map{"data": products}
}
