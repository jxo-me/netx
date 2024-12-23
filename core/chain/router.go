package chain

import (
	"context"
	"net"
	"time"

	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/core/resolver"
)

type SockOpts struct {
	Mark int
}

type RouterOptions struct {
	Retries    int
	Timeout    time.Duration
	IfceName   string
	Netns      string
	SockOpts   *SockOpts
	Chain      IChainer
	Resolver   resolver.IResolver
	HostMapper hosts.IHostMapper
	Recorders  []recorder.RecorderObject
	Logger     logger.ILogger
}

type RouterOption func(*RouterOptions)

func InterfaceRouterOption(ifceName string) RouterOption {
	return func(o *RouterOptions) {
		o.IfceName = ifceName
	}
}

func NetnsRouterOption(netns string) RouterOption {
	return func(o *RouterOptions) {
		o.Netns = netns
	}
}

func SockOptsRouterOption(so *SockOpts) RouterOption {
	return func(o *RouterOptions) {
		o.SockOpts = so
	}
}

func TimeoutRouterOption(timeout time.Duration) RouterOption {
	return func(o *RouterOptions) {
		o.Timeout = timeout
	}
}

func RetriesRouterOption(retries int) RouterOption {
	return func(o *RouterOptions) {
		o.Retries = retries
	}
}

func ChainRouterOption(chain IChainer) RouterOption {
	return func(o *RouterOptions) {
		o.Chain = chain
	}
}

func ResolverRouterOption(resolver resolver.IResolver) RouterOption {
	return func(o *RouterOptions) {
		o.Resolver = resolver
	}
}

func HostMapperRouterOption(m hosts.IHostMapper) RouterOption {
	return func(o *RouterOptions) {
		o.HostMapper = m
	}
}

func RecordersRouterOption(recorders ...recorder.RecorderObject) RouterOption {
	return func(o *RouterOptions) {
		o.Recorders = recorders
	}
}

func LoggerRouterOption(logger logger.ILogger) RouterOption {
	return func(o *RouterOptions) {
		o.Logger = logger
	}
}

type Router interface {
	Options() *RouterOptions
	Dial(ctx context.Context, network, address string) (net.Conn, error)
	Bind(ctx context.Context, network, address string, opts ...BindOption) (net.Listener, error)
}
