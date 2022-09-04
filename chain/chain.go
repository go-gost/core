package chain

import (
	"context"

	"github.com/go-gost/core/metadata"
	"github.com/go-gost/core/selector"
)

type Chainer interface {
	Route(ctx context.Context, network, address string) Route
}

type Chain struct {
	name     string
	groups   []*NodeGroup
	marker   selector.Marker
	metadata metadata.Metadata
}

func NewChain(name string, groups ...*NodeGroup) *Chain {
	return &Chain{
		name:   name,
		groups: groups,
		marker: selector.NewFailMarker(),
	}
}

func (c *Chain) AddNodeGroup(group *NodeGroup) {
	c.groups = append(c.groups, group)
}

func (c *Chain) WithMetadata(md metadata.Metadata) {
	c.metadata = md
}

// Metadata implements metadata.Metadatable interface.
func (c *Chain) Metadata() metadata.Metadata {
	return c.metadata
}

// Marker implements selector.Markable interface.
func (c *Chain) Marker() selector.Marker {
	return c.marker
}

func (c *Chain) Route(ctx context.Context, network, address string) Route {
	if c == nil || len(c.groups) == 0 {
		return nil
	}

	rt := newRoute().WithChain(c)
	for _, group := range c.groups {
		// hop level bypass test
		if group.bypass != nil && group.bypass.Contains(address) {
			break
		}

		node := group.FilterAddr(address).Next(ctx)
		if node == nil {
			return rt
		}
		if node.transport.Multiplex() {
			tr := node.transport.
				Copy().
				WithRoute(rt)
			node = node.Copy()
			node.transport = tr
			rt = newRoute()
		}

		rt.addNode(node)
	}
	return rt
}

type ChainGroup struct {
	chains   []Chainer
	selector selector.Selector[Chainer]
}

func NewChainGroup(chains ...Chainer) *ChainGroup {
	return &ChainGroup{chains: chains}
}

func (p *ChainGroup) WithSelector(s selector.Selector[Chainer]) *ChainGroup {
	p.selector = s
	return p
}

func (p *ChainGroup) Route(ctx context.Context, network, address string) Route {
	if chain := p.next(ctx); chain != nil {
		return chain.Route(ctx, network, address)
	}
	return nil
}

func (p *ChainGroup) next(ctx context.Context) Chainer {
	if p == nil || len(p.chains) == 0 {
		return nil
	}

	return p.selector.Select(ctx, p.chains...)
}
