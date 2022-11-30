package chain

import "context"

type SelectOptions struct {
	Addr     string
	Host     string
	Protocol string
}

type SelectOption func(*SelectOptions)

func AddrSelectOption(addr string) SelectOption {
	return func(o *SelectOptions) {
		o.Addr = addr
	}
}

func HostSelectOption(host string) SelectOption {
	return func(o *SelectOptions) {
		o.Host = host
	}
}

func ProtocolSelectOption(protocol string) SelectOption {
	return func(o *SelectOptions) {
		o.Protocol = protocol
	}
}

type Hop interface {
	Nodes() []*Node
	Select(ctx context.Context, opts ...SelectOption) *Node
}
