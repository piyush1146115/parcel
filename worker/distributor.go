package worker

import (
	"context"
	"github.com/piyush1146115/parcel/data"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskProcessOrder(
		ctx context.Context,
		payload *OrderProcessingPayload,
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

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
