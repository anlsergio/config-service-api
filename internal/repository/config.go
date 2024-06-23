package repository

import (
	"errors"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"strings"
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
	// Search fetches all configs that match the key/value combination in query.
	//
	// repository.Search(map[string]string{"metadata.monitoring", "true"})
	// TODO: generate Godoc example
	Search(query map[string]string) ([]domain.Config, error)
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

// Delete removes a given config from the in-memory datastore, based on its name.
func (i *InMemoryConfig) Delete(name string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// make sure the resource exists in the first place.
	_, ok := i.configs[name]
	if !ok {
		return ErrConfigNotFound
	}

	delete(i.configs, name)

	return nil
}

// Search gets all configs from the in-memory datastore that match the key/value pairs
// in query.
func (i *InMemoryConfig) Search(query map[string]string) ([]domain.Config, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	var configs []domain.Config

	// loop through all the stored configs
	// and for every key/value pair in query
	// check if the corresponding value for the given key
	// in metadata matches the expected value.
out:
	for _, c := range i.configs {
		for k, v := range query {
			// remove the metadata prefix because it's redundant
			// because the method already expects the search to be made
			// in metadata.
			k = strings.TrimPrefix(k, "metadata.")
			// TODO: extract a function for better code readability.
			foundValue := c.MetadataValue(k)
			// if any of the key/value combinations
			// doesn't find a match, skip adding the corresponding config.
			if foundValue == nil || foundValue.(string) != v {
				continue out
			}
		}
		// if the current config passes all query validations
		// added to the list.
		configs = append(configs, c)
	}

	return configs, nil
}
