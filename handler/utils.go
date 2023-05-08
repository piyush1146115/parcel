package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getCustomerID(r *http.Request) (int, error) {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["customer_id"])
	if err != nil {
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
		return -1, err
	}

	return id, nil
}

func getRiderID(r *http.Request) (int, error) {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["rider_id"])
	if err != nil {
		return -1, err
	}

	return id, nil
}
