package tap

import (
	"context"
	"net"
	"time"

	"github.com/jxo-me/netx/core/limiter"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	mdata "github.com/jxo-me/netx/core/metadata"
	xnet "github.com/jxo-me/netx/x/internal/net"
	traffic_limiter "github.com/jxo-me/netx/x/limiter/traffic"
	limiter_wrapper "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	mdx "github.com/jxo-me/netx/x/metadata"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
	stats "github.com/jxo-me/netx/x/observer/stats/wrapper"
)

type tapListener struct {
	addr    net.Addr
	cqueue  chan net.Conn
	closed  chan struct{}
	logger  logger.ILogger
	md      metadata
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.IListener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &tapListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *tapListener) Init(md mdata.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "udp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "udp4"
	}
	l.addr, err = net.ResolveUDPAddr(network, l.options.Addr)
	if err != nil {
		return
	}

	l.cqueue = make(chan net.Conn, 1)
	l.closed = make(chan struct{})

	go l.listenLoop()

	return
}

func (l *tapListener) listenLoop() {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		err := func() error {
			ifce, name, ip, err := l.createTap()
			if err != nil {
				if ifce != nil {
					ifce.Close()
				}
				return err
			}

			itf, err := net.InterfaceByName(name)
			if err != nil {
				return err
			}

			addrs, _ := itf.Addrs()
			l.logger.Infof("name: %s, net: %s, mtu: %d, addrs: %s",
				itf.Name, ip, itf.MTU, addrs)

			var c net.Conn
			c = &conn{
				ifce:   ifce,
				laddr:  l.addr,
				raddr:  &net.IPAddr{IP: ip},
				cancel: cancel,
			}
			c = metrics.WrapConn(l.options.Service, c)
			c = stats.WrapConn(c, l.options.Stats)
			c = limiter_wrapper.WrapConn(
				c,
				l.options.TrafficLimiter,
				traffic_limiter.ServiceLimitKey,
				limiter.ScopeOption(limiter.ScopeService),
				limiter.ServiceOption(l.options.Service),
				limiter.NetworkOption(c.LocalAddr().Network()),
			)
			c = withMetadata(mdx.NewMetadata(map[string]any{
				"config": l.md.config,
			}), c)

			l.cqueue <- c

			return nil
		}()
		if err != nil {
			l.logger.Error(err)
			cancel()
		}

		select {
		case <-ctx.Done():
		case <-l.closed:
			return
		}

		time.Sleep(time.Second)
	}
}

func (l *tapListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.cqueue:
		return conn, nil
	case <-l.closed:
	}

	return nil, listener.ErrClosed
}

func (l *tapListener) Addr() net.Addr {
	return l.addr
}

func (l *tapListener) Close() error {
	select {
	case <-l.closed:
		return net.ErrClosed
	default:
		close(l.closed)
	}
	return nil
}
