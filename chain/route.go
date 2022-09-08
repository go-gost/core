package chain

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/go-gost/core/common/net/dialer"
	"github.com/go-gost/core/common/net/udp"
	"github.com/go-gost/core/connector"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/core/metrics"
)

var (
	ErrEmptyRoute = errors.New("empty route")
)

type Route interface {
	Dial(ctx context.Context, network, address string, opts ...DialOption) (net.Conn, error)
	Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error)
	Len() int
	Path() []*Node
}

type route struct {
	chain *Chain
	nodes []*Node
}

func newRoute() *route {
	return &route{}
}

func (r *route) addNode(node *Node) {
	r.nodes = append(r.nodes, node)
}

func (r *route) WithChain(chain *Chain) *route {
	r.chain = chain
	return r
}

func (r *route) Dial(ctx context.Context, network, address string, opts ...DialOption) (net.Conn, error) {
	var options DialOptions
	for _, opt := range opts {
		opt(&options)
	}

	if r.Len() == 0 {
		netd := dialer.NetDialer{
			Timeout:   options.Timeout,
			Interface: options.Interface,
		}
		if options.SockOpts != nil {
			netd.Mark = options.SockOpts.Mark
		}
		if r != nil {
			netd.Logger = options.Logger
		}

		return netd.Dial(ctx, network, address)
	}

	conn, err := r.connect(ctx, options.Logger)
	if err != nil {
		return nil, err
	}

	cc, err := r.GetNode(r.Len()-1).transport.Connect(ctx, conn, network, address)
	if err != nil {
		if conn != nil {
			conn.Close()
		}
		return nil, err
	}
	return cc, nil
}

func (r *route) Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error) {
	var options BindOptions
	for _, opt := range opts {
		opt(&options)
	}

	if r.Len() == 0 {
		return r.bindLocal(ctx, network, address, &options)
	}

	conn, err := r.connect(ctx, options.Logger)
	if err != nil {
		return nil, err
	}

	ln, err := r.GetNode(r.Len()-1).transport.Bind(ctx,
		conn, network, address,
		connector.BacklogBindOption(options.Backlog),
		connector.MuxBindOption(options.Mux),
		connector.UDPConnTTLBindOption(options.UDPConnTTL),
		connector.UDPDataBufferSizeBindOption(options.UDPDataBufferSize),
		connector.UDPDataQueueSizeBindOption(options.UDPDataQueueSize),
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return ln, nil
}

func (r *route) connect(ctx context.Context, logger logger.Logger) (conn net.Conn, err error) {
	if r.Len() == 0 {
		return nil, ErrEmptyRoute
	}

	network := "ip"
	node := r.nodes[0]

	defer func() {
		if r.chain != nil {
			marker := r.chain.Marker()
			// chain error
			if err != nil {
				if marker != nil {
					marker.Mark()
				}
				if v := metrics.GetCounter(metrics.MetricChainErrorsCounter,
					metrics.Labels{"chain": r.chain.name, "node": node.Name}); v != nil {
					v.Inc()
				}
			} else {
				if marker != nil {
					marker.Reset()
				}
			}
		}
	}()

	addr, err := resolve(ctx, network, node.Addr, node.resolver, node.hostMapper, logger)
	marker := node.Marker()
	if err != nil {
		if marker != nil {
			marker.Mark()
		}
		return
	}

	start := time.Now()
	cc, err := node.transport.Dial(ctx, addr)
	if err != nil {
		if marker != nil {
			marker.Mark()
		}
		return
	}

	cn, err := node.transport.Handshake(ctx, cc)
	if err != nil {
		cc.Close()
		if marker != nil {
			marker.Mark()
		}
		return
	}
	if marker != nil {
		marker.Reset()
	}

	if r.chain != nil {
		if v := metrics.GetObserver(metrics.MetricNodeConnectDurationObserver,
			metrics.Labels{"chain": r.chain.name, "node": node.Name}); v != nil {
			v.Observe(time.Since(start).Seconds())
		}
	}

	preNode := node
	for _, node := range r.nodes[1:] {
		marker := node.Marker()
		addr, err = resolve(ctx, network, node.Addr, node.resolver, node.hostMapper, logger)
		if err != nil {
			cn.Close()
			if marker != nil {
				marker.Mark()
			}
			return
		}
		cc, err = preNode.transport.Connect(ctx, cn, "tcp", addr)
		if err != nil {
			cn.Close()
			if marker != nil {
				marker.Mark()
			}
			return
		}
		cc, err = node.transport.Handshake(ctx, cc)
		if err != nil {
			cn.Close()
			if marker != nil {
				marker.Mark()
			}
			return
		}
		if marker != nil {
			marker.Reset()
		}

		cn = cc
		preNode = node
	}

	conn = cn
	return
}

func (r *route) Len() int {
	if r == nil {
		return 0
	}
	return len(r.nodes)
}

func (r *route) GetNode(index int) *Node {
	if r.Len() == 0 || index < 0 || index >= len(r.nodes) {
		return nil
	}
	return r.nodes[index]
}

func (r *route) Path() (path []*Node) {
	if r == nil || len(r.nodes) == 0 {
		return nil
	}

	for _, node := range r.nodes {
		if node.transport != nil && node.transport.route != nil {
			path = append(path, node.transport.route.Path()...)
		}
		path = append(path, node)
	}
	return
}

func (r *route) bindLocal(ctx context.Context, network, address string, opts *BindOptions) (net.Listener, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		addr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return net.ListenTCP(network, addr)
	case "udp", "udp4", "udp6":
		addr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, err
		}
		conn, err := net.ListenUDP(network, addr)
		if err != nil {
			return nil, err
		}
		logger := logger.Default().WithFields(map[string]any{
			"network": network,
			"address": address,
		})
		ln := udp.NewListener(conn, &udp.ListenConfig{
			Backlog:        opts.Backlog,
			ReadQueueSize:  opts.UDPDataQueueSize,
			ReadBufferSize: opts.UDPDataBufferSize,
			TTL:            opts.UDPConnTTL,
			KeepAlive:      true,
			Logger:         logger,
		})
		return ln, err
	default:
		err := fmt.Errorf("network %s unsupported", network)
		return nil, err
	}
}

type DialOptions struct {
	Timeout   time.Duration
	Interface string
	SockOpts  *SockOpts
	Logger    logger.Logger
}

type DialOption func(opts *DialOptions)

func TimeoutDialOption(d time.Duration) DialOption {
	return func(opts *DialOptions) {
		opts.Timeout = d
	}
}

func InterfaceDialOption(ifName string) DialOption {
	return func(opts *DialOptions) {
		opts.Interface = ifName
	}
}

func SockOptsDialOption(so *SockOpts) DialOption {
	return func(opts *DialOptions) {
		opts.SockOpts = so
	}
}

func LoggerDialOption(logger logger.Logger) DialOption {
	return func(opts *DialOptions) {
		opts.Logger = logger
	}
}

type BindOptions struct {
	Mux               bool
	Backlog           int
	UDPDataQueueSize  int
	UDPDataBufferSize int
	UDPConnTTL        time.Duration
	Logger            logger.Logger
}

type BindOption func(opts *BindOptions)

func MuxBindOption(mux bool) BindOption {
	return func(opts *BindOptions) {
		opts.Mux = mux
	}
}

func BacklogBindOption(backlog int) BindOption {
	return func(opts *BindOptions) {
		opts.Backlog = backlog
	}
}

func UDPDataQueueSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataQueueSize = size
	}
}

func UDPDataBufferSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataBufferSize = size
	}
}

func UDPConnTTLBindOption(ttl time.Duration) BindOption {
	return func(opts *BindOptions) {
		opts.UDPConnTTL = ttl
	}
}

func LoggerBindOption(logger logger.Logger) BindOption {
	return func(opts *BindOptions) {
		opts.Logger = logger
	}
}
