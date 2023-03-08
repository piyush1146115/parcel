package data

type Customer struct {
	Id               int     `json:"id"`
	CurrentLongitude float64 `json:"current_Longitude,omitempty"`
	CurrentLatitude  float64 `json:"current_Latitude,omitempty"`
}

func IsCustomerAuthorized(Id int) bool {
	for i := range customers {
		if customers[i].Id == Id {
			return true
		}
	}

	return false
}

var customers = []*Customer{
	&Customer{
		Id:               11,
		CurrentLongitude: -122.4194,
		CurrentLatitude:  37.7749,
	},
	&Customer{
		Id:               12,
		CurrentLongitude: -121.8863,
		CurrentLatitude:  37.3362,
	},
	&Customer{
		Id:               13,
		CurrentLongitude: -122.4313,
		CurrentLatitude:  37.7699,
	},
	&Customer{
		Id:               14,
		CurrentLongitude: -122.0375,
		CurrentLatitude:  37.3688,
	},
	&Customer{
		Id:               15,
		CurrentLongitude: -122.4191,
		CurrentLatitude:  37.7749,
	},
}
