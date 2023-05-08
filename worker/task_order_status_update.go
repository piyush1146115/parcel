package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/piyush1146115/parcel/config"
	"github.com/piyush1146115/parcel/data"
	"github.com/rs/zerolog/log"
)

const TaskOrderStatusUpdate = "task:order_status_update"

func (distributor *RedisTaskDistributor) DistributeTaskOrderStatusUpdate(
	ctx context.Context,
	payload *OrderStatusPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskOrderStatusUpdate, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) TaskOrderStatusUpdate(ctx context.Context, task *asynq.Task) error {
	var payload OrderStatusPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	status, err := data.GetOrderStatus(payload.OrderId)
	if err != nil {
		log.Error().Str("type", task.Type()).Err(err)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if *status == data.COMPLETED || *status == data.CANCELLED {
		if err := data.UpdateRiderStatusFromOrder(payload.OrderId, data.Available); err != nil {
			log.Error().Str("type", task.Type()).Err(err)
			return fmt.Errorf("failed to update order status: %w", err)
		}

		return nil
	}

	if err := requeueOrderStatusUpdate(payload.OrderId); err != nil {
		log.Error().Str("type", task.Type()).Err(err)
		return fmt.Errorf("failed to requeue order status update task: %w", err)
	}

	return nil
}

func requeueOrderStatusUpdate(orderID int64) error {
	config, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := NewRedisTaskDistributor(redisOpt)
	statusPayload := &OrderStatusPayload{OrderId: orderID}

	opts := []asynq.Option{
		asynq.MaxRetry(1),
		asynq.ProcessIn(1 * time.Minute),
		asynq.Queue(QueueCritical),
	}

	return taskDistributor.DistributeTaskOrderStatusUpdate(context.Background(), statusPayload, opts...)
}
