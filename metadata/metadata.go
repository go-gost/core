// Package metadata defines a key-value metadata store attached to every
// component in the GOST framework. Components that carry metadata implement
// the Metadatable interface.
package metadata

// Metadatable marks a type as carrying Metadata. It is implemented by
// components like Node and HandleOptions to expose their metadata.
type Metadatable interface {
	Metadata() Metadata
}

// Metadata is a key-value store for arbitrary configuration or contextual
// data. It is used throughout GOST to pass settings that don't map to
// explicit struct fields. Keys are typically string constants defined
// alongside the parsing code.
type Metadata interface {
	// IsExists reports whether a key is present.
	IsExists(key string) bool
	// Set stores a value for the given key.
	Set(key string, value any)
	// Get retrieves the value for the given key, or nil if not present.
	Get(key string) any
}
