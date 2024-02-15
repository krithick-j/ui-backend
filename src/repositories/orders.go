package repositories

import (
	"fmt"
	"ui-back-end/configs"
	"ui-back-end/src/models"
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


func SaveOrderFooter(order *models.OrderFooter) error {

	result := configs.DB.Create(&order)
	if result.Error != nil {
		fmt.Printf("Error %v\n", result.Error.Error())
	}

	return nil
}