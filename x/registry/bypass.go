package registry

import (
	"context"

	"github.com/jxo-me/netx/core/bypass"
)

type BypassRegistry struct {
	registry[bypass.IBypass]
}

func (r *BypassRegistry) Register(name string, v bypass.IBypass) error {
	return r.registry.Register(name, v)
}

func (r *BypassRegistry) Get(name string) bypass.IBypass {
	if name != "" {
		return &bypassWrapper{name: name, r: r}
	}
	return nil
}

func (r *BypassRegistry) get(name string) bypass.IBypass {
	return r.registry.Get(name)
}

type bypassWrapper struct {
	name string
	r    *BypassRegistry
}

func (w *bypassWrapper) Contains(ctx context.Context, network, addr string, opts ...bypass.Option) bool {
	bp := w.r.get(w.name)
	if bp == nil {
		return false
	}
	return bp.Contains(ctx, network, addr, opts...)
}

func (w *bypassWrapper) IsWhitelist() bool {
	bp := w.r.get(w.name)
	if bp == nil {
		return false
	}
	return bp.IsWhitelist()
}
