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
		node := group.Next()
		if node == nil {
			return
		}
		if node.Bypass != nil && node.Bypass.Contains(address) {
			break
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
