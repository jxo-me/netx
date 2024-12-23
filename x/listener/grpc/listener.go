package grpc

import (
	"context"
	"net"
	"time"

	"github.com/jxo-me/netx/core/limiter"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	admission "github.com/jxo-me/netx/x/admission/wrapper"
	xnet "github.com/jxo-me/netx/x/internal/net"
	"github.com/jxo-me/netx/x/internal/net/proxyproto"
	pb "github.com/jxo-me/netx/x/internal/util/grpc/proto"
	limiter_util "github.com/jxo-me/netx/x/internal/util/limiter"
	climiter "github.com/jxo-me/netx/x/limiter/conn/wrapper"
	limiter_wrapper "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
	stats "github.com/jxo-me/netx/x/observer/stats/wrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type grpcListener struct {
	addr    net.Addr
	server  *grpc.Server
	cqueue  chan net.Conn
	errChan chan error
	md      metadata
	logger  logger.ILogger
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.IListener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &grpcListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *grpcListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	lc := net.ListenConfig{}
	if l.md.mptcp {
		lc.SetMultipathTCP(true)
		l.logger.Debugf("mptcp enabled: %v", lc.MultipathTCP())
	}
	ln, err := lc.Listen(context.Background(), network, l.options.Addr)
	if err != nil {
		return
	}
	ln = proxyproto.WrapListener(l.options.ProxyProtocol, ln, 10*time.Second)
	ln = metrics.WrapListener(l.options.Service, ln)
	ln = stats.WrapListener(ln, l.options.Stats)
	ln = admission.WrapListener(l.options.Admission, ln)
	ln = limiter_wrapper.WrapListener(
		l.options.Service,
		ln,
		limiter_util.NewCachedTrafficLimiter(l.options.TrafficLimiter, l.md.limiterRefreshInterval, 60*time.Second),
	)
	ln = climiter.WrapListener(l.options.ConnLimiter, ln)

	var opts []grpc.ServerOption
	if !l.md.insecure {
		opts = append(opts, grpc.Creds(credentials.NewTLS(l.options.TLSConfig)))
	}
	if l.md.keepalive {
		opts = append(opts,
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				MinTime:             l.md.keepaliveMinTime,
				PermitWithoutStream: l.md.keepalivePermitWithoutStream,
			}),
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionIdle: l.md.keepaliveMaxConnectionIdle,
				Time:              l.md.keepaliveTime,
				Timeout:           l.md.keepaliveTimeout,
			}),
		)
	}

	l.server = grpc.NewServer(opts...)
	l.addr = ln.Addr()
	l.cqueue = make(chan net.Conn, l.md.backlog)
	l.errChan = make(chan error, 1)

	pb.RegisterGostTunelServerX(l.server, &server{
		cqueue:    l.cqueue,
		localAddr: l.addr,
		logger:    l.options.Logger,
	}, l.md.path)

	go func() {
		err := l.server.Serve(ln)
		if err != nil {
			l.errChan <- err
		}
		close(l.errChan)
	}()

	return
}

func (l *grpcListener) Accept() (conn net.Conn, err error) {
	var ok bool
	select {
	case conn = <-l.cqueue:
		conn = limiter_wrapper.WrapConn(
			conn,
			limiter_util.NewCachedTrafficLimiter(l.options.TrafficLimiter, l.md.limiterRefreshInterval, 60*time.Second),
			conn.RemoteAddr().String(),
			limiter.ScopeOption(limiter.ScopeConn),
			limiter.ServiceOption(l.options.Service),
			limiter.NetworkOption(conn.LocalAddr().Network()),
			limiter.SrcOption(conn.RemoteAddr().String()),
		)
	case err, ok = <-l.errChan:
		if !ok {
			err = listener.ErrClosed
		}
	}
	return
}

func (l *grpcListener) Close() error {
	l.server.Stop()
	return nil
}

func (l *grpcListener) Addr() net.Addr {
	return l.addr
}
