package repository_test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemoryConfig_List(t *testing.T) {
	customData := test.GenerateInMemoryTestData(t)
	repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))

	configs, err := repo.List()
	require.NoError(t, err)

	t.Run("it returns the expected number of configs", func(t *testing.T) {
		wantLen := len(customData)
		require.Equal(t, wantLen, len(configs))
	})
}

func TestInMemoryConfig_Save(t *testing.T) {
	repo := repository.NewInMemoryConfig()

	toCreateConfig := domain.Config{
		Name:     "config 1",
		Metadata: []byte(`{"foo": "bar"}`),
	}
	require.NoError(t, repo.Save(toCreateConfig))

	t.Run("created config is the expected config", func(t *testing.T) {
		config, err := repo.Get(toCreateConfig.Name)
		require.NoError(t, err)
		assert.Equal(t, toCreateConfig, config)
	})
}

func TestInMemoryConfig_Get(t *testing.T) {
	customData := test.GenerateInMemoryTestData(t)
	repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))

	t.Run("config is found", func(t *testing.T) {
		wantName := test.ConfigName1

		config, err := repo.Get(wantName)
		require.NoError(t, err)

		t.Run("it returns the expected config", func(t *testing.T) {
			assert.Equal(t, wantName, config.Name)
		})
	})

	t.Run("config is not found", func(t *testing.T) {
		config, err := repo.Get("invalid")

		t.Run("it's the expected error type", func(t *testing.T) {
			assert.ErrorIs(t, err, repository.ErrConfigNotFound)
		})

		t.Run("config object is empty", func(t *testing.T) {
			assert.Empty(t, config)
		})
	})
}

func TestInMemoryConfig_Update(t *testing.T) {
	customData := test.GenerateInMemoryTestData(t)
	repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))

	t.Run("config is updated", func(t *testing.T) {
		wantName := test.ConfigName1
		wantMetadata := []byte(`{"got": "updated!"}`)

		// extract the exact config object to be updated
		// to keep references of it for comparison purposes.
		config1 := customData[wantName]
		metadataBeforeUpdate := config1.Metadata

		require.NoError(t, repo.Update(wantName, wantMetadata))

		gotConfig, err := repo.Get(wantName)
		require.NoError(t, err)

		t.Run("metadata is updated", func(t *testing.T) {
			assert.Equal(t, wantMetadata, gotConfig.Metadata)
			assert.NotEqual(t, metadataBeforeUpdate, gotConfig.Metadata)
		})
	})

	t.Run("config not found", func(t *testing.T) {
		err := repo.Update("nope", nil)

		t.Run("not found error", func(t *testing.T) {
			assert.ErrorIs(t, err, repository.ErrConfigNotFound)
		})
	})
}
