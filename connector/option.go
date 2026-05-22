package connector

import (
	"crypto/tls"
	"net/url"
	"time"

	xnet "github.com/go-gost/core/common/net"
	"github.com/go-gost/core/logger"
)

// Options holds the initialization parameters for a Connector.
type Options struct {
	// Auth is the authentication credentials for the connector.
	Auth *url.Userinfo
	// TLSConfig is the TLS configuration for the connector.
	TLSConfig *tls.Config
	// Logger is the logger for connector operations.
	Logger logger.Logger
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

// ConnectOptions holds the runtime parameters for the Connect operation.
type ConnectOptions struct {
	// Dialer is used to establish the underlying connection if needed.
	Dialer xnet.Dialer
}

// ConnectOption is a functional option for configuring ConnectOptions.
type ConnectOption func(opts *ConnectOptions)

// DialerConnectOption sets the dialer for the connect operation.
func DialerConnectOption(dialer xnet.Dialer) ConnectOption {
	return func(opts *ConnectOptions) {
		opts.Dialer = dialer
	}
}

// BindOptions holds the runtime parameters for the Bind operation.
type BindOptions struct {
	// Mux enables connection multiplexing on the bind listener.
	Mux bool
	// Backlog is the accept backlog size.
	Backlog int
	// UDPDataQueueSize is the queue size for UDP data packets.
	UDPDataQueueSize int
	// UDPDataBufferSize is the buffer size for UDP data packets.
	UDPDataBufferSize int
	// UDPConnTTL is the TTL for UDP connections.
	UDPConnTTL time.Duration
}

// BindOption is a functional option for configuring BindOptions.
type BindOption func(opts *BindOptions)

// MuxBindOption sets whether multiplexing is enabled.
func MuxBindOption(mux bool) BindOption {
	return func(opts *BindOptions) {
		opts.Mux = mux
	}
}

// BacklogBindOption sets the accept backlog size.
func BacklogBindOption(backlog int) BindOption {
	return func(opts *BindOptions) {
		opts.Backlog = backlog
	}
}

// UDPDataQueueSizeBindOption sets the UDP data queue size.
func UDPDataQueueSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataQueueSize = size
	}
}

// UDPDataBufferSizeBindOption sets the UDP data buffer size.
func UDPDataBufferSizeBindOption(size int) BindOption {
	return func(opts *BindOptions) {
		opts.UDPDataBufferSize = size
	}
}

// UDPConnTTLBindOption sets the UDP connection TTL.
func UDPConnTTLBindOption(ttl time.Duration) BindOption {
	return func(opts *BindOptions) {
		opts.UDPConnTTL = ttl
	}
}
