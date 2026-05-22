package dialer

import (
	"crypto/tls"
	"net/url"

	xnet "github.com/go-gost/core/common/net"
	"github.com/go-gost/core/logger"
)

// Options holds the initialization parameters for a Dialer.
type Options struct {
	// Auth is the authentication credentials for the dialer.
	Auth *url.Userinfo
	// TLSConfig is the TLS configuration for the dialer's outgoing connections.
	TLSConfig *tls.Config
	// Logger is the logger for dialer operations.
	Logger logger.Logger
	// ProxyProtocol is the proxy protocol version (0 = disabled).
	ProxyProtocol int
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// AuthOption sets the authentication credentials.
func AuthOption(auth *url.Userinfo) Option {
	return func(opts *Options) {
		opts.Auth = auth
	}
}

// TLSConfigOption sets the TLS configuration.
func TLSConfigOption(tlsConfig *tls.Config) Option {
	return func(opts *Options) {
		opts.TLSConfig = tlsConfig
	}
}

// LoggerOption sets the logger.
func LoggerOption(logger logger.Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

// ProxyProtocolOption sets the proxy protocol version.
func ProxyProtocolOption(ppv int) Option {
	return func(opts *Options) {
		opts.ProxyProtocol = ppv
	}
}

// DialOptions holds the runtime parameters for the Dial operation.
type DialOptions struct {
	// Host is the target hostname.
	Host string
	// Dialer is the underlying dialer to use.
	Dialer xnet.Dialer
}

// DialOption is a functional option for configuring DialOptions.
type DialOption func(opts *DialOptions)

// HostDialOption sets the target host.
func HostDialOption(host string) DialOption {
	return func(opts *DialOptions) {
		opts.Host = host
	}
}

// NetDialerDialOption sets the underlying dialer.
func NetDialerDialOption(dialer xnet.Dialer) DialOption {
	return func(opts *DialOptions) {
		opts.Dialer = dialer
	}
}

// HandshakeOptions holds the runtime parameters for the Handshake operation.
type HandshakeOptions struct {
	// Addr is the address of the proxy server.
	Addr string
}

// HandshakeOption is a functional option for configuring HandshakeOptions.
type HandshakeOption func(opts *HandshakeOptions)

// AddrHandshakeOption sets the proxy server address for the handshake.
func AddrHandshakeOption(addr string) HandshakeOption {
	return func(opts *HandshakeOptions) {
		opts.Addr = addr
	}
}
