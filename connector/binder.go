package connector

import (
	"context"
	"errors"
	"net"
)

var (
	// ErrBindUnsupported is returned when a connector does not support
	// the BIND operation (reverse tunnels).
	ErrBindUnsupported = errors.New("bind unsupported")
)

// Binder is an optional interface implemented by connectors that support
// reverse tunnels. In a reverse tunnel, the client binds a port and the
// proxy node listens for incoming connections from the destination side.
type Binder interface {
	Bind(ctx context.Context, conn net.Conn, network, address string, opts ...BindOption) (net.Listener, error)
}
