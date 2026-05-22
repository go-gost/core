// Package sd defines the SD (Service Discovery) interface for registering
// and discovering services in dynamic environments.
package sd

import (
	"context"
)

// Options holds the initialization parameters for an SD operation.
type Options struct{}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// Service represents a discovered service instance.
type Service struct {
	// ID is the unique instance identifier.
	ID string
	// Name is the logical service name.
	Name string
	// Node is the node (host) the service runs on.
	Node string
	// Network is the network type, e.g. "tcp", "udp".
	Network string
	// Address is the service address (host:port).
	Address string
}

// SD is the service discovery interface. It allows GOST services to register
// themselves with a service registry (Consul, etcd, DNS-SD, etc.) and
// discover other services at runtime. This enables dynamic proxy
// configurations where backend endpoints can come and go.
type SD interface {
	// Register adds a service instance to the registry.
	Register(ctx context.Context, service *Service, opts ...Option) error
	// Deregister removes a service instance from the registry.
	Deregister(ctx context.Context, service *Service) error
	// Renew refreshes the service registration lease.
	Renew(ctx context.Context, service *Service) error
	// Get returns all instances of a named service.
	Get(ctx context.Context, name string) ([]*Service, error)
}
