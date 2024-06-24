package dto

import (
	"encoding/json"
	"fmt"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
)

// Config is the data transfer object for the config controller request and response.
type Config struct {
	// Name is the name of the config.
	Name string `json:"name,omitempty"`
	// Metadata is the arbitrary key value pairs of metadata
	// that compose a config.
	Metadata map[string]any `json:"metadata"`
}

// ToDomainConfig converts the dto.Config into a domain.Config.
func (c Config) ToDomainConfig() (domain.Config, error) {
	bytes, err := json.Marshal(c.Metadata)
	if err != nil {
		return domain.Config{}, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return domain.Config{
		Name:     c.Name,
		Metadata: bytes,
	}, nil
}

// FromDomainConfig converts the domain.Config into a dto.Config.
func FromDomainConfig(d domain.Config) (Config, error) {
	var metadata map[string]any

	err := json.Unmarshal(d.Metadata, &metadata)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return Config{
		Name:     d.Name,
		Metadata: metadata,
	}, nil
}
