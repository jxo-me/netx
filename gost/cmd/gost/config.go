package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/service"
	"github.com/jxo-me/netx/x/api"
	"github.com/jxo-me/netx/x/app"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/config/parsing"
	xlogger "github.com/jxo-me/netx/x/logger"
	metrics "github.com/jxo-me/netx/x/metrics/service"
	"gopkg.in/natefinch/lumberjack.v2"
)

func buildService(cfg *config.Config) (services []service.IService) {
	if cfg == nil {
		return
	}

	log := logger.Default()

	for _, autherCfg := range cfg.Authers {
		if auther := parsing.ParseAuther(autherCfg); auther != nil {
			if err := app.Runtime.AutherRegistry().Register(autherCfg.Name, auther); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, admissionCfg := range cfg.Admissions {
		if adm := parsing.ParseAdmission(admissionCfg); adm != nil {
			if err := app.Runtime.AdmissionRegistry().Register(admissionCfg.Name, adm); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, bypassCfg := range cfg.Bypasses {
		if bp := parsing.ParseBypass(bypassCfg); bp != nil {
			if err := app.Runtime.BypassRegistry().Register(bypassCfg.Name, bp); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, resolverCfg := range cfg.Resolvers {
		r, err := parsing.ParseResolver(resolverCfg)
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
		if h := parsing.ParseHosts(hostsCfg); h != nil {
			if err := app.Runtime.HostsRegistry().Register(hostsCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, ingressCfg := range cfg.Ingresses {
		if h := parsing.ParseIngress(ingressCfg); h != nil {
			if err := app.Runtime.IngressRegistry().Register(ingressCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, recorderCfg := range cfg.Recorders {
		if h := parsing.ParseRecorder(recorderCfg); h != nil {
			if err := app.Runtime.RecorderRegistry().Register(recorderCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, limiterCfg := range cfg.Limiters {
		if h := parsing.ParseTrafficLimiter(limiterCfg); h != nil {
			if err := app.Runtime.TrafficLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, limiterCfg := range cfg.CLimiters {
		if h := parsing.ParseConnLimiter(limiterCfg); h != nil {
			if err := app.Runtime.ConnLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, limiterCfg := range cfg.RLimiters {
		if h := parsing.ParseRateLimiter(limiterCfg); h != nil {
			if err := app.Runtime.RateLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, hopCfg := range cfg.Hops {
		hop, err := parsing.ParseHop(hopCfg)
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
		c, err := parsing.ParseChain(chainCfg)
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
		svc, err := parsing.ParseService(svcCfg)
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

func logFromConfig(cfg *config.LogConfig) logger.ILogger {
	if cfg == nil {
		cfg = &config.LogConfig{}
	}
	opts := []xlogger.LoggerOption{
		xlogger.FormatLoggerOption(logger.LogFormat(cfg.Format)),
		xlogger.LevelLoggerOption(logger.LogLevel(cfg.Level)),
	}

	var out io.Writer = os.Stderr
	switch cfg.Output {
	case "none", "null":
		return xlogger.Nop()
	case "stdout":
		out = os.Stdout
	case "stderr", "":
		out = os.Stderr
	default:
		if cfg.Rotation != nil {
			out = &lumberjack.Logger{
				Filename:   cfg.Output,
				MaxSize:    cfg.Rotation.MaxSize,
				MaxAge:     cfg.Rotation.MaxAge,
				MaxBackups: cfg.Rotation.MaxBackups,
				LocalTime:  cfg.Rotation.LocalTime,
				Compress:   cfg.Rotation.Compress,
			}
		} else {
			os.MkdirAll(filepath.Dir(cfg.Output), 0755)
			f, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				logger.Default().Warn(err)
			} else {
				out = f
			}
		}
	}
	opts = append(opts, xlogger.OutputLoggerOption(out))

	return xlogger.NewLogger(opts...)
}

func buildAPIService(cfg *config.APIConfig) (service.IService, error) {
	auther := parsing.ParseAutherFromAuth(cfg.Auth)
	if cfg.Auther != "" {
		auther = app.Runtime.AutherRegistry().Get(cfg.Auther)
	}
	return api.NewService(
		cfg.Addr,
		api.PathPrefixOption(cfg.PathPrefix),
		api.AccessLogOption(cfg.AccessLog),
		api.AutherOption(auther),
	)
}

func buildMetricsService(cfg *config.MetricsConfig) (service.IService, error) {
	return metrics.NewService(
		cfg.Addr,
		metrics.PathOption(cfg.Path),
	)
}
