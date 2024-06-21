package service_test

import (
	"errors"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository/mocks"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_List(t *testing.T) {
	t.Run("listing is successful", func(t *testing.T) {
		mockRepo := mocks.NewConfig(t)
		stubs := generateConfigListStubs(t)
		mockRepo.On("List").Return(stubs, nil)

		svc := service.NewConfig(mockRepo)

		configs, err := svc.List()
		require.NoError(t, err)

		t.Run("it returns the expected number of configs", func(t *testing.T) {
			wantLen := len(stubs)
			gotLen := len(configs)

			assert.Equal(t, wantLen, gotLen)
		})
	})

	t.Run("listing returns error", func(t *testing.T) {
		mockRepo := mocks.NewConfig(t)
		stubs := generateConfigListStubs(t)
		mockRepo.On("List").Return(stubs, errors.New("oops"))

		svc := service.NewConfig(mockRepo)

		configs, err := svc.List()
		assert.Error(t, err)

		t.Run("configs list is empty", func(t *testing.T) {
			assert.Empty(t, configs)
		})
	})
}

// TODO: refactor it in a generic function
// to be used across multiple packages.
func generateConfigListStubs(t testing.TB) []domain.Config {
	t.Helper()

	return []domain.Config{
		{
			Name: "config 1",
			Metadata: []byte(`
				"foo": "bar",
				"abc": 123,
				"obj": {
					"aaa": "bbb",
				},
			`),
		},
		{
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
}
