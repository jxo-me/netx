package hosts

import (
	"context"
	"net"
)

type Options struct{}
type Option func(opts *Options)
// IHostMapper is a mapping from hostname to IP.
type IHostMapper interface {
	Lookup(ctx context.Context, network, host string, opts ...Option) ([]net.IP, bool)
}
