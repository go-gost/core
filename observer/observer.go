// Package observer defines the Observer interface for collecting
// observability events such as connection status and traffic statistics.
package observer

import "context"

// Options holds the initialization parameters for an Observer.
type Options struct{}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// Observer receives observability events from components. Events can
// include status changes, connection statistics snapshots, and other
// runtime diagnostics. Observers are typically used to feed data to
// monitoring dashboards or log systems.
type Observer interface {
	// Observe processes a batch of events.
	Observe(ctx context.Context, events []Event, opts ...Option) error
}

// EventType categorizes an event.
type EventType string

const (
	// EventStatus indicates a service or connection status change.
	EventStatus EventType = "status"
	// EventStats indicates a traffic statistics snapshot.
	EventStats EventType = "stats"
)

// Event is a generic observability event. Implementations carry type-specific
// data such as connection counts or byte throughput.
type Event interface {
	Type() EventType
}
