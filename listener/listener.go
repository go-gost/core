// Package listener defines the Listener interface for accepting incoming
// connections, analogous to net.Listener.
package listener

import (
	"errors"
	"net"

	"github.com/go-gost/core/metadata"
)

var (
	// ErrClosed is returned when Accept is called on a closed listener.
	ErrClosed = errors.New("accept on closed listener")
)

// Listener is a server socket that accepts incoming connections, analogous
// to net.Listener. Each listener is paired with a Handler to form a Service.
// The listener's Options carry a chain.Router so accepted connections can be
// routed upstream without an explicit handler router.
type Listener interface {
	// Init initializes the listener with metadata.
	Init(metadata.Metadata) error
	// Accept waits for and returns the next connection.
	Accept() (net.Conn, error)
	// Addr returns the listener's network address.
	Addr() net.Addr
	// Close stops the listener and releases resources.
	Close() error
}

// AcceptError wraps an error that occurred during Accept. It implements
// the net.Error interface and is always treated as temporary, preventing
// callers from giving up on transient accept failures.
type AcceptError struct {
	err error
}

// NewAcceptError wraps an error as an AcceptError.
func NewAcceptError(err error) error {
	return &AcceptError{err: err}
}

func (e *AcceptError) Error() string {
	return e.err.Error()
}

// Timeout always returns false for AcceptError.
func (e *AcceptError) Timeout() bool {
	return false
}

// Temporary always returns true for AcceptError, signaling that the caller
// should retry rather than close the listener.
func (e *AcceptError) Temporary() bool {
	return true
}

func (e *AcceptError) Unwrap() error {
	return e.err
}

// BindError wraps an error that occurred during a Bind (reverse tunnel)
// operation. Like AcceptError, it is always treated as temporary.
type BindError struct {
	err error
}

// NewBindError wraps an error as a BindError.
func NewBindError(err error) error {
	return &BindError{err: err}
}

func (e *BindError) Error() string {
	return e.err.Error()
}

// Timeout always returns false for BindError.
func (e *BindError) Timeout() bool {
	return false
}

// Temporary always returns true for BindError.
func (e *BindError) Temporary() bool {
	return true
}

func (e *BindError) Unwrap() error {
	return e.err
}
