package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/piyush1146115/parcel/config"
	"github.com/piyush1146115/parcel/data"
	"github.com/piyush1146115/parcel/utils"
	"github.com/rs/zerolog/log"
	"time"
)

const TaskProcessOrder = "task:process_order"

func (distributor *RedisTaskDistributor) DistributeTaskProcessOrder(
	ctx context.Context,
	payload *OrderProcessingPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskProcessOrder, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) TaskProcessOrder(ctx context.Context, task *asynq.Task) error {
	var payload OrderProcessingPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	var rider *data.Rider
	distance := utils.Haversine(payload.Parcel.PickupLatitude, payload.Parcel.PickupLongitude, payload.Parcel.DropOffLatitude, payload.Parcel.DropOffLongitude)
	if distance < 3 {
		rider = data.GetAvailableCyclist(payload.Parcel.PickupLongitude, payload.Parcel.PickupLatitude)
	} else {
		rider = data.GetAvailableBiker(payload.Parcel.PickupLongitude, payload.Parcel.PickupLatitude)
	}

	if rider == nil {
		err := data.UpdateOrderStatus(payload.Order.Id, data.SEARCHINGRIDER)
		if err != nil {
			log.Error().Str("type", task.Type()).Err(err)
			return fmt.Errorf("failed to update order status: %w", err)
		}

		if err := requeueOrderProcess(&payload); err != nil {
			log.Error().Str("type", task.Type()).Err(err)
			return fmt.Errorf("failed to requeue order processing task: %w", err)
		}

		return nil
	}

	if err := data.UpdateRiderInOrder(payload.Order.Id, rider.Id); err != nil {
		log.Error().Str("type", task.Type()).Err(err)
		return fmt.Errorf("failed to update rider information: %w", err)
	}

	if err := data.UpdateOrderStatus(payload.Order.Id, data.ACCEPTED); err != nil {
		log.Error().Str("type", task.Type()).Err(err)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if err := data.UpdateRiderStatus(rider.Id, data.OnTrip); err != nil {
		log.Error().Str("type", task.Type()).Err(err)
		return fmt.Errorf("failed to update rider information: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("processed task")
	return nil
}

func requeueOrderProcess(payload *OrderProcessingPayload) error {
	config, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := NewRedisTaskDistributor(redisOpt)

	opts := []asynq.Option{
		asynq.MaxRetry(5),
		asynq.ProcessIn(1 * time.Minute),
		asynq.Queue(QueueCritical),
	}

	return taskDistributor.DistributeTaskProcessOrder(context.Background(), payload, opts...)
}
