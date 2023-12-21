package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/hosts"
)

var (
	Hosts = hHosts{}
)

type hHosts struct{}

type CreateHostsReq struct {
	g.Meta `path:"/hosts" method:"post" tags:"Hosts" summary:"Create a new hosts, the name of the hosts must be unique in hosts list."`
	// in: body
	Data config.HostsConfig `json:"data"`
}

func (h *hHosts) CreateHosts(ctx context.Context, req *CreateHostsReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if req.Data.Name == "" {
		return nil, ErrInvalid
	}

	v := parser.ParseHostMapper(&req.Data)

	if err := app.Runtime.HostsRegistry().Register(req.Data.Name, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		c.Hosts = append(c.Hosts, &req.Data)
		return nil
	})

	return res, nil
}

type UpdateHostsReq struct {
	g.Meta `path:"/hosts/{hosts}" method:"put" tags:"Hosts" summary:"Update hosts by name, the hosts must already exist."`
	// in: path
	// required: true
	Hosts string `json:"hosts"`
	// in: body
	Data config.HostsConfig `json:"data"`
}

func (h *hHosts) UpdateHosts(ctx context.Context, req *UpdateHostsReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.HostsRegistry().IsRegistered(req.Hosts) {
		return nil, ErrNotFound
	}

	req.Data.Name = req.Hosts

	v := parser.ParseHostMapper(&req.Data)

	app.Runtime.HostsRegistry().Unregister(req.Hosts)

	if err := app.Runtime.HostsRegistry().Register(req.Hosts, v); err != nil {
		return nil, ErrDup
	}

	_ = config.OnUpdate(func(c *config.Config) error {
		for i := range c.Hosts {
			if c.Hosts[i].Name == req.Hosts {
				c.Hosts[i] = &req.Data
				break
			}
		}
		return nil
	})

	return res, nil
}

type DeleteHostsReq struct {
	g.Meta `path:"/hosts/{hosts}" method:"delete" tags:"Hosts" summary:"Delete hosts by name."`
	// in: path
	// required: true
	Hosts string `json:"hosts"`
}

func (h *hHosts) DeleteHosts(ctx context.Context, req *DeleteHostsReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	if !app.Runtime.HostsRegistry().IsRegistered(req.Hosts) {
		return nil, ErrNotFound
	}
	app.Runtime.HostsRegistry().Unregister(req.Hosts)

	_ = config.OnUpdate(func(c *config.Config) error {
		hosts := c.Hosts
		c.Hosts = nil
		for _, s := range hosts {
			if s.Name == req.Hosts {
				continue
			}
			c.Hosts = append(c.Hosts, s)
		}
		return nil
	})

	return res, nil
}
