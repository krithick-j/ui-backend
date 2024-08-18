package repositories

import (
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveOrderHeader(order *models.OrdersHeader, tx *gorm.DB) error {

	err :=
		tx.
			Debug().
			Create(&order).
			Error
	return err
}

func SaveOrderLiner(order *models.OrdersLiner, tx *gorm.DB) error {

	err := tx.Create(&order).Error

	return err
}

func GetOrderByDistribId(distrib_id string, tx *gorm.DB) ([]models.OrdersHeader, error) {
	var order []models.OrdersHeader
	err :=
		tx.
			Table("orders_headers").
			Order("created_at DESC").
			Preload("OrdersLiners").
			Where("distrib_id = ?", distrib_id).
			Find(&order).
			Error
	return order, err
}

func GetAllOrders(tx *gorm.DB) ([]models.OrdersHeader, error) {
	var order []models.OrdersHeader
	err :=
		tx.
			Order("created_at DESC").
			Preload("OrdersLiners").
			Find(&order).
			Error
	return order, err
}

func GetOrderDetailsByOrderId(orderId string, tx *gorm.DB) (models.OrdersHeader, error) {
	order := models.OrdersHeader{}
	err :=
		tx.
			Preload("OrdersLiners").
			Where("order_id = ?", orderId).
			Find(&order).
			Error
	return order, err
}
