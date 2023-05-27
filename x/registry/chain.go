package registry

import (
	"context"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/selector"
)

type chainRegistry struct {
	registry[chain.Chainer]
}

func (r *chainRegistry) Register(name string, v chain.Chainer) error {
	return r.registry.Register(name, v)
}

func (r *chainRegistry) Get(name string) chain.Chainer {
	if name != "" {
		return &chainWrapper{name: name, r: r}
	}
	return nil
}

func (r *chainRegistry) get(name string) chain.Chainer {
	return r.registry.Get(name)
}

type chainWrapper struct {
	name string
	r    *chainRegistry
}

func (w *chainWrapper) Marker() selector.Marker {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	if mi, ok := v.(selector.Markable); ok {
		return mi.Marker()
	}
	return nil
}

func (w *chainWrapper) Metadata() metadata.IMetaData {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}

	if mi, ok := v.(metadata.IMetaDatable); ok {
		return mi.Metadata()
	}
	return nil
}

func (w *chainWrapper) Route(ctx context.Context, network, address string) chain.Route {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Route(ctx, network, address)
}
