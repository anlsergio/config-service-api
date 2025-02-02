package service

import (
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
	return c.repo.List()
}

// Create creates a new config according to cfg.
func (c Config) Create(cfg domain.Config) error {
	return c.repo.Save(cfg)
}

// Get gets a config identified by its name.
func (c Config) Get(name string) (domain.Config, error) {
	return c.repo.Get(name)
}

// Update updates the config identified by name applying whatever is in metadata.
func (c Config) Update(name string, metadata []byte) error {
	return c.repo.Update(name, metadata)
}

// Delete removes the config identified by name.
func (c Config) Delete(name string) error {
	return c.repo.Delete(name)
}

// Search gets a list of configs matching the query key/value pairs,
// where key represents the nested property in metadata, and value is the
// value that should match.
func (c Config) Search(query map[string]string) ([]domain.Config, error) {
	return c.repo.Search(query)
}
