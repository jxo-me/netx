package tun

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/songgao/water/waterutil"
	"github.com/jxo-me/netx/sdk/core/chain"
	"github.com/jxo-me/netx/sdk/core/handler"
	tun_util "github.com/jxo-me/netx/sdk/internal/util/tun"
	md "github.com/jxo-me/netx/sdk/core/metadata"
)

var (
	ErrTun        = errors.New("tun device error")
	ErrInvalidNet = errors.New("invalid net IP")
)

type tunHandler struct {
	hop     chain.IHop
	routes  sync.Map
	router  *chain.Router
	md      metadata
	options handler.Options
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &tunHandler{
		options: options,
	}
}

func (h *tunHandler) Init(md md.IMetaData) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	h.router = h.options.Router
	if h.router == nil {
		h.router = chain.NewRouter(chain.LoggerRouterOption(h.options.Logger))
	}

	return
}

// Forward implements handler.IForwarder.
func (h *tunHandler) Forward(hop chain.IHop) {
	h.hop = hop
}

func (h *tunHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
	defer conn.Close()

	log := h.options.Logger

	v, _ := conn.(md.IMetaDatable)
	if v == nil {
		err := errors.New("tun: wrong connection type")
		log.Error(err)
		return err
	}
	config := v.Metadata().Get("config").(*tun_util.Config)

	start := time.Now()
	log = log.WithFields(map[string]any{
		"remote": conn.RemoteAddr().String(),
		"local":  conn.LocalAddr().String(),
	})

	log.Infof("%s <> %s", conn.RemoteAddr(), conn.LocalAddr())
	defer func() {
		log.WithFields(map[string]any{
			"duration": time.Since(start),
		}).Infof("%s >< %s", conn.RemoteAddr(), conn.LocalAddr())
	}()

	var target *chain.Node
	if h.hop != nil {
		target = h.hop.Select(ctx)
	}
	if target != nil {
		log = log.WithFields(map[string]any{
			"dst": fmt.Sprintf("%s/%s", target.Addr, "udp"),
		})
		log.Debugf("%s >> %s", conn.RemoteAddr(), target.Addr)

		if err := h.handleClient(ctx, conn, target.Addr, config, log); err != nil {
			log.Error(err)
		}
		return nil
	}

	return h.handleServer(ctx, conn, config, log)
}

func (h *tunHandler) findRouteFor(dst net.IP, routes ...tun_util.Route) net.Addr {
	if v, ok := h.routes.Load(ipToTunRouteKey(dst)); ok {
		return v.(net.Addr)
	}
	for _, route := range routes {
		if route.Net.Contains(dst) && route.Gateway != nil {
			if v, ok := h.routes.Load(ipToTunRouteKey(route.Gateway)); ok {
				return v.(net.Addr)
			}
		}
	}
	return nil
}

var mIPProts = map[waterutil.IPProtocol]string{
	waterutil.HOPOPT:     "HOPOPT",
	waterutil.ICMP:       "ICMP",
	waterutil.IGMP:       "IGMP",
	waterutil.GGP:        "GGP",
	waterutil.TCP:        "TCP",
	waterutil.UDP:        "UDP",
	waterutil.IPv6_Route: "IPv6-Route",
	waterutil.IPv6_Frag:  "IPv6-Frag",
	waterutil.IPv6_ICMP:  "IPv6-ICMP",
}

func ipProtocol(p waterutil.IPProtocol) string {
	if v, ok := mIPProts[p]; ok {
		return v
	}
	return fmt.Sprintf("unknown(%d)", p)
}

type tunRouteKey [16]byte

func ipToTunRouteKey(ip net.IP) (key tunRouteKey) {
	copy(key[:], ip.To16())
	return
}
