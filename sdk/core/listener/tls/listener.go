package tls

import (
	"crypto/tls"
	"net"
	"time"

	admission "github.com/jxo-me/netx/sdk/core/admission/wrapper"
	xnet "github.com/jxo-me/netx/sdk/internal/net"
	"github.com/jxo-me/netx/sdk/internal/net/proxyproto"
	climiter "github.com/jxo-me/netx/sdk/core/limiter/conn/wrapper"
	limiter "github.com/jxo-me/netx/sdk/core/limiter/traffic/wrapper"
	"github.com/jxo-me/netx/sdk/core/listener"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
	metrics "github.com/jxo-me/netx/sdk/core/metrics/wrapper"
)

type tlsListener struct {
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
	return &tlsListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *tlsListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	ln, err := net.Listen(network, l.options.Addr)
	if err != nil {
		return
	}
	ln = proxyproto.WrapListener(l.options.ProxyProtocol, ln, 10*time.Second)
	ln = metrics.WrapListener(l.options.Service, ln)
	ln = admission.WrapListener(l.options.Admission, ln)
	ln = limiter.WrapListener(l.options.TrafficLimiter, ln)
	ln = climiter.WrapListener(l.options.ConnLimiter, ln)

	l.ln = tls.NewListener(ln, l.options.TLSConfig)

	return
}

func (l *tlsListener) Accept() (conn net.Conn, err error) {
	return l.ln.Accept()
}

func (l *tlsListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *tlsListener) Close() error {
	return l.ln.Close()
}
