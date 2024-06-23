package domain

import (
	"encoding/json"
	"strings"
)

// Config represents a set of configs identified by its name.
type Config struct {
	// Name is the name of the config.
	Name string `json:"name"`
	// Metadata is the arbitrary key value pairs of metadata
	// that compose a config.
	Metadata []byte `json:"metadata"`
}

// MetadataValue traverses the metadata data structure
// to get the corresponding value for the key.
//
// A key is expected to have the following format: `aaa.bbb.ccc`,
// where each dot represents each key node in a different nest level.
//
// It returns nil if no matching value is found.
// TODO: generate Godoc example.
func (c Config) MetadataValue(key string) any {
	// break down the nested key format
	// into separate distinct keys.
	keys := strings.Split(key, ".")

	// it should be represented as map[string]any because
	// of the unstructured nature of the key/value pairs
	// expected.
	var m map[string]any
	if err := json.Unmarshal(c.Metadata, &m); err != nil {
		return nil
	}

	return traverseAndFind(keys, m)
}

// traverseAndFind takes in a key slice representing a nested key structure
// and the key value data that it's trying to match.
func traverseAndFind(keys []string, data any) any {
	// use type cast to know if the current data
	// is a key/value pair or if it's the final value.
	switch t := data.(type) {
	case map[string]any:
		data, ok := t[keys[0]]
		if !ok {
			// if it doesn't match any key in the current level
			// it's safe to assume the key doesn't match.
			return nil
		}
		// In the case that there's a corresponding data matching
		// the current key, proceed to traverse for the next keys
		// to make sure the entire set of nested keys is considered.
		return traverseAndFind(keys[1:], data)
	default:
		return data
	}
}
