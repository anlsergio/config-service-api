package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/config"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"log"
	"net/http"
)

func main() {
	// Load the application configuration params
	cfg := config.NewAppConfig(".")

	// set the config controller handlers injecting the dependency
	// in the router
	svc := service.NewConfig(repository.NewInMemoryConfig())
	configController := controller.NewConfig(svc)

	r := mux.NewRouter()
	configController.SetRouter(r)

	// start the HTTP server
	log.Printf("Starting server on port %d", cfg.ServerPort)
	// TODO: listen for syscalls to shutdown server gracefully
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r))
}
