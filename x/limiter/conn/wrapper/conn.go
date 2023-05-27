package wrapper

import (
	"errors"
	"net"
	"syscall"

	limiter "github.com/jxo-me/netx/core/limiter/conn"
	"github.com/jxo-me/netx/core/metadata"
)

var (
	errUnsupport = errors.New("unsupported operation")
)

// serverConn is a server side Conn with metrics supported.
type serverConn struct {
	net.Conn
	limiter limiter.ILimiter
}

func WrapConn(limiter limiter.ILimiter, c net.Conn) net.Conn {
	if limiter == nil {
		return c
	}
	return &serverConn{
		Conn:    c,
		limiter: limiter,
	}
}

func (c *serverConn) SyscallConn() (rc syscall.RawConn, err error) {
	if sc, ok := c.Conn.(syscall.Conn); ok {
		rc, err = sc.SyscallConn()
		return
	}
	err = errUnsupport
	return
}

func (c *serverConn) Close() error {
	c.limiter.Allow(-1)
	return c.Conn.Close()
}

func (c *serverConn) Metadata() metadata.IMetaData {
	if md, ok := c.Conn.(metadata.IMetaDatable); ok {
		return md.Metadata()
	}
	return nil
}
