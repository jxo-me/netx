package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/chain"
)

var (
	Chain = hChain{}
)

type hChain struct{}

type CreateChainReq struct {
	g.Meta `path:"/chains" method:"post" tags:"Chains" summary:"Create a new chain, the name of chain must be unique in chain list."`
	// in: body
	Data config.ChainConfig `json:"data"`
}

func (h *hChain) CreateChain(ctx context.Context, req *CreateChainReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v, err := parser.ParseChain(&req.Data, logger.Default())
	if err != nil {
		return nil, ErrCreate
	}

	if err := app.Runtime.ChainRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Chains = append(c.Chains, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateChainReq struct {
	g.Meta `path:"/chains/{chain}" method:"put" tags:"Chains" summary:"Update chain by name, the chain must already exist."`
	// in: path
	// required: true
	Chain string `json:"chain"`
	// in: body
	Data config.ChainConfig `json:"data"`
}

func (h *hChain) UpdateChain(ctx context.Context, req *UpdateChainReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.ChainRegistry().IsRegistered(req.Chain) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Chain

	v, err := parser.ParseChain(&req.Data, logger.Default())
	if err != nil {
		return nil, ErrCreate
	}

	app.Runtime.ChainRegistry().Unregister(req.Chain)

	if err := app.Runtime.ChainRegistry().Register(req.Chain, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Chains {
			if c.Chains[i].Name == req.Chain {
				c.Chains[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteChainReq struct {
	g.Meta `path:"/chains/{chain}" method:"delete" tags:"Chains" summary:"Delete chain by name."`
	// in: path
	// required: true
	Chain string `json:"chain"`
}

func (h *hChain) DeleteChain(ctx context.Context, req *DeleteChainReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.ChainRegistry().IsRegistered(req.Chain) {
		return nil, ErrNotFound
	}
	app.Runtime.ChainRegistry().Unregister(req.Chain)

	_ = config.OnUpdate(func(c *config.Config) error {
		chains := c.Chains
		c.Chains = nil
		for _, s := range chains {
			if s.Name == req.Chain {
				continue
			}
			c.Chains = append(c.Chains, s)
		}
		return nil
	})

	return res, nil
}
