// Package limiter defines shared types and constants used by all limiter
// sub-packages (conn, rate, traffic).
package limiter

// Limiter scopes define how a limiter key is interpreted.
const (
	// ScopeService applies the limiter at the service level.
	ScopeService string = "service"
	// ScopeConn applies the limiter at the connection level.
	ScopeConn string = "conn"
	// ScopeClient applies the limiter per client (source IP).
	ScopeClient string = "client"
)

// Options holds the parameters passed when requesting a limiter.
type Options struct {
	// Service is the service name.
	Service string
	// Scope is the limiter scope ("service", "conn", or "client").
	Scope string
	// Network is the connection network type.
	Network string
	// Addr is the destination address.
	Addr string
	// Client is the client identifier (typically the source IP).
	Client string
	// Src is an additional source identifier.
	Src string
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// ServiceOption sets the service name.
func ServiceOption(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

// ScopeOption sets the limiter scope.
func ScopeOption(scope string) Option {
	return func(opts *Options) {
		opts.Scope = scope
	}
}

// NetworkOption sets the network type.
func NetworkOption(network string) Option {
	return func(opts *Options) {
		opts.Network = network
	}
}

// AddrOption sets the destination address.
func AddrOption(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

// ClientOption sets the client identifier.
func ClientOption(client string) Option {
	return func(opts *Options) {
		opts.Client = client
	}
}

// SrcOption sets an additional source identifier.
func SrcOption(src string) Option {
	return func(opts *Options) {
		opts.Src = src
	}
}
