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
	// ErrConfigExists is used when there's already an existing resource with the same name.
	ErrConfigExists = errors.New("config already exists")
)

var (
	db         *inMemoryDBState
	initDBOnce sync.Once
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
	Search(query map[string]string) ([]domain.Config, error)
}

// NewInMemoryConfig returns a InMemoryConfig repository instance.
// Use InMemoryOption options to use custom settings.
func NewInMemoryConfig(opts ...InMemoryOption) Config {
	c := &InMemoryConfig{
		db: getInMemoryDBState(),
	}

	// apply options sent by the user if there's any.
	for _, opt := range opts {
		opt(c)
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
		c.db.configs = configs
	}
}

// InMemoryConfig defines the in-memory implementation of Config.
type InMemoryConfig struct {
	db *inMemoryDBState
}

// List fetches all available configs from an in-memory datastore.
func (i *InMemoryConfig) List() ([]domain.Config, error) {
	i.db.lock()
	defer i.db.unlock()

	var configs []domain.Config

	for _, c := range i.db.configs {
		configs = append(configs, c)
	}

	return configs, nil
}

// Save persists a config into an in-memory datastore.
// If there's a config with the same name, it won't be allowed
// to be created returning ErrConfigExists.
func (i *InMemoryConfig) Save(cfg domain.Config) error {
	i.db.lock()
	defer i.db.unlock()

	// make sure there's no existing resource with the same name.
	_, ok := i.db.configs[cfg.Name]
	if ok {
		return ErrConfigExists
	}

	i.db.configs[cfg.Name] = cfg

	return nil
}

// Get fetches a config from the in-memory datastore.
// If the resource is not found, it returns ErrConfigNotFound.
func (i *InMemoryConfig) Get(name string) (domain.Config, error) {
	i.db.lock()
	defer i.db.unlock()

	config, ok := i.db.configs[name]
	if !ok {
		return domain.Config{}, ErrConfigNotFound
	}

	return config, nil
}

// Update updates a config in the in-memory datastore, based on its name,
// applying what's defined in metadata.
// If the resource is not found, it returns ErrConfigNotFound.
func (i *InMemoryConfig) Update(name string, metadata []byte) error {
	i.db.lock()
	defer i.db.unlock()

	// make sure the resource exists in the first place.
	existingConfig, ok := i.db.configs[name]
	if !ok {
		return ErrConfigNotFound
	}

	// preserve existing config name
	// because it's the only identifier at this point.
	existingConfig.Metadata = metadata
	i.db.configs[name] = existingConfig

	return nil
}

// Delete removes a given config from the in-memory datastore, based on its name.
func (i *InMemoryConfig) Delete(name string) error {
	i.db.lock()
	defer i.db.unlock()

	// make sure the resource exists in the first place.
	_, ok := i.db.configs[name]
	if !ok {
		return ErrConfigNotFound
	}

	delete(i.db.configs, name)

	return nil
}

// Search gets all configs from the in-memory datastore that match the key/value pairs
// in query.
func (i *InMemoryConfig) Search(query map[string]string) ([]domain.Config, error) {
	i.db.lock()
	defer i.db.unlock()

	var configs []domain.Config

	// loop through all the stored configs
	// and for every key/value pair in query
	// check if the corresponding value for the given key
	// in metadata matches the expected value.
out:
	for _, c := range i.db.configs {
		for k, v := range query {
			// remove the metadata prefix because it's redundant
			// because the method already expects the search to be made
			// in metadata.
			k = strings.TrimPrefix(k, "metadata.")
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

// inMemoryDBState holds the in-memory DB state for the lifecycle
// of the application.
type inMemoryDBState struct {
	// used to protect the map from race conditions.
	mu      sync.Mutex
	configs map[string]domain.Config
}

// lock the operation on the db until the token is released.
func (i *inMemoryDBState) lock() {
	i.mu.Lock()
}

// unlock releases the token, allowing another go routine
// to manipulate data.
func (i *inMemoryDBState) unlock() {
	i.mu.Unlock()
}

// getInMemoryDBState is a singleton to get the same
// in-memory DB state.
func getInMemoryDBState() *inMemoryDBState {
	// ensures that the db is initialized only once.
	initDBOnce.Do(func() {
		db = &inMemoryDBState{
			configs: make(map[string]domain.Config),
		}
	})

	return db
}
