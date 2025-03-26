package v5

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jxo-me/netx/core/limiter"
	"github.com/jxo-me/netx/core/logger"
	"github.com/jxo-me/netx/core/observer/stats"
	"github.com/jxo-me/netx/gosocks5"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
	xnet "github.com/jxo-me/netx/x/internal/net"
	"github.com/jxo-me/netx/x/internal/util/mux"
	traffic_wrapper "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	stats_wrapper "github.com/jxo-me/netx/x/observer/stats/wrapper"
)

func (h *socks5Handler) handleMuxBind(ctx context.Context, conn net.Conn, network, address string, log logger.ILogger) error {
	log = log.WithFields(map[string]any{
		"dst": fmt.Sprintf("%s/%s", address, network),
		"cmd": "mbind",
	})

	log.Debugf("%s >> %s", conn.RemoteAddr(), address)

	if !h.md.enableBind {
		reply := gosocks5.NewReply(gosocks5.NotAllowed, nil)
		log.Trace(reply)
		log.Error("socks5: BIND is disabled")
		return reply.Write(conn)
	}

	{
		clientID := ctxvalue.ClientIDFromContext(ctx)
		rw := traffic_wrapper.WrapReadWriter(
			h.limiter,
			conn,
			string(clientID),
			limiter.ServiceOption(h.options.Service),
			limiter.ScopeOption(limiter.ScopeClient),
			limiter.NetworkOption(network),
			limiter.AddrOption(address),
			limiter.ClientOption(string(clientID)),
			limiter.SrcOption(conn.RemoteAddr().String()),
		)
		if h.options.Observer != nil {
			pstats := h.stats.Stats(string(clientID))
			pstats.Add(stats.KindTotalConns, 1)
			pstats.Add(stats.KindCurrentConns, 1)
			defer pstats.Add(stats.KindCurrentConns, -1)
			rw = stats_wrapper.WrapReadWriter(rw, pstats)
		}

		conn = xnet.NewReadWriteConn(rw, rw, conn)
	}

	return h.muxBindLocal(ctx, conn, network, address, log)
}

func (h *socks5Handler) muxBindLocal(ctx context.Context, conn net.Conn, network, address string, log logger.ILogger) error {
	lc := xnet.ListenConfig{
		Netns: h.options.Netns,
	}
	ln, err := lc.Listen(ctx, network, address) // strict mode: if the port already in use, it will return error
	if err != nil {
		log.Error(err)
		reply := gosocks5.NewReply(gosocks5.Failure, nil)
		log.Trace(reply)
		if err := reply.Write(conn); err != nil {
			log.Error(err)
		}
		return err
	}

	socksAddr := gosocks5.Addr{}
	err = socksAddr.ParseFrom(ln.Addr().String())
	if err != nil {
		log.Warn(err)
	}

	// Issue: may not reachable when host has multi-interface
	socksAddr.Host, _, _ = net.SplitHostPort(conn.LocalAddr().String())
	socksAddr.Type = 0
	reply := gosocks5.NewReply(gosocks5.Succeeded, &socksAddr)
	log.Trace(reply)
	if err := reply.Write(conn); err != nil {
		log.Error(err)
		ln.Close()
		return err
	}

	log = log.WithFields(map[string]any{
		"bind": fmt.Sprintf("%s/%s", ln.Addr(), ln.Addr().Network()),
	})

	log.Debugf("bind on %s OK", ln.Addr())

	return h.serveMuxBind(ctx, conn, ln, log)
}

func (h *socks5Handler) serveMuxBind(ctx context.Context, conn net.Conn, ln net.Listener, log logger.ILogger) error {
	// Upgrade connection to multiplex stream.
	session, err := mux.ClientSession(conn, h.md.muxCfg)
	if err != nil {
		log.Error(err)
		return err
	}
	defer session.Close()

	go func() {
		defer ln.Close()
		for {
			conn, err := session.Accept()
			if err != nil {
				log.Error(err)
				return
			}
			conn.Close() // we do not handle incoming connections.
		}
	}()

	for {
		rc, err := ln.Accept()
		if err != nil {
			log.Error(err)
			return err
		}
		log.Debugf("peer %s accepted", rc.RemoteAddr())

		go func(c net.Conn) {
			defer c.Close()

			log = log.WithFields(map[string]any{
				"local":  rc.LocalAddr().String(),
				"remote": rc.RemoteAddr().String(),
			})
			sc, err := session.GetConn()
			if err != nil {
				log.Error(err)
				return
			}
			defer sc.Close()

			// incompatible with GOST v2.x
			if !h.md.compatibilityMode {
				addr := gosocks5.Addr{}
				addr.ParseFrom(c.RemoteAddr().String())
				reply := gosocks5.NewReply(gosocks5.Succeeded, &addr)
				log.Trace(reply)
				if err := reply.Write(sc); err != nil {
					log.Error(err)
					return
				}
			}

			t := time.Now()
			log.Debugf("%s <-> %s", c.LocalAddr(), c.RemoteAddr())
			xnet.Transport(sc, c)
			log.WithFields(map[string]any{"duration": time.Since(t)}).
				Debugf("%s >-< %s", c.LocalAddr(), c.RemoteAddr())
		}(rc)
	}
}
