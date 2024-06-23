package repository_test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/test"
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

func TestInMemoryConfig_Delete(t *testing.T) {
	customData := test.GenerateInMemoryTestData(t)
	repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))

	t.Run("config is deleted", func(t *testing.T) {
		configs, err := repo.List()
		require.NoError(t, err)

		// the expected number of configs available
		// should drop by 1.
		wantLen := len(configs) - 1

		require.NoError(t, repo.Delete(test.ConfigName1))

		t.Run("it returns the expected number of configs", func(t *testing.T) {
			updatedConfigList, err := repo.List()
			require.NoError(t, err)

			assert.Len(t, updatedConfigList, wantLen)
		})
	})

	t.Run("config not found", func(t *testing.T) {
		err := repo.Delete("nope")

		t.Run("not found error", func(t *testing.T) {
			assert.ErrorIs(t, err, repository.ErrConfigNotFound)
		})
	})
}

func TestInMemoryConfig_Search(t *testing.T) {
	customData := test.GenerateInMemoryTestData(t)
	repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))

	tests := []struct {
		name           string
		query          map[string]string
		wantConfigsLen int
	}{
		{
			name:           "matching configs are found",
			query:          map[string]string{"abc": "123"},
			wantConfigsLen: 2,
		},
		{
			name:           "only one matching config is found",
			query:          map[string]string{"enabled": "true"},
			wantConfigsLen: 1,
		},
		{
			name:           "no matching config is found",
			query:          map[string]string{"enabled": "false"},
			wantConfigsLen: 0,
		},
		{
			name:           "only one matching config with nested keys is found",
			query:          map[string]string{"allergens.eggs": "true"},
			wantConfigsLen: 1,
		},
		{
			name: "only 1 config matching multiple key/value pairs",
			query: map[string]string{
				"abc":         "123",
				"obj.aaa.bbb": "ccc",
			},
			wantConfigsLen: 1,
		},
		{
			name: "not corresponding match multiple key/value pairs",
			query: map[string]string{
				"abc":         "123",
				"obj.aaa.bbb": "ccc",
				"enabled":     "false",
			},
			wantConfigsLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundConfigs, err := repo.Search(tt.query)
			require.NoError(t, err)

			assert.Len(t, foundConfigs, tt.wantConfigsLen)
		})
	}
}
