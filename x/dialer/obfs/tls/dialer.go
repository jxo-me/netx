package tls

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
)

type obfsTLSDialer struct {
	md     metadata
	logger logger.ILogger
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := &dialer.Options{}
	for _, opt := range opts {
		opt(options)
	}

	return &obfsTLSDialer{
		logger: options.Logger,
	}
}

func (d *obfsTLSDialer) Init(md md.IMetaData) (err error) {
	return d.parseMetadata(md)
}

func (d *obfsTLSDialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (net.Conn, error) {
	options := &dialer.DialOptions{}
	for _, opt := range opts {
		opt(options)
	}

	conn, err := options.NetDialer.Dial(ctx, "tcp", addr)
	if err != nil {
		d.logger.Error(err)
	}
	return conn, err
}

// Handshake implements dialer.Handshaker
func (d *obfsTLSDialer) Handshake(ctx context.Context, conn net.Conn, options ...dialer.HandshakeOption) (net.Conn, error) {
	opts := &dialer.HandshakeOptions{}
	for _, option := range options {
		option(opts)
	}

	host := d.md.host
	if host == "" {
		host = opts.Addr
	}

	return &obfsTLSConn{
		Conn:       conn,
		host:       host,
		handshaked: make(chan struct{}),
	}, nil
}
