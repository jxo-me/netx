package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/ingress"
)

var (
	Ingress = hIngress{}
)

type hIngress struct{}

type CreateIngressReq struct {
	g.Meta `path:"/ingresses" method:"post" tags:"Ingress" summary:"Create a new ingress, the name of the ingress must be unique in ingress list."`
	// in: body
	Data config.IngressConfig `json:"data"`
}

func (h *hIngress) CreateIngress(ctx context.Context, req *CreateIngressReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parser.ParseIngress(&req.Data)

	if err := app.Runtime.IngressRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Ingresses = append(c.Ingresses, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateIngressReq struct {
	g.Meta `path:"/ingresses/{ingress}" method:"put" tags:"Ingress" summary:"Update ingress by name, the ingress must already exist."`
	// in: path
	// required: true
	Ingress string `json:"ingress"`
	// in: body
	Data config.IngressConfig `json:"data"`
}

func (h *hIngress) UpdateIngress(ctx context.Context, req *UpdateIngressReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	if !app.Runtime.IngressRegistry().IsRegistered(req.Ingress) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Ingress

	v := parser.ParseIngress(&req.Data)

	app.Runtime.IngressRegistry().Unregister(req.Ingress)

	if err := app.Runtime.IngressRegistry().Register(req.Ingress, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Ingresses {
			if c.Ingresses[i].Name == req.Ingress {
				c.Ingresses[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteIngressReq struct {
	g.Meta `path:"/ingresses/{ingress}" method:"delete" tags:"Ingress" summary:"Delete ingress by name."`
	// in: path
	// required: true
	Ingress string `json:"ingress"`
}

func (h *hIngress) DeleteIngress(ctx context.Context, req *DeleteIngressReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.IngressRegistry().IsRegistered(req.Ingress) {
		return nil, ErrNotFound
	}
	app.Runtime.IngressRegistry().Unregister(req.Ingress)

	_ = config.OnUpdate(func(c *config.Config) error {
		ingresses := c.Ingresses
		c.Ingresses = nil
		for _, s := range ingresses {
			if s.Name == req.Ingress {
				continue
			}
			c.Ingresses = append(c.Ingresses, s)
		}
		return nil
	})

	return res, nil
}
