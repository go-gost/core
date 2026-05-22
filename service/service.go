// Package service defines the Service interface — the top-level abstraction
// that combines a Listener and a Handler to form a running proxy service.
package service

import (
	"net"
)

// Service is the top-level component in GOST. It wraps one Listener and one
// Handler: the Listener accepts inbound connections, and the Handler processes
// them by authenticating, routing, and forwarding traffic through a proxy chain.
// A running GOST instance is essentially a collection of Services.
type Service interface {
	// Serve starts the service and blocks until it is stopped or encounters
	// a fatal error.
	Serve() error
	// Addr returns the service's network address.
	Addr() net.Addr
	// Close stops the service and releases all resources.
	Close() error
}
