package registry

import (
	"context"

	"github.com/jxo-me/netx/core/router"
)

type RouterRegistry struct {
	registry[router.IRouter]
}

func (r *RouterRegistry) Register(name string, v router.IRouter) error {
	return r.registry.Register(name, v)
}

func (r *RouterRegistry) Get(name string) router.IRouter {
	if name != "" {
		return &routerWrapper{name: name, r: r}
	}
	return nil
}

func (r *RouterRegistry) get(name string) router.IRouter {
	return r.registry.Get(name)
}

type routerWrapper struct {
	name string
	r    *RouterRegistry
}

func (w *routerWrapper) GetRoute(ctx context.Context, dst string, opts ...router.Option) *router.Route {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.GetRoute(ctx, dst, opts...)
}
