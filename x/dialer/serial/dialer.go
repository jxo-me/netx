package serial

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	serial_util "github.com/jxo-me/netx/x/internal/util/serial"
	goserial "github.com/tarm/serial"
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

	cfg := serial_util.ParseConfigFromAddr(addr)
	port, err := goserial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}

	return serial_util.NewConn(port, &serial_util.Addr{Port: cfg.Name}, nil), nil
}
