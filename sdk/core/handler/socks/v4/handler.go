package v4

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/jxo-me/netx/sdk/core/chain"
	"github.com/jxo-me/netx/sdk/core/handler"
	netpkg "github.com/jxo-me/netx/sdk/core/internal/net"
	sx "github.com/jxo-me/netx/sdk/core/internal/util/selector"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
	"github.com/jxo-me/netx/sdk/gosocks4"
)

var (
	ErrUnknownCmd    = errors.New("socks4: unknown command")
	ErrUnimplemented = errors.New("socks4: unimplemented")
)

type socks4Handler struct {
	router  *chain.Router
	md      metadata
	options handler.Options
}

func NewHandler(opts ...handler.Option) handler.IHandler {
	options := handler.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &socks4Handler{
		options: options,
	}
}

func (h *socks4Handler) Init(md md.IMetaData) (err error) {
	if err := h.parseMetadata(md); err != nil {
		return err
	}

	h.router = h.options.Router
	if h.router == nil {
		h.router = chain.NewRouter(chain.LoggerRouterOption(h.options.Logger))
	}

	return nil
}

func (h *socks4Handler) Handle(ctx context.Context, conn net.Conn, opts ...handler.HandleOption) error {
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

	req, err := gosocks4.ReadRequest(conn)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Trace(req)

	conn.SetReadDeadline(time.Time{})

	if h.options.Auther != nil &&
		!h.options.Auther.Authenticate(ctx, string(req.Userid), "") {
		resp := gosocks4.NewReply(gosocks4.RejectedUserid, nil)
		log.Trace(resp)
		return resp.Write(conn)
	}

	switch req.Cmd {
	case gosocks4.CmdConnect:
		return h.handleConnect(ctx, conn, req, log)
	case gosocks4.CmdBind:
		return h.handleBind(ctx, conn, req)
	default:
		err = ErrUnknownCmd
		log.Error(err)
		return err
	}
}

func (h *socks4Handler) handleConnect(ctx context.Context, conn net.Conn, req *gosocks4.Request, log logger.ILogger) error {
	addr := req.Addr.String()

	log = log.WithFields(map[string]any{
		"dst": addr,
	})
	log.Debugf("%s >> %s", conn.RemoteAddr(), addr)

	if h.options.Bypass != nil && h.options.Bypass.Contains(ctx, addr) {
		resp := gosocks4.NewReply(gosocks4.Rejected, nil)
		log.Trace(resp)
		log.Debug("bypass: ", addr)
		return resp.Write(conn)
	}

	switch h.md.hash {
	case "host":
		ctx = sx.ContextWithHash(ctx, &sx.Hash{Source: addr})
	}

	cc, err := h.router.Dial(ctx, "tcp", addr)
	if err != nil {
		resp := gosocks4.NewReply(gosocks4.Failed, nil)
		log.Trace(resp)
		resp.Write(conn)
		return err
	}

	defer cc.Close()

	resp := gosocks4.NewReply(gosocks4.Granted, nil)
	log.Trace(resp)
	if err := resp.Write(conn); err != nil {
		log.Error(err)
		return err
	}

	t := time.Now()
	log.Infof("%s <-> %s", conn.RemoteAddr(), addr)
	netpkg.Transport(conn, cc)
	log.WithFields(map[string]any{
		"duration": time.Since(t),
	}).Infof("%s >-< %s", conn.RemoteAddr(), addr)

	return nil
}

func (h *socks4Handler) handleBind(ctx context.Context, conn net.Conn, req *gosocks4.Request) error {
	// TODO: bind
	return ErrUnimplemented
}

func (h *socks4Handler) checkRateLimit(addr net.Addr) bool {
	if h.options.RateLimiter == nil {
		return true
	}
	host, _, _ := net.SplitHostPort(addr.String())
	if limiter := h.options.RateLimiter.Limiter(host); limiter != nil {
		return limiter.Allow(1)
	}

	return true
}
