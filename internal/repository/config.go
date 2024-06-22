package repository

import (
	"errors"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"sync"
)

var (
	// ErrConfigNotFound is used when a given config doesn't exist.
	ErrConfigNotFound = errors.New("config not found")
)

// Config is the port defining the I/O operations
// for the domain.Config resource.
//
//go:generate mockery --name Config
type Config interface {
	// List gets a list of configs.
	List() ([]domain.Config, error)
	// Save persists a new config.
	Save(cfg domain.Config) error
	// Get gets a config identified by its name.
	Get(name string) (domain.Config, error)
	// Update updates a given config, applying what's in
	// metadata to the corresponding config identified by its name.
	Update(name string, metadata []byte) error
	// Delete deletes a given config by its name.
	Delete(name string) error
	// Search fetches all configs that match the property/value combination.
	//
	// repository.Search("metadata.monitoring", "true")
	// TODO: generate Godoc example
	Search(property, value string) ([]domain.Config, error)
}

// NewInMemoryConfig returns a InMemoryConfig repository instance.
// Use InMemoryOption options to use custom settings.
func NewInMemoryConfig(opts ...InMemoryOption) Config {
	c := &InMemoryConfig{}

	// apply options sent by the user if there's any.
	for _, opt := range opts {
		opt(c)
	}

	// if no custom data is provided
	// initialize the configs datasource.
	if c.configs == nil {
		c.configs = make(map[string]domain.Config)
	}

	return c
}

// InMemoryOption defines the optional params for the
// NewInMemoryConfig constructor.
type InMemoryOption func(c *InMemoryConfig)

// WithCustomData allows one to set custom data to initialize
// the InMemoryConfig repository.
func WithCustomData(configs map[string]domain.Config) InMemoryOption {
	return func(c *InMemoryConfig) {
		c.configs = configs
	}
}

// InMemoryConfig defines the in-memory implementation of Config.
type InMemoryConfig struct {
	// used to protect the map from race conditions.
	mu sync.Mutex
	// TODO: improve the hashmap to avoid the redundant name as the key
	configs map[string]domain.Config
}

// List fetches all available configs from an in-memory datastore.
func (i *InMemoryConfig) List() ([]domain.Config, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	var configs []domain.Config

	for _, c := range i.configs {
		configs = append(configs, c)
	}

	return configs, nil
}

// Save persists a config into an in-memory datastore.
func (i *InMemoryConfig) Save(cfg domain.Config) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.configs[cfg.Name] = cfg

	return nil
}

// Get fetches a config from the in-memory datastore.
// If the resource is not found, it returns ErrConfigNotFound.
func (i *InMemoryConfig) Get(name string) (domain.Config, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	config, ok := i.configs[name]
	if !ok {
		return domain.Config{}, ErrConfigNotFound
	}

	return config, nil
}

// Update updates a config in the in-memory datastore, based on its name,
// applying what's defined in metadata.
// If the resource is not found, it returns ErrConfigNotFound.
//
// TODO: perhaps the Update method should only expect a metadata
// not the whole domain.Config object.
func (i *InMemoryConfig) Update(name string, metadata []byte) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// make sure the resource exists in the first place.
	existingConfig, ok := i.configs[name]
	if !ok {
		return ErrConfigNotFound
	}

	// preserve existing config name
	// because it's the only identifier at this point.
	existingConfig.Metadata = metadata
	i.configs[name] = existingConfig

	return nil
}

func (i *InMemoryConfig) Delete(name string) error {
	//TODO implement me
	panic("implement me")
}

func (i *InMemoryConfig) Search(property, value string) ([]domain.Config, error) {
	//TODO implement me
	panic("implement me")
}
