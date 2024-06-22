package service_test

import (
	"errors"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository/mocks"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_List(t *testing.T) {
	t.Run("listing is successful", func(t *testing.T) {
		mockRepo := mocks.NewConfig(t)
		stubs := test.GenerateConfigListStubs(t)
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
		stubs := test.GenerateConfigListStubs(t)
		mockRepo.On("List").Return(stubs, errors.New("oops"))

		svc := service.NewConfig(mockRepo)

		configs, err := svc.List()
		assert.Error(t, err)

		t.Run("configs list is empty", func(t *testing.T) {
			assert.Empty(t, configs)
		})
	})
}

func TestConfig_Create(t *testing.T) {
	t.Run("creation is successful", func(t *testing.T) {
		toCreateConfig := domain.Config{
			Name:     "config 1",
			Metadata: []byte(`{"foo": "bar"}`),
		}

		mockRepo := mocks.NewConfig(t)
		mockRepo.On("Save", mock.Anything).
			Return(func(config domain.Config) error {
				assert.Equal(t, toCreateConfig, config)
				return nil
			})

		svc := service.NewConfig(mockRepo)
		require.NoError(t, svc.Create(toCreateConfig))
	})
}

func TestConfig_Get(t *testing.T) {
	t.Run("get is successful", func(t *testing.T) {
		mockRepo := mocks.NewConfig(t)
		wantName := test.ConfigName1
		mockRepo.On("Get", mock.Anything).
			Return(domain.Config{
				Name:     wantName,
				Metadata: []byte(`{"foo": "bar"}`),
			}, nil)

		svc := service.NewConfig(mockRepo)
		config, err := svc.Get(wantName)
		require.NoError(t, err)

		t.Run("returned config match expected name", func(t *testing.T) {
			assert.Equal(t, wantName, config.Name)
		})
	})
}
