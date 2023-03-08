package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"github.com/piyush1146115/parcel/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

func main() {

	sm := mux.NewRouter()

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/api/v1/parcel", handler.HandleNewParcelRequest)
	l := hclog.Default()

	// create a new server
	srv := http.Server{
		Addr:         ":9080",                                          //bindAddress,                                     // configure the bind address
		Handler:      sm,                                               // set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 // max time to write response to the client
		IdleTimeout:  120 * time.Second,                                // max time for connections using TCP Keep-Alive
	}

	ctx, cancel := context.WithCancel(context.Background())

	// start the server
	go func() {
		l.Info("Starting server...")

		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				l.Error("Error starting server", "error", err)
				os.Exit(1)
			}
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	sig := <-sigint

	log.Println("Got signal:", sig)
	l.Info("Shutting down server...")

	// First, cancel the context to signal the server to stop
	cancel()

	// Then, wait for the server to finish processing any requests
	timeout := 10 * time.Second
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown server: %v\n", err)
	}

	log.Println("Server stopped.")
}
