package chain

import (
	"context"

	"github.com/jxo-me/netx/sdk/core/logger"
	"github.com/jxo-me/netx/sdk/core/metadata"
	"github.com/jxo-me/netx/sdk/core/selector"
)

var (
	_ IChainer = (*chainGroup)(nil)
)

type ChainOptions struct {
	Metadata metadata.IMetaData
	Logger   logger.ILogger
}

type ChainOption func(*ChainOptions)

func MetadataChainOption(md metadata.IMetaData) ChainOption {
	return func(opts *ChainOptions) {
		opts.Metadata = md
	}
}

func LoggerChainOption(logger logger.ILogger) ChainOption {
	return func(opts *ChainOptions) {
		opts.Logger = logger
	}
}

type chainNamer interface {
	Name() string
}

type Chain struct {
	name     string
	hops     []IHop
	marker   selector.Marker
	metadata metadata.IMetaData
	logger   logger.ILogger
}

func NewChain(name string, opts ...ChainOption) *Chain {
	var options ChainOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	return &Chain{
		name:     name,
		metadata: options.Metadata,
		marker:   selector.NewFailMarker(),
		logger:   options.Logger,
	}
}

func (c *Chain) AddHop(hop IHop) {
	c.hops = append(c.hops, hop)
}

// Metadata implements metadata.Metadatable interface.
func (c *Chain) Metadata() metadata.IMetaData {
	return c.metadata
}

// Marker implements selector.Markable interface.
func (c *Chain) Marker() selector.Marker {
	return c.marker
}

func (c *Chain) Name() string {
	return c.name
}

func (c *Chain) Route(ctx context.Context, network, address string) IRoute {
	if c == nil || len(c.hops) == 0 {
		return nil
	}

	rt := NewRoute(ChainRouteOption(c))
	for _, hop := range c.hops {
		node := hop.Select(ctx, AddrSelectOption(address))
		if node == nil {
			return rt
		}
		if node.Options().Transport.Multiplex() {
			tr := node.Options().Transport.Copy()
			tr.Options().Route = rt
			node = node.Copy()
			node.Options().Transport = tr
			rt = NewRoute()
		}

		rt.addNode(node)
	}
	return rt
}

type chainGroup struct {
	chains   []IChainer
	selector selector.Selector[IChainer]
}

func NewChainGroup(chains ...IChainer) *chainGroup {
	return &chainGroup{chains: chains}
}

func (p *chainGroup) WithSelector(s selector.Selector[IChainer]) *chainGroup {
	p.selector = s
	return p
}

func (p *chainGroup) Route(ctx context.Context, network, address string) IRoute {
	if chain := p.next(ctx); chain != nil {
		return chain.Route(ctx, network, address)
	}
	return nil
}

func (p *chainGroup) next(ctx context.Context) IChainer {
	if p == nil || len(p.chains) == 0 {
		return nil
	}

	return p.selector.Select(ctx, p.chains...)
}
