package handler

import (
	"crypto/tls"
	"net/url"

	"github.com/go-gost/core/auth"
	"github.com/go-gost/core/bypass"
	"github.com/go-gost/core/chain"
	"github.com/go-gost/core/limiter/rate"
	"github.com/go-gost/core/limiter/traffic"
	"github.com/go-gost/core/logger"
	"github.com/go-gost/core/metadata"
	"github.com/go-gost/core/observer"
	"github.com/go-gost/core/recorder"
	"github.com/go-gost/core/rewriter"
)

// Options holds the initialization parameters for a Handler.
type Options struct {
	// Bypass determines which addresses skip the proxy chain.
	Bypass bypass.Bypass
	// Router is the chain router used for upstream forwarding.
	Router chain.Router
	// Auth is the authentication credentials.
	Auth *url.Userinfo
	// Auther is the authenticator for verifying client credentials.
	Auther auth.Authenticator
	// RateLimiter limits the rate of new connections.
	RateLimiter rate.RateLimiter
	// Limiter is the traffic (bandwidth) limiter.
	Limiter traffic.TrafficLimiter
	// TLSConfig is the TLS configuration for handler connections.
	TLSConfig *tls.Config
	// Logger is the logger for handler operations.
	Logger logger.Logger
	// Observer receives observability events.
	Observer observer.Observer
	// Recorders records traffic data.
	Recorders []recorder.RecorderObject
	// Rewriter rewrites traffic data.
	Rewriter rewriter.Rewriter
	// Service is the service name this handler belongs to.
	Service string
	// Netns is the network namespace name.
	Netns string
}

// Option is a functional option for configuring Options.
type Option func(opts *Options)

// BypassOption sets the Bypass rule.
func BypassOption(bypass bypass.Bypass) Option {
	return func(opts *Options) {
		opts.Bypass = bypass
	}
}

// RouterOption sets the chain Router.
func RouterOption(router chain.Router) Option {
	return func(opts *Options) {
		opts.Router = router
	}
}

// AuthOption sets the authentication credentials.
func AuthOption(auth *url.Userinfo) Option {
	return func(opts *Options) {
		opts.Auth = auth
	}
}

// AutherOption sets the Authenticator.
func AutherOption(auther auth.Authenticator) Option {
	return func(opts *Options) {
		opts.Auther = auther
	}
}

// RateLimiterOption sets the rate limiter.
func RateLimiterOption(limiter rate.RateLimiter) Option {
	return func(opts *Options) {
		opts.RateLimiter = limiter
	}
}

// TrafficLimiterOption sets the traffic limiter.
func TrafficLimiterOption(limiter traffic.TrafficLimiter) Option {
	return func(opts *Options) {
		opts.Limiter = limiter
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

// ObserverOption sets the Observer.
func ObserverOption(observer observer.Observer) Option {
	return func(opts *Options) {
		opts.Observer = observer
	}
}

// RecordersOption sets the traffic Recorders.
func RecordersOption(recorders ...recorder.RecorderObject) Option {
	return func(o *Options) {
		o.Recorders = recorders
	}
}

// RewriterOption sets the Rewriter.
func RewriterOption(rewriter rewriter.Rewriter) Option {
	return func(opts *Options) {
		opts.Rewriter = rewriter
	}
}

// ServiceOption sets the service name.
func ServiceOption(service string) Option {
	return func(opts *Options) {
		opts.Service = service
	}
}

// NetnsOption sets the network namespace.
func NetnsOption(netns string) Option {
	return func(opts *Options) {
		opts.Netns = netns
	}
}

// HandleOptions holds runtime parameters for a Handle call.
type HandleOptions struct {
	// Metadata is per-connection metadata.
	Metadata metadata.Metadata
}

// HandleOption is a functional option for configuring HandleOptions.
type HandleOption func(opts *HandleOptions)

// MetadataHandleOption sets the per-connection metadata.
func MetadataHandleOption(md metadata.Metadata) HandleOption {
	return func(opts *HandleOptions) {
		opts.Metadata = md
	}
}
