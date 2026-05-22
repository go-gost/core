// Package net provides a minimal shared Dialer interface used across the
// connector and dialer packages to avoid circular dependencies.
package net

import (
	"context"
	"net"
)

// Dialer is a minimal interface for establishing network connections.
// It is used as a shared dependency by connector and dialer option types
// rather than depending directly on the full dialer module.
type Dialer interface {
	Dial(ctx context.Context, network, addr string) (net.Conn, error)
}
