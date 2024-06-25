package chain

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/go-gost/core/common/net/dialer"
	"github.com/go-gost/core/common/net/udp"
	"github.com/go-gost/core/logger"
)

var (
	ErrEmptyRoute = errors.New("empty route")
)

var (
	DefaultRoute Route = &route{}
)

type Route interface {
	Dial(ctx context.Context, network, address string, opts ...DialOption) (net.Conn, error)
	Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error)
	Nodes() []*Node
}

// route is a Route without nodes.
type route struct{}

func (*route) Dial(ctx context.Context, network, address string, opts ...DialOption) (net.Conn, error) {
	var options DialOptions
	for _, opt := range opts {
		opt(&options)
	}

	netd := dialer.NetDialer{
		Interface: options.Interface,
		Netns:     options.Netns,
		Logger:    options.Logger,
	}
	if options.SockOpts != nil {
		netd.Mark = options.SockOpts.Mark
	}

	return netd.Dial(ctx, network, address)
}

func (*route) Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error) {
	var options BindOptions
	for _, opt := range opts {
		opt(&options)
	}

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
			Backlog:        options.Backlog,
			ReadQueueSize:  options.UDPDataQueueSize,
			ReadBufferSize: options.UDPDataBufferSize,
			TTL:            options.UDPConnTTL,
			KeepAlive:      true,
			Logger:         logger,
		})
		return ln, err
	default:
		err := fmt.Errorf("network %s unsupported", network)
		return nil, err
	}
}

func (r *route) Nodes() []*Node {
	return nil
}

type DialOptions struct {
	Interface string
	Netns     string
	SockOpts  *SockOpts
	Logger    logger.Logger
}

type DialOption func(opts *DialOptions)

func InterfaceDialOption(ifName string) DialOption {
	return func(opts *DialOptions) {
		opts.Interface = ifName
	}
}

func NetnsDialOption(netns string) DialOption {
	return func(opts *DialOptions) {
		opts.Netns = netns
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
