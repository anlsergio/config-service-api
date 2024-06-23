package service_test

import (
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

func TestConfig_Update(t *testing.T) {
	t.Run("update is successful", func(t *testing.T) {
		mockRepo := mocks.NewConfig(t)
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

		wantName := test.ConfigName1

		svc := service.NewConfig(mockRepo)
		err := svc.Update(wantName, []byte(`{"foo": "bar"}`))
		require.NoError(t, err)
	})
}

func TestConfig_Delete(t *testing.T) {
	t.Run("delete is successful", func(t *testing.T) {
		mockRepo := mocks.NewConfig(t)
		mockRepo.On("Delete", mock.Anything).Return(nil)
		svc := service.NewConfig(mockRepo)

		err := svc.Delete(test.ConfigName1)
		require.NoError(t, err)
	})
}

func TestConfig_Search(t *testing.T) {
	t.Run("search is successful", func(t *testing.T) {
		mockRepo := mocks.NewConfig(t)
		stubs := test.GenerateConfigListStubs(t)
		mockRepo.On("Search", mock.Anything).Return(stubs, nil)

		svc := service.NewConfig(mockRepo)

		configs, err := svc.Search(map[string]string{"foo": "bar"})
		require.NoError(t, err)

		t.Run("it returns the expected number of configs", func(t *testing.T) {
			wantLen := 2
			gotLen := len(configs)

			assert.Equal(t, wantLen, gotLen)
		})
	})
}
