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
