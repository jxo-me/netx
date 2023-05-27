package ftcp

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	"github.com/xtaci/tcpraw"
)

type ftcpDialer struct {
	md     metadata
	logger logger.ILogger
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := &dialer.Options{}
	for _, opt := range opts {
		opt(options)
	}

	return &ftcpDialer{
		logger: options.Logger,
	}
}

func (d *ftcpDialer) Init(md md.IMetaData) (err error) {
	return d.parseMetadata(md)
}

func (d *ftcpDialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (conn net.Conn, err error) {
	raddr, er := net.ResolveTCPAddr("tcp", addr)
	if er != nil {
		return nil, er
	}
	c, err := tcpraw.Dial("tcp", addr)
	if err != nil {
		return
	}
	return &fakeTCPConn{
		raddr: raddr,
		pc:    c,
	}, nil
}
