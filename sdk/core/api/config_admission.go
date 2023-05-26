package api

import (
	"net/http"
	"github.com/jxo-me/netx/sdk"

	"github.com/gin-gonic/gin"
	"github.com/jxo-me/netx/sdk/core/config"
	"github.com/jxo-me/netx/sdk/core/config/parsing"
)

// swagger:parameters createAdmissionRequest
type createAdmissionRequest struct {
	// in: body
	Data config.AdmissionConfig `json:"data"`
}

// successful operation.
// swagger:response createAdmissionResponse
type createAdmissionResponse struct {
	Data Response
}

func createAdmission(ctx *gin.Context) {
	// swagger:route POST /config/admissions Admission createAdmissionRequest
	//
	// Create a new admission, the name of admission must be unique in admission list.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: createAdmissionResponse

	var req createAdmissionRequest
	ctx.ShouldBindJSON(&req.Data)

	if req.Data.Name == "" {
		writeError(ctx, ErrInvalid)
		return
	}

	v := parsing.ParseAdmission(&req.Data)

	if err := sdk.Runtime.AdmissionRegistry().Register(req.Data.Name, v); err != nil {
		writeError(ctx, ErrDup)
		return
	}

	config.OnUpdate(func(c *config.Config) error {
		c.Admissions = append(c.Admissions, &req.Data)
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters updateAdmissionRequest
type updateAdmissionRequest struct {
	// in: path
	// required: true
	Admission string `uri:"admission" json:"admission"`
	// in: body
	Data config.AdmissionConfig `json:"data"`
}

// successful operation.
// swagger:response updateAdmissionResponse
type updateAdmissionResponse struct {
	Data Response
}

func updateAdmission(ctx *gin.Context) {
	// swagger:route PUT /config/admissions/{admission} Admission updateAdmissionRequest
	//
	// Update admission by name, the admission must already exist.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: updateAdmissionResponse

	var req updateAdmissionRequest
	ctx.ShouldBindUri(&req)
	ctx.ShouldBindJSON(&req.Data)

	if !sdk.Runtime.AdmissionRegistry().IsRegistered(req.Admission) {
		writeError(ctx, ErrNotFound)
		return
	}

	req.Data.Name = req.Admission

	v := parsing.ParseAdmission(&req.Data)

	sdk.Runtime.AdmissionRegistry().Unregister(req.Admission)

	if err := sdk.Runtime.AdmissionRegistry().Register(req.Admission, v); err != nil {
		writeError(ctx, ErrDup)
		return
	}

	config.OnUpdate(func(c *config.Config) error {
		for i := range c.Admissions {
			if c.Admissions[i].Name == req.Admission {
				c.Admissions[i] = &req.Data
				break
			}
		}
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}

// swagger:parameters deleteAdmissionRequest
type deleteAdmissionRequest struct {
	// in: path
	// required: true
	Admission string `uri:"admission" json:"admission"`
}

// successful operation.
// swagger:response deleteAdmissionResponse
type deleteAdmissionResponse struct {
	Data Response
}

func deleteAdmission(ctx *gin.Context) {
	// swagger:route DELETE /config/admissions/{admission} Admission deleteAdmissionRequest
	//
	// Delete admission by name.
	//
	//     Security:
	//       basicAuth: []
	//
	//     Responses:
	//       200: deleteAdmissionResponse

	var req deleteAdmissionRequest
	ctx.ShouldBindUri(&req)

	if !sdk.Runtime.AdmissionRegistry().IsRegistered(req.Admission) {
		writeError(ctx, ErrNotFound)
		return
	}
	sdk.Runtime.AdmissionRegistry().Unregister(req.Admission)

	config.OnUpdate(func(c *config.Config) error {
		admissiones := c.Admissions
		c.Admissions = nil
		for _, s := range admissiones {
			if s.Name == req.Admission {
				continue
			}
			c.Admissions = append(c.Admissions, s)
		}
		return nil
	})

	ctx.JSON(http.StatusOK, Response{
		Msg: "OK",
	})
}
