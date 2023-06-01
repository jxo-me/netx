package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/netx/api/handler"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/service"
	"net"
)

type options struct {
	accessLog  bool
	pathPrefix string
	auther     auth.IAuthenticator
}

type Option func(*options)

func PathPrefixOption(pathPrefix string) Option {
	return func(o *options) {
		o.pathPrefix = pathPrefix
	}
}

func AccessLogOption(enable bool) Option {
	return func(o *options) {
		o.accessLog = enable
	}
}

func AutherOption(auther auth.IAuthenticator) Option {
	return func(o *options) {
		o.auther = auther
	}
}

type server struct {
	s  *ghttp.Server
	ln net.Listener
}

func NewService(addr string, opts ...Option) (service.IService, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	var options options
	for _, opt := range opts {
		opt(&options)
	}
	s := g.Server()
	s.SetOpenApiPath("/api.json")
	s.SetSwaggerPath("/swagger")
	s.SetDumpRouterMap(false)
	s.SetLogStdout(false)
	err = s.SetListener(ln)
	if err != nil {
		return nil, err
	}
	s.BindStatusHandlerByMap(map[int]ghttp.HandlerFunc{
		403: func(r *ghttp.Request) { r.Response.ClearBuffer(); r.Response.Writeln("403") },
		404: func(r *ghttp.Request) { r.Response.ClearBuffer(); r.Response.Writeln("404") },
		500: func(r *ghttp.Request) { r.Response.ClearBuffer(); r.Response.Writeln("500") },
	})
	s.BindMiddlewareDefault(
		CORSMiddleware,
		Response,
	)
	err = registerApiDocument(s)
	if err != nil {
		return nil, err
	}
	s.Group("", func(root *ghttp.RouterGroup) {
		if options.pathPrefix != "" {
			root = root.Group(options.pathPrefix)
		}

		if options.accessLog {
			root.Middleware(mwLogger())
		}
		cfg := root.Group("/config").Middleware(
			mwBasicAuth(options.auther),
		)
		registerRouters(cfg)
	})

	return &server{
		s:  s,
		ln: ln,
	}, nil
}

func (s *server) Serve() error {
	s.s.Run()
	return nil
}

func (s *server) Addr() net.Addr {
	return s.ln.Addr()
}

func (s *server) Close() error {
	return s.s.Shutdown()
}

func registerRouters(r *ghttp.RouterGroup) {
	r.Bind(
		handler.Config,
		handler.Service,
		handler.Chain,
		handler.Hop,
		handler.Auther,
		handler.Admission,
		handler.Bypass,
		handler.Resolver,
		handler.Hosts,
		handler.Ingress,
		handler.Limiter,
		handler.ConnLimiter,
		handler.RateLimiter,
	)
}
