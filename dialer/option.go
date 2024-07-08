package dialer

import (
	"crypto/tls"
	"net/url"

	xnet "github.com/go-gost/core/common/net"
	"github.com/go-gost/core/logger"
)

type Options struct {
	Auth          *url.Userinfo
	TLSConfig     *tls.Config
	Logger        logger.Logger
	ProxyProtocol int
}

type Option func(opts *Options)

func AuthOption(auth *url.Userinfo) Option {
	return func(opts *Options) {
		opts.Auth = auth
	}
}

func TLSConfigOption(tlsConfig *tls.Config) Option {
	return func(opts *Options) {
		opts.TLSConfig = tlsConfig
	}
}

func LoggerOption(logger logger.Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

func ProxyProtocolOption(ppv int) Option {
	return func(opts *Options) {
		opts.ProxyProtocol = ppv
	}
}

type DialOptions struct {
	Host   string
	Dialer xnet.Dialer
}

type DialOption func(opts *DialOptions)

func HostDialOption(host string) DialOption {
	return func(opts *DialOptions) {
		opts.Host = host
	}
}

func NetDialerDialOption(dialer xnet.Dialer) DialOption {
	return func(opts *DialOptions) {
		opts.Dialer = dialer
	}
}

type HandshakeOptions struct {
	Addr string
}

type HandshakeOption func(opts *HandshakeOptions)

func AddrHandshakeOption(addr string) HandshakeOption {
	return func(opts *HandshakeOptions) {
		opts.Addr = addr
	}
}
