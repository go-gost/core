// Package connector defines interfaces for connecting to destination addresses
// from the final proxy node in a chain. It handles the application-level
// handshake (e.g. SOCKS5 CONNECT, HTTP CONNECT) that reaches the target.
package connector

import (
	"context"
	"net"

	"github.com/go-gost/core/metadata"
)

// Connector is responsible for connecting to the destination address through
// a proxy node. It performs the application-level negotiation that reaches
// the final target, such as a SOCKS5 CONNECT request or HTTP CONNECT tunnel.
type Connector interface {
	// Init initializes the connector with metadata.
	Init(metadata.Metadata) error
	// Connect establishes a connection to the target address through the
	// given connection (which is already established to the proxy node).
	Connect(ctx context.Context, conn net.Conn, network, address string, opts ...ConnectOption) (net.Conn, error)
}

// Handshaker performs an application-level handshake on an established
// connection. This is distinct from the transport-level Handshaker in the
// dialer package — connector.Handshaker handles protocol negotiation with
// the destination (e.g. SOCKS5 greeting), while dialer.Handshaker handles
// transport setup (e.g. TLS).
type Handshaker interface {
	Handshake(ctx context.Context, conn net.Conn) (net.Conn, error)
}
