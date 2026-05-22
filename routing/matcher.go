// Package routing defines the Matcher interface for request-based routing
// decisions using HTTP-level context (hostname, method, path, headers, etc.).
package routing

import (
	"net"
	"net/http"
	"net/url"
)

// Request carries the HTTP-level context for a routing decision. It is
// populated by protocol handlers (HTTP, SOCKS with hostname, TLS SNI) and
// used by Matchers to decide which route or node to select.
type Request struct {
	// ClientIP is the client's IP address.
	ClientIP net.IP
	// Host is the target hostname.
	Host string
	// Protocol is the proxy protocol.
	Protocol string
	// Method is the HTTP request method.
	Method string
	// Path is the HTTP request path.
	Path string
	// Query is the HTTP query parameters.
	Query url.Values
	// Header is the HTTP request headers.
	Header http.Header
}

// Matcher decides whether a Request matches a routing rule. It is used at
// the node level to filter which proxy nodes are eligible for a given request
// based on protocol-specific attributes.
type Matcher interface {
	Match(*Request) bool
}
