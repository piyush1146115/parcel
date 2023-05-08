package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/piyush1146115/parcel/data"
	"github.com/piyush1146115/parcel/worker"
)

type ParcelResponse struct {
	OrderId int64 `json:"order_id,omitempty"`
	Success bool  `json:"success,omitempty"`
}

type OrderHandler struct {
	taskDistributor worker.TaskDistributor
}

func NewOrderHandler(t worker.TaskDistributor) *OrderHandler {
	return &OrderHandler{
		taskDistributor: t,
	}
}

func (oh *OrderHandler) NewParcelRequest(w http.ResponseWriter, r *http.Request) {
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

	taskPayload := &worker.OrderProcessingPayload{
		Order: data.Order{
			Id: orderId,
		},
		Parcel: parcel,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(30 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	if err := oh.taskDistributor.DistributeTaskProcessOrder(context.Background(), taskPayload, opts...); err != nil {
		http.Error(w, fmt.Sprintf("failed to process order"), http.StatusInternalServerError)
		return
	}

	statusPayload := &worker.OrderStatusPayload{OrderId: orderId}
	if err := oh.taskDistributor.DistributeTaskOrderStatusUpdate(context.Background(), statusPayload, opts...); err != nil {
		http.Error(w, fmt.Sprintf("failed to process order"), http.StatusInternalServerError)
		return
	}

	response := ParcelResponse{OrderId: orderId, Success: true}
	json.NewEncoder(w).Encode(response)
	return
}
