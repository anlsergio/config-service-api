package repository

import "github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"

// Config is the port defining the I/O operations
// for the domain.Config resource.
//
//go:generate mockery --name Config
type Config interface {
	// List gets a list of configs.
	List() (config []domain.Config, error error)
	// Save persists a new config.
	Save(config []domain.Config) error
	// Get gets a config identified by its name.
	Get(name string) (domain.Config, error)
	// Update updates a given domain, applying what's in
	// form to the corresponding config identified by its name.
	Update(name string, form domain.Config) error
	// Delete deletes a given config by its name.
	Delete(name string) error
	// Search fetches all configs that match the property/value combination.
	//
	// repository.Search("metadata.monitoring", "true")
	// TODO: generate Godoc example
	Search(property, value string) ([]domain.Config, error)
}

// NewInMemoryConfig returns a InMemoryConfig repository instance.
func NewInMemoryConfig() Config {
	return &InMemoryConfig{
		configs: make(map[string]domain.Config),
	}
}

// InMemoryConfig defines the in-memory implementation of Config.
type InMemoryConfig struct {
	configs map[string]domain.Config
}

func (i *InMemoryConfig) List() (config []domain.Config, error error) {
	//TODO implement me
	panic("implement me")
}

func (i *InMemoryConfig) Save(config []domain.Config) error {
	//TODO implement me
	panic("implement me")
}

func (i *InMemoryConfig) Get(name string) (domain.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (i *InMemoryConfig) Update(name string, form domain.Config) error {
	//TODO implement me
	panic("implement me")
}

func (i *InMemoryConfig) Delete(name string) error {
	//TODO implement me
	panic("implement me")
}

func (i *InMemoryConfig) Search(property, value string) ([]domain.Config, error) {
	//TODO implement me
	panic("implement me")
}
