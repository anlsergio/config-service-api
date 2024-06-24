package dto

import (
	"encoding/json"
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
