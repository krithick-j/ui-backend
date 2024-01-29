package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

// has many
func GetAllProducts(product []models.Product) ([]models.Product, *gorm.DB) {
	result := configs.DB.Model(&models.Product{}).Preload("ProductImages").Find(&product)
	return product, result
}

func GetAllProductCategories(productCategories []models.ProductCategory) ([]models.ProductCategory, *gorm.DB) {
	result := configs.DB.Find(&productCategories) //select * from product_category
	return productCategories, result
}
