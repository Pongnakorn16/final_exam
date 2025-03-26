package dto

type CartItemRequest struct {
	CustomerID int    `json:"customer_id" binding:"required"`
	CartName   string `json:"cart_name" binding:"required"`
	ProductID  int    `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
}
