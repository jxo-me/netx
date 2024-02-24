package router

import (
	"crypto/tls"
	"net"
	"strings"

	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/router"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/internal/loader"
	"github.com/jxo-me/netx/x/internal/plugin"
	xrouter "github.com/jxo-me/netx/x/router"
	routerplugin "github.com/jxo-me/netx/x/router/plugin"
)

func ParseRouter(cfg *config.RouterConfig) router.IRouter {
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
			return routerplugin.NewHTTPPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TLSConfigOption(tlsCfg),
				plugin.TimeoutOption(cfg.Plugin.Timeout),
			)
		default:
			return routerplugin.NewGRPCPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TokenOption(cfg.Plugin.Token),
				plugin.TLSConfigOption(tlsCfg),
			)
		}
	}

	var routes []*router.Route
	for _, route := range cfg.Routes {
		_, ipNet, _ := net.ParseCIDR(route.Net)
		if ipNet == nil {
			continue
		}
		gw := net.ParseIP(route.Gateway)
		if gw == nil {
			continue
		}

		routes = append(routes, &router.Route{
			Net:     ipNet,
			Gateway: gw,
		})
	}
	opts := []xrouter.Option{
		xrouter.RoutesOption(routes),
		xrouter.ReloadPeriodOption(cfg.Reload),
		xrouter.LoggerOption(logger.Default().WithFields(map[string]any{
			"kind":   "router",
			"router": cfg.Name,
		})),
	}
	if cfg.File != nil && cfg.File.Path != "" {
		opts = append(opts, xrouter.FileLoaderOption(loader.FileLoader(cfg.File.Path)))
	}
	if cfg.Redis != nil && cfg.Redis.Addr != "" {
		switch cfg.Redis.Type {
		case "list": // rediss list
			opts = append(opts, xrouter.RedisLoaderOption(loader.RedisListLoader(
				cfg.Redis.Addr,
				loader.DBRedisLoaderOption(cfg.Redis.DB),
				loader.PasswordRedisLoaderOption(cfg.Redis.Password),
				loader.KeyRedisLoaderOption(cfg.Redis.Key),
			)))
		case "set": // redis set
			opts = append(opts, xrouter.RedisLoaderOption(loader.RedisSetLoader(
				cfg.Redis.Addr,
				loader.DBRedisLoaderOption(cfg.Redis.DB),
				loader.PasswordRedisLoaderOption(cfg.Redis.Password),
				loader.KeyRedisLoaderOption(cfg.Redis.Key),
			)))
		default: // redis hash
			opts = append(opts, xrouter.RedisLoaderOption(loader.RedisHashLoader(
				cfg.Redis.Addr,
				loader.DBRedisLoaderOption(cfg.Redis.DB),
				loader.PasswordRedisLoaderOption(cfg.Redis.Password),
				loader.KeyRedisLoaderOption(cfg.Redis.Key),
			)))
		}
	}
	if cfg.HTTP != nil && cfg.HTTP.URL != "" {
		opts = append(opts, xrouter.HTTPLoaderOption(loader.HTTPLoader(
			cfg.HTTP.URL,
			loader.TimeoutHTTPLoaderOption(cfg.HTTP.Timeout),
		)))
	}
	return xrouter.NewRouter(opts...)
}
