package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/limiter"
)

var (
	Limiter = hLimiter{}
)

type hLimiter struct{}

type CreateLimiterReq struct {
	g.Meta `path:"/limiters" method:"post" tags:"Limiter" summary:"Create a new limiter, the name of limiter must be unique in limiter list."`
	// in: body
	Data config.LimiterConfig `json:"data"`
}

func (h *hLimiter) CreateLimiter(ctx context.Context, req *CreateLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parser.ParseTrafficLimiter(&req.Data)

	if err := app.Runtime.TrafficLimiterRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Limiters = append(c.Limiters, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateLimiterReq struct {
	g.Meta `path:"/limiters/{limiter}" method:"put" tags:"Limiter" summary:"Update limiter by name, the limiter must already exist."`
	// in: path
	// required: true
	Limiter string `json:"limiter"`
	// in: body
	Data config.LimiterConfig `json:"data"`
}

func (h *hLimiter) UpdateLimiter(ctx context.Context, req *UpdateLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.TrafficLimiterRegistry().IsRegistered(req.Limiter) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Limiter

	v := parser.ParseTrafficLimiter(&req.Data)

	app.Runtime.TrafficLimiterRegistry().Unregister(req.Limiter)

	if err := app.Runtime.TrafficLimiterRegistry().Register(req.Limiter, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Limiters {
			if c.Limiters[i].Name == req.Limiter {
				c.Limiters[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteLimiterReq struct {
	g.Meta `path:"/limiters/{limiter}" method:"delete" tags:"Limiter" summary:"Delete limiter by name."`
	// in: path
	// required: true
	Limiter string `json:"limiter"`
}

func (h *hLimiter) DeleteLimiter(ctx context.Context, req *DeleteLimiterReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.TrafficLimiterRegistry().IsRegistered(req.Limiter) {
		return nil, ErrNotFound
	}
	app.Runtime.TrafficLimiterRegistry().Unregister(req.Limiter)

	_ = config.OnUpdate(func(c *config.Config) error {
		limiteres := c.Limiters
		c.Limiters = nil
		for _, s := range limiteres {
			if s.Name == req.Limiter {
				continue
			}
			c.Limiters = append(c.Limiters, s)
		}
		return nil
	})

	return res, nil
}
