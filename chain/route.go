package chain

import (
	"context"
	"net"
	"time"

	"github.com/go-gost/core/logger"
)

type Route interface {
	Dial(ctx context.Context, network, address string, opts ...DialOption) (net.Conn, error)
	Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error)
	Nodes() []*Node
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
