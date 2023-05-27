package registry

import (
	"context"

	"github.com/jxo-me/netx/core/bypass"
)

type bypassRegistry struct {
	registry[bypass.IBypass]
}

func (r *bypassRegistry) Register(name string, v bypass.IBypass) error {
	return r.registry.Register(name, v)
}

func (r *bypassRegistry) Get(name string) bypass.IBypass {
	if name != "" {
		return &bypassWrapper{name: name, r: r}
	}
	return nil
}

func (r *bypassRegistry) get(name string) bypass.IBypass {
	return r.registry.Get(name)
}

type bypassWrapper struct {
	name string
	r    *bypassRegistry
}

func (w *bypassWrapper) Contains(ctx context.Context, addr string) bool {
	bp := w.r.get(w.name)
	if bp == nil {
		return false
	}
	return bp.Contains(ctx, addr)
}
