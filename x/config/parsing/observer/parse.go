package observer

import (
	"crypto/tls"
	"strings"

	"github.com/jxo-me/netx/core/observer"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/internal/plugin"
	observerplugin "github.com/jxo-me/netx/x/observer/plugin"
)

func ParseObserver(cfg *config.ObserverConfig) observer.IObserver {
	if cfg == nil || cfg.Plugin == nil {
		return nil
	}

	var tlsCfg *tls.Config
	if cfg.Plugin.TLS != nil {
		tlsCfg = &tls.Config{
			ServerName:         cfg.Plugin.TLS.ServerName,
			InsecureSkipVerify: !cfg.Plugin.TLS.Secure,
		}
	}
	switch strings.ToLower(cfg.Plugin.Type) {
	case "http":
		return observerplugin.NewHTTPPlugin(
			cfg.Name, cfg.Plugin.Addr,
			plugin.TLSConfigOption(tlsCfg),
			plugin.TimeoutOption(cfg.Plugin.Timeout),
		)
	default:
		return observerplugin.NewGRPCPlugin(
			cfg.Name, cfg.Plugin.Addr,
			plugin.TokenOption(cfg.Plugin.Token),
			plugin.TLSConfigOption(tlsCfg),
		)
	}
}
