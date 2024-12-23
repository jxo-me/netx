package v5

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/limiter/traffic"
	md "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/observer/stats"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/gosocks5"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	limiter_util "github.com/jxo-me/netx/x/internal/util/limiter"
	"github.com/jxo-me/netx/x/internal/util/socks"
	stats_util "github.com/jxo-me/netx/x/internal/util/stats"
	tls_util "github.com/jxo-me/netx/x/internal/util/tls"
	rate_limiter "github.com/jxo-me/netx/x/limiter/rate"
	stats_wrapper "github.com/jxo-me/netx/x/observer/stats/wrapper"
	xrecorder "github.com/jxo-me/netx/x/recorder"
)

var (
	ErrUnknownCmd = errors.New("socks5: unknown command")
)

type socks5Handler struct {
	selector gosocks5.Selector
	md       metadata
	options  handler.Options
	stats    *stats_util.HandlerStats
	limiter  traffic.ITrafficLimiter
	cancel   context.CancelFunc
	recorder recorder.RecorderObject
	certPool tls_util.CertPool
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &socks5Handler{
		options: options,
	}
}

func (h *socks5Handler) Init(md md.IMetaData) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	h.selector = &serverSelector{
		Authenticator: h.options.Auther,
		TLSConfig:     h.options.TLSConfig,
		logger:        h.options.Logger,
		noTLS:         h.md.noTLS,
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel

	if h.options.Observer != nil {
		h.stats = stats_util.NewHandlerStats(h.options.Service)
		go h.observeStats(ctx)
	}

	if limiter := h.options.Limiter; limiter != nil {
		h.limiter = limiter_util.NewCachedTrafficLimiter(limiter, h.md.limiterRefreshInterval, 60*time.Second)
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

func (h *socks5Handler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) (err error) {
	defer conn.Close()

	start := time.Now()

	ro := &xrecorder.HandlerRecorderObject{
		Service:    h.options.Service,
		Network:    "tcp",
		RemoteAddr: conn.RemoteAddr().String(),
		LocalAddr:  conn.LocalAddr().String(),
		Time:       start,
		SID:        string(ctxvalue.SidFromContext(ctx)),
	}

	ro.ClientIP = conn.RemoteAddr().String()
	if clientAddr := ctxvalue.ClientAddrFromContext(ctx); clientAddr != "" {
		ro.ClientIP = string(clientAddr)
	}
	if h, _, _ := net.SplitHostPort(ro.ClientIP); h != "" {
		ro.ClientIP = h
	}

	log := h.options.Logger.WithFields(map[string]any{
		"remote": conn.RemoteAddr().String(),
		"local":  conn.LocalAddr().String(),
		"sid":    ctxvalue.SidFromContext(ctx),
		"client": ro.ClientIP,
	})
	log.Infof("%s <> %s", conn.RemoteAddr(), conn.LocalAddr())

	pStats := stats.Stats{}
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

	conn.SetReadDeadline(time.Now().Add(h.md.readTimeout))

	sc := gosocks5.ServerConn(conn, h.selector)
	req, err := gosocks5.ReadRequest(sc)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Trace(req)

	if clientID := sc.ID(); clientID != "" {
		ctx = ctxvalue.ContextWithClientID(ctx, ctxvalue.ClientID(clientID))
		log = log.WithFields(map[string]any{"user": clientID})
		ro.ClientID = clientID
	}

	conn = sc
	conn.SetReadDeadline(time.Time{})

	address := req.Addr.String()
	ro.Host = address

	switch req.Cmd {
	case gosocks5.CmdConnect:
		return h.handleConnect(ctx, conn, "tcp", address, ro, log)
	case gosocks5.CmdBind:
		return h.handleBind(ctx, conn, "tcp", address, log)
	case socks.CmdMuxBind:
		return h.handleMuxBind(ctx, conn, "tcp", address, log)
	case gosocks5.CmdUdp:
		ro.Network = "udp"
		return h.handleUDP(ctx, conn, ro, log)
	case socks.CmdUDPTun:
		ro.Network = "udp"
		return h.handleUDPTun(ctx, conn, "udp", address, ro, log)
	default:
		err = ErrUnknownCmd
		log.Error(err)
		resp := gosocks5.NewReply(gosocks5.CmdUnsupported, nil)
		log.Trace(resp)
		resp.Write(conn)
		return err
	}
}

func (h *socks5Handler) Close() error {
	if h.cancel != nil {
		h.cancel()
	}
	return nil
}

func (h *socks5Handler) checkRateLimit(addr net.Addr) bool {
	if h.options.RateLimiter == nil {
		return true
	}
	host, _, _ := net.SplitHostPort(addr.String())
	if limiter := h.options.RateLimiter.Limiter(host); limiter != nil {
		return limiter.Allow(1)
	}

	return true
}

func (h *socks5Handler) observeStats(ctx context.Context) {
	if h.options.Observer == nil {
		return
	}

	ticker := time.NewTicker(h.md.observePeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.options.Observer.Observe(ctx, h.stats.Events())
		case <-ctx.Done():
			return
		}
	}
}
