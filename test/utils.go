package test

import (
	"github.com/hellofreshdevtests/HFtest-platform-anlsergio/internal/domain"
	"testing"
)

const (
	// ConfigName1 represents a config name present from GenerateConfigListStubs.
	ConfigName1 = "config1"
	// ConfigName2 represents a config name present from GenerateConfigListStubs.
	ConfigName2 = "config2"
	// ConfigName3 represents a config name present from GenerateConfigListStubs
	ConfigName3 = "burger-nutrition"
)

// GenerateConfigListStubs generates a prepopulated list of configs
// to serve as test stubs.
func GenerateConfigListStubs(t testing.TB) []domain.Config {
	t.Helper()

	return []domain.Config{
		{
			Name: ConfigName1,
			Metadata: []byte(`
				"foo": "bar",
				"abc": 123,
				"obj": {
					"aaa": "bbb",
				},
			`),
		},
		{
			Name: ConfigName2,
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

// GenerateInMemoryTestData generates custom test data to populate the
// InMemory repository.
func GenerateInMemoryTestData(t testing.TB) map[string]domain.Config {
	t.Helper()

	return map[string]domain.Config{
		ConfigName1: {
			Name: ConfigName1,
			Metadata: []byte(`
			{
				"foo": "bar",
				"abc": "123",
				"obj": {
					"aaa": "bbb"
				}
			}`),
		},
		ConfigName2: {
			Name: ConfigName2,
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
		},
		ConfigName3: {
			Name: ConfigName3,
			Metadata: []byte(`{
				"calories": "230",
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
			`),
		},
	}
}
