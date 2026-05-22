// Package registry defines a generic Registry interface used throughout
// GOST for plugin-style component registration and lookup.
package registry

// Registry is a generic named container. Every component type in GOST
// (handlers, listeners, dialers, connectors, authenticators, etc.) has
// its own Registry instance. Components are registered by name (usually
// via init() side-effects and blank imports) and retrieved by name at
// config parsing time.
type Registry[T any] interface {
	// Register adds a value under the given name. Returns an error if the
	// name is already registered.
	Register(name string, v T) error
	// Unregister removes the value registered under the given name.
	Unregister(name string)
	// IsRegistered reports whether a value is registered under the name.
	IsRegistered(name string) bool
	// Get returns the value registered under the name, or the zero value
	// of T if not registered.
	Get(name string) T
	// GetAll returns all registered name-to-value mappings.
	GetAll() map[string]T
}
