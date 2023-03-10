package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/piyush1146115/parcel/data"
	"github.com/piyush1146115/parcel/logger"
	"github.com/piyush1146115/parcel/worker"
	"github.com/rs/zerolog"
	"math/rand"

	"github.com/hibiken/asynq"
	"github.com/piyush1146115/parcel/config"
	"github.com/piyush1146115/parcel/handler"

	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	logger := logger.NewLogger()

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	oh := handler.NewOrderHandler(taskDistributor)

	sm := getNewRouter(oh)
	srv := getNewServer(sm, config.HTTPServerAddress)

	// start the server
	go startServer(srv, logger)

	// Run the task processor for Asynq
	go runTaskProcessor(redisOpt)

	// Update rider's location periodically
	go updateRidersLocationPeriodically(time.Second * 30)

	// Update the order statuses periodically
	go updateOrderStatusPeriodically(time.Second * 45)

	// Wait for an interrupt signal to gracefully shut down the server
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	sig := <-sigint

	log.Print("Got signal:", sig)
	logger.Info("Shutting down server...")

	// Then, wait for the server to finish processing any requests
	timeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Could not gracefully shutdown server: %v\n", err)
	}

	logger.Info("Server stopped.")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt)
	log.Info().Msg("Start task processor...")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}

func getNewRouter(oh *handler.OrderHandler) *mux.Router {
	sm := mux.NewRouter()

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/", handler.Home)
	getR.HandleFunc("/api/v1/order/status/{order_id:[0-9]+}", handler.OrderStatus)
	getR.HandleFunc("/api/v1/rider/status/{rider_id:[0-9]+}", handler.RiderStatus)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/api/v1/parcel/{customer_id:[0-9]+}", oh.NewParcelRequest)

	return sm
}

func getNewServer(sm *mux.Router, address string) http.Server {
	return http.Server{
		Addr:         address,
		Handler:      sm,                // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

}

func startServer(srv http.Server, l *logger.Logger) {
	l.Info("Starting server...")

	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}
}
func updateRidersLocationPeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		totalRiders := data.GetTotalAvailableRiders()
		id := rand.Intn(totalRiders)
		long := rand.Float64()*(180-(-180)) + (-180)
		lat := rand.Float64()*(90-(-90)) + (-90)

		if err := data.UpdateRidersLocation(id, long, lat); err != nil {
			log.Err(fmt.Errorf("failed to update riders location with id %d: %w", id, err))
		}

		log.Info().Msg(fmt.Sprintf("Updated location of rider with id: %d", id))
	}
}

func updateOrderStatusPeriodically(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		n := rand.Intn(2)
		if n == 1 {
			id := data.GetAnOrderIdWithAcceptedStatus()
			if id != -1 {
				data.UpdateOrderStatus(id, data.INPROGRESS)
			}
		} else {
			id := data.GetAnOrderIdWithInProgressStatus()
			if id != -1 {
				data.UpdateOrderStatus(id, data.COMPLETED)
			}
		}
	}
}
