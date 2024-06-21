package main

import (
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller"
	"log"
	"net/http"
)

func main() {
	configController := controller.Config{}
	r := mux.NewRouter()
	configController.SetRouter(r)

	// TODO: port should be parsed from env
	// TODO: call log.Fatal if port is not provided
	log.Println("Starting server on port 8080")
	// TODO: listen for syscalls to shutdown server gracefully
	log.Fatal(http.ListenAndServe(":8080", r))
}
