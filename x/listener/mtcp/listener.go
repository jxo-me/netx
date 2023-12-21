package mtcp

import (
	"context"
	"net"
	"time"

	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	admission "github.com/jxo-me/netx/x/admission/wrapper"
	xnet "github.com/jxo-me/netx/x/internal/net"
	"github.com/jxo-me/netx/x/internal/net/proxyproto"
	"github.com/jxo-me/netx/x/internal/util/mux"
	climiter "github.com/jxo-me/netx/x/limiter/conn/wrapper"
	limiter "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
)

type mtcpListener struct {
	ln      net.Listener
	cqueue  chan net.Conn
	errChan chan error
	logger  logger.ILogger
	md      metadata
	options listener.Options
}

func NewListener(opts ...listener.Option) listener.IListener {
	options := listener.Options{}
	for _, opt := range opts {
		opt(&options)
	}
	return &mtcpListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *mtcpListener) Init(md md.IMetaData) (err error) {
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

	l.logger.Debugf("pp: %d", l.options.ProxyProtocol)

	ln = proxyproto.WrapListener(l.options.ProxyProtocol, ln, 10*time.Second)
	ln = metrics.WrapListener(l.options.Service, ln)
	ln = admission.WrapListener(l.options.Admission, ln)
	ln = limiter.WrapListener(l.options.TrafficLimiter, ln)
	ln = climiter.WrapListener(l.options.ConnLimiter, ln)
	l.ln = ln

	l.cqueue = make(chan net.Conn, l.md.backlog)
	l.errChan = make(chan error, 1)

	go l.listenLoop()

	return
}

func (l *mtcpListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *mtcpListener) Close() error {
	return l.ln.Close()
}

func (l *mtcpListener) Accept() (conn net.Conn, err error) {
	var ok bool
	select {
	case conn = <-l.cqueue:
	case err, ok = <-l.errChan:
		if !ok {
			err = listener.ErrClosed
		}
	}
	return
}

func (l *mtcpListener) listenLoop() {
	for {
		conn, err := l.ln.Accept()
		if err != nil {
			l.errChan <- err
			close(l.errChan)
			return
		}
		go l.mux(conn)
	}
}

func (l *mtcpListener) mux(conn net.Conn) {
	defer conn.Close()

	session, err := mux.ServerSession(conn, l.md.muxCfg)
	if err != nil {
		l.logger.Error(err)
		return
	}
	defer session.Close()

	for {
		stream, err := session.Accept()
		if err != nil {
			l.logger.Error("accept stream: ", err)
			return
		}

		select {
		case l.cqueue <- stream:
		default:
			stream.Close()
			l.logger.Warnf("connection queue is full, client %s discarded", stream.RemoteAddr())
		}
	}
}