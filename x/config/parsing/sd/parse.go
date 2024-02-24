package sd

import (
	"crypto/tls"
	"strings"

	"github.com/jxo-me/netx/core/sd"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/internal/plugin"
	sdplugin "github.com/jxo-me/netx/x/sd/plugin"
)

func ParseSD(cfg *config.SDConfig) sd.ISD {
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
		return sdplugin.NewHTTPPlugin(
			cfg.Name, cfg.Plugin.Addr,
			plugin.TLSConfigOption(tlsCfg),
			plugin.TimeoutOption(cfg.Plugin.Timeout),
		)
	default:
		return sdplugin.NewGRPCPlugin(
			cfg.Name, cfg.Plugin.Addr,
			plugin.TokenOption(cfg.Plugin.Token),
			plugin.TLSConfigOption(tlsCfg),
		)
	}
}
