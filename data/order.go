package data

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Order struct {
	Id         int64       `json:"id,omitempty"`
	RiderId    int         `json:"rider_id,omitempty"`
	CustomerID int         `json:"customer_id,omitempty"`
	Status     OrderStatus `json:"status,omitempty"`
}

type OrderStatus string

const (
	REQUESTED      OrderStatus = "Requested"
	COMPLETED      OrderStatus = "Completed"
	ACCEPTED       OrderStatus = "Accepted"
	INPROGRESS     OrderStatus = "In Progress"
	CANCELLED      OrderStatus = "Cancelled"
	SEARCHINGRIDER OrderStatus = "Searching Rider"
)

var orderIDCounter int64 = 1

func generateOrderID() int64 {
	// Get the current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Increment the order ID counter
	orderID := atomic.AddInt64(&orderIDCounter, 1)

	// Combine the timestamp and counter to generate a unique ID
	return timestamp*1000000 + orderID
}

var orderList []*Order

func CreateNewOrder(cid int) int64 {
	newId := generateOrderID()

	orderList = append(orderList, &Order{
		Id:         newId,
		CustomerID: cid,
		Status:     REQUESTED,
	})

	return newId
}

var ErrOrderNotFound = fmt.Errorf("Order not found")

func UpdateRiderInOrder(id int64, rid int) error {
	i := findOrderById(id)
	if i == -1 {
		return ErrOrderNotFound
	}

	orderList[i].RiderId = rid
	return nil
}

func UpdateOrderStatus(id int64, status OrderStatus) error {
	i := findOrderById(id)
	if i == -1 {
		return ErrOrderNotFound
	}

	orderList[i].Status = status
	return nil
}

func IsValidOrderId(id int64) bool {
	i := findOrderById(id)
	return i != -1
}

func GetOrderStatus(id int64) (*OrderStatus, error) {
	i := findOrderById(id)
	if i == -1 {
		return nil, ErrOrderNotFound
	}

	return &orderList[i].Status, nil
}

func GetAnOrderIdWithInProgressStatus() int64 {
	for i := range orderList {
		if orderList[i].Status == INPROGRESS {
			return orderList[i].Id
		}
	}

	return -1
}

func GetAnOrderIdWithAcceptedStatus() int64 {
	for i := range orderList {
		if orderList[i].Status == ACCEPTED {
			return orderList[i].Id
		}
	}

	return -1
}

func UpdateRiderStatusFromOrder(orderId int64, status RiderStatus) error {
	i := findOrderById(orderId)
	if i == -1 {
		return ErrOrderNotFound
	}

	rid := orderList[i].RiderId
	return UpdateRiderStatus(rid, status)
}

func findOrderById(id int64) int {
	for i, o := range orderList {
		if o.Id == id {
			return i
		}
	}

	return -1
}
