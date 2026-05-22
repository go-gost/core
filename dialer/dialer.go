// Package dialer defines the Dialer interface for establishing connections
// to the next-hop proxy server in a chain.
package dialer

import (
	"context"
	"net"

	"github.com/go-gost/core/metadata"
)

// Dialer is responsible for dialing to a proxy server. It establishes the
// raw network connection to the next hop in the proxy chain. After the
// connection is established, a connector handles the application-level
// protocol negotiation with the destination.
type Dialer interface {
	// Init initializes the dialer with metadata.
	Init(metadata.Metadata) error
	// Dial connects to the proxy server at the given address.
	Dial(ctx context.Context, addr string, opts ...DialOption) (net.Conn, error)
}

// Handshaker performs a transport-level handshake on a connection, such as
// TLS or other transport-layer protocol negotiation. This runs after Dial
// and before any application-level handshake (see connector.Handshaker).
type Handshaker interface {
	Handshake(ctx context.Context, conn net.Conn, opts ...HandshakeOption) (net.Conn, error)
}

// Multiplexer reports whether the dialed connection supports multiplexing
// multiple logical streams over a single physical connection.
type Multiplexer interface {
	Multiplex() bool
}
