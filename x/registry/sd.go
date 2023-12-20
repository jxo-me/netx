package registry

import (
	"context"

	"github.com/jxo-me/netx/core/sd"
)

type SdRegistry struct {
	registry[sd.ISD]
}

func (r *SdRegistry) Register(name string, v sd.ISD) error {
	return r.registry.Register(name, v)
}

func (r *SdRegistry) Get(name string) sd.ISD {
	if name != "" {
		return &sdWrapper{name: name, r: r}
	}
	return nil
}

func (r *SdRegistry) get(name string) sd.ISD {
	return r.registry.Get(name)
}

type sdWrapper struct {
	name string
	r    *SdRegistry
}

func (w *sdWrapper) Register(ctx context.Context, service *sd.Service, opts ...sd.Option) error {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Register(ctx, service, opts...)
}

func (w *sdWrapper) Deregister(ctx context.Context, service *sd.Service) error {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}

	return v.Deregister(ctx, service)
}

func (w *sdWrapper) Renew(ctx context.Context, service *sd.Service) error {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}

	return v.Renew(ctx, service)
}

func (w *sdWrapper) Get(ctx context.Context, name string) ([]*sd.Service, error) {
	v := w.r.get(w.name)
	if v == nil {
		return nil, nil
	}

	return v.Get(ctx, name)
}
