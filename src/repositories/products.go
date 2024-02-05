package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func SaveToCart(product models.CartItem) error {

	result := configs.DB.Create(&product)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return nil
}

func GetAllCartProductsByDistribID(distrib_id string, productsOut []dto.ProductsOut) ([]dto.ProductsOut, *gorm.DB) {
	result := configs.DB.Table("cart_items").Preload("Product").Select("product_id", "quantity").Where("distrib_id= ?", distrib_id).Find(&productsOut)
	return productsOut, result
}

func DeleteCartProduct(distrib_id string, product_id string, cartItem models.CartItem) (models.CartItem, *gorm.DB) {
	result := configs.DB.Clauses(clause.Returning{}).Unscoped().Where("distrib_id= ? AND product_id=?", distrib_id, product_id).Delete(&cartItem)
	return cartItem, result
}

func EditCartProducts(distrib_id string, product_id string, payload models.CartItem) (models.CartItem, *gorm.DB) {
	result := configs.DB.Model(models.CartItem{}).Where("distrib_id= ? AND product_id=?", distrib_id, product_id).Updates(payload)
	return payload, result
}
