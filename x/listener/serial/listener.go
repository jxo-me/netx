package serial

import (
	"context"
	"net"
	"time"

	"github.com/jxo-me/netx/core/limiter"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	serial "github.com/jxo-me/netx/x/internal/util/serial"
	traffic_limiter "github.com/jxo-me/netx/x/limiter/traffic"
	limiter_wrapper "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
	stats "github.com/jxo-me/netx/x/observer/stats/wrapper"
)

type serialListener struct {
	cqueue  chan net.Conn
	closed  chan struct{}
	addr    net.Addr
	logger  logger.ILogger
	md      metadata
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.IListener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.Addr == "" {
		options.Addr = serial.DefaultPort
	}

	return &serialListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *serialListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	l.addr = &serial.Addr{Port: l.options.Addr}

	l.cqueue = make(chan net.Conn)
	l.closed = make(chan struct{})

	go l.listenLoop()

	return
}

func (l *serialListener) Accept() (conn net.Conn, err error) {
	select {
	case conn := <-l.cqueue:
		return conn, nil
	case <-l.closed:
	}

	return nil, listener.ErrClosed
}

func (l *serialListener) Addr() net.Addr {
	return l.addr
}

func (l *serialListener) Close() error {
	select {
	case <-l.closed:
		return net.ErrClosed
	default:
		close(l.closed)
	}
	return nil
}

func (l *serialListener) listenLoop() {
	for {
		ctx, cancel := context.WithCancel(context.Background())
		err := func() error {
			cfg := serial.ParseConfigFromAddr(l.options.Addr)
			cfg.ReadTimeout = l.md.timeout
			port, err := serial.OpenPort(cfg)
			if err != nil {
				return err
			}

			conn := serial.NewConn(port, l.addr, cancel)
			conn = metrics.WrapConn(l.options.Service, conn)
			conn = stats.WrapConn(conn, l.options.Stats)
			conn = limiter_wrapper.WrapConn(
				conn,
				l.options.TrafficLimiter,
				traffic_limiter.ServiceLimitKey,
				limiter.ScopeOption(limiter.ScopeService),
				limiter.ServiceOption(l.options.Service),
				limiter.NetworkOption(conn.LocalAddr().Network()),
			)

			l.cqueue <- conn

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
