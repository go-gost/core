package chain

import (
	"sync/atomic"
	"time"

	"github.com/go-gost/core/bypass"
	"github.com/go-gost/core/hosts"
	"github.com/go-gost/core/resolver"
)

type Node struct {
	Name      string
	Addr      string
	Transport *Transport
	Bypass    bypass.Bypass
	Resolver  resolver.Resolver
	Hosts     hosts.HostMapper
	Marker    *FailMarker
}

func (node *Node) Copy() *Node {
	n := &Node{}
	*n = *node
	return n
}

type NodeGroup struct {
	nodes    []*Node
	selector Selector
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

func (g *NodeGroup) WithSelector(selector Selector) *NodeGroup {
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
		if node.Bypass == nil || !node.Bypass.Contains(addr) {
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
		s = DefaultSelector
	}

	return s.Select(g.nodes...)
}

type FailMarker struct {
	failTime  int64
	failCount int64
}

func (m *FailMarker) FailTime() int64 {
	if m == nil {
		return 0
	}

	return atomic.LoadInt64(&m.failTime)
}

func (m *FailMarker) FailCount() int64 {
	if m == nil {
		return 0
	}

	return atomic.LoadInt64(&m.failCount)
}

func (m *FailMarker) Mark() {
	if m == nil {
		return
	}

	atomic.AddInt64(&m.failCount, 1)
	atomic.StoreInt64(&m.failTime, time.Now().Unix())
}

func (m *FailMarker) Reset() {
	if m == nil {
		return
	}

	atomic.StoreInt64(&m.failCount, 0)
}
