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
	result := configs.DB.Model(&models.Product{}).Preload("ProductImage").Preload("ProductImages").Find(&product)
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

func GetAllEpProductByCategoryID(category_id string, product []models.Product) ([]models.Product, *gorm.DB) {
	result := configs.DB.Preload("ProductImage").Find(&product, "product_category_id = ? and ep IS NOT NULL", category_id)
	return product, result
}

// Get the first cart item using distrib ID in cart table
func GetFirstCartItem(distribId string) (models.CartItem, *gorm.DB) {
	var item models.CartItem
	fmt.Println("distrib id", distribId)
	result := configs.DB.Model(&models.CartItem{}).Where("distrib_id=?", distribId).First(&item)
	return item, result
}

func GetProductTypeByProductID(product_id uint) (string, *gorm.DB) {
	var prodType string
	result := configs.DB.Model(&models.Product{}).Select("product_type").First(&prodType, "id=?", product_id)
	return prodType, result
}

func GetAllProductByIDs(ids []uint) ([]models.Product, *gorm.DB) {
	var product []models.Product
	result := configs.DB.Find(&product, ids)
	return product, result
}

func SaveToCart(product models.CartItem) error {
	println("hello from save to cart")
	result := configs.DB.Table("cart_items").Create(&product)
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

func DeleteAllCartProduct(distrib_id string, cartItem []models.CartItem) ([]models.CartItem, *gorm.DB) {
	result := configs.DB.Clauses(clause.Returning{}).Unscoped().Where("distrib_id= ?", distrib_id).Delete(&cartItem)
	return cartItem, result
}

func EditCartProducts(distrib_id string, product_id string, payload models.CartItem) (models.CartItem, *gorm.DB) {
	result := configs.DB.Model(models.CartItem{}).Where("distrib_id=? AND product_id=?", distrib_id, product_id).Updates(payload)
	return payload, result
}

func SaveProductImage(productImage *models.ProductImage) error {

	result := configs.DB.Create(&productImage)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return nil
}

func SaveProduct(product *models.Product) *models.Product {
	result := configs.DB.Create(&product)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return product
}
