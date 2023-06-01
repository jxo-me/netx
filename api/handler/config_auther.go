package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
)

var (
	Auther = hAuther{}
)

type hAuther struct{}

type CreateAutherReq struct {
	g.Meta `path:"/authers" method:"post" tags:"Authers" summary:"Create a new auther, the name of the auther must be unique in auther list."`
	// in: body
	Data config.AutherConfig `json:"data"`
}

func (h *hAuther) CreateAuther(ctx context.Context, req *CreateAutherReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parsing.ParseAuther(&req.Data)
	if err := app.Runtime.AutherRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Authers = append(c.Authers, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateAutherReq struct {
	g.Meta `path:"/authers/{auther}" method:"put" tags:"Authers" summary:"Update auther by name, the auther must already exist."`
	// in: path
	// required: true
	Auther string `json:"auther"`
	// in: body
	Data config.AutherConfig `json:"data"`
}

func (h *hAuther) UpdateAuther(ctx context.Context, req *UpdateAutherReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	if !app.Runtime.AutherRegistry().IsRegistered(req.Auther) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Auther

	v := parsing.ParseAuther(&req.Data)
	app.Runtime.AutherRegistry().Unregister(req.Auther)

	if err := app.Runtime.AutherRegistry().Register(req.Auther, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Authers {
			if c.Authers[i].Name == req.Auther {
				c.Authers[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteAutherReq struct {
	g.Meta `path:"/authers/{auther}" method:"delete" tags:"Authers" summary:"Delete auther by name."`
	// in: path
	// required: true
	Auther string `json:"auther"`
}

func (h *hAuther) DeleteAuther(ctx context.Context, req *DeleteAutherReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.AutherRegistry().IsRegistered(req.Auther) {
		return nil, ErrNotFound
	}
	app.Runtime.AutherRegistry().Unregister(req.Auther)

	_ = config.OnUpdate(func(c *config.Config) error {
		authers := c.Authers
		c.Authers = nil
		for _, s := range authers {
			if s.Name == req.Auther {
				continue
			}
			c.Authers = append(c.Authers, s)
		}
		return nil
	})

	return res, nil
}
