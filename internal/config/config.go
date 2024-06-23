package config

import (
	"github.com/spf13/viper"
	"log"
)

// AppConfig represents the application configuration params.
type AppConfig struct {
	// ServerPort is the port where the API server will
	// listen for connections.
	ServerPort int `mapstructure:"SERVE_PORT"`
}

// checkRequired will check for missing required configuration
// and stop executing the application if there's any.
func (c AppConfig) checkRequired() {
	if c.ServerPort == 0 {
		log.Fatal("Missing required server port")
	}
}

// NewAppConfig loads the application configuration parameters
// and returns an instance of it.
func NewAppConfig(path string) *AppConfig {
	// setup viper to parse an .env file
	// containing the configuration parameters
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// automatically parse env. variables
	// which will take precedence over what's defined in
	// the config file.
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Unable to read configuration file: %s", err.Error())
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to parse configuration parameters for the application: %s", err.Error())
	}

	cfg.checkRequired()

	return &cfg
}
