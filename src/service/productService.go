package service

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/models"
	"ui-back-end/src/repositories"
	"ui-back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetProducts(tx *gorm.DB) (fiber.Map, int) {
	products, err := repositories.GetAllProducts(tx)
	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	if err != nil {
		utils.NotNilErrorMessage(err, "GetAllProducts", "GetProducts", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage(products, fiber.StatusOK)
}

func GetProductCategories(tx *gorm.DB) (fiber.Map, int) {
	var productCategories []models.ProductCategory
	productCategories, err := repositories.GetAllProductCategories(productCategories, tx)
	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllProductCategories", "GetProductCategories", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage(productCategories, fiber.StatusOK)
}

func GetProductByCategoryID(categoryIds []string, productTypes []string, tx *gorm.DB) (fiber.Map, int) {

    // Construct the query with both category IDs and product types
    product, err := repositories.GetAllProductByCategoryIdAndProductTypeFilter(categoryIds, productTypes, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllProductByCategoryIdAndProductTypeFilter", "GetProductByCategoryID", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage(product, fiber.StatusOK)
}

func GetEpProductsByCategoryId(category_id string, tx *gorm.DB) (fiber.Map, int) {

	product, err := repositories.GetAllEpProductByCategoryID(category_id, tx)

	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllEpProductByCategoryID", "GetEpProductsByCategoryId", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage(product, fiber.StatusOK)
}

func GetProductsByIds(ids []uint, tx *gorm.DB) (fiber.Map, int) {
	products, err := repositories.GetAllProductByIDs(ids, tx)

	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllProductByIDs", "GetProductsByIds", fiber.StatusInternalServerError, tx)
	}
	return utils.SuccessMessage(products, fiber.StatusOK)
}

func AddToCart(request dto.CartItemIn, tx *gorm.DB) (fiber.Map, int) {

	var firstProdType string = ""
	var existingProductFlag bool = false
	var existingProduct dto.ProductsOut

	firstCartItem, err := repositories.GetFirstCartItem(request.DistribID, tx)

	if err != gorm.ErrRecordNotFound {
		firstProdType, err = repositories.GetProductTypeByProductID(firstCartItem.ProductID, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetProductTypeByProductID", "AddToCart", fiber.StatusInternalServerError, tx)
		}
	}

	if firstProdType == "" || request.ProductType == firstProdType {

		existingProducts, err := repositories.GetAllCartProductsByDistribID(request.DistribID, tx)

		if err != nil {
			return utils.NotNilErrorMessage(err, "GetAllCartProductsByDistribID", "AddToCart", fiber.StatusBadRequest, tx)
		}

		ids := []uint{request.Product.ProductID}
		_, err = repositories.GetAllProductByIDs(ids, tx)
		if err == gorm.ErrRecordNotFound {
			return utils.RecordNotFoundMessage(err, tx)
		}
		//Checks existing product length is 0, so that error does not hit when looping
		if len(existingProducts) != 0 {
			//Checks product id already exist in cart
			for _, existingCartProduct := range existingProducts {
				if request.Product.ProductID == existingCartProduct.ProductID {
					existingProductFlag = true
					existingProduct = existingCartProduct
					break
				}
			}
			//if cart product already exist, quantity is updating
			if existingProductFlag {
				configs.Log.Infoln("Updating Cart Product Quantity of product id: ", existingProduct.ID)
				err := repositories.UpdateCartProductQuantityById(existingProduct.ID, existingProduct.Quantity+1, tx)
				if err != nil {
					return utils.NotNilErrorMessage(err, "UpdateCartProductQuantityById", "AddToCart", fiber.StatusInternalServerError, tx)
				}
			} else { //else saving as a new cart product
				product := models.CartItem{
					DistribID: request.DistribID,
					ProductID: request.Product.ProductID,
					Quantity:  request.Product.Quantity,
				}
				configs.Log.Infoln("Saving new Cart Product of product id: ", product.ProductID)
				err := repositories.SaveToCart(product, tx)
				if err != nil {
					return utils.NotNilErrorMessage(err, "SaveToCart", "AddToCart", fiber.StatusInternalServerError, tx)
				}
			}
		} else { //Saving as a new cart product when there is no existing cart product
			product := models.CartItem{
				DistribID: request.DistribID,
				ProductID: request.Product.ProductID,
				Quantity:  request.Product.Quantity,
			}
			configs.Log.Infoln("Saving new Cart Product when there are no existing product of product id: ", product.ProductID)
			err := repositories.SaveToCart(product, tx)
			if err != nil {
				return utils.NotNilErrorMessage(err, "SaveToCart", "AddToCart", fiber.StatusInternalServerError, tx)
			}
		}
		return fiber.Map{"success": "Added to Cart Successfully"}, fiber.StatusOK
	} else {
		return fiber.Map{"failed": "All products in cart must be same"}, fiber.StatusBadRequest
	}
}

func GetCartProductsByDistribId(user_id string, tx *gorm.DB) (fiber.Map, int) {

	var products []dto.ProductsOut

	configs.Log.Infoln("Retrieving all products from cart INIT")
	products, err := repositories.GetAllCartProductsByDistribID(user_id, tx)

	if err != nil {
		return utils.NotNilErrorMessage(err, "GetAllCartProductsByDistribID", "GetCartProductsByDistribId", fiber.StatusInternalServerError, tx)
	}

	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	configs.Log.Infoln("Retrieving all products from cart DONE")
	configs.Log.Infoln("validating cart products that exists in Product table")
	for _, product := range products {
		_, err := repositories.GetProductById(product.ProductID, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "GetProductById", "GetCartProductsByDistribId", fiber.StatusInternalServerError, tx)
		}
	}
	configs.Log.Infoln("validating cart products that exists in Product table DONE")

	return fiber.Map{"data": products}, fiber.StatusOK
}

func DeleteCartProduct(distrib_id string, product_id string, tx *gorm.DB) (fiber.Map, int) {
	_, err := repositories.DeleteCartProduct(distrib_id, product_id, tx)

	if err != nil {
		return utils.NotNilErrorMessage(err, "DeleteCartProduct", "DeleteCartProduct", fiber.StatusInternalServerError, tx)
	}

	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	return utils.SuccessMessage("Product Deleted Successfully", fiber.StatusOK)
}

func DeleteAllCartProduct(distrib_id string, tx *gorm.DB) (fiber.Map, int) {
	cartItem, err := repositories.DeleteAllCartProduct(distrib_id, tx)

	if err != nil {
		return utils.NotNilErrorMessage(err, "DeleteAllCartProduct", "DeleteAllCartProduct", fiber.StatusInternalServerError, tx)
	}

	if err == gorm.ErrRecordNotFound {
		return utils.RecordNotFoundMessage(err, tx)
	}
	return fiber.Map{"success": "All Products in Cart Deleted Successfully", "deleted_products": cartItem}, fiber.StatusOK
}

func EditCartProducts(payload models.CartItem, distrib_id string, product_id string, tx *gorm.DB) (fiber.Map, int) {

	cartItem, err := repositories.EditCartProducts(distrib_id, product_id, payload, tx)

	if err != nil {
		return utils.NotNilErrorMessage(err, "EditCartProducts", "EditCartProducts", fiber.StatusInternalServerError, tx)
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
func CreateProduct(c *fiber.Ctx, form *multipart.Form, tx *gorm.DB) (fiber.Map, int) {
	data := c.FormValue("data")
	product := models.Product{}
	json.Unmarshal([]byte(data), &product)
	prod, err := repositories.SaveProduct(&product, tx)
	if err != nil {
		return utils.NotNilErrorMessage(err, "SaveProduct", "CreateProduct", fiber.StatusInternalServerError, tx)
	}
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
				return utils.NotNilErrorMessage(err, "SaveFile", "CreateProduct", fiber.StatusInternalServerError, tx)
			}
			pm.Image = "media/" + fmt.Sprintf("%d", pid) + "-" + fs + extension
			pm.ProductID = pid
			pms = append(pms, pm)
		}
	}
	for _, pm := range pms {

		err := repositories.SaveProductImage(pm, tx)
		if err != nil {
			return utils.NotNilErrorMessage(err, "SaveProductImage", "CreateProduct", fiber.StatusInternalServerError, tx)
		}
	}
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
