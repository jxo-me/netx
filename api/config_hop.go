package api

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
)

var (
	Hop = hHop{}
)

type hHop struct{}

type CreateHopReq struct {
	g.Meta `path:"/hops" method:"post" tags:"Hops" summary:"Create a new hop, the name of hop must be unique in hop list."`
	// in: body
	Data config.HopConfig `json:"data"`
}

func (h *hHop) CreateHop(ctx context.Context, req *CreateHopReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v, err := parsing.ParseHop(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	if err := app.Runtime.HopRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Hops = append(c.Hops, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateHopReq struct {
	g.Meta `path:"/hops/{hop}" method:"put" tags:"Hops" summary:"Update hop by name, the hop must already exist."`
	// in: path
	// required: true
	Hop string `json:"hop"`
	// in: body
	Data config.HopConfig `json:"data"`
}

func (h *hHop) UpdateHop(ctx context.Context, req *UpdateHopReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	req.Data.Name = req.Hop

	v, err := parsing.ParseHop(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	app.Runtime.HopRegistry().Unregister(req.Hop)

	if err := app.Runtime.HopRegistry().Register(req.Hop, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Hops {
			if c.Hops[i].Name == req.Hop {
				c.Hops[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteHopReq struct {
	g.Meta `path:"/hops/{hop}" method:"delete" tags:"Hops" summary:"Delete hop by name."`
	// in: path
	// required: true
	Hop string `json:"hop"`
}

func (h *hHop) DeleteHop(ctx context.Context, req *DeleteHopReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.HopRegistry().IsRegistered(req.Hop) {
		return nil, ErrNotFound
	}
	app.Runtime.HopRegistry().Unregister(req.Hop)

	_ = config.OnUpdate(func(c *config.Config) error {
		hops := c.Hops
		c.Hops = nil
		for _, s := range hops {
			if s.Name == req.Hop {
				continue
			}
			c.Hops = append(c.Hops, s)
		}
		return nil
	})

	return res, nil
}
