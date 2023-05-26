package rtcp

import (
	"context"
	"net"

	admission "github.com/jxo-me/netx/sdk/core/admission/wrapper"
	"github.com/jxo-me/netx/sdk/core/chain"
	xnet "github.com/jxo-me/netx/sdk/core/internal/net"
	climiter "github.com/jxo-me/netx/sdk/core/limiter/conn/wrapper"
	limiter "github.com/jxo-me/netx/sdk/core/limiter/traffic/wrapper"
	"github.com/jxo-me/netx/sdk/core/listener"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
	metrics "github.com/jxo-me/netx/sdk/core/metrics/wrapper"
)

type rtcpListener struct {
	laddr   net.Addr
	ln      net.Listener
	router  *chain.Router
	logger  logger.ILogger
	closed  chan struct{}
	options listener.Options
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
