package recorder

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/x/config"
	"github.com/jxo-me/netx/x/internal/plugin"
	xrecorder "github.com/jxo-me/netx/x/recorder"
	recorderplugin "github.com/jxo-me/netx/x/recorder/plugin"
)

func ParseRecorder(cfg *config.RecorderConfig) (r recorder.IRecorder) {
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
			return recorderplugin.NewHTTPPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TLSConfigOption(tlsCfg),
				plugin.TimeoutOption(cfg.Plugin.Timeout),
			)
		default:
			return recorderplugin.NewGRPCPlugin(
				cfg.Name, cfg.Plugin.Addr,
				plugin.TokenOption(cfg.Plugin.Token),
				plugin.TLSConfigOption(tlsCfg),
			)
		}
	}

	if cfg.File != nil && cfg.File.Path != "" {
		return xrecorder.FileRecorder(cfg.File.Path,
			xrecorder.SepRecorderOption(cfg.File.Sep),
		)
	}

	if cfg.TCP != nil && cfg.TCP.Addr != "" {
		return xrecorder.TCPRecorder(cfg.TCP.Addr, xrecorder.TimeoutTCPRecorderOption(cfg.TCP.Timeout))
	}

	if cfg.HTTP != nil && cfg.HTTP.URL != "" {
		h := http.Header{}
		for k, v := range cfg.HTTP.Header {
			h.Add(k, v)
		}
		return xrecorder.HTTPRecorder(cfg.HTTP.URL,
			xrecorder.TimeoutHTTPRecorderOption(cfg.HTTP.Timeout),
			xrecorder.HeaderHTTPRecorderOption(h),
		)
	}

	if cfg.Redis != nil &&
		cfg.Redis.Addr != "" &&
		cfg.Redis.Key != "" {
		switch cfg.Redis.Type {
		case "list": // redis list
			return xrecorder.RedisListRecorder(cfg.Redis.Addr,
				xrecorder.DBRedisRecorderOption(cfg.Redis.DB),
				xrecorder.KeyRedisRecorderOption(cfg.Redis.Key),
				xrecorder.UsernameRedisRecorderOption(cfg.Redis.Username),
				xrecorder.PasswordRedisRecorderOption(cfg.Redis.Password),
			)
		case "sset": // sorted set
			return xrecorder.RedisSortedSetRecorder(cfg.Redis.Addr,
				xrecorder.DBRedisRecorderOption(cfg.Redis.DB),
				xrecorder.KeyRedisRecorderOption(cfg.Redis.Key),
				xrecorder.UsernameRedisRecorderOption(cfg.Redis.Username),
				xrecorder.PasswordRedisRecorderOption(cfg.Redis.Password),
			)
		default: // redis set
			return xrecorder.RedisSetRecorder(cfg.Redis.Addr,
				xrecorder.DBRedisRecorderOption(cfg.Redis.DB),
				xrecorder.KeyRedisRecorderOption(cfg.Redis.Key),
				xrecorder.UsernameRedisRecorderOption(cfg.Redis.Username),
				xrecorder.PasswordRedisRecorderOption(cfg.Redis.Password),
			)
		}
	}

	return
}
