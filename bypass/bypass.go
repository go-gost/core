package bypass

import (
	"context"
	"slices"
)

type Options struct {
	Host string
	Path string
}

type Option func(opts *Options)

func WithHostOpton(host string) Option {
	return func(opts *Options) {
		opts.Host = host
	}
}

func WithPathOption(path string) Option {
	return func(opts *Options) {
		opts.Path = path
	}
}

// Bypass is a filter of address (IP or domain).
type Bypass interface {
	// Contains reports whether the bypass includes addr.
	IsWhitelist() bool
	Contains(ctx context.Context, network, addr string, opts ...Option) bool
}

type bypassGroup struct {
	bypasses []Bypass
}

func BypassGroup(bypasses ...Bypass) Bypass {
	return &bypassGroup{
		bypasses: bypasses,
	}
}

func (p *bypassGroup) Contains(ctx context.Context, network, addr string, opts ...Option) bool {
	var whitelist, blacklist []bool
	for _, bypass := range p.bypasses {
		result := bypass.Contains(ctx, network, addr, opts...)
		if bypass.IsWhitelist() {
			whitelist = append(whitelist, result)
		} else {
			blacklist = append(blacklist, result)
		}
	}
	status := false
	if len(whitelist) > 0 {
		if slices.Contains(whitelist, false) {
			status = false
		} else {
			status = true
		}
	}
	if !status && len(blacklist) > 0 {
		if slices.Contains(blacklist, true) {
			status = true
		} else {
			status = false
		}
	}
	return status
}

func (p *bypassGroup) IsWhitelist() bool {
	return false
}
