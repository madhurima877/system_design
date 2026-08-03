package model

type Order struct {
	ID         int
	CustomerID string
	Status     string
}
type UpdateOrderStatusRequest struct {
	OrderID    int    `json:"order_id"`
	CustomerID string `json:"customer_id"`
	Status     string `json:"status"`
}
