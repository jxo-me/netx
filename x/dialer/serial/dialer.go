package serial

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	serial "github.com/jxo-me/netx/x/internal/util/serial"
)

type serialDialer struct {
	md     metadata
	logger logger.ILogger
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := &dialer.Options{}
	for _, opt := range opts {
		opt(options)
	}

	return &serialDialer{
		logger: options.Logger,
	}
}

func (d *serialDialer) Init(md md.IMetaData) (err error) {
	return d.parseMetadata(md)
}

func (d *serialDialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (net.Conn, error) {
	var options dialer.DialOptions
	for _, opt := range opts {
		opt(&options)
	}

	cfg := serial.ParseConfigFromAddr(addr)
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}

	return serial.NewConn(port, &serial.Addr{Port: cfg.Name}, nil), nil
}
