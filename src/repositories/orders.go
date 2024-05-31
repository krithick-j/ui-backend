package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveOrderHeader(order *models.OrdersHeader) error {

	result := configs.DB.Create(&order)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return nil
}

func SaveOrderLiner(order *models.OrdersLiner) error {

	result := configs.DB.Create(&order)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return nil
}

func GetOrderByDistribId(distrib_id string) ([]models.OrdersHeader, *gorm.DB) {
	var order []models.OrdersHeader
	result := configs.DB.Order("created_at DESC").Preload("OrdersLiner").Where("distrib_id = ?", distrib_id).Find(&order)
	return order, result
}

func GetAllOrders() ([]models.OrdersHeader, *gorm.DB) {
	var order []models.OrdersHeader
	result := configs.DB.Order("created_at DESC").Preload("OrdersLiner").Find(&order)
	return order, result
}
