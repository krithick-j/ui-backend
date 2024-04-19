package service

import (
	"fmt"
	"net/http"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

func GetProductByCategoryID(categoryId string, productType string) (fiber.Map, int) {

	var product []models.Product
	var result *gorm.DB

	if productType == "" {
		product, result = repositories.GetAllProductByCategoryID(categoryId, product)
	} else {
		product, result = repositories.GetAllProductByCategoryIdAndProductType(categoryId, productType)
	}

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": product}, fiber.StatusOK
}

func GetEpProductsByCategoryId(category_id string) (fiber.Map, int) {

	var product []models.Product
	var result *gorm.DB

	product, result = repositories.GetAllEpProductByCategoryID(category_id, product)

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}

	return fiber.Map{"data": product}, fiber.StatusOK
}

func GetProductsByIds(ids []uint) ([]models.Product, error) {
	products, result := repositories.GetAllProductByIDs(ids)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func AddToCart(request dto.CartItemIn) (fiber.Map, int) {
	var firstProdType string
	firstCartItem, err := repositories.GetFirstCartItem(request.DistribID)

	if err.Error != gorm.ErrRecordNotFound {
		firstProdType, err = repositories.GetProductTypeByProductID(firstCartItem.ProductID)
		if err.Error != nil {
			return fiber.Map{"message": err.Error}, fiber.StatusInternalServerError
		}
	}

	fmt.Println("firstprodtype", firstProdType, "req prod type", request)

	if firstProdType == "" || request.ProductType == firstProdType {

		productsOut, res := repositories.GetAllCartProductsByDistribID(request.DistribID)

		if res.Error != nil {
			fmt.Println(res.Error.Error())
			return fiber.Map{"error": res.Error.Error()}, fiber.StatusBadRequest
		}

		for _, item := range request.Items {
			ids := []uint{item.ProductID}
			_, err := repositories.GetAllProductByIDs(ids)
			if err.RowsAffected == 0 {
				return fiber.Map{"success": "Product Id does not exist"}, fiber.StatusBadRequest

			}
			if len(productsOut) != 0 {
				for _, existingCartProduct := range productsOut {
					if item.ProductID == existingCartProduct.ProductID {
						repositories.UpdateCartProductQuantityById(existingCartProduct.ID, existingCartProduct.Quantity+1)
					} else {
						product := models.CartItem{
							DistribID: request.DistribID,
							ProductID: item.ProductID,
							Quantity:  item.Quantity,
						}
						repositories.SaveToCart(product)
					}
				}
			} else {
				product := models.CartItem{
					DistribID: request.DistribID,
					ProductID: item.ProductID,
					Quantity:  item.Quantity,
				}
				repositories.SaveToCart(product)
				break
			}
		}
		return fiber.Map{"success": "Added to Cart Successfully"}, fiber.StatusOK
	} else {
		return fiber.Map{"failed": "All products in cart must be same"}, fiber.StatusBadRequest
	}
}

func GetCartProductsByDistribId(user_id string) fiber.Map {

	var products []dto.ProductsOut

	products, result := repositories.GetAllCartProductsByDistribID(user_id)

	if result.Error != nil {
		return fiber.Map{"error": result.Error}
	}

	if result.RowsAffected == 0 {
		return fiber.Map{"data": "No Products in Cart"}
	}

	return fiber.Map{"data": products}
}

func DeleteCartProduct(distrib_id string, product_id string) (fiber.Map, int) {
	var cartItem models.CartItem
	cartItem, result := repositories.DeleteCartProduct(distrib_id, product_id, cartItem)

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"success": "Product Not Found", "DeletedProduct": cartItem}, http.StatusNoContent
	}

	return fiber.Map{"success": "Product Deleted Successfully", "DeletedProduct": cartItem}, http.StatusOK
}

func DeleteAllCartProduct(distrib_id string) (fiber.Map, int) {
	cartItem, result := repositories.DeleteAllCartProduct(distrib_id)

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	if result.Error == gorm.ErrRecordNotFound {
		return fiber.Map{"success": "Product Not Found", "DeletedProduct": cartItem}, http.StatusNoContent
	}

	return fiber.Map{"success": "All Products in Cart Deleted Successfully", "deleted_products": cartItem}, http.StatusOK
}

func EditCartProducts(payload models.CartItem, distrib_id string, product_id string) (fiber.Map, int) {

	cartItem, result := repositories.EditCartProducts(distrib_id, product_id, payload)

	if result.Error != nil {
		log.Info("Error saving user to the database:", result.Error)
		return fiber.Map{"error": result.Error}, http.StatusBadGateway
	}

	if result.Error != nil {
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}
	return fiber.Map{"success": "Product Updated Successfully", "UpdatedProduct": cartItem}, http.StatusOK
}

func CreateProduct(payload dto.ProductIn, adminName string) (fiber.Map, int) {
	//saving product
	product := &models.Product{
		Name:              payload.Name,
		Quantity:          payload.Quantity,
		ShipmentTime:      payload.ShipmentTime,
		Price:             payload.Price,
		SandH:             payload.SandH,
		ProductCategoryID: payload.ProductCategoryID,
		TypeValue:         payload.TypeValue,
		ProductType:       payload.ProductType,
		AdminName:         adminName,
	}
	product = repositories.SaveProduct(product)
	//group of pictures stores in product image table
	for _, image := range payload.ProductImages {
		productImage := &models.ProductImage{
			Image:     image.Image,
			ProductID: product.ID,
		}
		repositories.SaveProductImage(productImage)
	}
	return fiber.Map{"data": "Product Successfully created"}, http.StatusCreated
}
