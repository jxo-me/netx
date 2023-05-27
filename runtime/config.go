package runtime

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jxo-me/netx/config"
	"github.com/jxo-me/netx/core/logger"
	xlogger "github.com/jxo-me/netx/core/logger"
	metrics "github.com/jxo-me/netx/core/metrics/service"
	"github.com/jxo-me/netx/core/service"
	"gopkg.in/natefinch/lumberjack.v2"
)

func (a *Application) BuildService(cfg *config.Config) (services []service.IService) {
	if cfg == nil {
		return
	}

	log := logger.Default()

	for _, autherCfg := range cfg.Authers {
		if auther := a.ParseAuther(autherCfg); auther != nil {
			if err := a.AutherRegistry().Register(autherCfg.Name, auther); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, admissionCfg := range cfg.Admissions {
		if adm := a.ParseAdmission(admissionCfg); adm != nil {
			if err := a.AdmissionRegistry().Register(admissionCfg.Name, adm); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, bypassCfg := range cfg.Bypasses {
		if bp := a.ParseBypass(bypassCfg); bp != nil {
			if err := a.BypassRegistry().Register(bypassCfg.Name, bp); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, resolverCfg := range cfg.Resolvers {
		r, err := a.ParseResolver(resolverCfg)
		if err != nil {
			log.Fatal(err)
		}
		if r != nil {
			if err := a.ResolverRegistry().Register(resolverCfg.Name, r); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, hostsCfg := range cfg.Hosts {
		if h := a.ParseHosts(hostsCfg); h != nil {
			if err := a.HostsRegistry().Register(hostsCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, ingressCfg := range cfg.Ingresses {
		if h := a.ParseIngress(ingressCfg); h != nil {
			if err := a.IngressRegistry().Register(ingressCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, recorderCfg := range cfg.Recorders {
		if h := a.ParseRecorder(recorderCfg); h != nil {
			if err := a.RecorderRegistry().Register(recorderCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, limiterCfg := range cfg.Limiters {
		if h := a.ParseTrafficLimiter(limiterCfg); h != nil {
			if err := a.TrafficLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, limiterCfg := range cfg.CLimiters {
		if h := a.ParseConnLimiter(limiterCfg); h != nil {
			if err := a.ConnLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, limiterCfg := range cfg.RLimiters {
		if h := a.ParseRateLimiter(limiterCfg); h != nil {
			if err := a.RateLimiterRegistry().Register(limiterCfg.Name, h); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, hopCfg := range cfg.Hops {
		hop, err := a.ParseHop(hopCfg)
		if err != nil {
			log.Fatal(err)
		}
		if hop != nil {
			if err := a.HopRegistry().Register(hopCfg.Name, hop); err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, chainCfg := range cfg.Chains {
		c, err := a.ParseChain(chainCfg)
		if err != nil {
			log.Fatal(err)
		}
		if c != nil {
			if err := a.ChainRegistry().Register(chainCfg.Name, c); err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, svcCfg := range cfg.Services {
		svc, err := a.ParseService(svcCfg)
		if err != nil {
			log.Fatal(err)
		}
		if svc != nil {
			if err := a.ServiceRegistry().Register(svcCfg.Name, svc); err != nil {
				log.Fatal(err)
			}
		}
		services = append(services, svc)
	}

	return
}

func (a *Application) LogFromConfig(cfg *config.LogConfig) logger.ILogger {
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

func (a *Application) BuildAPIService(cfg *config.APIConfig) (service.IService, error) {
	//auther := a.ParseAutherFromAuth(cfg.Auth)
	//if cfg.Auther != "" {
	//	auther = a.AutherRegistry().Get(cfg.Auther)
	//}
	//return api.NewService(
	//	cfg.Addr,
	//	api.PathPrefixOption(cfg.PathPrefix),
	//	api.AccessLogOption(cfg.AccessLog),
	//	api.AutherOption(auther),
	//)
	return nil, nil
}

func (a *Application) BuildMetricsService(cfg *config.MetricsConfig) (service.IService, error) {
	return metrics.NewService(
		cfg.Addr,
		metrics.PathOption(cfg.Path),
	)
}
