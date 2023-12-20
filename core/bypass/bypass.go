package bypass

import "context"

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
// IBypass is a filter of address (IP or domain).
type IBypass interface {
	// Contains reports whether the bypass includes addr.
	Contains(ctx context.Context, network, addr string, opts ...Option) bool
}

type bypassGroup struct {
	bypasses []IBypass
}

func BypassGroup(bypasses ...IBypass) IBypass {
	return &bypassGroup{
		bypasses: bypasses,
	}
}

func (p *bypassGroup) Contains(ctx context.Context, network, addr string, opts ...Option) bool {
	for _, bypass := range p.bypasses {
		if bypass != nil && bypass.Contains(ctx, network, addr, opts...) {
			return true
		}
	}
	return false
}
