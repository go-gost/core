// Package hop defines the Hop interface for selecting a node from a group
// of equivalent proxy nodes using load-balancing strategies.
package hop

import (
	"context"
	"net"
	"net/http"
	"net/url"

	"github.com/go-gost/core/chain"
)

// SelectOptions holds the runtime context used to select a node.
// These provide HTTP-level visibility into the request, allowing
// routing decisions based on hostname, path, headers, etc.
type SelectOptions struct {
	// ClientIP is the client's IP address.
	ClientIP net.IP
	// Network is the connection network type ("tcp", "udp").
	Network string
	// Addr is the destination address.
	Addr string
	// Protocol is the proxy protocol in use.
	Protocol string
	// Host is the target hostname (from HTTP Host header or TLS SNI).
	Host string
	// Method is the HTTP request method.
	Method string
	// Path is the HTTP request path.
	Path string
	// Query is the HTTP request query parameters.
	Query url.Values
	// Header is the HTTP request headers.
	Header http.Header
}

// SelectOption is a functional option for configuring SelectOptions.
type SelectOption func(*SelectOptions)

// ClientIPSelectOption sets the client IP.
func ClientIPSelectOption(clientIP net.IP) SelectOption {
	return func(o *SelectOptions) {
		o.ClientIP = clientIP
	}
}

// NetworkSelectOption sets the network type.
func NetworkSelectOption(network string) SelectOption {
	return func(so *SelectOptions) {
		so.Network = network
	}
}

// AddrSelectOption sets the destination address.
func AddrSelectOption(addr string) SelectOption {
	return func(o *SelectOptions) {
		o.Addr = addr
	}
}

// ProtocolSelectOption sets the proxy protocol.
func ProtocolSelectOption(protocol string) SelectOption {
	return func(o *SelectOptions) {
		o.Protocol = protocol
	}
}

// HostSelectOption sets the target hostname.
func HostSelectOption(host string) SelectOption {
	return func(o *SelectOptions) {
		o.Host = host
	}
}

// MethodSelectOption sets the HTTP method.
func MethodSelectOption(method string) SelectOption {
	return func(o *SelectOptions) {
		o.Method = method
	}
}

// PathSelectOption sets the HTTP request path.
func PathSelectOption(path string) SelectOption {
	return func(o *SelectOptions) {
		o.Path = path
	}
}

// QuerySelectOption sets the HTTP query parameters.
func QuerySelectOption(query url.Values) SelectOption {
	return func(o *SelectOptions) {
		o.Query = query
	}
}

// HeaderSelectOption sets the HTTP headers.
func HeaderSelectOption(header http.Header) SelectOption {
	return func(o *SelectOptions) {
		o.Header = header
	}
}

// Hop represents a group of functionally equivalent proxy nodes. It uses a
// selector.Strategy to pick one node from the group, enabling load balancing
// (round-robin, random, weighted, etc.) across redundant proxy servers.
type Hop interface {
	// Select picks a Node from the group based on the selection options.
	// It returns nil if no node is available.
	Select(ctx context.Context, opts ...SelectOption) *chain.Node
}

// NodeList provides access to all nodes in a group. Implemented by hops
// that expose their full node list for introspection or custom routing.
type NodeList interface {
	Nodes() []*chain.Node
}
