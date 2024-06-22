package domain

// Config represents a set of configs identified by its name.
type Config struct {
	// Name is the name of the config.
	Name string `json:"name"`
	// Metadata is the arbitrary key value pairs of metadata
	// that compose a config.
	Metadata []byte `json:"metadata"`
}
