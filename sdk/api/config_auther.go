package api

import (
	"github.com/jxo-me/netx/sdk"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jxo-me/netx/sdk/config"
)

// swagger:parameters createAutherRequest
type createAutherRequest struct {
	// in: body
	Data config.AutherConfig `json:"data"`
}

// successful operation.
// swagger:response createAutherResponse
type createAutherResponse struct {
	Data Response
}

func createAuther(ctx *gin.Context) {
	// swagger:route POST /config/authers Auther createAutherRequest
	//
	// Create a new auther, the name of the auther must be unique in auther list.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: createAutherResponse

	var req createAutherRequest
	ctx.ShouldBindJSON(&req.Data)

	if req.Data.Name == "" {
		writeError(ctx, ErrInvalid)
		return
	}

	v := sdk.Runtime.ParseAuther(&req.Data)
	if err := sdk.Runtime.AutherRegistry().Register(req.Data.Name, v); err != nil {
		writeError(ctx, ErrDup)
		return
	}

	config.OnUpdate(func(c *config.Config) error {
		c.Authers = append(c.Authers, &req.Data)
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters updateAutherRequest
type updateAutherRequest struct {
	// in: path
	// required: true
	Auther string `uri:"auther" json:"auther"`
	// in: body
	Data config.AutherConfig `json:"data"`
}

// successful operation.
// swagger:response updateAutherResponse
type updateAutherResponse struct {
	Data Response
}

func updateAuther(ctx *gin.Context) {
	// swagger:route PUT /config/authers/{auther} Auther updateAutherRequest
	//
	// Update auther by name, the auther must already exist.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: updateAutherResponse

	var req updateAutherRequest
	ctx.ShouldBindUri(&req)
	ctx.ShouldBindJSON(&req.Data)

	if !sdk.Runtime.AutherRegistry().IsRegistered(req.Auther) {
		writeError(ctx, ErrNotFound)
		return
	}

	req.Data.Name = req.Auther

	v := sdk.Runtime.ParseAuther(&req.Data)
	sdk.Runtime.AutherRegistry().Unregister(req.Auther)

	if err := sdk.Runtime.AutherRegistry().Register(req.Auther, v); err != nil {
		writeError(ctx, ErrDup)
		return
	}

	config.OnUpdate(func(c *config.Config) error {
		for i := range c.Authers {
			if c.Authers[i].Name == req.Auther {
				c.Authers[i] = &req.Data
				break
			}
		}
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters deleteAutherRequest
type deleteAutherRequest struct {
	// in: path
	// required: true
	Auther string `uri:"auther" json:"auther"`
}

// successful operation.
// swagger:response deleteAutherResponse
type deleteAutherResponse struct {
	Data Response
}

func deleteAuther(ctx *gin.Context) {
	// swagger:route DELETE /config/authers/{auther} Auther deleteAutherRequest
	//
	// Delete auther by name.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: deleteAutherResponse

	var req deleteAutherRequest
	ctx.ShouldBindUri(&req)

	if !sdk.Runtime.AutherRegistry().IsRegistered(req.Auther) {
		writeError(ctx, ErrNotFound)
		return
	}
	sdk.Runtime.AutherRegistry().Unregister(req.Auther)

	config.OnUpdate(func(c *config.Config) error {
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

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}
