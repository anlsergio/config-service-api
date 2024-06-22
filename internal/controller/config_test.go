package controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository/mocks"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("list configs", func(t *testing.T) {
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

		t.Run("listing is successful", func(t *testing.T) {
			repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))
			svc := service.NewConfig(repo)
			configController := controller.NewConfig(svc)

			r := mux.NewRouter()
			configController.SetRouter(r)

			req := httptest.NewRequest(http.MethodGet, "/configs", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("http status is OK", func(t *testing.T) {
				assert.Equal(t, http.StatusOK, rr.Code)
			})

			t.Run("it returns the expected number of configs", func(t *testing.T) {
				wantLen := len(customData)

				var responseConfigs []domain.Config
				err := json.Unmarshal(rr.Body.Bytes(), &responseConfigs)
				require.NoError(t, err)

				assert.Equal(t, wantLen, len(responseConfigs))
			})
		})

		t.Run("service errors out", func(t *testing.T) {
			mockRepo := mocks.NewConfig(t)
			mockRepo.On("List").Return(nil, errors.New("oops"))

			svc := service.NewConfig(mockRepo)
			configController := controller.NewConfig(svc)

			r := mux.NewRouter()
			configController.SetRouter(r)

			req := httptest.NewRequest(http.MethodGet, "/configs", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("status is InternalServerError", func(t *testing.T) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			})
		})
	})

	t.Run("create config", func(t *testing.T) {
		t.Run("creation is successful", func(t *testing.T) {
			repo := repository.NewInMemoryConfig()
			svc := service.NewConfig(repo)
			configController := controller.NewConfig(svc)

			r := mux.NewRouter()
			configController.SetRouter(r)

			requestBody := `
				{
					"name": "burger-nutrition",
					"metadata": {
					  "calories": 230,
					  "fats": {
						"saturated-fat": "0g",
						"trans-fat": "1g"
					  },
					  "carbohydrates": {
						  "dietary-fiber": "4g",
						  "sugars": "1g"
					  },
					  "allergens": {
						"nuts": "false",
						"seafood": "false",
						"eggs": "true"
					  }
					}
				}
`

			req := httptest.NewRequest(http.MethodPost, "/configs", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("http status is created", func(t *testing.T) {
				assert.Equal(t, http.StatusCreated, rr.Code)
			})
		})

		t.Run("invalid request body", func(t *testing.T) {
			repo := repository.NewInMemoryConfig()
			svc := service.NewConfig(repo)
			configController := controller.NewConfig(svc)

			r := mux.NewRouter()
			configController.SetRouter(r)

			requestBody := `
				{
					"name": "burger-nutrition",
					"metadata": {
					  "calories": 230,
					  "fats": {
						"saturated-fat": "0g",
						"trans-fat": "1g"
					  },
				}
`

			req := httptest.NewRequest(http.MethodPost, "/configs", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("http status is bad request", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})
		})

		t.Run("service errors out", func(t *testing.T) {
			mockRepo := mocks.NewConfig(t)
			mockRepo.On("Save", mock.Anything).Return(errors.New("oops"))

			svc := service.NewConfig(mockRepo)
			configController := controller.NewConfig(svc)

			r := mux.NewRouter()
			configController.SetRouter(r)

			requestBody := `
				{
					"name": "burger-nutrition",
					"metadata": {
					  "calories": 230,
					  "fats": {
						"saturated-fat": "0g",
						"trans-fat": "1g"
					  },
					  "carbohydrates": {
						  "dietary-fiber": "4g",
						  "sugars": "1g"
					  },
					  "allergens": {
						"nuts": "false",
						"seafood": "false",
						"eggs": "true"
					  }
					}
				}
`

			req := httptest.NewRequest(http.MethodPost, "/configs", strings.NewReader(requestBody))
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("http status is internal server error", func(t *testing.T) {
				assert.Equal(t, http.StatusInternalServerError, rr.Code)
			})
		})
	})

	t.Run("get config", func(t *testing.T) {
		configController := controller.Config{}

		r := mux.NewRouter()
		configController.SetRouter(r)

		req := httptest.NewRequest(http.MethodGet, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("update config using the PUT HTTP verb", func(t *testing.T) {
		configController := controller.Config{}

		r := mux.NewRouter()
		configController.SetRouter(r)

		req := httptest.NewRequest(http.MethodPut, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("update config using the PATCH HTTP verb", func(t *testing.T) {
		configController := controller.Config{}

		r := mux.NewRouter()
		configController.SetRouter(r)

		req := httptest.NewRequest(http.MethodPatch, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("delete config", func(t *testing.T) {
		configController := controller.Config{}

		r := mux.NewRouter()
		configController.SetRouter(r)

		req := httptest.NewRequest(http.MethodDelete, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("search configs using query params", func(t *testing.T) {
		configController := controller.Config{}

		r := mux.NewRouter()
		configController.SetRouter(r)

		wantKey := "metadata.foo"
		wantValue := "bar"

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/search?%s=%s", wantKey, wantValue), nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		assert.Contains(t, rr.Body.String(), wantKey)
		assert.Contains(t, rr.Body.String(), wantValue)
	})
}
