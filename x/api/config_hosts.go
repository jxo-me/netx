package api

import (
	"fmt"
	"github.com/jxo-me/netx/x/app"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jxo-me/netx/x/config"
	parser "github.com/jxo-me/netx/x/config/parsing/hosts"
)

// swagger:parameters createHostsRequest
type createHostsRequest struct {
	// in: body
	Data config.HostsConfig `json:"data"`
}

// successful operation.
// swagger:response createHostsResponse
type createHostsesponse struct {
	Data Response
}

func createHosts(ctx *gin.Context) {
	// swagger:route POST /config/hosts Hosts createHostsRequest
	//
	// Create a new hosts, the name of the hosts must be unique in hosts list.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: createHostsResponse

	var req createHostsRequest
	ctx.ShouldBindJSON(&req.Data)

	name := strings.TrimSpace(req.Data.Name)
	if name == "" {
		writeError(ctx, NewError(http.StatusBadRequest, ErrCodeInvalid, "hosts name is required"))
		return
	}
	req.Data.Name = name

	if app.Runtime.HostsRegistry().IsRegistered(name) {
		writeError(ctx, NewError(http.StatusBadRequest, ErrCodeDup, fmt.Sprintf("hosts %s already exists", name)))
		return
	}

	v := parser.ParseHostMapper(&req.Data)

	if err := app.Runtime.HostsRegistry().Register(name, v); err != nil {
		writeError(ctx, NewError(http.StatusBadRequest, ErrCodeDup, fmt.Sprintf("hosts %s already exists", name)))
		return
	}

	config.OnUpdate(func(c *config.Config) error {
		c.Hosts = append(c.Hosts, &req.Data)
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters updateHostsRequest
type updateHostsRequest struct {
	// in: path
	// required: true
	Hosts string `uri:"hosts" json:"hosts"`
	// in: body
	Data config.HostsConfig `json:"data"`
}

// successful operation.
// swagger:response updateHostsResponse
type updateHostsResponse struct {
	Data Response
}

func updateHosts(ctx *gin.Context) {
	// swagger:route PUT /config/hosts/{hosts} Hosts updateHostsRequest
	//
	// Update hosts by name, the hosts must already exist.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: updateHostsResponse

	var req updateHostsRequest
	ctx.ShouldBindUri(&req)
	ctx.ShouldBindJSON(&req.Data)

	name := strings.TrimSpace(req.Hosts)

	if !app.Runtime.HostsRegistry().IsRegistered(name) {
		writeError(ctx, NewError(http.StatusBadRequest, ErrCodeNotFound, fmt.Sprintf("hosts %s not found", name)))
		return
	}

	req.Data.Name = name

	v := parser.ParseHostMapper(&req.Data)

	app.Runtime.HostsRegistry().Unregister(name)

	if err := app.Runtime.HostsRegistry().Register(name, v); err != nil {
		writeError(ctx, NewError(http.StatusBadRequest, ErrCodeDup, fmt.Sprintf("hosts %s already exists", name)))
		return
	}

	config.OnUpdate(func(c *config.Config) error {
		for i := range c.Hosts {
			if c.Hosts[i].Name == name {
				c.Hosts[i] = &req.Data
				break
			}
		}
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters deleteHostsRequest
type deleteHostsRequest struct {
	// in: path
	// required: true
	Hosts string `uri:"hosts" json:"hosts"`
}

// successful operation.
// swagger:response deleteHostsResponse
type deleteHostsResponse struct {
	Data Response
}

func deleteHosts(ctx *gin.Context) {
	// swagger:route DELETE /config/hosts/{hosts} Hosts deleteHostsRequest
	//
	// Delete hosts by name.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: deleteHostsResponse

	var req deleteHostsRequest
	ctx.ShouldBindUri(&req)

	name := strings.TrimSpace(req.Hosts)

	if !app.Runtime.HostsRegistry().IsRegistered(name) {
		writeError(ctx, NewError(http.StatusBadRequest, ErrCodeNotFound, fmt.Sprintf("hosts %s not found", name)))
		return
	}
	app.Runtime.HostsRegistry().Unregister(name)

	config.OnUpdate(func(c *config.Config) error {
		hosts := c.Hosts
		c.Hosts = nil
		for _, s := range hosts {
			if s.Name == name {
				continue
			}
			c.Hosts = append(c.Hosts, s)
		}
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}
