package remote

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/hop"
	mdata "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/observer/stats"
	"github.com/jxo-me/netx/core/recorder"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	xnet "github.com/jxo-me/netx/x/internal/net"
	"github.com/jxo-me/netx/x/internal/net/proxyproto"
	"github.com/jxo-me/netx/x/internal/util/forwarder"
	"github.com/jxo-me/netx/x/internal/util/sniffing"
	tls_util "github.com/jxo-me/netx/x/internal/util/tls"
	rate_limiter "github.com/jxo-me/netx/x/limiter/rate"
	mdutil "github.com/jxo-me/netx/x/metadata/util"
	xstats "github.com/jxo-me/netx/x/observer/stats"
	stats_wrapper "github.com/jxo-me/netx/x/observer/stats/wrapper"
	xrecorder "github.com/jxo-me/netx/x/recorder"
)

type forwardHandler struct {
	hop      hop.IHop
	md       metadata
	options  handler.Options
	recorder recorder.RecorderObject
	certPool tls_util.CertPool
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &forwardHandler{
		options: options,
	}
}

func (h *forwardHandler) Init(md mdata.IMetaData) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	for _, ro := range h.options.Recorders {
		if ro.Record == xrecorder.RecorderServiceHandler {
			h.recorder = ro
			break
		}
	}

	if h.md.certificate != nil && h.md.privateKey != nil {
		h.certPool = tls_util.NewMemoryCertPool()
	}

	return
}

// Forward implements handler.Forwarder.
func (h *forwardHandler) Forward(hop hop.IHop) {
	h.hop = hop
}

func (h *forwardHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) (err error) {
	defer conn.Close()

	start := time.Now()

	ro := &xrecorder.HandlerRecorderObject{
		Service:    h.options.Service,
		RemoteAddr: conn.RemoteAddr().String(),
		LocalAddr:  conn.LocalAddr().String(),
		Network:    "tcp",
		Time:       start,
		SID:        string(ctxvalue.SidFromContext(ctx)),
	}

	ro.ClientIP = conn.RemoteAddr().String()
	if clientAddr := ctxvalue.ClientAddrFromContext(ctx); clientAddr != "" {
		ro.ClientIP = string(clientAddr)
	} else {
		ctx = ctxvalue.ContextWithClientAddr(ctx, ctxvalue.ClientAddr(conn.RemoteAddr().String()))
	}

	if h, _, _ := net.SplitHostPort(ro.ClientIP); h != "" {
		ro.ClientIP = h
	}

	log := h.options.Logger.WithFields(map[string]any{
		"remote": conn.RemoteAddr().String(),
		"local":  conn.LocalAddr().String(),
		"sid":    ro.SID,
		"client": ro.ClientIP,
	})
	log.Infof("%s <> %s", conn.RemoteAddr(), conn.LocalAddr())

	network := "tcp"
	if _, ok := conn.(net.PacketConn); ok {
		network = "udp"
	}
	ro.Network = network

	pStats := xstats.Stats{}
	conn = stats_wrapper.WrapConn(conn, &pStats)

	defer func() {
		if err != nil {
			ro.Err = err.Error()
		}
		ro.InputBytes = pStats.Get(stats.KindInputBytes)
		ro.OutputBytes = pStats.Get(stats.KindOutputBytes)
		ro.Duration = time.Since(start)
		if err := ro.Record(ctx, h.recorder.Recorder); err != nil {
			log.Errorf("record: %v", err)
		}

		log.WithFields(map[string]any{
			"duration":    time.Since(start),
			"inputBytes":  ro.InputBytes,
			"outputBytes": ro.OutputBytes,
		}).Infof("%s >< %s", conn.RemoteAddr(), conn.LocalAddr())
	}()

	if !h.checkRateLimit(conn.RemoteAddr()) {
		return rate_limiter.ErrRateLimit
	}

	var host string
	if md, ok := conn.(mdata.IMetaDatable); ok {
		if v := mdutil.GetString(md.Metadata(), "host"); v != "" {
			host = v
		}
	}

	var proto string
	if network == "tcp" && h.md.sniffing {
		if h.md.sniffingTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(h.md.sniffingTimeout))
		}

		br := bufio.NewReader(conn)
		proto, _ = sniffing.Sniff(ctx, br)
		ro.Proto = proto

		if h.md.sniffingTimeout > 0 {
			conn.SetReadDeadline(time.Time{})
		}

		dial := func(ctx context.Context, network, address string) (net.Conn, error) {
			var buf bytes.Buffer
			cc, err := h.options.Router.Dial(ctxvalue.ContextWithBuffer(ctx, &buf), "tcp", address)
			ro.Route = buf.String()
			return cc, err
		}
		sniffer := &forwarder.Sniffer{
			Websocket:           h.md.sniffingWebsocket,
			WebsocketSampleRate: h.md.sniffingWebsocketSampleRate,
			Recorder:            h.recorder.Recorder,
			RecorderOptions:     h.recorder.Options,
			Certificate:         h.md.certificate,
			PrivateKey:          h.md.privateKey,
			NegotiatedProtocol:  h.md.alpn,
			CertPool:            h.certPool,
			MitmBypass:          h.md.mitmBypass,
			ReadTimeout:         h.md.readTimeout,
		}

		conn = xnet.NewReadWriteConn(br, conn, conn)
		switch proto {
		case sniffing.ProtoHTTP:
			return sniffer.HandleHTTP(ctx, conn,
				forwarder.WithDial(dial),
				forwarder.WithHop(h.hop),
				forwarder.WithBypass(h.options.Bypass),
				forwarder.WithHTTPKeepalive(h.md.httpKeepalive),
				forwarder.WithRecorderObject(ro),
				forwarder.WithLog(log),
			)
		case sniffing.ProtoTLS:
			return sniffer.HandleTLS(ctx, conn,
				forwarder.WithDial(dial),
				forwarder.WithHop(h.hop),
				forwarder.WithBypass(h.options.Bypass),
				forwarder.WithRecorderObject(ro),
				forwarder.WithLog(log),
			)
		}
	}

	var target *chain.Node
	if host != "" {
		target = &chain.Node{
			Addr: host,
		}
	}
	if h.hop != nil {
		target = h.hop.Select(ctx,
			hop.ProtocolSelectOption(proto),
		)
	}
	if target == nil {
		err := errors.New("node not available")
		log.Error(err)
		return err
	}

	if opts := target.Options(); opts != nil {
		switch opts.Network {
		case "unix":
			network = opts.Network
		default:
		}
	}

	ro.Network = network
	ro.Host = target.Addr

	log = log.WithFields(map[string]any{
		"node": target.Name,
		"dst":  fmt.Sprintf("%s/%s", target.Addr, network),
	})

	log.Debugf("%s >> %s", conn.RemoteAddr(), target.Addr)

	var buf bytes.Buffer
	cc, err := h.options.Router.Dial(ctxvalue.ContextWithBuffer(ctx, &buf), network, target.Addr)
	ro.Route = buf.String()
	if err != nil {
		log.Error(err)
		// TODO: the router itself may be failed due to the failed node in the router,
		// the dead marker may be a wrong operation.
		if marker := target.Marker(); marker != nil {
			marker.Mark()
		}
		return err
	}
	defer cc.Close()
	if marker := target.Marker(); marker != nil {
		marker.Reset()
	}

	cc = proxyproto.WrapClientConn(h.md.proxyProtocol, conn.RemoteAddr(), convertAddr(conn.LocalAddr()), cc)

	t := time.Now()
	log.Infof("%s <-> %s", conn.RemoteAddr(), target.Addr)
	xnet.Transport(conn, cc)
	log.WithFields(map[string]any{
		"duration": time.Since(t),
	}).Infof("%s >-< %s", conn.RemoteAddr(), target.Addr)

	return nil
}

func (h *forwardHandler) checkRateLimit(addr net.Addr) bool {
	if h.options.RateLimiter == nil {
		return true
	}
	host, _, _ := net.SplitHostPort(addr.String())
	if limiter := h.options.RateLimiter.Limiter(host); limiter != nil {
		return limiter.Allow(1)
	}

	return true
}

func convertAddr(addr net.Addr) net.Addr {
	host, sp, _ := net.SplitHostPort(addr.String())
	ip := net.ParseIP(host)
	port, _ := strconv.Atoi(sp)

	if ip == nil || ip.Equal(net.IPv6zero) {
		ip = net.IPv4zero
	}

	switch addr.Network() {
	case "tcp", "tcp4", "tcp6":
		return &net.TCPAddr{
			IP:   ip,
			Port: port,
		}

	default:
		return &net.UDPAddr{
			IP:   ip,
			Port: port,
		}
	}
}
