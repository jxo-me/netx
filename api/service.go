package api

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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
	r := g.Server()
	r.SetOpenApiPath("/openapi")
	r.SetDumpRouterMap(false)
	r.SetLogStdout(false)
	err = r.SetListener(ln)
	if err != nil {
		return nil, err
	}
	r.BindMiddlewareDefault(CORSMiddleware)

	router := r.Group("")
	if options.pathPrefix != "" {
		router = router.Group(options.pathPrefix)
	}

	if options.accessLog {
		router.Middleware(mwLogger())
	}
	_ = InitDoc(r)
	config := router.Group("/config")
	config.Middleware(Response)
	config.Middleware(mwBasicAuth(options.auther))
	registerConfig(config)

	return &server{
		s:  r,
		ln: ln,
	}, nil
}

func (s *server) Serve() error {
	return s.s.Start()
}

func (s *server) Addr() net.Addr {
	return s.ln.Addr()
}

func (s *server) Close() error {
	return s.s.Shutdown()
}

func registerConfig(config *ghttp.RouterGroup) {
	config.Bind(
		Config,
		Service,
		Chain,
		Hop,
		Auther,
		Admission,
		Bypass,
		Resolver,
		Hosts,
		Ingress,
		Limiter,
		ConnLimiter,
		RateLimiter,
	)
}
