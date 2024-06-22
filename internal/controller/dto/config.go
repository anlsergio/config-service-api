package dto

import (
	"encoding/json"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
)

// Config is the data transfer object for the config controller request and response.
// TODO: there's should be a validation method to validate if the metadata is a valid JSON.
type Config struct {
	// Name is the name of the config.
	Name string `json:"name,omitempty"`
	// Metadata is the arbitrary key value pairs of metadata
	// that compose a config.
	Metadata any `json:"metadata"`
}

// ToDomainConfig converts the dto.Config into a domain.Config.
func ToDomainConfig(c Config) (domain.Config, error) {
	bytes, err := json.Marshal(c.Metadata)
	if err != nil {
		return domain.Config{}, err
	}

	return domain.Config{
		Name:     c.Name,
		Metadata: bytes,
	}, nil
}
