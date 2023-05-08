package data

import (
	"fmt"
	"github.com/piyush1146115/parcel/utils"
	"math"
)

type Rider struct {
	Id               int         `json:"id,omitempty"`
	CurrentLongitude float64     `json:"current_Longitude,omitempty"`
	CurrentLatitude  float64     `json:"current_Latitude,omitempty"`
	Type             RiderType   `json:"rider_type,omitempty"`
	Status           RiderStatus `json:"status,omitempty"`
}

type RiderType int8

const (
	BIKER RiderType = iota
	CYCLIST
)

type RiderStatus string

const (
	Available RiderStatus = "Available"
	Offline   RiderStatus = "Offline"
	OnTrip    RiderStatus = "On Trip"
)

func GetAvailableCyclist(long, lat float64) *Rider {
	var rider *Rider
	minDis := math.MaxFloat64

	for _, r := range riderList {
		if r.Type == CYCLIST && r.Status == Available {
			//return r
			dis := utils.Haversine(r.CurrentLatitude, r.CurrentLongitude, lat, long)
			if dis < minDis {
				rider = r
				minDis = dis
			}
		}
	}

	return rider
}

func GetAvailableBiker(long, lat float64) *Rider {
	var rider *Rider
	minDis := math.MaxFloat64

	for _, r := range riderList {
		if r.Type == BIKER && r.Status == Available {
			//return r
			dis := utils.Haversine(r.CurrentLatitude, r.CurrentLongitude, lat, long)
			if dis < minDis {
				rider = r
				minDis = dis
			}
		}
	}

	return rider
}

var ErrRiderNotFound = fmt.Errorf("Rider not found")

func UpdateRiderStatus(id int, status RiderStatus) error {
	i := findRiderByID(id)
	if i == -1 {
		return ErrRiderNotFound
	}

	riderList[i].Status = status
	return nil
}

func GetRiderStatus(id int) (*RiderStatus, error) {
	i := findRiderByID(id)
	if i == -1 {
		return nil, ErrRiderNotFound
	}

	return &riderList[i].Status, nil
}

func GetTotalAvailableRiders() int {
	return len(riderList)
}

func IsValidRiderId(id int) bool {
	i := findRiderByID(id)
	return i != -1
}

func GetRidersCurrentLocation(id int) (*float64, *float64, error) {
	i := findRiderByID(id)
	if i == -1 {
		return nil, nil, ErrRiderNotFound
	}

	return &riderList[i].CurrentLongitude, &riderList[i].CurrentLatitude, nil
}

func UpdateRidersLocation(id int, long float64, lat float64) error {
	i := findRiderByID(id)
	if i == -1 {
		return ErrRiderNotFound
	}

	riderList[i].CurrentLongitude = long
	riderList[i].CurrentLatitude = lat

	return nil
}

func findRiderByID(id int) int {
	for i, r := range riderList {
		if r.Id == id {
			return i
		}
	}

	return -1
}

var riderList = []*Rider{
	&Rider{
		Id:               1,
		CurrentLongitude: -122.4194,
		CurrentLatitude:  37.7749,
		Type:             BIKER,
		Status:           Available,
	},
	&Rider{
		Id:               2,
		CurrentLongitude: -121.8863,
		CurrentLatitude:  37.3362,
		Type:             CYCLIST,
		Status:           OnTrip,
	},
	&Rider{
		Id:               3,
		CurrentLongitude: -122.4313,
		CurrentLatitude:  37.7699,
		Type:             BIKER,
		Status:           Available,
	},
	&Rider{
		Id:               4,
		CurrentLongitude: -122.4313,
		CurrentLatitude:  37.7699,
		Type:             BIKER,
		Status:           Offline,
	},
	&Rider{
		Id:               5,
		CurrentLongitude: -122.4191,
		CurrentLatitude:  37.7749,
		Type:             CYCLIST,
		Status:           Available,
	},
	&Rider{
		Id:               6,
		CurrentLongitude: -122.4056,
		CurrentLatitude:  37.7886,
		Type:             BIKER,
		Status:           Available,
	},
	&Rider{
		Id:               7,
		CurrentLongitude: -122.4314,
		CurrentLatitude:  37.7739,
		Type:             CYCLIST,
		Status:           Available,
	},
	&Rider{
		Id:               8,
		CurrentLongitude: -122.4383,
		CurrentLatitude:  37.7613,
		Type:             BIKER,
		Status:           OnTrip,
	},
	&Rider{
		Id:               9,
		CurrentLongitude: -122.4727,
		CurrentLatitude:  37.7510,
		Type:             CYCLIST,
		Status:           Available,
	},
	&Rider{
		Id:               10,
		CurrentLongitude: -122.4014,
		CurrentLatitude:  37.7887,
		Type:             CYCLIST,
		Status:           Offline,
	},
}
