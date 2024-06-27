package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
)

var (
	// ErrFailedValidation is used when the validation check has failed.
	ErrFailedValidation = errors.New("failed validation")
)

// Config is the data transfer object for the config controller request and response.
type Config struct {
	// Name is the name of the config.
	Name string `json:"name,omitempty"`
	// Metadata is the arbitrary key value pairs of metadata
	// that compose a config.
	Metadata Metadata `json:"metadata"`
}

// Validate returns an error ErrFailedValidation if Config
// doesn't pass validation of the schema.
func (c Config) Validate() (err error) {
	if c.Name == "" {
		err = errors.Join(ErrFailedValidation, errors.New("name is required"))
	}

	if c.Metadata != nil {
		if metadataErr := c.Metadata.Validate(); metadataErr != nil {
			err = errors.Join(ErrFailedValidation, metadataErr)
		}
	}

	return err
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
