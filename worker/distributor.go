package worker

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/piyush1146115/parcel/data"
)

type TaskDistributor interface {
	DistributeTaskProcessOrder(
		ctx context.Context,
		payload *OrderProcessingPayload,
		opts ...asynq.Option,
	) error

	DistributeTaskOrderStatusUpdate(
		ctx context.Context,
		payload *OrderStatusPayload,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

type OrderProcessingPayload struct {
	Order  data.Order
	Parcel data.Parcel
}

type OrderStatusPayload struct {
	OrderId int64 `json:"order_id"`
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
