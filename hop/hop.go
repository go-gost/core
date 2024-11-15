package hop

import (
	"context"
	"net"
	"net/http"
	"net/url"

	"github.com/go-gost/core/chain"
)

type SelectOptions struct {
	ClientIP net.IP
	Network  string
	Addr     string
	Protocol string
	Host     string
	Method   string
	Path     string
	Query    url.Values
	Header   http.Header
}

type SelectOption func(*SelectOptions)

func ClientIPSelectOption(clientIP net.IP) SelectOption {
	return func(o *SelectOptions) {
		o.ClientIP = clientIP
	}
}

func NetworkSelectOption(network string) SelectOption {
	return func(so *SelectOptions) {
		so.Network = network
	}
}

func AddrSelectOption(addr string) SelectOption {
	return func(o *SelectOptions) {
		o.Addr = addr
	}
}

func ProtocolSelectOption(protocol string) SelectOption {
	return func(o *SelectOptions) {
		o.Protocol = protocol
	}
}

func HostSelectOption(host string) SelectOption {
	return func(o *SelectOptions) {
		o.Host = host
	}
}

func MethodSelectOption(method string) SelectOption {
	return func(o *SelectOptions) {
		o.Method = method
	}
}

func PathSelectOption(path string) SelectOption {
	return func(o *SelectOptions) {
		o.Path = path
	}
}

func QuerySelectOption(query url.Values) SelectOption {
	return func(o *SelectOptions) {
		o.Query = query
	}
}

func HeaderSelectOption(header http.Header) SelectOption {
	return func(o *SelectOptions) {
		o.Header = header
	}
}

type Hop interface {
	Select(ctx context.Context, opts ...SelectOption) *chain.Node
}

type NodeList interface {
	Nodes() []*chain.Node
}
