package handler

import (
	"encoding/json"
	"fmt"
	"github.com/piyush1146115/parcel/data"
	"net/http"
)

type OrderResponse struct {
	OrderId     int64            `json:"order_id,omitempty"`
	OrderStatus data.OrderStatus `json:"order_status,omitempty"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to the Homepage of parcel simulator!\n")
}

func OrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orderId, err := getOrderID(r)
	if err != nil {
		http.Error(w, "Could not get Order Id from the URL", http.StatusBadRequest)
		return
	}

	if !data.IsValidOrderId(int64(orderId)) {
		http.Error(w, fmt.Sprintf("Invalid order ID: %d", orderId), http.StatusBadRequest)
		return
	}

	status, err := data.GetOrderStatus(int64(orderId))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get status for order id: %d", orderId), http.StatusInternalServerError)
		return
	}

	response := OrderResponse{OrderId: int64(orderId), OrderStatus: *status}
	json.NewEncoder(w).Encode(response)
	return
}
