package chain

import (
	"context"
	"net"
	"time"

	"github.com/go-gost/core/logger"
)

// Route represents a path through a chain of proxy nodes to a destination.
// A Route knows how to Dial (forward) and Bind (reverse tunnel) traffic
// through the proxy chain.
type Route interface {
	// Dial connects to the target address through the proxy chain.
	Dial(ctx context.Context, network, address string, opts ...DialOption) (net.Conn, error)
	// Bind sets up a reverse listener through the proxy chain, allowing
	// remote clients to connect back.
	Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error)
	// Nodes returns the list of nodes that make up this route.
	Nodes() []*Node
}

// DialOptions holds the runtime parameters for Dial operations.
type DialOptions struct {
	// Interface is the network interface to bind to.
	Interface string
	// Netns is the network namespace name.
	Netns string
	// SockOpts holds socket options.
	SockOpts *SockOpts
	// Logger is the logger for this dial operation.
	Logger logger.Logger
}

// DialOption is a functional option for configuring DialOptions.
type DialOption func(opts *DialOptions)

// InterfaceDialOption sets the network interface for the dial.
func InterfaceDialOption(ifName string) DialOption {
	return func(opts *DialOptions) {
		opts.Interface = ifName
	}
}

// NetnsDialOption sets the network namespace for the dial.
func NetnsDialOption(netns string) DialOption {
	return func(opts *DialOptions) {
		opts.Netns = netns
	}
}

// SockOptsDialOption sets the socket options for the dial.
func SockOptsDialOption(so *SockOpts) DialOption {
	return func(opts *DialOptions) {
		opts.SockOpts = so
	}
}

// LoggerDialOption sets the logger for the dial.
func LoggerDialOption(logger logger.Logger) DialOption {
	return func(opts *DialOptions) {
		opts.Logger = logger
	}
}

// BindOptions holds the runtime parameters for Bind operations.
type BindOptions struct {
	// Mux enables multiplexing on the bind listener.
	Mux bool
	// Backlog is the accept backlog size.
	Backlog int
	// UDPDataQueueSize is the queue size for UDP data packets.
	UDPDataQueueSize int
	// UDPDataBufferSize is the buffer size for UDP data packets.
	UDPDataBufferSize int
	// UDPConnTTL is the TTL for UDP connections.
	UDPConnTTL time.Duration
	// Logger is the logger for this bind operation.
	Logger logger.Logger
}

// BindOption is a functional option for configuring BindOptions.
type BindOption func(opts *BindOptions)

// MuxBindOption sets whether multiplexing is enabled.
func MuxBindOption(mux bool) BindOption {
	return func(opts *BindOptions) {
		opts.Mux = mux
	}
}

// BacklogBindOption sets the accept backlog size.
func BacklogBindOption(backlog int) BindOption {
	return func(opts *BindOptions) {
		opts.Backlog = backlog
	}
}

// UDPDataQueueSizeBindOption sets the UDP data queue size.
func UDPDataQueueSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataQueueSize = size
	}
}

// UDPDataBufferSizeBindOption sets the UDP data buffer size.
func UDPDataBufferSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataBufferSize = size
	}
}

// UDPConnTTLBindOption sets the UDP connection TTL.
func UDPConnTTLBindOption(ttl time.Duration) BindOption {
	return func(opts *BindOptions) {
		opts.UDPConnTTL = ttl
	}
}

// LoggerBindOption sets the logger for the bind operation.
func LoggerBindOption(logger logger.Logger) BindOption {
	return func(opts *BindOptions) {
		opts.Logger = logger
	}
}
