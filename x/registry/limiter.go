package registry

import (
	"github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/limiter/rate"
	"github.com/jxo-me/netx/core/limiter/traffic"
)

type TrafficLimiterRegistry struct {
	registry[traffic.ITrafficLimiter]
}

func (r *TrafficLimiterRegistry) Register(name string, v traffic.ITrafficLimiter) error {
	return r.registry.Register(name, v)
}

func (r *TrafficLimiterRegistry) Get(name string) traffic.ITrafficLimiter {
	if name != "" {
		return &trafficLimiterWrapper{name: name, r: r}
	}
	return nil
}

func (r *TrafficLimiterRegistry) get(name string) traffic.ITrafficLimiter {
	return r.registry.Get(name)
}

type trafficLimiterWrapper struct {
	name string
	r    *TrafficLimiterRegistry
}

func (w *trafficLimiterWrapper) In(key string) traffic.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.In(key)
}

func (w *trafficLimiterWrapper) Out(key string) traffic.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Out(key)
}

type ConnLimiterRegistry struct {
	registry[conn.IConnLimiter]
}

func (r *ConnLimiterRegistry) Register(name string, v conn.IConnLimiter) error {
	return r.registry.Register(name, v)
}

func (r *ConnLimiterRegistry) Get(name string) conn.IConnLimiter {
	if name != "" {
		return &connLimiterWrapper{name: name, r: r}
	}
	return nil
}

func (r *ConnLimiterRegistry) get(name string) conn.IConnLimiter {
	return r.registry.Get(name)
}

type connLimiterWrapper struct {
	name string
	r    *ConnLimiterRegistry
}

func (w *connLimiterWrapper) Limiter(key string) conn.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Limiter(key)
}

type RateLimiterRegistry struct {
	registry[rate.IRateLimiter]
}

func (r *RateLimiterRegistry) Register(name string, v rate.IRateLimiter) error {
	return r.registry.Register(name, v)
}

func (r *RateLimiterRegistry) Get(name string) rate.IRateLimiter {
	if name != "" {
		return &rateLimiterWrapper{name: name, r: r}
	}
	return nil
}

func (r *RateLimiterRegistry) get(name string) rate.IRateLimiter {
	return r.registry.Get(name)
}

type rateLimiterWrapper struct {
	name string
	r    *RateLimiterRegistry
}

func (w *rateLimiterWrapper) Limiter(key string) rate.ILimiter {
	v := w.r.get(w.name)
	if v == nil {
		return nil
	}
	return v.Limiter(key)
}
