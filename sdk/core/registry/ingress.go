package registry

import (
	"context"

	"github.com/jxo-me/netx/sdk/core/ingress"
)

type IngressRegistry struct {
	registry[ingress.IIngress]
}

func (r *IngressRegistry) Register(name string, v ingress.IIngress) error {
	return r.registry.Register(name, v)
}

func (r *IngressRegistry) Get(name string) ingress.IIngress {
	if name != "" {
		return &ingressWrapper{name: name, r: r}
	}
	return nil
}

func (r *IngressRegistry) get(name string) ingress.IIngress {
	return r.registry.Get(name)
}

type ingressWrapper struct {
	name string
	r    *IngressRegistry
}

func (w *ingressWrapper) Get(ctx context.Context, host string) string {
	v := w.r.get(w.name)
	if v == nil {
		return ""
	}
	return v.Get(ctx, host)
}
