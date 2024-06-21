package repository_test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInMemoryConfig_List(t *testing.T) {
	customData := map[string]domain.Config{
		"config 1": {
			Name: "config 1",
			Metadata: []byte(`
				"foo": "bar",
				"abc": 123,
				"obj": {
					"aaa": "bbb",
				},
			`),
		},
		"config 2": {
			Name: "config 2",
			Metadata: []byte(`
				"enabled": "true",
				"abc": 123,
				"obj": {
					"aaa": {
						"bbb": "ccc"
					},
				},
			`),
		},
	}

	repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))
	configs, err := repo.List()
	require.NoError(t, err)

	t.Run("it returns the expected number of configs", func(t *testing.T) {
		wantLen := len(customData)
		require.Equal(t, wantLen, len(configs))
	})
}
