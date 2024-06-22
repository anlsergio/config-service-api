package middleware_test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetJSONContent(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Howdy!"))
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	middleware.SetJSONContent(handler)(rr, req)

	t.Run("content type is set", func(t *testing.T) {
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("http status is ok", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("it contains a response body", func(t *testing.T) {
		assert.Equal(t, "Howdy!", rr.Body.String())
	})
}
