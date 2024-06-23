package domain_test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_MetadataValue(t *testing.T) {
	c := domain.Config{
		Name: test.ConfigName1,
		Metadata: []byte(`
			{
				"enabled": "true",
				"abc": "123",
				"obj": {
					"aaa": {
						"bbb": "ccc"
					}
				}
			}`),
	}

	t.Run("not matching key", func(t *testing.T) {
		got := c.MetadataValue("nope")
		assert.Empty(t, got)
	})

	t.Run("matching key", func(t *testing.T) {
		got := c.MetadataValue("enabled")
		assert.Equal(t, "true", got)
	})

	t.Run("nested matching key", func(t *testing.T) {
		got := c.MetadataValue("obj.aaa.bbb")
		assert.Equal(t, "ccc", got)
	})
}
