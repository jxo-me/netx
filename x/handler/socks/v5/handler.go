package v5

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/handler"
	md "github.com/jxo-me/netx/core/metadata"
	"github.com/jxo-me/netx/gosocks5"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	"github.com/jxo-me/netx/x/internal/util/socks"
	stats_util "github.com/jxo-me/netx/x/internal/util/stats"
)

var (
	ErrUnknownCmd = errors.New("socks5: unknown command")
)

type socks5Handler struct {
	selector gosocks5.Selector
	router   *chain.Router
	md       metadata
	options  handler.Options
	stats    *stats_util.HandlerStats
	cancel   context.CancelFunc
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &socks5Handler{
		options: options,
		stats:   stats_util.NewHandlerStats(options.Service),
	}
}

func (h *socks5Handler) Init(md md.IMetaData) (err error) {
	if err = h.parseMetadata(md); err != nil {
		return
	}

	h.router = h.options.Router
	if h.router == nil {
		h.router = chain.NewRouter(chain.LoggerRouterOption(h.options.Logger))
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
		go h.observeStats(ctx)
	}

	return
}

func (h *socks5Handler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
	defer conn.Close()

	start := time.Now()

	log := h.options.Logger.WithFields(map[string]any{
		"remote": conn.RemoteAddr().String(),
		"local":  conn.LocalAddr().String(),
	})

	log.Infof("%s <> %s", conn.RemoteAddr(), conn.LocalAddr())
	defer func() {
		log.WithFields(map[string]any{
			"duration": time.Since(start),
		}).Infof("%s >< %s", conn.RemoteAddr(), conn.LocalAddr())
	}()

	if !h.checkRateLimit(conn.RemoteAddr()) {
		return nil
	}

	if h.md.readTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(h.md.readTimeout))
	}

	sc := gosocks5.ServerConn(conn, h.selector)
	req, err := gosocks5.ReadRequest(sc)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Trace(req)

	if clientID := sc.ID(); clientID != "" {
		ctx = ctxvalue.ContextWithClientID(ctx, ctxvalue.ClientID(clientID))
	}

	conn = sc
	conn.SetReadDeadline(time.Time{})

	address := req.Addr.String()

	switch req.Cmd {
	case gosocks5.CmdConnect:
		return h.handleConnect(ctx, conn, "tcp", address, log)
	case gosocks5.CmdBind:
		return h.handleBind(ctx, conn, "tcp", address, log)
	case socks.CmdMuxBind:
		return h.handleMuxBind(ctx, conn, "tcp", address, log)
	case gosocks5.CmdUdp:
		return h.handleUDP(ctx, conn, log)
	case socks.CmdUDPTun:
		return h.handleUDPTun(ctx, conn, "udp", address, log)
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

	d := h.md.observePeriod
	if d < time.Millisecond {
		d = 5 * time.Second
	}
	ticker := time.NewTicker(d)
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
