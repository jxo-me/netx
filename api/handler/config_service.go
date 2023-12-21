package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/service"
)

var (
	Service = hService{}
)

type hService struct{}

type CreateServiceReq struct {
	g.Meta `path:"/services" method:"post" tags:"Services" summary:"Create a new service, the name of the service must be unique in service list."`
	// in: body
	Data config.ServiceConfig `json:"data"`
}

func (h *hService) CreateService(ctx context.Context, req *CreateServiceReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}

	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	if app.Runtime.ServiceRegistry().IsRegistered(req.Data.Name) {
		return nil, ErrDup
	}

	svc, err := parser.ParseService(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	if err := app.Runtime.ServiceRegistry().Register(req.Data.Name, svc); err != nil {
		svc.Close()
		return nil, ErrDup
	}

	go svc.Serve()

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Services = append(c.Services, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateServiceReq struct {
	g.Meta `path:"/services/{service}" method:"put" tags:"Services" summary:"Update service by name, the service must already exist."`
	// in: path
	// required: true
	Service string `json:"service"`
	// in: body
	Data config.ServiceConfig `json:"data"`
}

func (h *hService) UpdateService(ctx context.Context, req *UpdateServiceReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	old := app.Runtime.ServiceRegistry().Get(req.Service)
	if old == nil {
		return nil, ErrInvalid
	}
	_ = old.Close()

	req.Data.Name = req.Service

	svc, err := parser.ParseService(&req.Data)
	if err != nil {
		return nil, ErrCreate
	}

	app.Runtime.ServiceRegistry().Unregister(req.Service)

	if err := app.Runtime.ServiceRegistry().Register(req.Service, svc); err != nil {
		_ = svc.Close()
		return nil, ErrDup
	}

	go svc.Serve()

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Services {
			if c.Services[i].Name == req.Service {
				c.Services[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteServiceReq struct {
	g.Meta `path:"/services/{service}" method:"delete" tags:"Services" summary:"Delete service by name."`
	// in: path
	// required: true
	Service string `json:"service"`
}

func (h *hService) DeleteService(ctx context.Context, req *DeleteServiceReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	svc := app.Runtime.ServiceRegistry().Get(req.Service)
	if svc == nil {
		return nil, ErrNotFound
	}

	app.Runtime.ServiceRegistry().Unregister(req.Service)
	svc.Close()

	_ = config.OnUpdate(func(c *config.Config) error {
		services := c.Services
		c.Services = nil
		for _, s := range services {
			if s.Name == req.Service {
				continue
			}
			c.Services = append(c.Services, s)
		}
		return nil
	})

	return res, nil
}
