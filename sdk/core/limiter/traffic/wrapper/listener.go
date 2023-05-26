package wrapper

import (
	"net"

	limiter "github.com/jxo-me/netx/sdk/core/limiter/traffic"
)

type listener struct {
	net.Listener
	limiter limiter.ITrafficLimiter
}

func WrapListener(limiter limiter.ITrafficLimiter, ln net.Listener) net.Listener {
	if limiter == nil {
		return ln
	}

	return &listener{
		limiter:  limiter,
		Listener: ln,
	}
}

func (ln *listener) Accept() (net.Conn, error) {
	c, err := ln.Listener.Accept()
	if err != nil {
		return nil, err
	}

	return WrapConn(ln.limiter, c), nil
}
