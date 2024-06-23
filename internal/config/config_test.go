package config_test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewAppConfig(t *testing.T) {
	t.Run("server port is populated", func(t *testing.T) {
		os.Setenv("SERVE_PORT", "8080")
		defer os.Unsetenv("SERVE_PORT")

		cfg := config.NewAppConfig()

		assert.Equal(t, 8080, cfg.ServerPort)
	})
}
