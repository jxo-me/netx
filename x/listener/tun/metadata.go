package tun

import (
	"github.com/jxo-me/netx/x/app"
	"net"
	"strings"
	"time"

	"github.com/jxo-me/netx/core/logger"
	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/router"
	tun_util "github.com/jxo-me/netx/x/internal/util/tun"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
	xrouter "github.com/jxo-me/netx/x/router"
)

const (
	defaultMTU            = 1350
	defaultReadBufferSize = 4096
)

type metadata struct {
	config                 *tun_util.Config
	readBufferSize         int
	limiterRefreshInterval time.Duration
}

func (l *tunListener) parseMetadata(md mdata.IMetaData) (err error) {
	const (
		name    = "name"
		netKey  = "net"
		peer    = "peer"
		mtu     = "mtu"
		route   = "route"
		routes  = "routes"
		gateway = "gw"
	)

	l.md.readBufferSize = mdutil.GetInt(md, "tun.rbuf", "rbuf", "readBufferSize")
	if l.md.readBufferSize <= 0 {
		l.md.readBufferSize = defaultReadBufferSize
	}

	config := &tun_util.Config{
		Name:   mdutil.GetString(md, name),
		Peer:   mdutil.GetString(md, peer),
		MTU:    mdutil.GetInt(md, mtu),
		Router: app.Runtime.RouterRegistry().Get(mdutil.GetString(md, "router")),
	}
	if config.MTU <= 0 {
		config.MTU = defaultMTU
	}
	if gw := mdutil.GetString(md, gateway); gw != "" {
		config.Gateway = net.ParseIP(gw)
	}

	for _, s := range strings.Split(mdutil.GetString(md, netKey), ",") {
		if s = strings.TrimSpace(s); s == "" {
			continue
		}
		ip, ipNet, err := net.ParseCIDR(s)
		if err != nil {
			continue
		}
		config.Net = append(config.Net, net.IPNet{
			IP:   ip,
			Mask: ipNet.Mask,
		})
	}

	for _, s := range strings.Split(mdutil.GetString(md, route), ",") {
		_, ipNet, _ := net.ParseCIDR(strings.TrimSpace(s))
		if ipNet == nil {
			continue
		}

		l.routes = append(l.routes, &router.Route{
			Net:     ipNet,
			Gateway: config.Gateway,
		})
	}

	for _, s := range mdutil.GetStrings(md, routes) {
		ss := strings.SplitN(s, " ", 2)
		if len(ss) == 2 {
			var route router.Route
			_, ipNet, _ := net.ParseCIDR(strings.TrimSpace(ss[0]))
			if ipNet == nil {
				continue
			}
			route.Net = ipNet
			gw := net.ParseIP(ss[1])
			if gw == nil {
				gw = config.Gateway
			}

			l.routes = append(l.routes, &router.Route{
				Net:     ipNet,
				Gateway: gw,
			})
		}
	}

	if config.Router == nil && len(l.routes) > 0 {
		config.Router = xrouter.NewRouter(
			xrouter.RoutesOption(l.routes),
			xrouter.LoggerOption(logger.Default().WithFields(map[string]any{
				"kind":   "router",
				"router": "@internal",
			})),
		)
	}

	l.md.config = config

	l.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	if l.md.limiterRefreshInterval == 0 {
		l.md.limiterRefreshInterval = 30 * time.Second
	}
	if l.md.limiterRefreshInterval < time.Second {
		l.md.limiterRefreshInterval = time.Second
	}

	return
}
