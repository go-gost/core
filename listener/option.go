package listener

import (
	"crypto/tls"
	"net/url"

	"github.com/go-gost/core/admission"
	"github.com/go-gost/core/auth"
	"github.com/go-gost/core/chain"
	"github.com/go-gost/core/limiter/conn"
	"github.com/go-gost/core/limiter/traffic"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/core/observer/stats"
)

// Options holds the initialization parameters for a Listener.
type Options struct {
	// Addr is the address the listener binds to.
	Addr string
	// Auther is the authenticator for verifying client credentials.
	Auther auth.Authenticator
	// Auth is the authentication credentials required from clients.
	Auth *url.Userinfo
	// TLSConfig is the TLS configuration for the listener.
	TLSConfig *tls.Config
	// Admission controls whether to accept or reject connections.
	Admission admission.Admission
	// TrafficLimiter limits data transfer bandwidth.
	TrafficLimiter traffic.TrafficLimiter
	// ConnLimiter limits the number of concurrent connections.
	ConnLimiter conn.ConnLimiter
	// Chain provides the route chain for upstream forwarding.
	Chain chain.Chainer
	// Stats tracks connection and traffic statistics.
	Stats stats.Stats
	// Logger is the logger for listener operations.
	Logger logger.Logger
	// Service is the service name this listener belongs to.
	Service string
	// ProxyProtocol is the proxy protocol version (0 = disabled).
	ProxyProtocol int
	// Netns is the network namespace name.
	Netns string
	// Router is the chain router for upstream forwarding.
	Router chain.Router
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// AddrOption sets the bind address.
func AddrOption(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

// AutherOption sets the authenticator.
func AutherOption(auther auth.Authenticator) Option {
	return func(opts *Options) {
		opts.Auther = auther
	}
}

// AuthOption sets the required client credentials.
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

// AdmissionOption sets the admission controller.
func AdmissionOption(admission admission.Admission) Option {
	return func(opts *Options) {
		opts.Admission = admission
	}
}

// TrafficLimiterOption sets the traffic limiter.
func TrafficLimiterOption(limiter traffic.TrafficLimiter) Option {
	return func(opts *Options) {
		opts.TrafficLimiter = limiter
	}
}

// ConnLimiterOption sets the connection limiter.
func ConnLimiterOption(limiter conn.ConnLimiter) Option {
	return func(opts *Options) {
		opts.ConnLimiter = limiter
	}
}

// StatsOption sets the statistics tracker.
func StatsOption(stats stats.Stats) Option {
	return func(opts *Options) {
		opts.Stats = stats
	}
}

// LoggerOption sets the logger.
func LoggerOption(logger logger.Logger) Option {
	return func(opts *Options) {
		opts.Logger = logger
	}
}

// ServiceOption sets the service name.
func ServiceOption(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

// ProxyProtocolOption sets the proxy protocol version.
func ProxyProtocolOption(ppv int) Option {
	return func(opts *Options) {
		opts.ProxyProtocol = ppv
	}
}

// NetnsOption sets the network namespace.
func NetnsOption(netns string) Option {
	return func(opts *Options) {
		opts.Netns = netns
	}
}

// RouterOption sets the chain router.
func RouterOption(router chain.Router) Option {
	return func(opts *Options) {
		opts.Router = router
	}
}
