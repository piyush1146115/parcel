package data

type Parcel struct {
	ReceiverName     string  `json:"receiver_name,omitempty"`
	ReceiverNumber   string  `json:"receiver_number,omitempty"`
	PickupAddress    string  `json:"pickup_address,omitempty"`
	DropOffAddress   string  `json:"dropoff_address,omitempty"`
	PickupLatitude   float64 `json:"pickup_latitude"`
	PickupLongitude  float64 `json:"pickup_longitude"`
	DropOffLatitude  float64 `json:"dropoff_latitude"`
	DropOffLongitude float64 `json:"dropoff_longitude"`
	OrderId          int64   `json:"order_id,omitempty"`
}

var parcelList []*Parcel
