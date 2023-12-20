package admission

import (
	"crypto/tls"
	"strings"

	"github.com/jxo-me/netx/core/admission"
	"github.com/jxo-me/netx/core/logger"
	xadmission "github.com/jxo-me/netx/x/admission"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/internal/loader"
	"github.com/jxo-me/netx/x/internal/plugin"
	"github.com/jxo-me/netx/x/registry"
)

func ParseAdmission(cfg *config.AdmissionConfig) admission.IAdmission {
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
			return xadmission.NewHTTPPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TLSConfigOption(tlsCfg),
				plugin.TimeoutOption(cfg.Plugin.Timeout),
			)
		default:
			return xadmission.NewGRPCPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TokenOption(cfg.Plugin.Token),
				plugin.TLSConfigOption(tlsCfg),
			)
		}
	}

	opts := []xadmission.Option{
		xadmission.MatchersOption(cfg.Matchers),
		xadmission.WhitelistOption(cfg.Reverse || cfg.Whitelist),
		xadmission.ReloadPeriodOption(cfg.Reload),
		xadmission.LoggerOption(logger.Default().WithFields(map[string]any{
			"kind":      "admission",
			"admission": cfg.Name,
		})),
	}
	if cfg.File != nil && cfg.File.Path != "" {
		opts = append(opts, xadmission.FileLoaderOption(loader.FileLoader(cfg.File.Path)))
	}
	if cfg.Redis != nil && cfg.Redis.Addr != "" {
		opts = append(opts, xadmission.RedisLoaderOption(loader.RedisSetLoader(
			cfg.Redis.Addr,
			loader.DBRedisLoaderOption(cfg.Redis.DB),
			loader.PasswordRedisLoaderOption(cfg.Redis.Password),
			loader.KeyRedisLoaderOption(cfg.Redis.Key),
		)))
	}
	if cfg.HTTP != nil && cfg.HTTP.URL != "" {
		opts = append(opts, xadmission.HTTPLoaderOption(loader.HTTPLoader(
			cfg.HTTP.URL,
			loader.TimeoutHTTPLoaderOption(cfg.HTTP.Timeout),
		)))
	}

	return xadmission.NewAdmission(opts...)
}

func List(name string, names ...string) []admission.IAdmission {
	var admissions []admission.IAdmission
	if adm := registry.AdmissionRegistry().Get(name); adm != nil {
		admissions = append(admissions, adm)
	}
	for _, s := range names {
		if adm := registry.AdmissionRegistry().Get(s); adm != nil {
			admissions = append(admissions, adm)
		}
	}

	return admissions
}
