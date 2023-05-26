package registry

import (
	"context"
	"net"

	"github.com/jxo-me/netx/sdk/core/resolver"
)

type ResolverRegistry struct {
	registry[resolver.IResolver]
}

func (r *ResolverRegistry) Register(name string, v resolver.IResolver) error {
	return r.registry.Register(name, v)
}

func (r *ResolverRegistry) Get(name string) resolver.IResolver {
	if name != "" {
		return &resolverWrapper{name: name, r: r}
	}
	return nil
}

func (r *ResolverRegistry) get(name string) resolver.IResolver {
	return r.registry.Get(name)
}

type resolverWrapper struct {
	name string
	r    *ResolverRegistry
}

func (w *resolverWrapper) Resolve(ctx context.Context, network, host string) ([]net.IP, error) {
	r := w.r.get(w.name)
	if r == nil {
		return nil, resolver.ErrInvalid
	}
	return r.Resolve(ctx, network, host)
}
