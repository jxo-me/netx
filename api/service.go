package api

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/jxo-me/netx/api/bot"
	"github.com/jxo-me/netx/api/handler"
	"github.com/jxo-me/netx/core/api"
	"net"
)

type Server api.Server

func (s *Server) HttpServer() *ghttp.Server {
	return s.Srv
}

func (s *Server) TGBot() *api.TGBot {
	return s.Bot
}

func (s *Server) Serve() error {
	if s.Bot != nil {
		go func() {
			s.Bot.Bot.Start()
		}()
	}
	s.Srv.Run()
	return nil
}

func (s *Server) Addr() net.Addr {
	return s.Listener.Addr()
}

func (s *Server) Close() error {
	if s.Bot != nil {
		s.Bot.Bot.Stop()
	}
	return s.Srv.Shutdown()
}

func NewService(network, addr string, opts ...Option) (api.IApi, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	var option options
	for _, opt := range opts {
		opt(&option)
	}
	ctx := context.Background()
	var b *api.TGBot
	// bot
	if option.botEnable {
		b, err = botService(ctx, option)
	}
	// api
	s, err := apiService(ln, option, b)
	if err != nil {
		return nil, err
	}

	return &Server{
		Srv:      s,
		Listener: ln,
		Bot:      b,
	}, nil
}

func apiService(ln net.Listener, options options, b *api.TGBot) (s *ghttp.Server, err error) {
	s = g.Server()
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
	//s.BindMiddlewareDefault(
	//	CORSMiddleware,
	//	Response,
	//)
	err = registerApiDocument(s)
	if err != nil {
		return nil, err
	}
	s.Group("", func(root *ghttp.RouterGroup) {
		if options.pathPrefix != "" {
			root = root.Group(options.pathPrefix)
		}
		// register TG bot
		if options.botEnable && b != nil {
			// TG bot Hook api
			root.POST("", bot.Bot.Hook)
			// Web apps api
			root.GET("/", bot.WebApp.Index)
			root.GET("/validate", bot.WebApp.Validate)
			// Login API
			root.GET("/check_authorization", bot.WebApp.CheckAuthorization)
			root.GET("/login", bot.WebApp.Login)
			if b.Bot != nil {
				// bot routers
				for key, handlerFunc := range bot.Router().List {
					b.Bot.Handle(key, handlerFunc)
				}
				for key, handlerFunc := range bot.Router().Btns {
					b.Bot.Handle(key, handlerFunc)
				}
			}
		}
		// Middlewares
		root.Middleware(
			CORSMiddleware,
			Response,
		)
		if options.accessLog {
			root.Middleware(mwLogger())
		}
		cfg := root.Group("/config").Middleware(
			mwBasicAuth(options.auther),
		)
		registerRouters(cfg)
	})

	return s, nil
}

func botService(ctx context.Context, op options) (s *api.TGBot, err error) {
	s, err = NewBot(ctx, op.domain, op.botToken, op.pathPrefix)
	return
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
		handler.Router,
		handler.Hosts,
		handler.Ingress,
		handler.Limiter,
		handler.ConnLimiter,
		handler.RateLimiter,
	)
}
