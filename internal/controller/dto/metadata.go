package dto

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Metadata represents an independent representation of
// Config.Metadata, for purposes where the entire Config
// object should not be used or is misleading.
type Metadata map[string]any

// ToByteSlice converts m into bytes.
func (m Metadata) ToByteSlice() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return bytes, nil
}

// Validate returns an error ErrFailedValidation if Config.Metadata
// doesn't pass validation of the schema.
func (m Metadata) Validate() error {
	// traverse metadata to find any non-string value.
	for _, v := range m {
		if isThereNonStringMetadataKeyValues(v) {
			return errors.Join(ErrFailedValidation, errors.New("metadata value is invalid"))
		}
	}
	return nil
}

// isThereNonStringMetadataKeyValues traverses v to check if there's
// any non-string key/value pair in the data structure.
func isThereNonStringMetadataKeyValues(v any) bool {
	switch t := v.(type) {
	// in case v is already a string, we've hit the end of the branch
	case string:
		return false
	// in case it's a map, we still have a way to traverse
	case map[string]any:
		for _, v := range t {
			if isThereNonStringMetadataKeyValues(v) {
				return true
			}
		}
		return false
	// if it's neither a string or a map, then it's a non-string type.
	default:
		return true
	}
}
