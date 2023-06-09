package chain

import (
	"context"
	"net"
	"strings"

	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/selector"
)

type HopOptions struct {
	bypass   bypass.IBypass
	selector selector.ISelector[*chain.Node]
	logger   logger.ILogger
}

type HopOption func(*HopOptions)

func BypassHopOption(bp bypass.IBypass) HopOption {
	return func(o *HopOptions) {
		o.bypass = bp
	}
}

func SelectorHopOption(s selector.ISelector[*chain.Node]) HopOption {
	return func(o *HopOptions) {
		o.selector = s
	}
}

func LoggerHopOption(logger logger.ILogger) HopOption {
	return func(opts *HopOptions) {
		opts.logger = logger
	}
}

type chainHop struct {
	nodes   []*chain.Node
	options HopOptions
}

func NewChainHop(nodes []*chain.Node, opts ...HopOption) chain.IHop {
	var options HopOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	hop := &chainHop{
		nodes:   nodes,
		options: options,
	}

	return hop
}

func (p *chainHop) Nodes() []*chain.Node {
	return p.nodes
}

func (p *chainHop) Select(ctx context.Context, opts ...chain.SelectOption) *chain.Node {
	var options chain.SelectOptions
	for _, opt := range opts {
		opt(&options)
	}

	if p == nil || len(p.nodes) == 0 {
		return nil
	}

	// hop level bypass
	if p.options.bypass != nil &&
		p.options.bypass.Contains(ctx, options.Addr) {
		return nil
	}

	filters := p.nodes
	if host := options.Host; host != "" {
		filters = nil
		if v, _, _ := net.SplitHostPort(host); v != "" {
			host = v
		}
		var nodes []*chain.Node
		for _, node := range p.nodes {
			if node == nil {
				continue
			}
			vhost := node.Options().Host
			if vhost == "" {
				nodes = append(nodes, node)
				continue
			}
			if vhost == host ||
				vhost[0] == '.' && strings.HasSuffix(host, vhost[1:]) {
				filters = append(filters, node)
			}
		}
		if len(filters) == 0 {
			filters = nodes
		}
	} else if protocol := options.Protocol; protocol != "" {
		filters = nil
		for _, node := range p.nodes {
			if node == nil {
				continue
			}
			if node.Options().Protocol == protocol {
				filters = append(filters, node)
			}
		}
	}

	var nodes []*chain.Node
	for _, node := range filters {
		if node == nil {
			continue
		}
		// node level bypass
		if node.Options().Bypass != nil &&
			node.Options().Bypass.Contains(ctx, options.Addr) {
			continue
		}

		nodes = append(nodes, node)
	}
	if len(nodes) == 0 {
		return nil
	}

	if s := p.options.selector; s != nil {
		return s.Select(ctx, nodes...)
	}
	return nodes[0]
}
