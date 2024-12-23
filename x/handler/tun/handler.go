package tun

import (
	"context"
	"errors"
	"fmt"
	tun_util "github.com/jxo-me/netx/x/internal/util/tun"
	"net"
	"sync"
	"time"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/hop"
	md "github.com/jxo-me/netx/core/metadata"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	"github.com/songgao/water/waterutil"
)

var (
	ErrTun        = errors.New("tun device error")
	ErrInvalidNet = errors.New("invalid net IP")
)

type tunHandler struct {
	hop     hop.IHop
	routes  sync.Map
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

	return
}

// Forward implements handler.Forwarder.
func (h *tunHandler) Forward(hop hop.IHop) {
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
		"sid":    ctxvalue.SidFromContext(ctx),
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
