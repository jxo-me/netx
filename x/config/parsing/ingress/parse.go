package ingress

import (
	"crypto/tls"
	"strings"

	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/x/config"
	xingress "github.com/jxo-me/netx/x/ingress"
	ingressplugin "github.com/jxo-me/netx/x/ingress/plugin"
	"github.com/jxo-me/netx/x/internal/loader"
	"github.com/jxo-me/netx/x/internal/plugin"
)

func ParseIngress(cfg *config.IngressConfig) ingress.IIngress {
	if cfg == nil {
		return nil
	}

	if cfg.Plugin != nil {
		var tlsCfg *tls.Config
		if cfg.Plugin.TLS != nil {
			tlsCfg = &tls.Config{
				ServerName:         cfg.Plugin.TLS.ServerName,
				InsecureSkipVerify: !cfg.Plugin.TLS.Secure,
			}
		}
		switch strings.ToLower(cfg.Plugin.Type) {
		case "http":
			return ingressplugin.NewHTTPPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TLSConfigOption(tlsCfg),
				plugin.TimeoutOption(cfg.Plugin.Timeout),
			)
		default:
			return ingressplugin.NewGRPCPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TokenOption(cfg.Plugin.Token),
				plugin.TLSConfigOption(tlsCfg),
			)
		}
	}

	var rules []*ingress.Rule
	for _, rule := range cfg.Rules {
		if rule.Hostname == "" || rule.Endpoint == "" {
			continue
		}

		rules = append(rules, &ingress.Rule{
			Hostname: rule.Hostname,
			Endpoint: rule.Endpoint,
		})
	}
	opts := []xingress.Option{
		xingress.RulesOption(rules),
		xingress.ReloadPeriodOption(cfg.Reload),
		xingress.LoggerOption(logger.Default().WithFields(map[string]any{
			"kind":    "ingress",
			"ingress": cfg.Name,
		})),
	}
	if cfg.File != nil && cfg.File.Path != "" {
		opts = append(opts, xingress.FileLoaderOption(loader.FileLoader(cfg.File.Path)))
	}
	if cfg.Redis != nil && cfg.Redis.Addr != "" {
		switch cfg.Redis.Type {
		case "set": // redis set
			opts = append(opts, xingress.RedisLoaderOption(loader.RedisSetLoader(
				cfg.Redis.Addr,
				loader.DBRedisLoaderOption(cfg.Redis.DB),
				loader.UsernameRedisLoaderOption(cfg.Redis.Username),
				loader.PasswordRedisLoaderOption(cfg.Redis.Password),
				loader.KeyRedisLoaderOption(cfg.Redis.Key),
			)))
		default: // redis hash
			opts = append(opts, xingress.RedisLoaderOption(loader.RedisHashLoader(
				cfg.Redis.Addr,
				loader.DBRedisLoaderOption(cfg.Redis.DB),
				loader.UsernameRedisLoaderOption(cfg.Redis.Username),
				loader.PasswordRedisLoaderOption(cfg.Redis.Password),
				loader.KeyRedisLoaderOption(cfg.Redis.Key),
			)))
		}
	}
	if cfg.HTTP != nil && cfg.HTTP.URL != "" {
		opts = append(opts, xingress.HTTPLoaderOption(loader.HTTPLoader(
			cfg.HTTP.URL,
			loader.TimeoutHTTPLoaderOption(cfg.HTTP.Timeout),
		)))
	}
	return xingress.NewIngress(opts...)
}
