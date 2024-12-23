package tcp

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
	limiter_util "github.com/jxo-me/netx/x/internal/util/limiter"
	climiter "github.com/jxo-me/netx/x/limiter/conn/wrapper"
	limiter_wrapper "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
	stats "github.com/jxo-me/netx/x/observer/stats/wrapper"
)

type redirectListener struct {
	ln      net.Listener
	logger  logger.ILogger
	md      metadata
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.IListener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &redirectListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *redirectListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	lc := net.ListenConfig{}
	if l.md.tproxy {
		lc.Control = l.control
	}
	if l.md.mptcp {
		lc.SetMultipathTCP(true)
		l.logger.Debugf("mptcp enabled: %v", lc.MultipathTCP())
	}
	ln, err := lc.Listen(context.Background(), network, l.options.Addr)
	if err != nil {
		return err
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
	l.ln = ln
	return
}

func (l *redirectListener) Accept() (conn net.Conn, err error) {
	conn, err = l.ln.Accept()
	if err != nil {
		return
	}

	conn = limiter_wrapper.WrapConn(
		conn,
		limiter_util.NewCachedTrafficLimiter(l.options.TrafficLimiter, l.md.limiterRefreshInterval, 60*time.Second),
		conn.RemoteAddr().String(),
		limiter.ScopeOption(limiter.ScopeConn),
		limiter.ServiceOption(l.options.Service),
		limiter.NetworkOption(conn.LocalAddr().Network()),
		limiter.SrcOption(conn.RemoteAddr().String()),
	)

	return
}

func (l *redirectListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *redirectListener) Close() error {
	return l.ln.Close()
}
