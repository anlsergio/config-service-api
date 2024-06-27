package controller_test

import (
	"github.com/gorilla/mux"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/controller"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	hc := controller.NewHealthCheck()

	r := mux.NewRouter()
	hc.SetRouter(r)

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "healthz endpoint returns OK",
			path:       "/healthz",
			wantStatus: http.StatusOK,
		},
		{
			name:       "readyz endpoint returns OK",
			path:       "/readyz",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}
