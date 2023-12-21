package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/bypass"
)

var (
	Bypass = hBypass{}
)

type hBypass struct{}

type CreateBypassReq struct {
	g.Meta `path:"/bypasses" method:"post" tags:"Bypass" summary:"Create a new bypass, the name of bypass must be unique in bypass list."`
	// in: body
	Data config.BypassConfig `json:"data"`
}

func (h *hBypass) CreateBypass(ctx context.Context, req *CreateBypassReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parser.ParseBypass(&req.Data)

	if err := app.Runtime.BypassRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Bypasses = append(c.Bypasses, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateBypassReq struct {
	g.Meta `path:"/bypasses/{bypass}" method:"put" tags:"Bypass" summary:"Update bypass by name, the bypass must already exist."`
	// in: path
	// required: true
	Bypass string `json:"bypass"`
	// in: body
	Data config.BypassConfig `json:"data"`
}

func (h *hBypass) UpdateBypass(ctx context.Context, req *UpdateBypassReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.BypassRegistry().IsRegistered(req.Bypass) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Bypass

	v := parser.ParseBypass(&req.Data)

	app.Runtime.BypassRegistry().Unregister(req.Bypass)

	if err := app.Runtime.BypassRegistry().Register(req.Bypass, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Bypasses {
			if c.Bypasses[i].Name == req.Bypass {
				c.Bypasses[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteBypassReq struct {
	g.Meta `path:"/bypasses/{bypass}" method:"delete" tags:"Bypass" summary:"Delete bypass by name."`
	// in: path
	// required: true
	Bypass string `json:"bypass"`
}

func (h *hBypass) DeleteBypass(ctx context.Context, req *DeleteBypassReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	if !app.Runtime.BypassRegistry().IsRegistered(req.Bypass) {
		return nil, ErrNotFound
	}
	app.Runtime.BypassRegistry().Unregister(req.Bypass)

	_ = config.OnUpdate(func(c *config.Config) error {
		bypasses := c.Bypasses
		c.Bypasses = nil
		for _, s := range bypasses {
			if s.Name == req.Bypass {
				continue
			}
			c.Bypasses = append(c.Bypasses, s)
		}
		return nil
	})

	return res, nil
}
