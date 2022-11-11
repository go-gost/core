package chain

import "context"

type SelectOptions struct {
	Addr string
	Host string
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

type Hop interface {
	Nodes() []*Node
	Select(ctx context.Context, opts ...SelectOption) *Node
}
