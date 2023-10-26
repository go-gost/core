package hosts

import (
	"context"
	"net"
)

type Options struct{}

type Option func(opts *Options)

// HostMapper is a mapping from hostname to IP.
type HostMapper interface {
	Lookup(ctx context.Context, network, host string, opts ...Option) ([]net.IP, bool)
}
