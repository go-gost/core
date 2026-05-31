// Package stats defines the Stats interface for tracking connection and
// traffic counters (total connections, current connections, bytes, errors).
package stats

// Kind identifies the type of statistic being tracked.
type Kind int

const (
	// KindTotalConns is the total number of connections established.
	KindTotalConns Kind = 1
	// KindCurrentConns is the current number of active connections.
	KindCurrentConns Kind = 2
	// KindInputBytes is the total number of bytes received.
	KindInputBytes Kind = 3
	// KindOutputBytes is the total number of bytes sent.
	KindOutputBytes Kind = 4
	// KindTotalErrs is the total number of errors that occurred.
	KindTotalErrs Kind = 5
)

// Stats tracks connection and traffic counters. It is typically attached
// to listeners and services to monitor runtime activity.
type Stats interface {
	// Add increments the given stat by n.
	Add(kind Kind, n int64)
	// Get returns the current value for the given stat.
	Get(kind Kind) uint64
	// IsUpdated reports whether any stat has changed since the last Reset.
	IsUpdated() bool
	// Reset clears all counters to zero and resets the IsUpdated flag.
	Reset()
}
