package data

type Parcel struct {
	CustomerID       int     `json:"customer_id"`
	ReceiverName     string  `json:"receiver_name,omitempty"`
	ReceiverNumber   string  `json:"receiver_number,omitempty"`
	PickupAddress    string  `json:"pickup_address"`
	DropOffAddress   string  `json:"dropoff_address"`
	PickupLatitude   float64 `json:"pickup_latitude"`
	PickupLongitude  float64 `json:"pickup_longitude"`
	DropOffLatitude  float64 `json:"dropoff_latitude"`
	DropOffLongitude float64 `json:"dropoff_longitude"`
}

type Response struct {
	Success bool `json:"success"`
	OrderID int  `json:"order_id,omitempty"`
}
