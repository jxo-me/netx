package api

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
)

var (
	ConnLimiter = hConnLimiter{}
)

type hConnLimiter struct{}

type CreateConnLimiterReq struct {
	g.Meta `path:"/climiters" method:"post" tags:"ConnLimiter" summary:"Create a new conn limiter, the name of limiter must be unique in limiter list."`
	// in: body
	Data config.LimiterConfig `json:"data"`
}

func (h *hConnLimiter) CreateConnLimiter(ctx context.Context, req *CreateConnLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parsing.ParseConnLimiter(&req.Data)

	if err := app.Runtime.ConnLimiterRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.CLimiters = append(c.CLimiters, &req.Data)
		return nil
	})
	return res, nil
}

type UpdateConnLimiterReq struct {
	g.Meta `path:"/climiters/{limiter}" method:"put" tags:"ConnLimiter" summary:"Update conn limiter by name, the limiter must already exist."`
	// in: path
	// required: true
	Limiter string `json:"limiter"`
	// in: body
	Data config.LimiterConfig `json:"data"`
}

func (h *hConnLimiter) UpdateConnLimiter(ctx context.Context, req *UpdateConnLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.ConnLimiterRegistry().IsRegistered(req.Limiter) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Limiter

	v := parsing.ParseConnLimiter(&req.Data)

	app.Runtime.ConnLimiterRegistry().Unregister(req.Limiter)

	if err := app.Runtime.ConnLimiterRegistry().Register(req.Limiter, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.CLimiters {
			if c.CLimiters[i].Name == req.Limiter {
				c.CLimiters[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteConnLimiterReq struct {
	g.Meta `path:"/climiters/{limiter}" method:"delete" tags:"ConnLimiter" summary:"Delete conn limiter by name."`
	// in: path
	// required: true
	Limiter string `json:"limiter"`
}

func (h *hConnLimiter) DeleteConnLimiter(ctx context.Context, req *DeleteConnLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.ConnLimiterRegistry().IsRegistered(req.Limiter) {
		return nil, ErrNotFound
	}
	app.Runtime.ConnLimiterRegistry().Unregister(req.Limiter)

	_ = config.OnUpdate(func(c *config.Config) error {
		limiteres := c.CLimiters
		c.CLimiters = nil
		for _, s := range limiteres {
			if s.Name == req.Limiter {
				continue
			}
			c.CLimiters = append(c.CLimiters, s)
		}
		return nil
	})

	return res, nil
}
