package registry

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/hosts"
)

type HostsRegistry struct {
	registry[hosts.IHostMapper]
}

func (r *HostsRegistry) Register(name string, v hosts.IHostMapper) error {
	return r.registry.Register(name, v)
}

func (r *HostsRegistry) Get(name string) hosts.IHostMapper {
	if name != "" {
		return &hostsWrapper{name: name, r: r}
	}
	return nil
}

func (r *HostsRegistry) get(name string) hosts.IHostMapper {
	return r.registry.Get(name)
}

type hostsWrapper struct {
	name string
	r    *HostsRegistry
}

func (w *hostsWrapper) Lookup(ctx context.Context, network, host string, opts ...hosts.Option) ([]net.IP, bool) {
	v := w.r.get(w.name)
	if v == nil {
		return nil, false
	}
	return v.Lookup(ctx, network, host, opts...)
}
