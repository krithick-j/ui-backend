package dto

type ProductCart struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type CartItemIn struct {
	UserID uint          `json:"user_id"`
	Items  []ProductCart `json:"items"`
}
