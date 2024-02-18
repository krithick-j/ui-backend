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

func GetOrderByDistribId(distrib_id string, order []models.OrdersHeader) ([]models.OrdersHeader, *gorm.DB) {
	result := configs.DB.Preload("OrdersLiner").Find(&order, "distrib_id = ?", distrib_id)
	return order, result
}
