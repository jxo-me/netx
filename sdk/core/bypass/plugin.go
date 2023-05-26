package bypass

import (
	"context"

	xlogger "github.com/jxo-me/netx/sdk/core/logger"
	"github.com/jxo-me/netx/sdk/plugin/bypass/proto"
)

type pluginBypass struct {
	client  proto.BypassClient
	options options
}

// NewPluginBypass creates a plugin bypass.
func NewPluginBypass(opts ...Option) IBypass {
	var options options
	for _, opt := range opts {
		opt(&options)
	}
	if options.logger == nil {
		options.logger = xlogger.Nop()
	}

	p := &pluginBypass{
		options: options,
	}
	if options.client != nil {
		p.client = proto.NewBypassClient(options.client)
	}
	return p
}

func (p *pluginBypass) Contains(ctx context.Context, addr string) bool {
	if p.client == nil {
		return false
	}

	r, err := p.client.Bypass(ctx,
		&proto.BypassRequest{
			Addr: addr,
		})
	if err != nil {
		p.options.logger.Error(err)
		return false
	}
	return r.Ok
}

func (p *pluginBypass) Close() error {
	if p.options.client != nil {
		return p.options.client.Close()
	}
	return nil
}
