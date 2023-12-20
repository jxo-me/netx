package registry

import (
	"context"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/hop"
)

type HopRegistry struct {
	registry[hop.IHop]
}

func (r *HopRegistry) Register(name string, v hop.IHop) error {
	return r.registry.Register(name, v)
}

func (r *HopRegistry) Get(name string) hop.IHop {
	if name != "" {
		return &hopWrapper{name: name, r: r}
	}
	return nil
}

func (r *HopRegistry) get(name string) hop.IHop {
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
	if nl, ok := v.(hop.NodeList); ok {
		return nl.Nodes()
	}
	return nil
}

func (w *hopWrapper) Select(ctx context.Context, opts ...hop.SelectOption) *chain.Node {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}

	return v.Select(ctx, opts...)
}
