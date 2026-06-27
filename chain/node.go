package chain

import (
	"regexp"

	"github.com/go-gost/core/auth"
	"github.com/go-gost/core/bypass"
	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/metadata"
	"github.com/go-gost/core/resolver"
	"github.com/go-gost/core/rewriter"
	"github.com/go-gost/core/routing"
	"github.com/go-gost/core/selector"
)

// NodeFilterSettings defines filtering criteria used to select nodes during routing.
type NodeFilterSettings struct {
	// Protocol filters by proxy protocol name (e.g. "http", "socks5").
	Protocol string
	// Host filters by hostname pattern.
	Host string
	// Path filters by URL path.
	Path string
}

// HTTPURLRewriteSetting defines an HTTP URL rewrite rule.
type HTTPURLRewriteSetting struct {
	// Pattern is the regex to match against the URL.
	Pattern *regexp.Regexp
	// Replacement is the replacement string.
	Replacement string
}

// HTTPBodyRewriteSettings defines an HTTP body rewrite rule.
type HTTPBodyRewriteSettings struct {
	// Type is the MIME type to match, e.g. "text/html".
	Type string
	// Pattern is the regex to match against the response body.
	Pattern *regexp.Regexp
	// Replacement is the replacement bytes.
	Replacement []byte
	// Rewriter is an optional plugin-based rewriter.
	// When set, Rewrite delegates to the plugin instead of using Pattern.
	Rewriter rewriter.Rewriter
}

// HTTPNodeSettings holds HTTP-level configuration for a node.
type HTTPNodeSettings struct {
	// Host is the HTTP Host header value.
	Host string
	// RequestHeader contains custom request headers to inject.
	RequestHeader map[string]string
	// ResponseHeader contains custom response headers to inject.
	ResponseHeader map[string]string
	// Auther is the HTTP authenticator for this node.
	Auther auth.Authenticator
	// RewriteURL holds the URL rewrite rules.
	RewriteURL []HTTPURLRewriteSetting
	// RewriteResponseBody holds the response body rewrite rules.
	RewriteResponseBody []HTTPBodyRewriteSettings
	// RewriteRequestBody holds the request body rewrite rules.
	RewriteRequestBody []HTTPBodyRewriteSettings
}

// TLSNodeSettings holds TLS configuration for a node.
type TLSNodeSettings struct {
	// ServerName is the TLS SNI server name.
	ServerName string
	// Secure indicates whether a secure connection is used.
	Secure bool
	// Options specifies TLS version and cipher suite settings.
	Options struct {
		MinVersion   string
		MaxVersion   string
		CipherSuites []string
		ALPN         []string
	}
}

// NodeOptions holds the initialization parameters for a Node.
type NodeOptions struct {
	Network    string
	Transport  Transporter
	Bypass     bypass.Bypass
	Resolver   resolver.Resolver
	HostMapper hosts.HostMapper
	Filter     *NodeFilterSettings
	HTTP       *HTTPNodeSettings
	TLS        *TLSNodeSettings
	Metadata   metadata.Metadata
	Matcher    routing.Matcher
	Priority   int
}

// NodeOption is a functional option for configuring NodeOptions.
type NodeOption func(*NodeOptions)

// TransportNodeOption sets the Transporter for the node.
func TransportNodeOption(tr Transporter) NodeOption {
	return func(o *NodeOptions) {
		o.Transport = tr
	}
}

// BypassNodeOption sets the Bypass rule for the node.
func BypassNodeOption(bp bypass.Bypass) NodeOption {
	return func(o *NodeOptions) {
		o.Bypass = bp
	}
}

// ResolverNodeOption sets the Resolver for the node.
func ResolverNodeOption(resolver resolver.Resolver) NodeOption {
	return func(o *NodeOptions) {
		o.Resolver = resolver
	}
}

// HostMapperNodeOption sets the HostMapper for the node.
func HostMapperNodeOption(m hosts.HostMapper) NodeOption {
	return func(o *NodeOptions) {
		o.HostMapper = m
	}
}

// NetworkNodeOption sets the network type for the node.
func NetworkNodeOption(network string) NodeOption {
	return func(o *NodeOptions) {
		o.Network = network
	}
}

// NodeFilterOption sets the filter settings for the node.
func NodeFilterOption(filter *NodeFilterSettings) NodeOption {
	return func(o *NodeOptions) {
		o.Filter = filter
	}
}

// HTTPNodeOption sets the HTTP settings for the node.
func HTTPNodeOption(httpSettings *HTTPNodeSettings) NodeOption {
	return func(o *NodeOptions) {
		o.HTTP = httpSettings
	}
}

// TLSNodeOption sets the TLS settings for the node.
func TLSNodeOption(tlsSettings *TLSNodeSettings) NodeOption {
	return func(o *NodeOptions) {
		o.TLS = tlsSettings
	}
}

// MetadataNodeOption sets the Metadata for the node.
func MetadataNodeOption(md metadata.Metadata) NodeOption {
	return func(o *NodeOptions) {
		o.Metadata = md
	}
}

// MatcherNodeOption sets the routing Matcher for the node.
func MatcherNodeOption(matcher routing.Matcher) NodeOption {
	return func(o *NodeOptions) {
		o.Matcher = matcher
	}
}

// PriorityNodeOption sets the priority for the node. Higher priority nodes
// are selected first during routing.
func PriorityNodeOption(priority int) NodeOption {
	return func(o *NodeOptions) {
		o.Priority = priority
	}
}

// Node represents a single hop in a proxy chain. Each node has a name,
// address, transport layer, and optional filters for protocol-specific routing.
// Nodes implement selector.Markable so that failed nodes can be temporarily
// excluded from selection.
type Node struct {
	Name    string
	Addr    string
	marker  selector.Marker
	options NodeOptions
}

// NewNode creates a new Node with the given name and address.
func NewNode(name string, addr string, opts ...NodeOption) *Node {
	var options NodeOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return &Node{
		Name:    name,
		Addr:    addr,
		marker:  selector.NewFailMarker(),
		options: options,
	}
}

// Options returns the Node's configuration.
func (node *Node) Options() *NodeOptions {
	return &node.options
}

// Metadata returns the Node's metadata. Implements the metadata.Metadatable interface.
func (node *Node) Metadata() metadata.Metadata {
	return node.options.Metadata
}

// Marker returns the Node's failure marker. Implements the selector.Markable interface.
func (node *Node) Marker() selector.Marker {
	return node.marker
}

// Copy returns a shallow copy of the Node.
func (node *Node) Copy() *Node {
	n := &Node{}
	*n = *node
	return n
}
