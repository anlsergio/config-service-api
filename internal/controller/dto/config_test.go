package dto_test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller/dto"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_ToDomainConfig(t *testing.T) {
	dtoConfig := dto.Config{
		Name:     "config name",
		Metadata: map[string]any{"foo": "bar"},
	}

	domainConfig, err := dtoConfig.ToDomainConfig()
	require.NoError(t, err)

	assert.Equal(t, "config name", domainConfig.Name)
	assert.Contains(t, string(domainConfig.Metadata), "foo")
}

func TestFromDomainConfig(t *testing.T) {
	domainConfig := domain.Config{
		Name:     "config name",
		Metadata: []byte(`{"foo":"bar"}`),
	}

	dtoConfig, err := dto.FromDomainConfig(domainConfig)
	require.NoError(t, err)

	assert.Equal(t, "config name", dtoConfig.Name)
	assert.Equal(t, "bar", dtoConfig.Metadata["foo"])
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  dto.Config
		wantErr error
	}{
		{
			name: "valid config",
			config: dto.Config{
				Name:     "config name",
				Metadata: map[string]any{"foo": "bar"},
			},
			wantErr: nil,
		},
		{
			name: "valid config with nested metadata",
			config: dto.Config{
				Name: "config name",
				Metadata: map[string]any{
					"nested": map[string]any{
						"foo": "bar",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "missing config name",
			config: dto.Config{
				Name:     "",
				Metadata: map[string]any{"foo": "bar"},
			},
			wantErr: dto.ErrFailedValidation,
		},
		{
			name: "invalid metadata",
			config: dto.Config{
				Name:     "config name",
				Metadata: map[string]any{"foo": 8},
			},
			wantErr: dto.ErrFailedValidation,
		},
		{
			name: "nested non-string metadata value",
			config: dto.Config{
				Name: "config name",
				Metadata: map[string]any{
					"nest": map[string]any{
						"foo": 8,
					},
				},
			},
			wantErr: dto.ErrFailedValidation,
		},
		{
			name: "nested non-string metadata key",
			config: dto.Config{
				Name: "config name",
				Metadata: map[string]any{
					"nest": map[any]any{
						8: "hey",
					},
				},
			},
			wantErr: dto.ErrFailedValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
