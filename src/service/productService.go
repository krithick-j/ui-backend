package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetProducts() (fiber.Map, int) {
	var products []models.Product
	var result *gorm.DB
	products, result = repositories.GetAllProducts(products)
	println(result)
	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Infoln("RecordNotFound", result.Error.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorln("Error on calling GetAllProducts repositories fn from GetProducts service fn",result.Error.Error())
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": products},fiber.StatusOK
}

func GetProductCategories() (fiber.Map,int) {
	var productCategories []models.ProductCategory
	var result *gorm.DB
	productCategories, result = repositories.GetAllProductCategories(productCategories)
	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Errorln("RecordNotFound",result.Error.Error())
		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorln("Error on calling GetAllProductCategories repositories fn from GetProductCategories service fn",result.Error.Error())
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}
	return fiber.Map{"data": productCategories}, fiber.StatusOK
}

func GetProductByCategoryID(categoryId string, productType string) (fiber.Map,int) {

	var product []models.Product
	var result *gorm.DB

	if productType == "" {
		product, result = repositories.GetAllProductByCategoryID(categoryId, product)

		if result.Error == gorm.ErrRecordNotFound {
			configs.Log.Errorln("RecordNotFound",result.Error.Error())
			return fiber.Map{"data": "Not Found"},fiber.StatusNotFound
		}
		if result.Error != nil {
			configs.Log.Errorln("Error on calling GetAllProductByCategoryID repositories fn from GetProductByCategoryID service fn",result.Error.Error())
			return fiber.Map{"error": result.Error},fiber.StatusInternalServerError
		}
	} else {
		product, result = repositories.GetAllProductByCategoryIdAndProductType(categoryId, productType)
		
		if result.Error == gorm.ErrRecordNotFound {
			configs.Log.Errorln("RecordNotFound",result.Error.Error())
			return fiber.Map{"data": "Not Found"},fiber.StatusNotFound
		}
		if result.Error != nil {
			configs.Log.Errorln("Error on calling GetAllProductByCategoryIdAndProductType repositories fn from GetProductByCategoryID service fn",result.Error.Error())
			return fiber.Map{"error": result.Error},fiber.StatusInternalServerError
		}
	}

	return fiber.Map{"data": product},fiber.StatusOK
}

func GetEpProductsByCategoryId(category_id string) (fiber.Map, int) {

	var product []models.Product
	var result *gorm.DB

	product, result = repositories.GetAllEpProductByCategoryID(category_id, product)

	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Errorln("RecordNotFound",result.Error.Error())

		return fiber.Map{"data": "Not Found"}, fiber.StatusNotFound
	}
	if result.Error != nil {
		configs.Log.Errorln("Error on calling GetAllEpProductByCategoryID repositories fn from GetEpProductsByCategoryId service fn",result.Error.Error())
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

func GetCartProductsByDistribId(user_id string) (fiber.Map, int) {

	var products []dto.ProductsOut

	configs.Log.Infoln("Retrieving all products from cart INIT")
	products, result := repositories.GetAllCartProductsByDistribID(user_id)

	if result.Error != nil {
		configs.Log.Errorf("Error while retrieving the cart products, %s", result.Error.Error())
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}

	if result.RowsAffected == 0 {
		configs.Log.Infoln("No Products in Cart")
		return fiber.Map{"data": "No Products in Cart"}, fiber.StatusOK
	}
	configs.Log.Infoln("Retrieving all products from cart DONE")
	configs.Log.Infoln("validating cart products that exists in Product table")
	for _, product := range products {
		_, result := repositories.GetProductById(product.ProductID)
		if result.Error != nil {
			configs.Log.Errorf("Error while retrieving the cart products, %s", result.Error.Error())
			return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
		}
	}
	configs.Log.Infoln("validating cart products that exists in Product table DONE")

	return fiber.Map{"data": products}, fiber.StatusOK
}

func DeleteCartProduct(distrib_id string, product_id string) (fiber.Map, int) {
	var cartItem models.CartItem
	cartItem, result := repositories.DeleteCartProduct(distrib_id, product_id, cartItem)

	if result.Error != nil {
		configs.Log.Errorln("Error on calling DeleteCartProduct repositories fn from DeleteCartProduct service fn ", result.Error.Error())
		return fiber.Map{"error": result.Error}, http.StatusInternalServerError
	}

	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Errorln("RecordNotFound ", result.Error.Error())
		return fiber.Map{"success": "Product Not Found", "DeletedProduct": cartItem}, http.StatusNoContent
	}

	return fiber.Map{"success": "Product Deleted Successfully", "DeletedProduct": cartItem}, http.StatusOK
}

func DeleteAllCartProduct(distrib_id string) (fiber.Map, int) {
	cartItem, result := repositories.DeleteAllCartProduct(distrib_id)

	if result.Error != nil {
		configs.Log.Errorln("Error on calling DeleteAllCartProduct repositories fn from DeleteAllCartProduct service fn",result.Error.Error())
		return fiber.Map{"error": result.Error}, fiber.StatusInternalServerError
	}

	if result.Error == gorm.ErrRecordNotFound {
		configs.Log.Errorln("RecordNotFound",result.Error.Error())

		return fiber.Map{"success": "Product Not Found", "DeletedProduct": cartItem}, fiber.StatusNoContent
	}

	return fiber.Map{"success": "All Products in Cart Deleted Successfully", "deleted_products": cartItem}, fiber.StatusOK
}

func EditCartProducts(payload models.CartItem, distrib_id string, product_id string) (fiber.Map, int) {

	cartItem, result := repositories.EditCartProducts(distrib_id, product_id, payload)

	if result.Error != nil {
		configs.Log.Errorln("Error saving user to the database:", result.Error.Error())
		return fiber.Map{"error": result.Error.Error()}, fiber.StatusInternalServerError
	}

	return fiber.Map{"success": "Product Updated Successfully", "UpdatedProduct": cartItem}, fiber.StatusOK
}

// func EditProduct(payload models.Product, distrib_id string, product_id string) (fiber.Map, int) {

// 	cartItem, result := repositories.EditProduct(distrib_id, product_id, payload)

// 	if result.Error != nil {
// 		configs.Log.Errorln("Error saving user to the database:", result.Error.Error())
// 		return fiber.Map{"error": result.Error.Error()}, http.StatusBadGateway
// 	}

// 	return fiber.Map{"success": "Product Updated Successfully", "UpdatedProduct": cartItem}, http.StatusOK
// }

// func CreateProduct(payload dto.ProductIn, adminName string) (fiber.Map, int) {
func CreateProduct(c *fiber.Ctx, form *multipart.Form) (fiber.Map, int) {
	data := c.FormValue("data")
	product := models.Product{}
	json.Unmarshal([]byte(data), &product)
	prod := repositories.SaveProduct(&product)
	var pms []models.ProductImage
	pid := prod.ID
	for fs, fhs := range form.File {
		for _, fh := range fhs {
			pm := models.ProductImage{}
			extension := filepath.Ext(fh.Filename)
			fmt.Println(extension, fh.Filename)
			fmt.Println(fs)
			fullPath := "./assets/" + fmt.Sprintf("%d", pid) + "-" + fs + extension
			err := c.SaveFile(fh, fullPath)
			if err != nil {
				configs.Log.Errorln("Error on SaveFile fn from CreateProduct service fn")
				fmt.Println(err.Error())
			}
			pm.Image = "media/" + fmt.Sprintf("%d", pid) + "-" + fs + extension
			pm.ProductID = pid
			pms = append(pms, pm)
		}
	}
	repositories.SaveProductImage(pms)
	return fiber.Map{"data": "Product Successfully created"}, fiber.StatusCreated
}

// func EditProduct(payload models.Product, distrib_id string, product_id string) (fiber.Map, int) {
// 	func CreateProduct(c *fiber.Ctx, form *multipart.Form) (fiber.Map, int) {
// 							fullPath := "./assets/" + fmt.Sprintf("%d", pid) + "-" + fs + extension
// 							err := c.SaveFile(fh, fullPath)
// 							if err != nil {
// 	                               configs.Log.Errorln("Error on SaveFile fn from CreateProduct service fn")
// 									fmt.Println(err.Error())
// 							}
// 							pm.Image = "media/" + fmt.Sprintf("%d", pid) + "-" + fs + extension
// 	func CreateProduct(c *fiber.Ctx, form *multipart.Form) (fiber.Map, int) {
// 					}
// 			}
// 			repositories.SaveProductImage(pms)
// 	       return fiber.Map{"data": "Product Successfully created"}, fiber.StatusCreated
// 	 }