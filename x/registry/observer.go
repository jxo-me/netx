package registry

import (
	"context"

	"github.com/jxo-me/netx/core/observer"
)

type ObserverRegistry struct {
	registry[observer.IObserver]
}

func (r *ObserverRegistry) Register(name string, v observer.IObserver) error {
	return r.registry.Register(name, v)
}

func (r *ObserverRegistry) Get(name string) observer.IObserver {
	if name != "" {
		return &observerWrapper{name: name, r: r}
	}
	return nil
}

func (r *ObserverRegistry) get(name string) observer.IObserver {
	return r.registry.Get(name)
}

type observerWrapper struct {
	name string
	r    *ObserverRegistry
}

func (w *observerWrapper) Observe(ctx context.Context, events []observer.Event, opts ...observer.Option) error {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Observe(ctx, events, opts...)
}
