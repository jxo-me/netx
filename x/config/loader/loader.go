package loader

import (
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
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
)

var (
	defaultLoader *loader = &loader{}
)

func Load(cfg *config.Config) error {
	return defaultLoader.Load(cfg)
}

type loader struct{}

func (l *loader) Load(cfg *config.Config) error {
	logCfg := cfg.Log
	if logCfg == nil {
		logCfg = &config.LogConfig{}
	}
	logger.SetDefault(logger_parser.ParseLogger(&config.LoggerConfig{Log: logCfg}))

	tlsCfg, err := parsing.BuildDefaultTLSConfig(cfg.TLS)
	if err != nil {
		return err
	}
	parsing.SetDefaultTLSConfig(tlsCfg)

	if err := register(cfg); err != nil {
		return err
	}

	return nil
}

func register(cfg *config.Config) error {
	if cfg == nil {
		return nil
	}

	for name := range app.Runtime.LoggerRegistry().GetAll() {
		app.Runtime.LoggerRegistry().Unregister(name)
	}
	for _, loggerCfg := range cfg.Loggers {
		if err := app.Runtime.LoggerRegistry().Register(loggerCfg.Name, logger_parser.ParseLogger(loggerCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.AutherRegistry().GetAll() {
		app.Runtime.AutherRegistry().Unregister(name)
	}
	for _, autherCfg := range cfg.Authers {
		if err := app.Runtime.AutherRegistry().Register(autherCfg.Name, auth_parser.ParseAuther(autherCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.AdmissionRegistry().GetAll() {
		app.Runtime.AdmissionRegistry().Unregister(name)
	}
	for _, admissionCfg := range cfg.Admissions {
		if err := app.Runtime.AdmissionRegistry().Register(admissionCfg.Name, admission_parser.ParseAdmission(admissionCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.BypassRegistry().GetAll() {
		app.Runtime.BypassRegistry().Unregister(name)
	}
	for _, bypassCfg := range cfg.Bypasses {
		if err := app.Runtime.BypassRegistry().Register(bypassCfg.Name, bypass_parser.ParseBypass(bypassCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.ResolverRegistry().GetAll() {
		app.Runtime.ResolverRegistry().Unregister(name)
	}
	for _, resolverCfg := range cfg.Resolvers {
		r, err := resolver_parser.ParseResolver(resolverCfg)
		if err != nil {
			return err
		}
		if err := app.Runtime.ResolverRegistry().Register(resolverCfg.Name, r); err != nil {
			return err
		}
	}

	for name := range app.Runtime.HostsRegistry().GetAll() {
		app.Runtime.HostsRegistry().Unregister(name)
	}
	for _, hostsCfg := range cfg.Hosts {
		if err := app.Runtime.HostsRegistry().Register(hostsCfg.Name, hosts_parser.ParseHostMapper(hostsCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.IngressRegistry().GetAll() {
		app.Runtime.IngressRegistry().Unregister(name)
	}
	for _, ingressCfg := range cfg.Ingresses {
		if err := app.Runtime.IngressRegistry().Register(ingressCfg.Name, ingress_parser.ParseIngress(ingressCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.RouterRegistry().GetAll() {
		app.Runtime.RouterRegistry().Unregister(name)
	}
	for _, routerCfg := range cfg.Routers {
		if err := app.Runtime.RouterRegistry().Register(routerCfg.Name, router_parser.ParseRouter(routerCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.SDRegistry().GetAll() {
		app.Runtime.SDRegistry().Unregister(name)
	}
	for _, sdCfg := range cfg.SDs {
		if err := app.Runtime.SDRegistry().Register(sdCfg.Name, sd_parser.ParseSD(sdCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.ObserverRegistry().GetAll() {
		app.Runtime.ObserverRegistry().Unregister(name)
	}
	for _, observerCfg := range cfg.Observers {
		if err := app.Runtime.ObserverRegistry().Register(observerCfg.Name, observer_parser.ParseObserver(observerCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.RecorderRegistry().GetAll() {
		app.Runtime.RecorderRegistry().Unregister(name)
	}
	for _, recorderCfg := range cfg.Recorders {
		if err := app.Runtime.RecorderRegistry().Register(recorderCfg.Name, recorder_parser.ParseRecorder(recorderCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.TrafficLimiterRegistry().GetAll() {
		app.Runtime.TrafficLimiterRegistry().Unregister(name)
	}
	for _, limiterCfg := range cfg.Limiters {
		if err := app.Runtime.TrafficLimiterRegistry().Register(limiterCfg.Name, limiter_parser.ParseTrafficLimiter(limiterCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.ConnLimiterRegistry().GetAll() {
		app.Runtime.ConnLimiterRegistry().Unregister(name)
	}
	for _, limiterCfg := range cfg.CLimiters {
		if err := app.Runtime.ConnLimiterRegistry().Register(limiterCfg.Name, limiter_parser.ParseConnLimiter(limiterCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.RateLimiterRegistry().GetAll() {
		app.Runtime.RateLimiterRegistry().Unregister(name)
	}
	for _, limiterCfg := range cfg.RLimiters {
		if err := app.Runtime.RateLimiterRegistry().Register(limiterCfg.Name, limiter_parser.ParseRateLimiter(limiterCfg)); err != nil {
			return err
		}
	}

	for name := range app.Runtime.HopRegistry().GetAll() {
		app.Runtime.HopRegistry().Unregister(name)
	}
	for _, hopCfg := range cfg.Hops {
		hop, err := hop_parser.ParseHop(hopCfg, logger.Default())
		if err != nil {
			return err
		}
		if err := app.Runtime.HopRegistry().Register(hopCfg.Name, hop); err != nil {
			return err
		}
	}

	for name := range app.Runtime.ChainRegistry().GetAll() {
		app.Runtime.ChainRegistry().Unregister(name)
	}
	for _, chainCfg := range cfg.Chains {
		c, err := chain_parser.ParseChain(chainCfg, logger.Default())
		if err != nil {
			return err
		}
		if err := app.Runtime.ChainRegistry().Register(chainCfg.Name, c); err != nil {
			return err
		}
	}

	for name := range app.Runtime.ServiceRegistry().GetAll() {
		app.Runtime.ServiceRegistry().Unregister(name)
	}
	for _, svcCfg := range cfg.Services {
		svc, err := service_parser.ParseService(svcCfg)
		if err != nil {
			return err
		}
		if svc != nil {
			if err := app.Runtime.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
				return err
			}
		}
	}

	return nil
}
