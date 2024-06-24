package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/hellofreshdevtests/HFtest-platform-anlsergio/api"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/config"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"log"
	"net/http"
)

// @title Config Service API
// @version 1.0
// @description A really nice description

// @contact.name Config API Support
// @contact.email foo@bar.com

// @host config-service
// @BasePath /
func main() {
	// Load the application configuration params
	cfg := config.NewAppConfig()

	// set the config controller handlers injecting the dependency
	// in the router
	svc := service.NewConfig(repository.NewInMemoryConfig())
	configController := controller.NewConfig(svc)

	r := mux.NewRouter()
	configController.SetRouter(r)

	// TODO: Swagger should have it's own controller object.
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	// start the HTTP server
	log.Printf("Starting server on port %d", cfg.ServerPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: r,
	}

	// Start up the HTTP server in a Go routine
	// to not block the execution so that the Signal listener can
	// take it from there.
	go func() {
		// When the server exits, make sure the error states that the server
		// was closed normally, meaning there's no unexpected error.
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Server is shutting down")
	}()

	// Listen to OS termination signals to allow for a graceful shutdown
	// (especially important in Kubernetes runtimes)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	// define the main context with cancel to release
	// associated resources upon shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call Shutdown for gracefully shut it down.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully shutdown complete.")
}
