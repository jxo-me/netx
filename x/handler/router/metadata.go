package router

import (
	"github.com/jxo-me/netx/x/app"
	"math"
	"time"

	"github.com/jxo-me/netx/core/ingress"
	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/router"
	"github.com/jxo-me/netx/core/sd"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
)

const (
	MaxMessageSize         = math.MaxUint16
	defaultTTL             = 15 * time.Second
	defaultCacheExpiration = time.Second
)

type metadata struct {
	readTimeout time.Duration

	entryPoint        string
	ingress           ingress.IIngress
	sd                sd.ISD
	sdCacheExpiration time.Duration
	sdRenewInterval   time.Duration

	router                router.IRouter
	routerCacheEnabled    bool
	routerCacheExpiration time.Duration

	observerPeriod       time.Duration
	observerResetTraffic bool

	limiterRefreshInterval time.Duration
	limiterCleanupInterval time.Duration
}

func (h *routerHandler) parseMetadata(md mdata.IMetaData) (err error) {
	h.md.readTimeout = mdutil.GetDuration(md, "readTimeout")

	h.md.entryPoint = mdutil.GetString(md, "entrypoint")
	h.md.ingress = app.Runtime.IngressRegistry().Get(mdutil.GetString(md, "ingress"))

	h.md.sd = app.Runtime.SDRegistry().Get(mdutil.GetString(md, "sd"))
	h.md.sdCacheExpiration = mdutil.GetDuration(md, "sd.cache.expiration")
	if h.md.sdCacheExpiration <= 0 {
		h.md.sdCacheExpiration = defaultCacheExpiration
	}
	h.md.sdRenewInterval = mdutil.GetDuration(md, "sd.renewInterval")
	if h.md.sdRenewInterval < time.Second {
		h.md.sdRenewInterval = defaultTTL
	}

	h.md.router = app.Runtime.RouterRegistry().Get(mdutil.GetString(md, "router"))
	h.md.routerCacheEnabled = mdutil.GetBool(md, "router.cache")
	h.md.routerCacheExpiration = mdutil.GetDuration(md, "router.cache.expiration")
	if h.md.routerCacheExpiration <= 0 {
		h.md.routerCacheExpiration = defaultCacheExpiration
	}

	h.md.observerPeriod = mdutil.GetDuration(md, "observePeriod", "observer.period", "observer.observePeriod")
	if h.md.observerPeriod == 0 {
		h.md.observerPeriod = 5 * time.Second
	}
	if h.md.observerPeriod < time.Second {
		h.md.observerPeriod = time.Second
	}
	h.md.observerResetTraffic = mdutil.GetBool(md, "observer.resetTraffic")

	h.md.limiterRefreshInterval = mdutil.GetDuration(md, "limiter.refreshInterval")
	h.md.limiterCleanupInterval = mdutil.GetDuration(md, "limiter.cleanupInterval")

	return
}
