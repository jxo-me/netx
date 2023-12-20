package registry

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/resolver"
)

type resolverRegistry struct {
	registry[resolver.IResolver]
}

func (r *resolverRegistry) Register(name string, v resolver.IResolver) error {
	return r.registry.Register(name, v)
}

func (r *resolverRegistry) Get(name string) resolver.IResolver {
	if name != "" {
		return &resolverWrapper{name: name, r: r}
	}
	return nil
}

func (r *resolverRegistry) get(name string) resolver.IResolver {
	return r.registry.Get(name)
}

type resolverWrapper struct {
	name string
	r    *resolverRegistry
}

func (w *resolverWrapper) Resolve(ctx context.Context, network, host string, opts ...resolver.Option) ([]net.IP, error) {
	r := w.r.get(w.name)
	if r == nil {
		return nil, resolver.ErrInvalid
	}
	return r.Resolve(ctx, network, host, opts...)
}
