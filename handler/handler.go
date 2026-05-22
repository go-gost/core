// Package handler defines the Handler interface for processing inbound
// connections — authenticating, routing, and forwarding traffic.
package handler

import (
	"context"
	"net"

	"github.com/go-gost/core/hop"
	"github.com/go-gost/core/metadata"
)

// Handler processes an inbound connection. It authenticates the client,
// selects a route through the proxy chain, and forwards traffic between
// the client and the destination. A Handler is typically paired with a
// Listener to form a Service.
type Handler interface {
	// Init initializes the handler with metadata.
	Init(metadata.Metadata) error
	// Handle processes an inbound connection. The connection is closed when
	// Handle returns.
	Handle(context.Context, net.Conn, ...HandleOption) error
}

// Forwarder is an optional interface implemented by handlers that support
// dynamic forwarding. It allows setting or changing the forwarding hop at
// runtime, which is used in reverse tunnel scenarios where the target is
// determined after the connection is established.
type Forwarder interface {
	Forward(hop.Hop)
}
