package config

import (
	"log"
	"os"
	"strconv"
)

// AppConfig represents the application configuration params.
type AppConfig struct {
	// ServerPort is the port where the API server will
	// listen for connections.
	ServerPort int
}

// NewAppConfig loads the application configuration parameters
// and returns an instance of it.
func NewAppConfig() *AppConfig {
	serverPort, err := strconv.Atoi(os.Getenv("SERVE_PORT"))
	if err != nil || serverPort == 0 {
		log.Fatal("Missing required server port")
	}

	return &AppConfig{
		ServerPort: serverPort,
	}
}
