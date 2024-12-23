package main

import (
	"github.com/jxo-me/netx/api"
	iApi "github.com/jxo-me/netx/core/api"
	"github.com/jxo-me/netx/core/auth"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/service"
	"github.com/jxo-me/netx/x/app"
	xauth "github.com/jxo-me/netx/x/auth"
	"github.com/jxo-me/netx/x/config"
	admission_parser "github.com/jxo-me/netx/x/config/parsing/admission"
	auth_parser "github.com/jxo-me/netx/x/config/parsing/auth"
	bypass_parser "github.com/jxo-me/netx/x/config/parsing/bypass"
	chain_parser "github.com/jxo-me/netx/x/config/parsing/chain"
	hop_parser "github.com/jxo-me/netx/x/config/parsing/hop"
	hosts_parser "github.com/jxo-me/netx/x/config/parsing/hosts"
	ingress_parser "github.com/jxo-me/netx/x/config/parsing/ingress"
	limiter_parser "github.com/jxo-me/netx/x/config/parsing/limiter"
	logger_parser "github.com/jxo-me/netx/x/config/parsing/logger"
	observer_parser "github.com/jxo-me/netx/x/config/parsing/observer"
	recorder_parser "github.com/jxo-me/netx/x/config/parsing/recorder"
	resolver_parser "github.com/jxo-me/netx/x/config/parsing/resolver"
	router_parser "github.com/jxo-me/netx/x/config/parsing/router"
	sd_parser "github.com/jxo-me/netx/x/config/parsing/sd"
	service_parser "github.com/jxo-me/netx/x/config/parsing/service"
	metrics "github.com/jxo-me/netx/x/metrics/service"
	"strings"
)

func buildService(cfg *config.Config) (services []service.IService) {
	if cfg == nil {
		return
	}

	log := logger.Default()

	for _, loggerCfg := range cfg.Loggers {
		if err := app.Runtime.LoggerRegistry().Register(loggerCfg.Name, logger_parser.ParseLogger(loggerCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, autherCfg := range cfg.Authers {
		if err := app.Runtime.AutherRegistry().Register(autherCfg.Name, auth_parser.ParseAuther(autherCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, admissionCfg := range cfg.Admissions {
		if err := app.Runtime.AdmissionRegistry().Register(admissionCfg.Name, admission_parser.ParseAdmission(admissionCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, bypassCfg := range cfg.Bypasses {
		if err := app.Runtime.BypassRegistry().Register(bypassCfg.Name, bypass_parser.ParseBypass(bypassCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, resolverCfg := range cfg.Resolvers {
		r, err := resolver_parser.ParseResolver(resolverCfg)
		if err != nil {
			log.Fatal(err)
		}
		if err := app.Runtime.ResolverRegistry().Register(resolverCfg.Name, r); err != nil {
			log.Fatal(err)
		}
	}

	for _, hostsCfg := range cfg.Hosts {
		if err := app.Runtime.HostsRegistry().Register(hostsCfg.Name, hosts_parser.ParseHostMapper(hostsCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, ingressCfg := range cfg.Ingresses {
		if err := app.Runtime.IngressRegistry().Register(ingressCfg.Name, ingress_parser.ParseIngress(ingressCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, routerCfg := range cfg.Routers {
		if err := app.Runtime.RouterRegistry().Register(routerCfg.Name, router_parser.ParseRouter(routerCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, sdCfg := range cfg.SDs {
		if err := app.Runtime.SDRegistry().Register(sdCfg.Name, sd_parser.ParseSD(sdCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, observerCfg := range cfg.Observers {
		if err := app.Runtime.ObserverRegistry().Register(observerCfg.Name, observer_parser.ParseObserver(observerCfg)); err != nil {
			log.Fatal(err)
		}
	}
	for _, recorderCfg := range cfg.Recorders {
		if err := app.Runtime.RecorderRegistry().Register(recorderCfg.Name, recorder_parser.ParseRecorder(recorderCfg)); err != nil {
			log.Fatal(err)
		}
	}

	for _, limiterCfg := range cfg.Limiters {
		if err := app.Runtime.TrafficLimiterRegistry().Register(limiterCfg.Name, limiter_parser.ParseTrafficLimiter(limiterCfg)); err != nil {
			log.Fatal(err)
		}
	}
	for _, limiterCfg := range cfg.CLimiters {
		if err := app.Runtime.ConnLimiterRegistry().Register(limiterCfg.Name, limiter_parser.ParseConnLimiter(limiterCfg)); err != nil {
			log.Fatal(err)
		}
	}
	for _, limiterCfg := range cfg.RLimiters {
		if err := app.Runtime.RateLimiterRegistry().Register(limiterCfg.Name, limiter_parser.ParseRateLimiter(limiterCfg)); err != nil {
			log.Fatal(err)
		}
	}
	for _, hopCfg := range cfg.Hops {
		hop, err := hop_parser.ParseHop(hopCfg, log)
		if err != nil {
			log.Fatal(err)
		}
		if err := app.Runtime.HopRegistry().Register(hopCfg.Name, hop); err != nil {
			log.Fatal(err)
		}
	}
	for _, chainCfg := range cfg.Chains {
		c, err := chain_parser.ParseChain(chainCfg, log)
		if err != nil {
			log.Fatal(err)
		}
		if err := app.Runtime.ChainRegistry().Register(chainCfg.Name, c); err != nil {
			log.Fatal(err)
		}
	}

	for _, svcCfg := range cfg.Services {
		svc, err := service_parser.ParseService(svcCfg)
		if err != nil {
			log.Fatal(err)
		}
		if svc != nil {
			if err := app.Runtime.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
				log.Fatal(err)
			}
		}
		services = append(services, svc)
	}

	return
}

func buildAPIService(cfg *config.APIConfig) (iApi.IApi, error) {
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
	return api.NewService(
		network,
		cfg.Addr,
		api.PathPrefixOption(cfg.PathPrefix),
		api.AccessLogOption(cfg.AccessLog),
		api.AutherOption(auther),
		api.BotEnableOption(cfg.BotEnable),
		api.DomainOption(cfg.Domain),
		api.TokenOption(cfg.BotToken),
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
		network,
		cfg.Addr,
		metrics.PathOption(cfg.Path),
		metrics.AutherOption(auther),
	)
}
