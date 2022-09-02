package chain

import (
	"github.com/go-gost/core/metadata"
)

type Chainer interface {
	Route(network, address string) *Route
}

type SelectableChainer interface {
	Chainer
	Selectable
}

type Chain struct {
	name     string
	groups   []*NodeGroup
	marker   Marker
	metadata metadata.Metadata
}

func NewChain(name string, groups ...*NodeGroup) *Chain {
	return &Chain{
		name:   name,
		groups: groups,
		marker: NewFailMarker(),
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

func (c *Chain) Marker() Marker {
	return c.marker
}

func (c *Chain) Route(network, address string) (r *Route) {
	if c == nil || len(c.groups) == 0 {
		return
	}

	r = &Route{
		chain: c,
	}
	for _, group := range c.groups {
		// hop level bypass test
		if group.bypass != nil && group.bypass.Contains(address) {
			break
		}

		node := group.FilterAddr(address).Next()
		if node == nil {
			return
		}
		if node.transport.Multiplex() {
			tr := node.transport.Copy().
				WithRoute(r)
			node = node.Copy()
			node.transport = tr
			r = &Route{}
		}

		r.addNode(node)
	}
	return r
}

type ChainGroup struct {
	Chains   []SelectableChainer
	Selector Selector[SelectableChainer]
}

func (p *ChainGroup) Route(network, address string) *Route {
	if chain := p.next(); chain != nil {
		return chain.Route(network, address)
	}
	return nil
}

func (p *ChainGroup) next() Chainer {
	if p == nil || len(p.Chains) == 0 {
		return nil
	}

	s := p.Selector
	if s == nil {
		s = DefaultChainSelector
	}
	return s.Select(p.Chains...)
}
