package chain

type Chainer interface {
	Route(network, address string) *Route
}

type Chain struct {
	name   string
	groups []*NodeGroup
}

func NewChain(name string, groups ...*NodeGroup) *Chain {
	return &Chain{
		name:   name,
		groups: groups,
	}
}

func (c *Chain) AddNodeGroup(group *NodeGroup) {
	c.groups = append(c.groups, group)
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
		if node.Transport.Multiplex() {
			tr := node.Transport.Copy().
				WithRoute(r)
			node = node.Copy()
			node.Transport = tr
			r = &Route{}
		}

		r.addNode(node)
	}
	return r
}
