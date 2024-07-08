package net

import (
	"context"
	"net"
)

type Dialer interface {
	Dial(ctx context.Context, network, addr string) (net.Conn, error)
}
