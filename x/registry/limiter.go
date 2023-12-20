package registry

import (
	"context"

	"github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/limiter/traffic"
)

type trafficLimiterRegistry struct {
	registry[traffic.ITrafficLimiter]
}

func (r *trafficLimiterRegistry) Register(name string, v traffic.ITrafficLimiter) error {
	return r.registry.Register(name, v)
}

func (r *trafficLimiterRegistry) Get(name string) traffic.ITrafficLimiter {
	if name != "" {
		return &trafficLimiterWrapper{name: name, r: r}
	}
	return nil
}

func (r *trafficLimiterRegistry) get(name string) traffic.ITrafficLimiter {
	return r.registry.Get(name)
}

type trafficLimiterWrapper struct {
	name string
	r    *trafficLimiterRegistry
}

func (w *trafficLimiterWrapper) In(ctx context.Context, key string, opts ...traffic.Option) traffic.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.In(ctx, key, opts...)
}

func (w *trafficLimiterWrapper) Out(ctx context.Context, key string, opts ...traffic.Option) traffic.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Out(ctx, key, opts...)
}

type connLimiterRegistry struct {
	registry[conn.IConnLimiter]
}

func (r *connLimiterRegistry) Register(name string, v conn.IConnLimiter) error {
	return r.registry.Register(name, v)
}

func (r *connLimiterRegistry) Get(name string) conn.IConnLimiter {
	if name != "" {
		return &connLimiterWrapper{name: name, r: r}
	}
	return nil
}

func (r *connLimiterRegistry) get(name string) conn.IConnLimiter {
	return r.registry.Get(name)
}

type connLimiterWrapper struct {
	name string
	r    *connLimiterRegistry
}

func (w *connLimiterWrapper) Limiter(key string) conn.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Limiter(key)
}

type rateLimiterRegistry struct {
	registry[rate.IRateLimiter]
}

func (r *rateLimiterRegistry) Register(name string, v rate.IRateLimiter) error {
	return r.registry.Register(name, v)
}

func (r *rateLimiterRegistry) Get(name string) rate.IRateLimiter {
	if name != "" {
		return &rateLimiterWrapper{name: name, r: r}
	}
	return nil
}

func (r *rateLimiterRegistry) get(name string) rate.IRateLimiter {
	return r.registry.Get(name)
}

type rateLimiterWrapper struct {
	name string
	r    *rateLimiterRegistry
}

func (w *rateLimiterWrapper) Limiter(key string) rate.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Limiter(key)
}
