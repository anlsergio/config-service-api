package service

import (
	"fmt"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
)

// NewConfig creates a new Config service instance.
func NewConfig(repo repository.Config) *Config {
	return &Config{repo}
}

// Config abstracts away the complexity of interacting
// with repositories to serve the config resources.
type Config struct {
	repo repository.Config
}

// List gets a list of configs.
func (c Config) List() ([]domain.Config, error) {
	configs, err := c.repo.List()
	if err != nil {
		return nil, fmt.Errorf("service failed to list configs: %w", err)
	}

	return configs, nil
}
