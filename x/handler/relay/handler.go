package relay

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/jxo-me/netx/core/handler"
	"github.com/jxo-me/netx/core/hop"
	"github.com/jxo-me/netx/core/limiter/traffic"
	md "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/core/observer/stats"
	"github.com/jxo-me/netx/core/recorder"
	"github.com/jxo-me/netx/relay"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	limiter_util "github.com/jxo-me/netx/x/internal/util/limiter"
	stats_util "github.com/jxo-me/netx/x/internal/util/stats"
	tls_util "github.com/jxo-me/netx/x/internal/util/tls"
	rate_limiter "github.com/jxo-me/netx/x/limiter/rate"
	stats_wrapper "github.com/jxo-me/netx/x/observer/stats/wrapper"
	xrecorder "github.com/jxo-me/netx/x/recorder"
)

var (
	ErrBadVersion   = errors.New("relay: bad version")
	ErrUnknownCmd   = errors.New("relay: unknown command")
	ErrUnauthorized = errors.New("relay: unauthorized")
)

type relayHandler struct {
	hop      hop.IHop
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

	return &relayHandler{
		options: options,
	}
}

func (h *relayHandler) Init(md md.IMetaData) (err error) {
	if err := h.parseMetadata(md); err != nil {
		return err
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

	return nil
}

// Forward implements handler.Forwarder.
func (h *relayHandler) Forward(hop hop.IHop) {
	h.hop = hop
}

func (h *relayHandler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) (err error) {
	defer conn.Close()

	start := time.Now()

	ro := &xrecorder.HandlerRecorderObject{
		Service:    h.options.Service,
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

	if h.md.readTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(h.md.readTimeout))
	}

	req := relay.Request{}
	if _, err := req.ReadFrom(conn); err != nil {
		return err
	}

	conn.SetReadDeadline(time.Time{})

	resp := relay.Response{
		Version: relay.Version1,
		Status:  relay.StatusOK,
	}

	if req.Version != relay.Version1 {
		resp.Status = relay.StatusBadRequest
		resp.WriteTo(conn)
		return ErrBadVersion
	}

	var user, pass string
	var address string
	var networkID relay.NetworkID
	for _, f := range req.Features {
		switch f.Type() {
		case relay.FeatureUserAuth:
			if feature, _ := f.(*relay.UserAuthFeature); feature != nil {
				user, pass = feature.Username, feature.Password
			}
		case relay.FeatureAddr:
			if feature, _ := f.(*relay.AddrFeature); feature != nil {
				address = net.JoinHostPort(feature.Host, strconv.Itoa(int(feature.Port)))
			}
		case relay.FeatureNetwork:
			if feature, _ := f.(*relay.NetworkFeature); feature != nil {
				networkID = feature.Network
			}
		}
	}

	if user != "" {
		ro.ClientID = user
		log = log.WithFields(map[string]any{"user": user})
	}

	if h.options.Auther != nil {
		clientID, ok := h.options.Auther.Authenticate(ctx, user, pass)
		if !ok {
			resp.Status = relay.StatusUnauthorized
			resp.WriteTo(conn)
			return ErrUnauthorized
		}
		ctx = ctxvalue.ContextWithClientID(ctx, ctxvalue.ClientID(clientID))
	}

	network := networkID.String()
	if (req.Cmd & relay.FUDP) == relay.FUDP {
		network = "udp"
	}
	ro.Network = network
	ro.Host = address

	if h.hop != nil {
		// forward mode
		return h.handleForward(ctx, conn, network, log)
	}

	switch req.Cmd & relay.CmdMask {
	case 0, relay.CmdConnect:
		return h.handleConnect(ctx, conn, network, address, ro, log)
	case relay.CmdBind:
		return h.handleBind(ctx, conn, network, address, log)
	default:
		resp.Status = relay.StatusBadRequest
		resp.WriteTo(conn)
		return ErrUnknownCmd
	}
}

// Close implements io.Closer interface.
func (h *relayHandler) Close() error {
	if h.cancel != nil {
		h.cancel()
	}
	return nil
}

func (h *relayHandler) checkRateLimit(addr net.Addr) bool {
	if h.options.RateLimiter == nil {
		return true
	}
	host, _, _ := net.SplitHostPort(addr.String())
	if limiter := h.options.RateLimiter.Limiter(host); limiter != nil {
		return limiter.Allow(1)
	}

	return true
}

func (h *relayHandler) observeStats(ctx context.Context) {
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
