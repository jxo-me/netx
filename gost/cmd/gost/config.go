package main

import (
	"github.com/jxo-me/netx/api"
	iApi "github.com/jxo-me/netx/core/api"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/service"
	"github.com/jxo-me/netx/x/app"
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
)

func buildService(cfg *config.Config) (services []service.IService) {
	if cfg == nil {
		return
	}

	log := logger.Default()

	for _, loggerCfg := range cfg.Loggers {
		if lg := logger_parser.ParseLogger(loggerCfg); lg != nil {
			if err := app.Runtime.LoggerRegistry().Register(loggerCfg.Name, lg); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, autherCfg := range cfg.Authers {
		if auther := auth_parser.ParseAuther(autherCfg); auther != nil {
			if err := app.Runtime.AutherRegistry().Register(autherCfg.Name, auther); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, admissionCfg := range cfg.Admissions {
		if adm := admission_parser.ParseAdmission(admissionCfg); adm != nil {
			if err := app.Runtime.AdmissionRegistry().Register(admissionCfg.Name, adm); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, bypassCfg := range cfg.Bypasses {
		if bp := bypass_parser.ParseBypass(bypassCfg); bp != nil {
			if err := app.Runtime.BypassRegistry().Register(bypassCfg.Name, bp); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, resolverCfg := range cfg.Resolvers {
		r, err := resolver_parser.ParseResolver(resolverCfg)
		if err != nil {
			log.Fatal(err)
		}
		if r != nil {
			if err := app.Runtime.ResolverRegistry().Register(resolverCfg.Name, r); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, hostsCfg := range cfg.Hosts {
		if h := hosts_parser.ParseHostMapper(hostsCfg); h != nil {
			if err := app.Runtime.HostsRegistry().Register(hostsCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, ingressCfg := range cfg.Ingresses {
		if h := ingress_parser.ParseIngress(ingressCfg); h != nil {
			if err := app.Runtime.IngressRegistry().Register(ingressCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, routerCfg := range cfg.Routers {
		if h := router_parser.ParseRouter(routerCfg); h != nil {
			if err := app.Runtime.RouterRegistry().Register(routerCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, sdCfg := range cfg.SDs {
		if h := sd_parser.ParseSD(sdCfg); h != nil {
			if err := app.Runtime.SDRegistry().Register(sdCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, observerCfg := range cfg.Observers {
		if h := observer_parser.ParseObserver(observerCfg); h != nil {
			if err := app.Runtime.ObserverRegistry().Register(observerCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, recorderCfg := range cfg.Recorders {
		if h := recorder_parser.ParseRecorder(recorderCfg); h != nil {
			if err := app.Runtime.RecorderRegistry().Register(recorderCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, limiterCfg := range cfg.Limiters {
		if h := limiter_parser.ParseTrafficLimiter(limiterCfg); h != nil {
			if err := app.Runtime.TrafficLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, limiterCfg := range cfg.CLimiters {
		if h := limiter_parser.ParseConnLimiter(limiterCfg); h != nil {
			if err := app.Runtime.ConnLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, limiterCfg := range cfg.RLimiters {
		if h := limiter_parser.ParseRateLimiter(limiterCfg); h != nil {
			if err := app.Runtime.RateLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, hopCfg := range cfg.Hops {
		hop, err := hop_parser.ParseHop(hopCfg, log)
		if err != nil {
			log.Fatal(err)
		}
		if hop != nil {
			if err := app.Runtime.HopRegistry().Register(hopCfg.Name, hop); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, chainCfg := range cfg.Chains {
		c, err := chain_parser.ParseChain(chainCfg, log)
		if err != nil {
			log.Fatal(err)
		}
		if c != nil {
			if err := app.Runtime.ChainRegistry().Register(chainCfg.Name, c); err != nil {
				log.Fatal(err)
			}
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
	auther := auth_parser.ParseAutherFromAuth(cfg.Auth)
	if cfg.Auther != "" {
		auther = app.Runtime.AutherRegistry().Get(cfg.Auther)
	}
	return api.NewService(
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
	return metrics.NewService(
		cfg.Addr,
		metrics.PathOption(cfg.Path),
		metrics.AutherOption(auther),
	)
}
