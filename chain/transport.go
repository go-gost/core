package chain

import (
	"context"
	"net"

	"github.com/go-gost/core/connector"
)

// TransportOptions holds the runtime parameters for a Transporter.
type TransportOptions struct {
	// Addr is the address to dial.
	Addr string
	// IfceName is the network interface to bind to.
	IfceName string
	// Netns is the network namespace name.
	Netns string
	// SockOpts holds socket options.
	SockOpts *SockOpts
	// Route is the route this transporter belongs to.
	Route Route
}

// TransportOption is a functional option for configuring TransportOptions.
type TransportOption func(*TransportOptions)

// AddrTransportOption sets the dial address.
func AddrTransportOption(addr string) TransportOption {
	return func(o *TransportOptions) {
		o.Addr = addr
	}
}

// InterfaceTransportOption sets the network interface.
func InterfaceTransportOption(ifceName string) TransportOption {
	return func(o *TransportOptions) {
		o.IfceName = ifceName
	}
}

// NetnsTransportOption sets the network namespace.
func NetnsTransportOption(netns string) TransportOption {
	return func(o *TransportOptions) {
		o.Netns = netns
	}
}

// SockOptsTransportOption sets the socket options.
func SockOptsTransportOption(so *SockOpts) TransportOption {
	return func(o *TransportOptions) {
		o.SockOpts = so
	}
}

// RouteTransportOption sets the parent Route.
func RouteTransportOption(route Route) TransportOption {
	return func(o *TransportOptions) {
		o.Route = route
	}
}

// Transporter combines a dialer and connector to handle one node in a proxy
// chain. It establishes a connection to the next proxy hop (Dial), negotiates
// the transport protocol (Handshake), and reaches the final destination
// through the hop (Connect). It also supports reverse tunnels (Bind).
type Transporter interface {
	// Dial establishes a connection to the proxy node's address.
	Dial(ctx context.Context, addr string) (net.Conn, error)
	// Handshake performs a transport-level handshake (e.g. TLS) on the connection.
	Handshake(ctx context.Context, conn net.Conn) (net.Conn, error)
	// Connect reaches the target destination through the proxy node.
	Connect(ctx context.Context, conn net.Conn, network, address string) (net.Conn, error)
	// Bind sets up a reverse listener through the proxy node.
	Bind(ctx context.Context, conn net.Conn, network, address string, opts ...connector.BindOption) (net.Listener, error)
	// Multiplex reports whether the underlying connection supports multiplexing.
	Multiplex() bool
	// Options returns the transporter's configuration.
	Options() *TransportOptions
	// Copy returns a copy of this Transporter.
	Copy() Transporter
}
