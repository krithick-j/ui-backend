package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

// has many
func GetAllProducts(product []models.Product) ([]models.Product, *gorm.DB) {
	result := configs.DB.Model(&models.Product{}).Preload("ProductImage").Find(&product)
	return product, result
}

func GetAllProductCategories(productCategories []models.ProductCategory) ([]models.ProductCategory, *gorm.DB) {
	result := configs.DB.Find(&productCategories) //select * from product_category
	return productCategories, result
}

func GetAllProductByCategoryID(category_id string, product []models.Product) ([]models.Product, *gorm.DB) {
	result := configs.DB.Preload("ProductImage").Find(&product, "product_category_id = ?", category_id)
	return product, result
}

func GetAllProductByIDs(ids []uint, product []models.Product) ([]models.Product, *gorm.DB) {
	result := configs.DB.Find(&product, ids)
	return product, result
}

// func SaveToCart(products []models.Product) {
// 	for _, product := range products {
// 		product.AddToCart = true
// 		configs.DB.Save(&product)
// 	}
// }
