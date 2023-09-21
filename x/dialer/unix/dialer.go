package unix

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
)

type unixDialer struct {
	md     metadata
	logger logger.ILogger
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := &dialer.Options{}
	for _, opt := range opts {
		opt(options)
	}

	return &unixDialer{
		logger: options.Logger,
	}
}

func (d *unixDialer) Init(md md.IMetaData) (err error) {
	return d.parseMetadata(md)
}

func (d *unixDialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (net.Conn, error) {
	var options dialer.DialOptions
	for _, opt := range opts {
		opt(&options)
	}

	conn, err := (&net.Dialer{}).DialContext(ctx, "unix", addr)
	if err != nil {
		d.logger.Error(err)
	}
	return conn, err
}
