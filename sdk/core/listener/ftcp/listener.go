package ftcp

import (
	"net"

	"github.com/xtaci/tcpraw"
	admission "github.com/jxo-me/netx/sdk/core/admission/wrapper"
	"github.com/jxo-me/netx/sdk/core/common/net/udp"
	xnet "github.com/jxo-me/netx/sdk/internal/net"
	limiter "github.com/jxo-me/netx/sdk/core/limiter/traffic/wrapper"
	"github.com/jxo-me/netx/sdk/core/listener"
	"github.com/jxo-me/netx/sdk/core/logger"
	md "github.com/jxo-me/netx/sdk/core/metadata"
	metrics "github.com/jxo-me/netx/sdk/core/metrics/wrapper"
)

type ftcpListener struct {
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
	return &ftcpListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *ftcpListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	var conn net.PacketConn
	network := "tcp"
	if xnet.IsIPv4(l.options.Addr) {
		network = "tcp4"
	}
	conn, err = tcpraw.Listen(network, l.options.Addr)
	if err != nil {
		return
	}
	conn = metrics.WrapPacketConn(l.options.Service, conn)
	conn = admission.WrapPacketConn(l.options.Admission, conn)
	conn = limiter.WrapPacketConn(l.options.TrafficLimiter, conn)

	l.ln = udp.NewListener(
		conn,
		&udp.ListenConfig{
			Backlog:        l.md.backlog,
			ReadQueueSize:  l.md.readQueueSize,
			ReadBufferSize: l.md.readBufferSize,
			TTL:            l.md.ttl,
			KeepAlive:      true,
			Logger:         l.logger,
		})
	return
}

func (l *ftcpListener) Accept() (conn net.Conn, err error) {
	return l.ln.Accept()
}

func (l *ftcpListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *ftcpListener) Close() error {
	return l.ln.Close()
}
