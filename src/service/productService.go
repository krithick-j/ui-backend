package service

import (
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

	if firstProdType == "" || request.ProductType == firstProdType {

		for _, item := range request.Items {
			ids := []uint{item.ProductID}
			_, err := repositories.GetAllProductByIDs(ids)
			if err.RowsAffected == 0 {
				return fiber.Map{"success": "Product Id does not exist"}, fiber.StatusBadRequest
			}
			product := models.CartItem{
				DistribID: request.DistribID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}
			repositories.SaveToCart(product)
		}
		return fiber.Map{"success": "Added to Cart Successfully"}, fiber.StatusOK
	} else {
		return fiber.Map{"failed": "All products in cart must be same"}, fiber.StatusBadRequest
	}
}

func GetCartProductsByDistribId(user_id string) fiber.Map {

	var products []dto.ProductsOut

	products, result := repositories.GetAllCartProductsByDistribID(user_id, products)

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

func GetOrderDetails(distrib_id string) (dto.OrderDetailsOut, int) {
	var subTotal float64 = 0.0
	var totalSandH float64 = 0.0
	var quantity uint = 0
	var orderProductArray []dto.OrderProduct
	var orderDetails dto.OrderDetailsOut
	var cartItems []dto.ProductsOut
	var userData models.User
	var TotalTypeValue float64
	//1. Retrieving All Products in Cart
	cartItems, result := repositories.GetAllCartProductsByDistribID(distrib_id, cartItems)

	if result.Error != nil {
		return orderDetails, http.StatusInternalServerError
	}

	//2. Populating OrderProduct Array field
	for _, item := range cartItems {
		orderProduct := dto.OrderProduct{
			Name:        item.Product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   uint64(item.Product.Price),
			SandH:       item.Product.SandH,
			SubTotal:    item.Product.Price * float64(item.Quantity),
			ProductType: item.Product.ProductType,
			TypeValue:   item.Product.TypeValue,
		}

		orderProductArray = append(orderProductArray, orderProduct)
		subTotal += orderProduct.SubTotal
		totalSandH += orderProduct.SandH
		quantity += item.Quantity
		TotalTypeValue += orderProduct.TypeValue * float64(item.Quantity)
	}
	//Retrieving User Data for Delivery Address
	userData, result = repositories.GetUserByID(distrib_id, userData)
	if result.Error != nil {
		println(fiber.Map{"error": result.Error})
		return orderDetails, http.StatusInternalServerError
	}

	deliveryAddress := dto.DeliveryAddress{
		ContactName:   userData.Name,
		ContactEmail:  userData.EmailAddress,
		Address:       userData.Address1,
		City:          userData.TownOrCity,
		District:      userData.District,
		State:         userData.StateOrProvince,
		ZipCode:       userData.PinOrZipCode,
		Country:       userData.Country,
		HomePhoneNo:   userData.HomePhoneNo,
		MobilePhoneNo: userData.MobilePhoneNo,
	}
	orderDetails = dto.OrderDetailsOut{
		Items:           orderProductArray,
		SubTotal:        subTotal,
		TotalSandH:      totalSandH,
		TotalAmount:     subTotal + totalSandH,
		DeliveryAddress: deliveryAddress,
		TotalQuantity:   float64(quantity),
		DistribId:       distrib_id,
		TotalTypeValue: TotalTypeValue,
	}

	if result.Error != nil {
		print(fiber.Map{"error": result.Error})
		return orderDetails, http.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		print(fiber.Map{"data": "No Products in Cart"})
		return orderDetails, http.StatusNoContent
	}

	return orderDetails, http.StatusOK
}
