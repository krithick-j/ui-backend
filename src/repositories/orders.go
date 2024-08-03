package repositories

import (
	"ui-back-end/configs"
	"ui-back-end/src/models"

	"gorm.io/gorm"
)

func SaveOrderHeader(order *models.OrdersHeader, tx *gorm.DB) error {

	err := configs.DB.Debug().Create(&order).Error
	return err
}

func SaveOrderLiner(order *models.OrdersLiner) error {

	err := configs.DB.Create(&order).Error

	return err
}

func GetOrderByDistribId(distrib_id string) ([]models.OrdersHeader, error) {
	var order []models.OrdersHeader
	err := configs.DB.Table("orders_headers").Order("created_at DESC").Preload("OrdersLiners").Where("distrib_id = ?", distrib_id).Find(&order).Error
	return order, err
}

// func GetOrderByDistribId(distrib_id string) ([]models.OrdersHeader, *gorm.DB) {
// 	var order []models.OrdersHeader
// 	query := `
// 	SELECT oh.*, ol.*
// 	FROM orders_headers oh
// 	LEFT JOIN orders_liners ol ON ol.orders_header_id = oh.id
// 	WHERE distrib_id = ?`
// 	result := configs.DB.Raw(query, distrib_id).Find(&order)
// 	return order, result
// }

// func GetOrdersByDistribId(distribId string) ([]models.OrdersHeader, error) {
// 	var orders_headers []models.OrdersHeader
// 	var orders_liners []models.OrdersLiner

// 	// Query OrdersHeader
// 	err := configs.DB.Where("distrib_id = ?", distribId).Find(&orders_headers).Error
// 	if err != nil {
// 		return orders_headers, err
// 	}

// 	// Query OrdersLiner
// 	var orders_headers_id []uint
// 	for _, orders_header := range orders_headers {
// 		orders_headers_id = append(orders_headers_id, orders_header.ID)
// 	}
// 	fmt.Println("orders header ids--->", orders_headers_id)
// 	err = configs.DB.Where("orders_header_id IN ?", orders_headers_id).Find(&orders_liners).Error
// 	fmt.Print("Liners--->", orders_liners)
// 	// Attach OrdersLiners to the corresponding OrdersHeader
// 	linerMap := make(map[uint][]models.OrdersLiner)
// 	for _, liner := range orders_liners {
// 		linerMap[liner.OrdersHeaderID] = append(linerMap[liner.OrdersHeaderID], liner)
// 	}

// 	for i, orders_header := range orders_headers {
// 		orders_headers[i].OrdersLiners = linerMap[orders_header.ID]
// 	}

// 	return orders_headers, err
// }

func GetAllOrders() ([]models.OrdersHeader, *gorm.DB) {
	var order []models.OrdersHeader
	result := configs.DB.Order("created_at DESC").Preload("OrdersLiners").Find(&order)
	return order, result
}

func GetOrderDetailsByOrderId(orderId string) (models.OrdersHeader, error) {
	order := models.OrdersHeader{}
	err := configs.DB.Preload("OrdersLiners").Where("order_id = ?", orderId).Find(&order).Error
	return order, err
}
