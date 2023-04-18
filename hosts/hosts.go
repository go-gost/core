package hosts

import (
	"context"
	"net"
)

// HostMapper is a mapping from hostname to IP.
type HostMapper interface {
	Lookup(ctx context.Context, network, host string) ([]net.IP, bool)
}
