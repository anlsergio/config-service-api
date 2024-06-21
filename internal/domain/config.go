package domain

// Config represents a set of configs identified by its name.
type Config struct {
	Name     string         `json:"name"`
	Metadata map[string]any `json:"metadata"`
}
