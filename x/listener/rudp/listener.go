package rudp

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

type rudpListener struct {
	laddr   net.Addr
	ln      net.Listener
	closed  chan struct{}
	logger  logger.ILogger
	md      metadata
	options listener.Options
	mu      sync.Mutex
}

func NewListener(opts ...listener.Option) listener.IListener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &rudpListener{
		closed:  make(chan struct{}),
		logger:  options.Logger,
		options: options,
	}
}

func (l *rudpListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "udp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "udp4"
	}
	if laddr, _ := net.ResolveUDPAddr(network, l.options.Addr); laddr != nil {
		l.laddr = laddr
	}
	if l.laddr == nil {
		l.laddr = &bindAddr{addr: l.options.Addr}
	}

	return
}

func (l *rudpListener) Accept() (conn net.Conn, err error) {
	select {
	case <-l.closed:
		return nil, net.ErrClosed
	default:
	}

	ln := l.getListener()
	if ln == nil {
		ln, err = l.options.Router.Bind(
			context.Background(), "udp", l.laddr.String(),
			chain.BacklogBindOption(l.md.backlog),
			chain.UDPConnTTLBindOption(l.md.ttl),
			chain.UDPDataBufferSizeBindOption(l.md.readBufferSize),
			chain.UDPDataQueueSizeBindOption(l.md.readQueueSize),
		)
		if err != nil {
			return nil, listener.NewAcceptError(err)
		}

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

	conn, err = l.ln.Accept()
	if err != nil {
		l.ln.Close()
		l.setListener(nil)
		return nil, listener.NewAcceptError(err)
	}

	if pc, ok := conn.(net.PacketConn); ok {
		uc := metrics.WrapUDPConn(l.options.Service, pc)
		uc = stats.WrapUDPConn(uc, l.options.Stats)
		uc = admission.WrapUDPConn(l.options.Admission, uc)
		conn = limiter_wrapper.WrapUDPConn(
			uc,
			limiter_util.NewCachedTrafficLimiter(l.options.TrafficLimiter, l.md.limiterRefreshInterval, 60*time.Second),
			"",
			limiter.ScopeOption(limiter.ScopeConn),
			limiter.ServiceOption(l.options.Service),
			limiter.NetworkOption(conn.LocalAddr().Network()),
		)
	}

	return
}

func (l *rudpListener) Addr() net.Addr {
	return l.laddr
}

func (l *rudpListener) Close() error {
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

func (l *rudpListener) setListener(ln net.Listener) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.ln = ln
}

func (l *rudpListener) getListener() net.Listener {
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
