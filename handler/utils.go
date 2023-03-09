package handler

import (
	"github.com/gorilla/mux"
	"math"
	"net/http"
	"strconv"
)

const (
	earthRadiusKm = 6371
)

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)
	lat1 = degreesToRadians(lat1)
	lat2 = degreesToRadians(lat2)

	a := math.Pow(math.Sin(dLat/2), 2) +
		math.Pow(math.Sin(dLon/2), 2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func getCustomerID(r *http.Request) (int, error) {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["customer_id"])
	if err != nil {
		// should never happen
		return -1, err
	}

	return id, nil
}

func getOrderID(r *http.Request) (int, error) {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["order_id"])
	if err != nil {
		// should never happen
		return -1, err
	}

	return id, nil
}
