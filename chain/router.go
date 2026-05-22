package chain

import (
	"context"
	"net"
	"time"

	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/core/recorder"
	"github.com/go-gost/core/resolver"
)

// SockOpts holds socket-level options applied to connections.
type SockOpts struct {
	// Mark is the SO_MARK value for the socket.
	Mark int
}

// RouterOptions holds the initialization parameters for a Router.
type RouterOptions struct {
	// Retries is the number of dial retries on failure.
	Retries int
	// Timeout is the dial timeout.
	Timeout time.Duration
	// IfceName is the network interface name to bind to.
	IfceName string
	// Netns is the network namespace name.
	Netns string
	// SockOpts holds socket options for connections.
	SockOpts *SockOpts
	// Chain is the Chainer that selects routes.
	Chain Chainer
	// Resolver resolves domain names to IP addresses.
	Resolver resolver.Resolver
	// HostMapper maps hostnames to static IP addresses.
	HostMapper hosts.HostMapper
	// Recorders records traffic data.
	Recorders []recorder.RecorderObject
	// Logger is the logger for router operations.
	Logger logger.Logger
}

// RouterOption is a functional option for configuring RouterOptions.
type RouterOption func(*RouterOptions)

// InterfaceRouterOption sets the network interface.
func InterfaceRouterOption(ifceName string) RouterOption {
	return func(o *RouterOptions) {
		o.IfceName = ifceName
	}
}

// NetnsRouterOption sets the network namespace.
func NetnsRouterOption(netns string) RouterOption {
	return func(o *RouterOptions) {
		o.Netns = netns
	}
}

// SockOptsRouterOption sets the socket options.
func SockOptsRouterOption(so *SockOpts) RouterOption {
	return func(o *RouterOptions) {
		o.SockOpts = so
	}
}

// TimeoutRouterOption sets the dial timeout.
func TimeoutRouterOption(timeout time.Duration) RouterOption {
	return func(o *RouterOptions) {
		o.Timeout = timeout
	}
}

// RetriesRouterOption sets the number of dial retries.
func RetriesRouterOption(retries int) RouterOption {
	return func(o *RouterOptions) {
		o.Retries = retries
	}
}

// ChainRouterOption sets the Chainer.
func ChainRouterOption(chain Chainer) RouterOption {
	return func(o *RouterOptions) {
		o.Chain = chain
	}
}

// ResolverRouterOption sets the Resolver.
func ResolverRouterOption(resolver resolver.Resolver) RouterOption {
	return func(o *RouterOptions) {
		o.Resolver = resolver
	}
}

// HostMapperRouterOption sets the HostMapper.
func HostMapperRouterOption(m hosts.HostMapper) RouterOption {
	return func(o *RouterOptions) {
		o.HostMapper = m
	}
}

// RecordersRouterOption sets the traffic Recorders.
func RecordersRouterOption(recorders ...recorder.RecorderObject) RouterOption {
	return func(o *RouterOptions) {
		o.Recorders = recorders
	}
}

// LoggerRouterOption sets the Logger.
func LoggerRouterOption(logger logger.Logger) RouterOption {
	return func(o *RouterOptions) {
		o.Logger = logger
	}
}

// Router is the top-level routing interface used by both Listeners and Handlers
// for upstream communication. It combines chain-based routing with DNS
// resolution, host mapping, and retry/timeout logic.
type Router interface {
	// Options returns the router's configuration.
	Options() *RouterOptions
	// Dial connects to the target address through the configured chain,
	// applying resolution, host mapping, recording, and retry logic.
	Dial(ctx context.Context, network, address string, opts ...DialOption) (net.Conn, error)
	// Bind creates a reverse listener through the configured chain.
	Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error)
}
