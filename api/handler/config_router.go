package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/router"
)

var (
	Router = hRouter{}
)

type hRouter struct{}

type CreateRouterReq struct {
	g.Meta `path:"/routers" method:"post" tags:"Router" summary:"Create a new router, the name of the router must be unique in router list."`
	// in: body
	Data config.RouterConfig `json:"data"`
}

func (h *hRouter) CreateRouter(ctx context.Context, req *CreateRouterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parser.ParseRouter(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	if err := app.Runtime.RouterRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	config.OnUpdate(func(c *config.Config) error {
		c.Routers = append(c.Routers, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateRouterReq struct {
	g.Meta `path:"/routers/{router}" method:"put" tags:"Router" summary:"Update router by name, the router must already exist."`
	// in: path
	// required: true
	Router string `json:"router"`
	// in: body
	Data config.RouterConfig `json:"data"`
}

func (h *hRouter) UpdateResolver(ctx context.Context, req *UpdateRouterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.RouterRegistry().IsRegistered(req.Router) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Router

	v := parser.ParseRouter(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	app.Runtime.RouterRegistry().Unregister(req.Router)

	if err := app.Runtime.RouterRegistry().Register(req.Router, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Routers {
			if c.Routers[i].Name == req.Router {
				c.Routers[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteRouterReq struct {
	g.Meta `path:"/routers/{router}" method:"delete" tags:"Router" summary:"Delete router by name."`
	// in: path
	// required: true
	Router string `json:"router"`
}

func (h *hRouter) DeleteResolver(ctx context.Context, req *DeleteRouterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.RouterRegistry().IsRegistered(req.Router) {
		return nil, ErrNotFound
	}
	app.Runtime.RouterRegistry().Unregister(req.Router)

	_ = config.OnUpdate(func(c *config.Config) error {
		routers := c.Routers
		c.Resolvers = nil
		for _, s := range routers {
			if s.Name == req.Router {
				continue
			}
			c.Routers = append(c.Routers, s)
		}
		return nil
	})

	return res, nil
}
