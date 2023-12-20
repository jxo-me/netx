package v5

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jxo-me/netx/core/limiter/traffic"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/gosocks5"
	ctxvalue "github.com/jxo-me/netx/x/internal/ctx"
	netpkg "github.com/jxo-me/netx/x/internal/net"
	"github.com/jxo-me/netx/x/limiter/traffic/wrapper"
)

func (h *socks5Handler) handleConnect(ctx context.Context, conn net.Conn, network, address string, log logger.ILogger) error {
	log = log.WithFields(map[string]any{
		"dst": fmt.Sprintf("%s/%s", address, network),
		"cmd": "connect",
	})
	log.Debugf("%s >> %s", conn.RemoteAddr(), address)

	if h.options.Bypass != nil && h.options.Bypass.Contains(ctx, network, address) {
		resp := gosocks5.NewReply(gosocks5.NotAllowed, nil)
		log.Trace(resp)
		log.Debug("bypass: ", address)
		return resp.Write(conn)
	}

	switch h.md.hash {
	case "host":
		ctx = ctxvalue.ContextWithHash(ctx, &ctxvalue.Hash{Source: address})
	}

	cc, err := h.router.Dial(ctx, network, address)
	if err != nil {
		resp := gosocks5.NewReply(gosocks5.NetUnreachable, nil)
		log.Trace(resp)
		resp.Write(conn)
		return err
	}

	defer cc.Close()

	resp := gosocks5.NewReply(gosocks5.Succeeded, nil)
	log.Trace(resp)
	if err := resp.Write(conn); err != nil {
		log.Error(err)
		return err
	}

	rw := wrapper.WrapReadWriter(h.options.Limiter, conn, conn.RemoteAddr().String(),
		traffic.NetworkOption(network),
		traffic.AddrOption(address),
		traffic.ClientOption(string(ctxvalue.ClientIDFromContext(ctx))),
		traffic.SrcOption(conn.RemoteAddr().String()),
	)

	t := time.Now()
	log.Infof("%s <-> %s", conn.RemoteAddr(), address)
	netpkg.Transport(rw, cc)
	log.WithFields(map[string]any{
		"duration": time.Since(t),
	}).Infof("%s >-< %s", conn.RemoteAddr(), address)

	return nil
}
