package wrapper

import (
	"net"

	"github.com/jxo-me/netx/core/limiter"
	"github.com/jxo-me/netx/core/limiter/traffic"
)

type listener struct {
	net.Listener
	limiter traffic.ITrafficLimiter
	service string
}

func WrapListener(service string, ln net.Listener, limiter traffic.ITrafficLimiter) net.Listener {
	if limiter == nil {
		return ln
	}

	return &listener{
		Listener: ln,
		limiter:  limiter,
		service:  service,
	}
}

func (ln *listener) Accept() (net.Conn, error) {
	c, err := ln.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return WrapConn(c, ln.limiter, "",
		limiter.ScopeOption(limiter.ScopeService),
		limiter.ServiceOption(ln.service),
		limiter.NetworkOption(ln.Addr().Network()),
	), nil
}
