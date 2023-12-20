package traffic

import "context"

type ILimiter interface {
	// Wait blocks with the requested n and returns the result value,
	// the returned value is less or equal to n.
	Wait(ctx context.Context, n int) int
	Limit() int
	Set(n int)
}

type Options struct {
	Network string
	Addr    string
	Client  string
	Src     string
}

type Option func(opts *Options)

func NetworkOption(network string) Option {
	return func(opts *Options) {
		opts.Network = network
	}
}

func AddrOption(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

func ClientOption(client string) Option {
	return func(opts *Options) {
		opts.Client = client
	}
}

func SrcOption(src string) Option {
	return func(opts *Options) {
		opts.Src = src
	}
}

type ITrafficLimiter interface {
	In(ctx context.Context, key string, opts ...Option) ILimiter
	Out(ctx context.Context, key string, opts ...Option) ILimiter
}
