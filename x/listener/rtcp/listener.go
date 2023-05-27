package rtcp

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/chain"
	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	admission "github.com/jxo-me/netx/x/admission/wrapper"
	xnet "github.com/jxo-me/netx/x/internal/net"
	climiter "github.com/jxo-me/netx/x/limiter/conn/wrapper"
	limiter "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
	"github.com/jxo-me/netx/x/registry"
)

func init() {
	registry.ListenerRegistry().Register("rtcp", NewListener)
}

type rtcpListener struct {
	laddr   net.Addr
	ln      net.Listener
	router  *chain.Router
	logger  logger.Logger
	closed  chan struct{}
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.Listener {
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

func (l *rtcpListener) Init(md md.Metadata) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	laddr, err := net.ResolveTCPAddr(network, l.options.Addr)
	if err != nil {
		return
	}

	l.laddr = laddr
	l.router = chain.NewRouter(
		chain.ChainRouterOption(l.options.Chain),
		chain.LoggerRouterOption(l.logger),
	)

	return
}

func (l *rtcpListener) Accept() (conn net.Conn, err error) {
	select {
	case <-l.closed:
		return nil, net.ErrClosed
	default:
	}

	if l.ln == nil {
		l.ln, err = l.router.Bind(
			context.Background(), "tcp", l.laddr.String(),
			chain.MuxBindOption(true),
		)
		if err != nil {
			return nil, listener.NewAcceptError(err)
		}
		l.ln = metrics.WrapListener(l.options.Service, l.ln)
		l.ln = admission.WrapListener(l.options.Admission, l.ln)
		l.ln = limiter.WrapListener(l.options.TrafficLimiter, l.ln)
		l.ln = climiter.WrapListener(l.options.ConnLimiter, l.ln)
	}
	conn, err = l.ln.Accept()
	if err != nil {
		l.ln.Close()
		l.ln = nil
		return nil, listener.NewAcceptError(err)
	}
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
		if l.ln != nil {
			l.ln.Close()
			// l.ln = nil
		}
	}

	return nil
}
