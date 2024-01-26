package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"
)

func GetAllProducts(product []models.Product) ([]models.Product, error) {
	result := configs.DB.Model(&models.Product{}).Preload("ProductImages").Find(&product).Error
	return product, result
}
