package dtls

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/pion/dtls/v2"
	admission "github.com/jxo-me/netx/sdk/core/admission/wrapper"
	xnet "github.com/jxo-me/netx/sdk/internal/net"
	"github.com/jxo-me/netx/sdk/internal/net/proxyproto"
	xdtls "github.com/jxo-me/netx/sdk/internal/util/dtls"
	climiter "github.com/jxo-me/netx/sdk/core/limiter/conn/wrapper"
	limiter "github.com/jxo-me/netx/sdk/core/limiter/traffic/wrapper"
	"github.com/jxo-me/netx/sdk/core/listener"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
	metrics "github.com/jxo-me/netx/sdk/core/metrics/wrapper"
)

type dtlsListener struct {
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
	return &dtlsListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *dtlsListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	network := "udp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "udp4"
	}
	laddr, err := net.ResolveUDPAddr(network, l.options.Addr)
	if err != nil {
		return
	}

	tlsCfg := l.options.TLSConfig
	if tlsCfg == nil {
		tlsCfg = &tls.Config{}
	}
	config := dtls.Config{
		Certificates:         tlsCfg.Certificates,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		// Create timeout context for accepted connection.
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(context.Background(), 30*time.Second)
		},
		ClientCAs:      tlsCfg.ClientCAs,
		ClientAuth:     dtls.ClientAuthType(tlsCfg.ClientAuth),
		FlightInterval: l.md.flightInterval,
		MTU:            l.md.mtu,
	}

	ln, err := dtls.Listen(network, laddr, &config)
	if err != nil {
		return
	}
	ln = proxyproto.WrapListener(l.options.ProxyProtocol, ln, 10*time.Second)
	ln = metrics.WrapListener(l.options.Service, ln)
	ln = admission.WrapListener(l.options.Admission, ln)
	ln = limiter.WrapListener(l.options.TrafficLimiter, ln)
	ln = climiter.WrapListener(l.options.ConnLimiter, ln)

	l.ln = ln

	return
}

func (l *dtlsListener) Accept() (conn net.Conn, err error) {
	c, err := l.ln.Accept()
	if err != nil {
		return
	}
	conn = xdtls.Conn(c, l.md.bufferSize)
	return
}

func (l *dtlsListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *dtlsListener) Close() error {
	return l.ln.Close()
}
