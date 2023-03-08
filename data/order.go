package data

type Order struct {
	OrderId    int         `json:"order_id,omitempty"`
	RiderId    int         `json:"rider_id,omitempty"`
	CustomerID int         `json:"customer_id,omitempty"`
	Status     OrderStatus `json:"status,omitempty"`
}

type OrderStatus string

const (
	COMPLETED  OrderStatus = "Completed"
	ACCEPTED   OrderStatus = "Accepted"
	INPROGRESS OrderStatus = "In Progress"
	CANCELLED  OrderStatus = "Cancelled"
)
