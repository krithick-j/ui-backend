package repositories

import (
	"ui-back-end/src/dto"
	"ui-back-end/src/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// has many
func GetAllProducts(tx *gorm.DB) ([]models.Product, error) {
	var product []models.Product
	err :=
		tx.
			Model(&models.Product{}).
			Preload("ProductImages").
			Find(&product).
			Error
	return product, err
}

func GetProductById(productId uint, tx *gorm.DB) (models.Product, error) {
	var product models.Product
	err :=
		tx.
			Model(&models.Product{}).
			Find(&product, productId).
			Error
	return product, err
}

func GetAllProductCategories(productCategories []models.ProductCategory, tx *gorm.DB) ([]models.ProductCategory, error) {
	err := tx.Find(&productCategories).Error //select * from product_category
	return productCategories, err
}

func GetAllProductByCategoryID(category_id string, product []models.Product, tx *gorm.DB) ([]models.Product, error) {
	err :=
		tx.
			Find(&product, "product_category_id = ?", category_id).
			Error
	return product, err
}

func GetAllProductByCategoryIdAndProductType(category_id string, productType string, tx *gorm.DB) ([]models.Product, error) {
	var product []models.Product
	err :=
		tx.
			Find(&product, "product_category_id = ? AND product_type=?", category_id, productType).
			Error
	return product, err
}

func GetAllProductByCategoryIdAndProductTypeFilter(categoryIds []string, productTypes []string, tx *gorm.DB) ([]models.Product, error) {

	var products []models.Product
	if len(categoryIds) > 0 {
		tx = tx.Where("product_category_id IN (?)", categoryIds)
	}
	if len(productTypes) > 0 {
		tx = tx.Where("product_type IN (?)", productTypes)
	}

	// Execute the query
	err := tx.Preload("ProductImages").Find(&products).Error
	return products, err
}

func GetAllEpProductByCategoryID(category_id string, tx *gorm.DB) ([]models.Product, error) {
	var product []models.Product
	err :=
		tx.
			Preload("ProductImages").
			Find(&product, "product_category_id = ? and ep IS NOT NULL", category_id).
			Error
	return product, err
}

// Get the first cart item using distrib ID in cart table
func GetFirstCartItem(distribId string, tx *gorm.DB) (models.CartItem, error) {
	var item models.CartItem
	err :=
		tx.
			Model(&models.CartItem{}).
			Where("distrib_id=?", distribId).
			First(&item).
			Error
	return item, err
}

func GetProductTypeByProductID(product_id uint, tx *gorm.DB) (string, error) {
	var prodType string
	err :=
		tx.
			Model(&models.Product{}).
			Select("product_type").
			First(&prodType, "id=?", product_id).
			Error
	return prodType, err
}

func GetAllProductByIDs(ids []uint, tx *gorm.DB) ([]models.Product, error) {
	var product []models.Product
	err :=
		tx.
			Find(&product, ids).
			Error
	return product, err
}

func SaveToCart(product models.CartItem, tx *gorm.DB) error {
	err :=
		tx.
			Table("cart_items").
			Create(&product).
			Error
	return err
}

func GetAllCartProductsByDistribID(distrib_id string, tx *gorm.DB) ([]dto.ProductsOut, error) {
	var productsOut []dto.ProductsOut
	err :=
		tx.
			Table("cart_items").
			Preload("Product").
			Select("product_id", "quantity", "id").
			Where("distrib_id= ?", distrib_id).
			Find(&productsOut).
			Error
	return productsOut, err
}

func UpdateCartProductQuantityById(cartId uint, quantity uint, tx *gorm.DB) error {
	err :=
		tx.
			Model(models.CartItem{}).
			Where("id=?", cartId).
			Update("quantity", quantity).
			Error
	return err
}

func DeleteCartProduct(distrib_id string, product_id string, tx *gorm.DB) (models.CartItem, error) {
	var cartItem models.CartItem
	err :=
		tx.
			Clauses(clause.Returning{}).
			Unscoped().
			Where("distrib_id= ? AND product_id=?", distrib_id, product_id).
			Delete(&cartItem).
			Error
	return cartItem, err
}

func DeleteAllCartProduct(distrib_id string, tx *gorm.DB) ([]models.CartItem, error) {
	var cartItem []models.CartItem
	err :=
		tx.
			Clauses(clause.Returning{}).
			Unscoped().
			Where("distrib_id= ?", distrib_id).
			Delete(&cartItem).
			Error
	return cartItem, err
}

func EditCartProducts(distrib_id string, product_id string, payload models.CartItem, tx *gorm.DB) (models.CartItem, error) {
	err :=
		tx.
			Model(models.CartItem{}).
			Where("distrib_id=? AND product_id=?", distrib_id, product_id).
			Updates(payload).
			Error
	return payload, err
}

// func EditProduct(distrib_id string, product_id string, payload models.Product) (models.Product, *gorm.DB) {
// 	err := tx.Model(models.Product{}).Where("product_id=?", distrib_id, product_id).Updates(payload)
// 	return payload, err
// }

func SaveProductImage(productImage models.ProductImage, tx *gorm.DB) error {

	err :=
		tx.
			Table("product_images").
			Create(&productImage).
			Error
	return err
}

func GetProductImage(productId uint, tx *gorm.DB) (string, error) {
	var image string
	err :=
		tx.
			Table("product_images").
			Select("image").
			Where("product_id=?", productId).
			Take(&image).
			Error
	return image, err
}

func SaveProduct(product *models.Product, tx *gorm.DB) (*models.Product, error) {
	err := tx.Create(&product).Error

	return product, err
}
