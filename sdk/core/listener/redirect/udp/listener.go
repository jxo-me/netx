package udp

import (
	"net"

	admission "github.com/jxo-me/netx/sdk/core/admission/wrapper"
	limiter "github.com/jxo-me/netx/sdk/core/limiter/traffic/wrapper"
	"github.com/jxo-me/netx/sdk/core/listener"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
	metrics "github.com/jxo-me/netx/sdk/core/metrics/wrapper"
)

type redirectListener struct {
	ln      *net.UDPConn
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

	ln, err := l.listenUDP(l.options.Addr)
	if err != nil {
		return
	}

	l.ln = ln
	return
}

func (l *redirectListener) Accept() (conn net.Conn, err error) {
	conn, err = l.accept()
	if err != nil {
		return
	}
	conn = metrics.WrapConn(l.options.Service, conn)
	conn = admission.WrapConn(l.options.Admission, conn)
	conn = limiter.WrapConn(l.options.TrafficLimiter, conn)
	return
}

func (l *redirectListener) Addr() net.Addr {
	return l.ln.LocalAddr()
}

func (l *redirectListener) Close() error {
	return l.ln.Close()
}
