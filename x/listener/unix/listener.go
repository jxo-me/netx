package unix

import (
	"net"

	"github.com/jxo-me/netx/core/listener"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	admission "github.com/jxo-me/netx/x/admission/wrapper"
	climiter "github.com/jxo-me/netx/x/limiter/conn/wrapper"
	limiter "github.com/jxo-me/netx/x/limiter/traffic/wrapper"
	metrics "github.com/jxo-me/netx/x/metrics/wrapper"
	stats "github.com/jxo-me/netx/x/stats/wrapper"
)

type unixListener struct {
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
	return &unixListener{
		logger:  options.Logger,
		options: options,
	}
}

func (l *unixListener) Init(md md.IMetaData) (err error) {
	if err = l.parseMetadata(md); err != nil {
		return
	}

	ln, err := net.Listen("unix", l.options.Addr)
	if err != nil {
		return
	}

	ln = metrics.WrapListener(l.options.Service, ln)
	ln = stats.WrapListener(ln, l.options.Stats)
	ln = admission.WrapListener(l.options.Admission, ln)
	ln = limiter.WrapListener(l.options.TrafficLimiter, ln)
	ln = climiter.WrapListener(l.options.ConnLimiter, ln)
	l.ln = ln

	return
}

func (l *unixListener) Accept() (conn net.Conn, err error) {
	return l.ln.Accept()
}

func (l *unixListener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *unixListener) Close() error {
	return l.ln.Close()
}
