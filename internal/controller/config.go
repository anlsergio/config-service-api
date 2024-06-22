package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller/dto"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"log"
	"net/http"
	"strings"
)

// NewConfig creates a new Config controller instance.
// It expects a service as a dependency.
func NewConfig(svc *service.Config) *Config {
	return &Config{service: svc}
}

// Config is the config controller.
// It defines routes and handlers for the config resources.
type Config struct {
	service *service.Config
}

// SetRouter returns the router r with all the necessary routes for the
// Config controller setup.
func (c Config) SetRouter(r *mux.Router) {
	// TODO: create an object Route to better organize these
	// routes
	r.HandleFunc("/configs", c.list).
		Methods(http.MethodGet)
	r.HandleFunc("/configs", c.create).
		Methods(http.MethodPost)
	r.HandleFunc("/configs/{name}", c.get).
		Methods(http.MethodGet)
	r.HandleFunc("/configs/{name}", c.update).
		Methods(http.MethodPut, http.MethodPatch)
	r.HandleFunc("/configs/{name}", c.delete).
		Methods(http.MethodDelete)
	r.HandleFunc("/search", c.query).
		Methods(http.MethodGet)
}

func (c Config) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	configs, err := c.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// convert the configs slice into a json byte array
	// so that it can be written as a JSON response back to the requester.
	bytes, err := json.Marshal(configs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Failed to write response: %s", err.Error())
	}
}

func (c Config) create(w http.ResponseWriter, r *http.Request) {
	// TODO: create a middleware to set the content type for all handlers.
	w.Header().Set("Content-Type", "application/json")

	var requestBody dto.Config
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	config, err := dto.ToDomainConfig(requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := c.service.Create(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (c Config) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("get"))
}

func (c Config) update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("update"))
}

func (c Config) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete"))
}

func (c Config) query(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	metadataParam := struct {
		key   string
		value string
	}{}

	// extract the query param from the URL
	// if it starts with "metadata"
	for k, v := range queryParams {
		if len(v) > 0 && strings.HasPrefix(k, "metadata") {
			metadataParam.key = k
			metadataParam.value = v[0]
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("(query) key: %s, value: %s", metadataParam.key, metadataParam.value)))
}
