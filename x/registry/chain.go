package registry

import (
	"context"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/selector"
)

type ChainRegistry struct {
	registry[chain.IChainer]
}

func (r *ChainRegistry) Register(name string, v chain.IChainer) error {
	return r.registry.Register(name, v)
}

func (r *ChainRegistry) Get(name string) chain.IChainer {
	if name != "" {
		return &chainWrapper{name: name, r: r}
	}
	return nil
}

func (r *ChainRegistry) get(name string) chain.IChainer {
	return r.registry.Get(name)
}

type chainWrapper struct {
	name string
	r    *ChainRegistry
}

func (w *chainWrapper) Marker() selector.IMarker {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	if mi, ok := v.(selector.IMarkable); ok {
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

func (w *chainWrapper) Route(ctx context.Context, network, address string) chain.IRoute {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Route(ctx, network, address)
}
