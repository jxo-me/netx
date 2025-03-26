package main

import (
	"context"
	"errors"
	"github.com/jxo-me/netx/x/app"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/judwhite/go-svc"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/service"
	api_service "github.com/jxo-me/netx/x/api/service"
	xauth "github.com/jxo-me/netx/x/auth"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/loader"
	auth_parser "github.com/jxo-me/netx/x/config/parsing/auth"
	"github.com/jxo-me/netx/x/config/parsing/parser"
	xmetrics "github.com/jxo-me/netx/x/metrics"
	metrics "github.com/jxo-me/netx/x/metrics/service"
)

type program struct {
	srvApi       service.IService
	srvMetrics   service.IService
	srvProfiling *http.Server

	cancel context.CancelFunc
}

func (p *program) Init(env svc.Environment) error {
	parser.Init(parser.Args{
		CfgFile:     cfgFile,
		Services:    services,
		Nodes:       nodes,
		Debug:       debug,
		Trace:       trace,
		ApiAddr:     apiAddr,
		MetricsAddr: metricsAddr,
	})

	return nil
}

func (p *program) Start() error {
	cfg, err := parser.Parse()
	if err != nil {
		return err
	}

	if outputFormat != "" {
		if err := cfg.Write(os.Stdout, outputFormat); err != nil {
			return err
		}
		os.Exit(0)
	}

	config.Set(cfg)

	if err := loader.Load(cfg); err != nil {
		return err
	}

	if err := p.run(cfg); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go p.reload(ctx)

	return nil
}

func (p *program) run(cfg *config.Config) error {
	for _, svc := range app.Runtime.ServiceRegistry().GetAll() {
		svc := svc
		go func() {
			svc.Serve()
		}()
	}

	if p.srvApi != nil {
		p.srvApi.Close()
		p.srvApi = nil
	}
	if cfg.API != nil {
		s, err := buildApiService(cfg.API)
		if err != nil {
			return err
		}

		p.srvApi = s

		go func() {
			defer s.Close()

			log := logger.Default().WithFields(map[string]any{"kind": "service", "service": "@api"})

			log.Info("listening on ", s.Addr())
			if err := s.Serve(); !errors.Is(err, http.ErrServerClosed) {
				log.Error(err)
			}
		}()
	}

	xmetrics.Enable(false)
	if p.srvMetrics != nil {
		p.srvMetrics.Close()
		p.srvMetrics = nil
	}
	if cfg.Metrics != nil && cfg.Metrics.Addr != "" {
		s, err := buildMetricsService(cfg.Metrics)
		if err != nil {
			return err
		}

		p.srvMetrics = s

		xmetrics.Enable(true)

		go func() {
			defer s.Close()

			log := logger.Default().WithFields(map[string]any{"kind": "service", "service": "@metrics"})

			log.Info("listening on ", s.Addr())
			if err := s.Serve(); !errors.Is(err, http.ErrServerClosed) {
				log.Error(err)
			}
		}()
	}

	if p.srvProfiling != nil {
		p.srvProfiling.Close()
		p.srvProfiling = nil
	}
	if cfg.Profiling != nil {
		addr := cfg.Profiling.Addr
		if addr == "" {
			addr = ":6060"
		}
		s := &http.Server{
			Addr: addr,
		}
		p.srvProfiling = s

		go func() {
			defer s.Close()

			log := logger.Default().WithFields(map[string]any{"kind": "service", "service": "@profiling"})

			log.Info("listening on ", addr)
			if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Error(err)
			}
		}()
	}

	return nil
}

func (p *program) Stop() error {
	if p.cancel != nil {
		p.cancel()
	}

	for name, srv := range app.Runtime.ServiceRegistry().GetAll() {
		srv.Close()
		logger.Default().Debugf("service %s shutdown", name)
	}

	if p.srvApi != nil {
		p.srvApi.Close()
		logger.Default().Debug("service @api shutdown")
	}
	if p.srvMetrics != nil {
		p.srvMetrics.Close()
		logger.Default().Debug("service @metrics shutdown")
	}
	if p.srvProfiling != nil {
		p.srvProfiling.Close()
		logger.Default().Debug("service @profiling shutdown")
	}

	return nil
}

func (p *program) reload(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)

	for {
		select {
		case <-c:
			if err := p.reloadConfig(); err != nil {
				logger.Default().Error(err)
			} else {
				logger.Default().Info("config reloaded")
			}

		case <-ctx.Done():
			return
		}
	}
}

func (p *program) reloadConfig() error {
	cfg, err := parser.Parse()
	if err != nil {
		return err
	}
	config.Set(cfg)

	if err := loader.Load(cfg); err != nil {
		return err
	}

	if err := p.run(cfg); err != nil {
		return err
	}

	return nil
}

func buildApiService(cfg *config.APIConfig) (service.IService, error) {
	var authers []auth.IAuthenticator
	if auther := auth_parser.ParseAutherFromAuth(cfg.Auth); auther != nil {
		authers = append(authers, auther)
	}
	if cfg.Auther != "" {
		authers = append(authers, app.Runtime.AutherRegistry().Get(cfg.Auther))
	}

	var auther auth.IAuthenticator
	if len(authers) > 0 {
		auther = xauth.AuthenticatorGroup(authers...)
	}

	network := "tcp"
	addr := cfg.Addr
	if strings.HasPrefix(addr, "unix://") {
		network = "unix"
		addr = strings.TrimPrefix(addr, "unix://")
	}
	return api_service.NewService(
		network, addr,
		api_service.PathPrefixOption(cfg.PathPrefix),
		api_service.AccessLogOption(cfg.AccessLog),
		api_service.AutherOption(auther),
	)
}

func buildMetricsService(cfg *config.MetricsConfig) (service.IService, error) {
	auther := auth_parser.ParseAutherFromAuth(cfg.Auth)
	if cfg.Auther != "" {
		auther = app.Runtime.AutherRegistry().Get(cfg.Auther)
	}

	network := "tcp"
	addr := cfg.Addr
	if strings.HasPrefix(addr, "unix://") {
		network = "unix"
		addr = strings.TrimPrefix(addr, "unix://")
	}
	return metrics.NewService(
		network, addr,
		metrics.PathOption(cfg.Path),
		metrics.AutherOption(auther),
	)
}
