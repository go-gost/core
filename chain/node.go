package chain

import (
	"regexp"

	"github.com/go-gost/core/auth"
	"github.com/go-gost/core/bypass"
	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/metadata"
	"github.com/go-gost/core/resolver"
	"github.com/go-gost/core/routing"
	"github.com/go-gost/core/selector"
)

type NodeFilterSettings struct {
	Protocol string
	Host     string
	Path     string
}

type HTTPURLRewriteSetting struct {
	Pattern     *regexp.Regexp
	Replacement string
}

type HTTPBodyRewriteSettings struct {
	Type        string
	Pattern     *regexp.Regexp
	Replacement []byte
}

type HTTPNodeSettings struct {
	Host           string
	Header         map[string]string
	ResponseHeader map[string]string
	Auther         auth.Authenticator
	RewriteURL     []HTTPURLRewriteSetting
	RewriteBody    []HTTPBodyRewriteSettings
}

type TLSNodeSettings struct {
	ServerName string
	Secure     bool
	Options    struct {
		MinVersion   string
		MaxVersion   string
		CipherSuites []string
		ALPN         []string
	}
}

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

type NodeOption func(*NodeOptions)

func TransportNodeOption(tr Transporter) NodeOption {
	return func(o *NodeOptions) {
		o.Transport = tr
	}
}

func BypassNodeOption(bp bypass.Bypass) NodeOption {
	return func(o *NodeOptions) {
		o.Bypass = bp
	}
}

func ResoloverNodeOption(resolver resolver.Resolver) NodeOption {
	return func(o *NodeOptions) {
		o.Resolver = resolver
	}
}

func HostMapperNodeOption(m hosts.HostMapper) NodeOption {
	return func(o *NodeOptions) {
		o.HostMapper = m
	}
}

func NetworkNodeOption(network string) NodeOption {
	return func(o *NodeOptions) {
		o.Network = network
	}
}

func NodeFilterOption(filter *NodeFilterSettings) NodeOption {
	return func(o *NodeOptions) {
		o.Filter = filter
	}
}

func HTTPNodeOption(httpSettings *HTTPNodeSettings) NodeOption {
	return func(o *NodeOptions) {
		o.HTTP = httpSettings
	}
}

func TLSNodeOption(tlsSettings *TLSNodeSettings) NodeOption {
	return func(o *NodeOptions) {
		o.TLS = tlsSettings
	}
}

func MetadataNodeOption(md metadata.Metadata) NodeOption {
	return func(o *NodeOptions) {
		o.Metadata = md
	}
}

func MatcherNodeOption(matcher routing.Matcher) NodeOption {
	return func(o *NodeOptions) {
		o.Matcher = matcher
	}
}

func PriorityNodeOption(priority int) NodeOption {
	return func(o *NodeOptions) {
		o.Priority = priority
	}
}

type Node struct {
	Name    string
	Addr    string
	marker  selector.Marker
	options NodeOptions
}

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

func (node *Node) Options() *NodeOptions {
	return &node.options
}

// Metadata implements metadadta.Metadatable interface.
func (node *Node) Metadata() metadata.Metadata {
	return node.options.Metadata
}

// Marker implements selector.Markable interface.
func (node *Node) Marker() selector.Marker {
	return node.marker
}

func (node *Node) Copy() *Node {
	n := &Node{}
	*n = *node
	return n
}
