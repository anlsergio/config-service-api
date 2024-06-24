package controller

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller/dto"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller/middleware"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"log"
	"net/http"
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
	r.HandleFunc("/configs", middleware.SetJSONContent(c.list)).
		Methods(http.MethodGet)
	r.HandleFunc("/configs", middleware.SetJSONContent(c.create)).
		Methods(http.MethodPost)
	r.HandleFunc("/configs/{name}", middleware.SetJSONContent(c.get)).
		Methods(http.MethodGet)
	r.HandleFunc("/configs/{name}", middleware.SetJSONContent(c.update)).
		Methods(http.MethodPut, http.MethodPatch)
	r.HandleFunc("/configs/{name}", middleware.SetJSONContent(c.delete)).
		Methods(http.MethodDelete)
	r.HandleFunc("/search", middleware.SetJSONContent(c.query)).
		Methods(http.MethodGet)
}

// @Summary List configs
// @Description Lists all available configs
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {array} dto.Config
// @Failure 500 {string} string "Error message"
// @Router /configs [get]
func (c Config) list(w http.ResponseWriter, r *http.Request) {
	configs, err := c.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseConfigs []dto.Config
	for _, config := range configs {
		dtoConfig, err := dto.FromDomainConfig(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseConfigs = append(responseConfigs, dtoConfig)
	}

	bytes, err := json.Marshal(responseConfigs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Failed to write response: %s", err.Error())
		return
	}
}

// @Summary Create a new config
// @Description Creates a new config resource
// @Tags config
// @Accept json
// @Produce json
// @Param config body dto.Config true "Config object to be created"
// @Success 201
// @Failure 400 {object} string "Error message"
// @Failure 500 {object} string "Error message"
// @Router /configs [post]
func (c Config) create(w http.ResponseWriter, r *http.Request) {
	var requestBody dto.Config
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := requestBody.ToDomainConfig()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.Create(config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Get a config by name
// @Description Gets a config resource by its name
// @Tags config
// @Accept json
// @Produce json
// @Param name path string true "Name of the config"
// @Success 200 {object} dto.Config
// @Failure 404 {object} string "Error message"
// @Failure 500 {object} string "Error message"
// @Router /configs/{name} [get]
func (c Config) get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	config, err := c.service.Get(name)
	if err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseConfig, err := dto.FromDomainConfig(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(responseConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Failed to write response: %s", err.Error())
		return
	}
}

// @Summary Update a config by name
// @Description Updates a config resource by its name
// @Tags config
// @Accept json
// @Produce json
// @Param name path string true "Name of the config"
// @Param config body dto.Metadata true "Metadata"
// @Success 200
// @Failure 400 {object} string "Error message"
// @Failure 404 {object} string "Error message"
// @Failure 500 {object} string "Error message"
// @Router /configs/{name} [put]
// @Router /configs/{name} [patch]
func (c Config) update(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	var requestBody dto.Metadata
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metadataBytes, err := requestBody.ToByteSlice()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.service.Update(name, metadataBytes); err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a config by name
// @Description Deletes a config resource by its name
// @Tags config
// @Accept json
// @Produce json
// @Param name path string true "Name of the config"
// @Success 200
// @Failure 404 {object} string "Error message"
// @Failure 500 {object} string "Error message"
// @Router /configs/{name} [delete]
func (c Config) delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	err := c.service.Delete(name)
	if err != nil {
		if errors.Is(err, repository.ErrConfigNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Query configs based on criteria
// @Description Query all available configs based on query parameters
// @Tags config
// @Accept json
// @Produce json
// @Param keyValuePairs query object true "Metadata filters not represented appropriately, due to limitations in OpenAPI 2.x. But it's a free key/value pair of strings"
// @Success 200 {array} dto.Config
// @Failure 500 {object} string "Error message"
// @Router /search [get]
func (c Config) query(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()

	// convert the query params into map[string]string
	query := make(map[string]string)
	for k, v := range urlQuery {
		if len(v) > 0 {
			query[k] = v[0]
		}
	}

	configs, err := c.service.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseConfigs []dto.Config
	for _, config := range configs {
		dtoConfig, err := dto.FromDomainConfig(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseConfigs = append(responseConfigs, dtoConfig)
	}

	bytes, err := json.Marshal(responseConfigs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		log.Printf("Failed to write response: %s", err.Error())
		return
	}
}
