package handler

import (
	"encoding/json"
	"github.com/piyush1146115/parcel/data"
	"net/http"
)

type OrderResponse struct {
	OrderId int64 `json:"order_id,omitempty"`
	Success bool  `json:"success,omitempty"`
}

func HandleNewParcelRequest(w http.ResponseWriter, r *http.Request) {
	var parcel data.Parcel
	if err := json.NewDecoder(r.Body).Decode(&parcel); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !data.IsCustomerAuthorized(parcel.CustomerID) {
		http.Error(w, "User is unauthorized to place this order", http.StatusUnauthorized)
		return
	}

	orderId := data.CreateNewOrder(parcel.CustomerID)
	response := OrderResponse{OrderId: orderId}
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
		w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return

}
