package chain

import (
	"github.com/go-gost/core/auth"
	"github.com/go-gost/core/bypass"
	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/metadata"
	"github.com/go-gost/core/resolver"
	"github.com/go-gost/core/selector"
)

type HTTPNodeSettings struct {
	Host   string
	Header map[string]string
}

type TLSNodeSettings struct {
	ServerName string
	Secure     bool
	Options    struct {
		MinVersion   string
		MaxVersion   string
		CipherSuites []string
	}
}

type NodeOptions struct {
	Transport  *Transport
	Bypass     bypass.Bypass
	Resolver   resolver.Resolver
	HostMapper hosts.HostMapper
	Metadata   metadata.Metadata
	Host       string
	Network    string
	Protocol   string
	Path       string
	HTTP       *HTTPNodeSettings
	TLS        *TLSNodeSettings
	Auther     auth.Authenticator
}

type NodeOption func(*NodeOptions)

func TransportNodeOption(tr *Transport) NodeOption {
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

func HostNodeOption(host string) NodeOption {
	return func(o *NodeOptions) {
		o.Host = host
	}
}

func NetworkNodeOption(network string) NodeOption {
	return func(o *NodeOptions) {
		o.Network = network
	}
}

func ProtocolNodeOption(protocol string) NodeOption {
	return func(o *NodeOptions) {
		o.Protocol = protocol
	}
}

func PathNodeOption(path string) NodeOption {
	return func(o *NodeOptions) {
		o.Path = path
	}
}

func MetadataNodeOption(md metadata.Metadata) NodeOption {
	return func(o *NodeOptions) {
		o.Metadata = md
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

func AutherNodeOption(auther auth.Authenticator) NodeOption {
	return func(o *NodeOptions) {
		o.Auther = auther
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
