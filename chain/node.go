package chain

import (
	"github.com/go-gost/core/bypass"
	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/metadata"
	"github.com/go-gost/core/resolver"
)

type Node struct {
	Name       string
	Addr       string
	transport  *Transport
	bypass     bypass.Bypass
	resolver   resolver.Resolver
	hostMapper hosts.HostMapper
	marker     Marker
	metadata   metadata.Metadata
}

func NewNode(name, addr string) *Node {
	return &Node{
		Name:   name,
		Addr:   addr,
		marker: NewFailMarker(),
	}
}

func (node *Node) WithTransport(tr *Transport) *Node {
	node.transport = tr
	return node
}

func (node *Node) WithBypass(bypass bypass.Bypass) *Node {
	node.bypass = bypass
	return node
}

func (node *Node) WithResolver(reslv resolver.Resolver) *Node {
	node.resolver = reslv
	return node
}

func (node *Node) WithHostMapper(m hosts.HostMapper) *Node {
	node.hostMapper = m
	return node
}

func (node *Node) WithMetadata(md metadata.Metadata) *Node {
	node.metadata = md
	return node
}

func (node *Node) Marker() Marker {
	return node.marker
}

func (node *Node) Metadata() metadata.Metadata {
	return node.metadata
}

func (node *Node) Copy() *Node {
	n := &Node{}
	*n = *node
	return n
}

type NodeGroup struct {
	nodes    []*Node
	selector Selector[*Node]
	bypass   bypass.Bypass
}

func NewNodeGroup(nodes ...*Node) *NodeGroup {
	return &NodeGroup{
		nodes: nodes,
	}
}

func (g *NodeGroup) AddNode(node *Node) {
	g.nodes = append(g.nodes, node)
}

func (g *NodeGroup) Nodes() []*Node {
	return g.nodes
}

func (g *NodeGroup) WithSelector(selector Selector[*Node]) *NodeGroup {
	g.selector = selector
	return g
}

func (g *NodeGroup) WithBypass(bypass bypass.Bypass) *NodeGroup {
	g.bypass = bypass
	return g
}

func (g *NodeGroup) FilterAddr(addr string) *NodeGroup {
	var nodes []*Node
	for _, node := range g.nodes {
		if node.bypass == nil || !node.bypass.Contains(addr) {
			nodes = append(nodes, node)
		}
	}
	return &NodeGroup{
		nodes:    nodes,
		selector: g.selector,
		bypass:   g.bypass,
	}
}

func (g *NodeGroup) Next() *Node {
	if g == nil || len(g.nodes) == 0 {
		return nil
	}

	s := g.selector
	if s == nil {
		s = DefaultNodeSelector
	}

	return s.Select(g.nodes...)
}
