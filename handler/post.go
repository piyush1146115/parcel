package handler

import (
	"encoding/json"
	"fmt"
	"github.com/piyush1146115/parcel/data"
	"net/http"
)

type ParcelResponse struct {
	OrderId int64 `json:"order_id,omitempty"`
	Success bool  `json:"success,omitempty"`
}

func NewParcelRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := getCustomerID(r)
	if err != nil {
		http.Error(w, "Could not get Customer Id from the URL", http.StatusBadRequest)
		return
	}

	var parcel data.Parcel
	if err := json.NewDecoder(r.Body).Decode(&parcel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !data.IsCustomerAuthorized(id) {
		http.Error(w, fmt.Sprintf("No authorized customer found with id: %d", id), http.StatusUnauthorized)
		return
	}

	orderId := data.CreateNewOrder(id)
	response := ParcelResponse{OrderId: orderId}
	var rider *data.Rider

	distance := haversine(parcel.PickupLatitude, parcel.PickupLongitude, parcel.DropOffLatitude, parcel.DropOffLongitude)
	if distance < 3 {
		rider = data.GetAvailableCyclist()
	} else {
		rider = data.GetAvailableBiker()
	}

	if rider == nil {
		err := data.UpdateOrderStatus(orderId, data.CANCELLED)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Success = false

		json.NewEncoder(w).Encode(response)
		return
	}

	if err := data.UpdateRiderInOrder(orderId, rider.Id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := data.UpdateOrderStatus(orderId, data.ACCEPTED); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := data.UpdateRiderStatus(rider.Id, data.OnTrip); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Success = true

	json.NewEncoder(w).Encode(response)
	return
}
