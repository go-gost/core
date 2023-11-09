package hop

import (
	"context"

	"github.com/go-gost/core/chain"
)

type SelectOptions struct {
	Network  string
	Addr     string
	Protocol string
	Host     string
	Path     string
}

type SelectOption func(*SelectOptions)

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

func PathSelectOption(path string) SelectOption {
	return func(o *SelectOptions) {
		o.Path = path
	}
}

type Hop interface {
	Select(ctx context.Context, opts ...SelectOption) *chain.Node
}

type NodeList interface {
	Nodes() []*chain.Node
}
