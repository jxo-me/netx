package hosts

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/hosts"
	"github.com/jxo-me/netx/plugin/hosts/proto"
	xlogger "github.com/jxo-me/netx/x/logger"
)

type pluginHostMapper struct {
	client  proto.HostMapperClient
	options options
}

// NewPluginHostMapper creates a plugin IHostMapper.
func NewPluginHostMapper(opts ...Option) hosts.IHostMapper {
	var options options
	for _, opt := range opts {
		opt(&options)
	}
	if options.logger == nil {
		options.logger = xlogger.Nop()
	}

	p := &pluginHostMapper{
		options: options,
	}
	if options.client != nil {
		p.client = proto.NewHostMapperClient(options.client)
	}
	return p
}

func (p *pluginHostMapper) Lookup(ctx context.Context, network, host string) (ips []net.IP, ok bool) {
	p.options.logger.Debugf("lookup %s/%s", host, network)

	if p.client == nil {
		return
	}

	r, err := p.client.Lookup(ctx,
		&proto.LookupRequest{
			Network: network,
			Host:    host,
		})
	if err != nil {
		p.options.logger.Error(err)
		return
	}
	for _, s := range r.Ips {
		if ip := net.ParseIP(s); ip != nil {
			ips = append(ips, ip)
		}
	}
	ok = r.Ok
	return
}

func (p *pluginHostMapper) Close() error {
	if p.options.client != nil {
		return p.options.client.Close()
	}
	return nil
}
