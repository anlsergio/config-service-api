package controller_test

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfig(t *testing.T) {
	configController := controller.NewConfig()

	r := mux.NewRouter()
	configController.SetRouter(r)

	t.Run("list configs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/configs", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("create config", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/configs", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("get config", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("update config using the PUT HTTP verb", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("update config using the PATCH HTTP verb", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("delete config", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/configs/foo", nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("search configs using query params", func(t *testing.T) {
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
