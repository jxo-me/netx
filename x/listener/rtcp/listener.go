package rtcp

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/limiter"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	admission "github.com/jxo-me/netx/x/admission/wrapper"
	xnet "github.com/jxo-me/netx/x/internal/net"
	limiter_util "github.com/jxo-me/netx/x/internal/util/limiter"
	climiter "github.com/jxo-me/netx/x/limiter/conn/wrapper"
	limiter_wrapper "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
	stats "github.com/jxo-me/netx/x/observer/stats/wrapper"
)

type rtcpListener struct {
	laddr   net.Addr
	ln      net.Listener
	logger  logger.ILogger
	closed  chan struct{}
	md      metadata
	options listener.Options
	mu      sync.Mutex
}

func NewListener(opts ...listener.Option) listener.IListener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &rtcpListener{
		closed:  make(chan struct{}),
		logger:  options.Logger,
		options: options,
	}
}

func (l *rtcpListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	if laddr, _ := net.ResolveTCPAddr(network, l.options.Addr); laddr != nil {
		l.laddr = laddr
	}
	if l.laddr == nil {
		l.laddr = &bindAddr{addr: l.options.Addr}
	}

	return
}

func (l *rtcpListener) Accept() (conn net.Conn, err error) {
	select {
	case <-l.closed:
		return nil, net.ErrClosed
	default:
	}

	ln := l.getListener()
	if ln == nil {
		ln, err = l.options.Router.Bind(
			context.Background(), "tcp", l.laddr.String(),
			chain.MuxBindOption(true),
		)
		if err != nil {
			return nil, listener.NewAcceptError(err)
		}
		ln = metrics.WrapListener(l.options.Service, ln)
		ln = stats.WrapListener(ln, l.options.Stats)
		ln = admission.WrapListener(l.options.Admission, ln)
		ln = limiter_wrapper.WrapListener(
			l.options.Service,
			ln,
			limiter_util.NewCachedTrafficLimiter(l.options.TrafficLimiter, l.md.limiterRefreshInterval, 60*time.Second),
		)
		ln = climiter.WrapListener(l.options.ConnLimiter, ln)
		l.setListener(ln)
	}

	select {
	case <-l.closed:
		ln.Close()
		return nil, net.ErrClosed
	default:
	}

	conn, err = ln.Accept()
	if err != nil {
		ln.Close()
		l.setListener(nil)
		return nil, listener.NewAcceptError(err)
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

func (l *rtcpListener) Addr() net.Addr {
	return l.laddr
}

func (l *rtcpListener) Close() error {
	select {
	case <-l.closed:
	default:
		close(l.closed)
		if ln := l.getListener(); ln != nil {
			ln.Close()
		}
	}

	return nil
}

func (l *rtcpListener) setListener(ln net.Listener) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ln = ln
}

func (l *rtcpListener) getListener() net.Listener {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.ln
}

type bindAddr struct {
	addr string
}

func (p *bindAddr) Network() string {
	return "tcp"
}

func (p *bindAddr) String() string {
	return p.addr
}
