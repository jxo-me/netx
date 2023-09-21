package api

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/netx/api/handler"
	"github.com/jxo-me/netx/core/logger"
	"net/http"
	"time"

	"github.com/jxo-me/netx/core/auth"
)

func mwLogger() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		// start time
		startTime := time.Now()
		// Processing request
		r.Middleware.Next()
		duration := time.Since(startTime)
		logger.Default().WithFields(map[string]any{
			"kind":     "api",
			"method":   r.Request.Method,
			"uri":      r.Request.RequestURI,
			"code":     r.Response.Writer.Status,
			"client":   r.GetClientIp(),
			"duration": duration,
		}).Infof("| %3d | %13v | %15s | %-7s %s",
			r.Response.Writer.Status, duration, r.GetClientIp(), r.Request.Method, r.Request.RequestURI)

	}
}

func mwBasicAuth(auther auth.IAuthenticator) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		if auther == nil {
			r.Middleware.Next()
			return
		}
		u, p, _ := r.Request.BasicAuth()
		if _, ok := auther.Authenticate(r.GetCtx(), u, p); !ok {
			JsonExit(r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}
		r.Middleware.Next()
	}
}

// CORSMiddleware 允许跨域请求中间件
func CORSMiddleware(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowMethods = "GET,PUT,POST,DELETE,OPTIONS"
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

// Response is the default middleware handling handler response object and its error.
func Response(r *ghttp.Request) {
	r.Middleware.Next()
	//glog.Warning(r.GetCtx(), "Response start ......")
	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}
	var (
		msg  string
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = gerror.Code(err)
	)
	contentType := r.Get("format")
	//glog.Debug(r.GetCtx(), "Content-Type:", contentType)
	if contentType.String() == "yaml" {
		r.Response.Header().Set("Content-Type", "text/x-yaml")
		err = handler.Write(r.Response.ResponseWriter, res, "yaml")
		r.Exit()
	}

	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		JsonExit(r, code.Code(), err.Error())
	} else if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		msg = http.StatusText(r.Response.Status)
		switch r.Response.Status {
		case http.StatusNotFound:
			code = gcode.CodeNotFound
		case http.StatusForbidden:
			code = gcode.CodeNotAuthorized
		default:
			code = gcode.CodeUnknown
		}
	} else {
		code = gcode.CodeOK
		msg = "ok"
	}
	//glog.Warning(r.GetCtx(), "Response end ......", res)
	JsonExit(r, code.Code(), msg, res)
}
