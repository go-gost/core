package chain

import (
	"context"
	"net"

	"github.com/go-gost/core/connector"
)

type TransportOptions struct {
	Addr     string
	IfceName string
	Netns    string
	SockOpts *SockOpts
	Route    Route
}

type TransportOption func(*TransportOptions)

func AddrTransportOption(addr string) TransportOption {
	return func(o *TransportOptions) {
		o.Addr = addr
	}
}

func InterfaceTransportOption(ifceName string) TransportOption {
	return func(o *TransportOptions) {
		o.IfceName = ifceName
	}
}

func NetnsTransportOption(netns string) TransportOption {
	return func(o *TransportOptions) {
		o.Netns = netns
	}
}

func SockOptsTransportOption(so *SockOpts) TransportOption {
	return func(o *TransportOptions) {
		o.SockOpts = so
	}
}

func RouteTransportOption(route Route) TransportOption {
	return func(o *TransportOptions) {
		o.Route = route
	}
}

type Transporter interface {
	Dial(ctx context.Context, addr string) (net.Conn, error)
	Handshake(ctx context.Context, conn net.Conn) (net.Conn, error)
	Connect(ctx context.Context, conn net.Conn, network, address string) (net.Conn, error)
	Bind(ctx context.Context, conn net.Conn, network, address string, opts ...connector.BindOption) (net.Listener, error)
	Multiplex() bool
	Options() *TransportOptions
	Copy() Transporter
}
