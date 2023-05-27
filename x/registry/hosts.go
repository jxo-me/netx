package registry

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/hosts"
)

type hostsRegistry struct {
	registry[hosts.IHostMapper]
}

func (r *hostsRegistry) Register(name string, v hosts.IHostMapper) error {
	return r.registry.Register(name, v)
}

func (r *hostsRegistry) Get(name string) hosts.IHostMapper {
	if name != "" {
		return &hostsWrapper{name: name, r: r}
	}
	return nil
}

func (r *hostsRegistry) get(name string) hosts.IHostMapper {
	return r.registry.Get(name)
}

type hostsWrapper struct {
	name string
	r    *hostsRegistry
}

func (w *hostsWrapper) Lookup(ctx context.Context, network, host string) ([]net.IP, bool) {
	v := w.r.get(w.name)
	if v == nil {
		return nil, false
	}
	return v.Lookup(ctx, network, host)
}
