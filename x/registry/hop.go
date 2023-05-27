package registry

import (
	"context"

	"github.com/jxo-me/netx/core/chain"
)

type HopRegistry struct {
	registry[chain.IHop]
}

func (r *HopRegistry) Register(name string, v chain.IHop) error {
	return r.registry.Register(name, v)
}

func (r *HopRegistry) Get(name string) chain.IHop {
	if name != "" {
		return &hopWrapper{name: name, r: r}
	}
	return nil
}

func (r *HopRegistry) get(name string) chain.IHop {
	return r.registry.Get(name)
}

type hopWrapper struct {
	name string
	r    *HopRegistry
}

func (w *hopWrapper) Nodes() []*chain.Node {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Nodes()
}

func (w *hopWrapper) Select(ctx context.Context, opts ...chain.SelectOption) *chain.Node {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}

	return v.Select(ctx, opts...)
}
