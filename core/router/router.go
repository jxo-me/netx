package router

import (
	"context"
	"net"
)

type Options struct {
	ID string
}

type Option func(opts *Options)

func IDOption(id string) Option {
	return func(opts *Options) {
		opts.ID = id
	}
}

type Route struct {
	// Net is the destination network, e.g. 192.168.0.0/16, 172.10.10.0/24.
	// Deprecated by Dst.
	Net *net.IPNet
	// Dst is the destination.
	Dst string
	// Gateway is the gateway for the destination network.
	Gateway string
}

type IRouter interface {
	// GetRoute queries a route by destination.
	GetRoute(ctx context.Context, dst string, opts ...Option) *Route
}
