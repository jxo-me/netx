package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/resolver"
)

var (
	Resolver = hResolver{}
)

type hResolver struct{}

type CreateResolverReq struct {
	g.Meta `path:"/resolvers" method:"post" tags:"Resolver" summary:"Create a new resolver, the name of the resolver must be unique in resolver list."`
	// in: body
	Data config.ResolverConfig `json:"data"`
}

func (h *hResolver) CreateResolver(ctx context.Context, req *CreateResolverReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v, err := parser.ParseResolver(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	if err := app.Runtime.ResolverRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Resolvers = append(c.Resolvers, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateResolverReq struct {
	g.Meta `path:"/resolvers/{resolver}" method:"put" tags:"Resolver" summary:"Update resolver by name, the resolver must already exist."`
	// in: path
	// required: true
	Resolver string `json:"resolver"`
	// in: body
	Data config.ResolverConfig `json:"data"`
}

func (h *hResolver) UpdateResolver(ctx context.Context, req *UpdateResolverReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.ResolverRegistry().IsRegistered(req.Resolver) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Resolver

	v, err := parser.ParseResolver(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	app.Runtime.ResolverRegistry().Unregister(req.Resolver)

	if err := app.Runtime.ResolverRegistry().Register(req.Resolver, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Resolvers {
			if c.Resolvers[i].Name == req.Resolver {
				c.Resolvers[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteResolverReq struct {
	g.Meta `path:"/resolvers/{resolver}" method:"delete" tags:"Resolver" summary:"Delete resolver by name."`
	// in: path
	// required: true
	Resolver string `json:"resolver"`
}

func (h *hResolver) DeleteResolver(ctx context.Context, req *DeleteResolverReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.ResolverRegistry().IsRegistered(req.Resolver) {
		return nil, ErrNotFound
	}
	app.Runtime.ResolverRegistry().Unregister(req.Resolver)

	_ = config.OnUpdate(func(c *config.Config) error {
		resolvers := c.Resolvers
		c.Resolvers = nil
		for _, s := range resolvers {
			if s.Name == req.Resolver {
				continue
			}
			c.Resolvers = append(c.Resolvers, s)
		}
		return nil
	})

	return res, nil
}
