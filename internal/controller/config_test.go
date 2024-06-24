package controller_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller/dto"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/repository/mocks"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/service"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("list configs", func(t *testing.T) {
		customData := test.GenerateInMemoryTestData(t)

		t.Run("listing is successful", func(t *testing.T) {
			repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))
			svc := service.NewConfig(repo)
			configController := controller.NewConfig(svc)

			r := mux.NewRouter()
			configController.SetRouter(r)

			req := httptest.NewRequest(http.MethodGet, "/configs", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			var responseConfigs []dto.Config
			err := json.Unmarshal(rr.Body.Bytes(), &responseConfigs)
			require.NoError(t, err)

			t.Run("http status is OK", func(t *testing.T) {
				assert.Equal(t, http.StatusOK, rr.Code)
			})

			t.Run("it returns the expected number of configs", func(t *testing.T) {
				wantLen := len(customData)

				assert.Equal(t, wantLen, len(responseConfigs))
			})

			t.Run("metadata is properly serialized", func(t *testing.T) {
				for _, c := range responseConfigs {
					assert.NotEmpty(t, c.Metadata)
				}
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
		customData := test.GenerateInMemoryTestData(t)
		repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))
		svc := service.NewConfig(repo)
		configController := controller.NewConfig(svc)

		r := mux.NewRouter()
		configController.SetRouter(r)

		t.Run("gets config successfully", func(t *testing.T) {
			wantName := test.ConfigName1

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/configs/%s", wantName), nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("http status is OK", func(t *testing.T) {
				assert.Equal(t, http.StatusOK, rr.Code)
			})

			t.Run("corresponding config is present in the response", func(t *testing.T) {
				assert.Contains(t, rr.Body.String(), wantName)
			})
		})

		t.Run("config not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/configs/nope", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			t.Run("status is not found", func(t *testing.T) {
				assert.Equal(t, http.StatusNotFound, rr.Code)
			})
		})
	})

	t.Run("update config", func(t *testing.T) {
		fineRequestBody := `
			{
				"metadata": {
				  "calories": 230,
				  "fats": {
					"saturated-fat": "0g",
					"trans-fat": "1g"
				  }
				}
			}
			`

		invalidRequestBody := `
			{
				"metadata": {
				  "calories": 230,
				  "fats": {
					"saturated-fat": "0g",
					"trans-fat": "1g"
				  }
			}
			`

		// using table tests approach to avoid test code duplication
		tests := []struct {
			name           string
			configName     string
			method         string
			requestBody    io.Reader
			wantHTTPStatus int
		}{
			{
				name:           "using the PUT HTTP verb",
				configName:     test.ConfigName1,
				method:         http.MethodPut,
				requestBody:    strings.NewReader(fineRequestBody),
				wantHTTPStatus: http.StatusOK,
			},
			{
				name:           "using the PATCH HTTP verb",
				configName:     test.ConfigName1,
				method:         http.MethodPatch,
				requestBody:    strings.NewReader(fineRequestBody),
				wantHTTPStatus: http.StatusOK,
			},
			{
				name:           "invalid request body",
				configName:     test.ConfigName1,
				method:         http.MethodPut,
				requestBody:    strings.NewReader(invalidRequestBody),
				wantHTTPStatus: http.StatusBadRequest,
			},
			{
				name:           "resource doesn't exist",
				configName:     "nope",
				method:         http.MethodPut,
				requestBody:    strings.NewReader(fineRequestBody),
				wantHTTPStatus: http.StatusNotFound,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				customData := test.GenerateInMemoryTestData(t)
				repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))
				svc := service.NewConfig(repo)
				configController := controller.NewConfig(svc)

				r := mux.NewRouter()
				configController.SetRouter(r)

				req := httptest.NewRequest(tt.method,
					fmt.Sprintf("/configs/%s", tt.configName),
					tt.requestBody)
				rr := httptest.NewRecorder()
				r.ServeHTTP(rr, req)

				t.Run("https status is the expected one", func(t *testing.T) {
					assert.Equal(t, tt.wantHTTPStatus, rr.Code)
				})
			})
		}
	})

	t.Run("delete config", func(t *testing.T) {
		customData := test.GenerateInMemoryTestData(t)
		repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))
		svc := service.NewConfig(repo)
		configController := controller.NewConfig(svc)

		r := mux.NewRouter()
		configController.SetRouter(r)

		t.Run("delete successfully", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/configs/%s", test.ConfigName1), nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
		})

		t.Run("config not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/configs/nope", nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusNotFound, rr.Code)
		})
	})

	t.Run("search configs using query params", func(t *testing.T) {
		customData := test.GenerateInMemoryTestData(t)

		repo := repository.NewInMemoryConfig(repository.WithCustomData(customData))
		svc := service.NewConfig(repo)
		configController := controller.NewConfig(svc)

		r := mux.NewRouter()
		configController.SetRouter(r)

		target := fmt.Sprintf("/search?metadata.allergens.eggs=true&metadata.fats.saturated-fat=0g")
		req := httptest.NewRequest(http.MethodGet, target, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		t.Run("http status is OK", func(t *testing.T) {
			assert.Equal(t, http.StatusOK, rr.Code)
		})

		t.Run("it returns the expected number of configs", func(t *testing.T) {
			wantLen := 1

			var responseConfigs []dto.Config
			err := json.Unmarshal(rr.Body.Bytes(), &responseConfigs)
			require.NoError(t, err)

			assert.Equal(t, wantLen, len(responseConfigs))
		})
	})
}
