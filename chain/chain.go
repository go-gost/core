package chain

import (
	"github.com/go-gost/core/metadata"
	"github.com/go-gost/core/selector"
)

type Chainer interface {
	Route(network, address string) Route
}

type SelectableChainer interface {
	Chainer
	selector.Selectable
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

func (c *Chain) Metadata() metadata.Metadata {
	return c.metadata
}

func (c *Chain) Marker() selector.Marker {
	return c.marker
}

func (c *Chain) Route(network, address string) Route {
	if c == nil || len(c.groups) == 0 {
		return nil
	}

	rt := newRoute().WithChain(c)
	for _, group := range c.groups {
		// hop level bypass test
		if group.bypass != nil && group.bypass.Contains(address) {
			break
		}

		node := group.FilterAddr(address).Next()
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
	chains   []SelectableChainer
	selector selector.Selector[SelectableChainer]
}

func NewChainGroup(chains ...SelectableChainer) *ChainGroup {
	return &ChainGroup{chains: chains}
}

func (p *ChainGroup) WithSelector(s selector.Selector[SelectableChainer]) *ChainGroup {
	p.selector = s
	return p
}

func (p *ChainGroup) Route(network, address string) Route {
	if chain := p.next(); chain != nil {
		return chain.Route(network, address)
	}
	return nil
}

func (p *ChainGroup) next() Chainer {
	if p == nil || len(p.chains) == 0 {
		return nil
	}

	return p.selector.Select(p.chains...)
}
