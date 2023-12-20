package bypass

import (
	"crypto/tls"
	"strings"

	"github.com/jxo-me/netx/core/bypass"
	"github.com/jxo-me/netx/core/logger"
	xbypass "github.com/jxo-me/netx/x/bypass"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/internal/loader"
	"github.com/jxo-me/netx/x/internal/plugin"
	"github.com/jxo-me/netx/x/registry"
)

func ParseBypass(cfg *config.BypassConfig) bypass.IBypass {
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
			return xbypass.NewHTTPPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TLSConfigOption(tlsCfg),
				plugin.TimeoutOption(cfg.Plugin.Timeout),
			)
		default:
			return xbypass.NewGRPCPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TokenOption(cfg.Plugin.Token),
				plugin.TLSConfigOption(tlsCfg),
			)
		}
	}

	opts := []xbypass.Option{
		xbypass.MatchersOption(cfg.Matchers),
		xbypass.WhitelistOption(cfg.Reverse || cfg.Whitelist),
		xbypass.ReloadPeriodOption(cfg.Reload),
		xbypass.LoggerOption(logger.Default().WithFields(map[string]any{
			"kind":   "bypass",
			"bypass": cfg.Name,
		})),
	}
	if cfg.File != nil && cfg.File.Path != "" {
		opts = append(opts, xbypass.FileLoaderOption(loader.FileLoader(cfg.File.Path)))
	}
	if cfg.Redis != nil && cfg.Redis.Addr != "" {
		opts = append(opts, xbypass.RedisLoaderOption(loader.RedisSetLoader(
			cfg.Redis.Addr,
			loader.DBRedisLoaderOption(cfg.Redis.DB),
			loader.PasswordRedisLoaderOption(cfg.Redis.Password),
			loader.KeyRedisLoaderOption(cfg.Redis.Key),
		)))
	}
	if cfg.HTTP != nil && cfg.HTTP.URL != "" {
		opts = append(opts, xbypass.HTTPLoaderOption(loader.HTTPLoader(
			cfg.HTTP.URL,
			loader.TimeoutHTTPLoaderOption(cfg.HTTP.Timeout),
		)))
	}

	return xbypass.NewBypass(opts...)
}

func List(name string, names ...string) []bypass.IBypass {
	var bypasses []bypass.IBypass
	if bp := registry.BypassRegistry().Get(name); bp != nil {
		bypasses = append(bypasses, bp)
	}
	for _, s := range names {
		if bp := registry.BypassRegistry().Get(s); bp != nil {
			bypasses = append(bypasses, bp)
		}
	}
	return bypasses
}
