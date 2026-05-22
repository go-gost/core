// Package hosts defines the HostMapper interface for static host-to-IP mappings.
package hosts

import (
	"context"
	"net"
)

// Options holds the initialization parameters for a HostMapper.
type Options struct{}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// HostMapper maps hostnames to IP addresses. It acts as a local hosts file,
// providing static resolution that overrides DNS for matched hostnames.
type HostMapper interface {
	// Lookup resolves a hostname to a list of IP addresses. The boolean
	// return value indicates whether a matching mapping was found — if
	// false, the caller should fall back to DNS resolution.
	Lookup(ctx context.Context, network, host string, opts ...Option) ([]net.IP, bool)
}
