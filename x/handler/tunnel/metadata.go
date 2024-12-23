package tunnel

import (
	"github.com/jxo-me/netx/x/app"
	"strings"
	"time"

	"github.com/jxo-me/netx/core/ingress"
	"github.com/jxo-me/netx/core/logger"
	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/sd"
	"github.com/jxo-me/netx/relay"
	xingress "github.com/jxo-me/netx/x/ingress"
	"github.com/jxo-me/netx/x/internal/util/mux"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	defaultTTL = 15 * time.Second
)

type metadata struct {
	readTimeout time.Duration

	entryPoint                  string
	entryPointID                relay.TunnelID
	entryPointProxyProtocol     int
	entryPointKeepalive         bool
	entryPointCompression       bool
	entryPointReadTimeout       time.Duration
	sniffingWebsocket           bool
	sniffingWebsocketSampleRate float64

	directTunnel           bool
	tunnelTTL              time.Duration
	ingress                ingress.IIngress
	sd                     sd.ISD
	muxCfg                 *mux.Config
	observePeriod          time.Duration
	limiterRefreshInterval time.Duration
}

func (h *tunnelHandler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.readTimeout = mdutil.GetDuration(md, "readTimeout")

	h.md.entryPoint = mdutil.GetString(md, "entrypoint")
	h.md.entryPointID = parseTunnelID(mdutil.GetString(md, "entrypoint.id"))
	h.md.entryPointProxyProtocol = mdutil.GetInt(md, "entrypoint.ProxyProtocol")

	h.md.entryPointKeepalive = mdutil.GetBool(md, "entrypoint.keepalive")
	h.md.entryPointCompression = mdutil.GetBool(md, "entrypoint.compression")

	h.md.entryPointReadTimeout = mdutil.GetDuration(md, "entrypoint.readTimeout")
	if h.md.entryPointReadTimeout <= 0 {
		h.md.entryPointReadTimeout = 15 * time.Second
	}

	h.md.sniffingWebsocket = mdutil.GetBool(md, "sniffing.websocket")
	h.md.sniffingWebsocketSampleRate = mdutil.GetFloat(md, "sniffing.websocket.sampleRate")

	h.md.tunnelTTL = mdutil.GetDuration(md, "tunnel.ttl")
	if h.md.tunnelTTL <= 0 {
		h.md.tunnelTTL = defaultTTL
	}
	h.md.directTunnel = mdutil.GetBool(md, "tunnel.direct")

	h.md.ingress = app.Runtime.IngressRegistry().Get(mdutil.GetString(md, "ingress"))
	if h.md.ingress == nil {
		var rules []*ingress.Rule
		for _, s := range strings.Split(mdutil.GetString(md, "tunnel"), ",") {
			ss := strings.SplitN(s, ":", 2)
			if len(ss) != 2 {
				continue
			}
			rules = append(rules, &ingress.Rule{
				Hostname: ss[0],
				Endpoint: ss[1],
			})
		}
		if len(rules) > 0 {
			h.md.ingress = xingress.NewIngress(
				xingress.RulesOption(rules),
				xingress.LoggerOption(logger.Default().WithFields(map[string]any{
					"kind":    "ingress",
					"ingress": "@internal",
				})),
			)
		}
	}
	h.md.sd = app.Runtime.SDRegistry().Get(mdutil.GetString(md, "sd"))

	h.md.muxCfg = &mux.Config{
		Version:           mdutil.GetInt(md, "mux.version"),
		KeepAliveInterval: mdutil.GetDuration(md, "mux.keepaliveInterval"),
		KeepAliveDisabled: mdutil.GetBool(md, "mux.keepaliveDisabled"),
		KeepAliveTimeout:  mdutil.GetDuration(md, "mux.keepaliveTimeout"),
		MaxFrameSize:      mdutil.GetInt(md, "mux.maxFrameSize"),
		MaxReceiveBuffer:  mdutil.GetInt(md, "mux.maxReceiveBuffer"),
		MaxStreamBuffer:   mdutil.GetInt(md, "mux.maxStreamBuffer"),
	}
	if h.md.muxCfg.Version == 0 {
		h.md.muxCfg.Version = 2
	}
	if h.md.muxCfg.MaxStreamBuffer == 0 {
		h.md.muxCfg.MaxStreamBuffer = 1048576
	}

	h.md.observePeriod = mdutil.GetDuration(md, "observePeriod", "observer.observePeriod")
	if h.md.observePeriod == 0 {
		h.md.observePeriod = 5 * time.Second
	}
	if h.md.observePeriod < time.Second {
		h.md.observePeriod = time.Second
	}

	h.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	if h.md.limiterRefreshInterval == 0 {
		h.md.limiterRefreshInterval = 30 * time.Second
	}
	if h.md.limiterRefreshInterval < time.Second {
		h.md.limiterRefreshInterval = time.Second
	}

	return
}
