package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
)

var (
	RateLimiter = hRateLimiter{}
)

type hRateLimiter struct{}

type CreateRateLimiterReq struct {
	g.Meta `path:"/rlimiters" method:"post" tags:"RateLimiter" summary:"Create a new rate limiter, the name of limiter must be unique in limiter list."`
	// in: body
	Data config.LimiterConfig `json:"data"`
}

func (h *hRateLimiter) CreateRateLimiter(ctx context.Context, req *CreateRateLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parsing.ParseRateLimiter(&req.Data)

	if err := app.Runtime.RateLimiterRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.RLimiters = append(c.RLimiters, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateRateLimiterReq struct {
	g.Meta `path:"/rlimiters/{limiter}" method:"put" tags:"RateLimiter" summary:"Update rate limiter by name, the limiter must already exist."`
	// in: path
	// required: true
	Limiter string `json:"limiter"`
	// in: body
	Data config.LimiterConfig `json:"data"`
}

func (h *hRateLimiter) UpdateRateLimiter(ctx context.Context, req *UpdateRateLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.RateLimiterRegistry().IsRegistered(req.Limiter) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Limiter

	v := parsing.ParseRateLimiter(&req.Data)

	app.Runtime.RateLimiterRegistry().Unregister(req.Limiter)

	if err := app.Runtime.RateLimiterRegistry().Register(req.Limiter, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.RLimiters {
			if c.RLimiters[i].Name == req.Limiter {
				c.RLimiters[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteRateLimiterReq struct {
	g.Meta `path:"/rlimiters/{limiter}" method:"delete" tags:"RateLimiter" summary:"Delete rate limiter by name."`
	// in: path
	// required: true
	Limiter string `json:"limiter"`
}

func (h *hRateLimiter) DeleteRateLimiter(ctx context.Context, req *DeleteRateLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.RateLimiterRegistry().IsRegistered(req.Limiter) {
		return nil, ErrNotFound
	}
	app.Runtime.RateLimiterRegistry().Unregister(req.Limiter)

	_ = config.OnUpdate(func(c *config.Config) error {
		limiteres := c.RLimiters
		c.RLimiters = nil
		for _, s := range limiteres {
			if s.Name == req.Limiter {
				continue
			}
			c.RLimiters = append(c.RLimiters, s)
		}
		return nil
	})

	return res, nil
}
