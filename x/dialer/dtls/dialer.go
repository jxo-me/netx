package dtls

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/jxo-me/netx/core/dialer"
	"github.com/jxo-me/netx/core/logger"
	md "github.com/jxo-me/netx/core/metadata"
	xdtls "github.com/jxo-me/netx/x/internal/util/dtls"
	"github.com/pion/dtls/v2"
)

type dtlsDialer struct {
	md      metadata
	logger  logger.ILogger
	options dialer.Options
}

func NewDialer(opts ...dialer.Option) dialer.IDialer {
	options := dialer.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &dtlsDialer{
		logger:  options.Logger,
		options: options,
	}
}

func (d *dtlsDialer) Init(md md.IMetaData) (err error) {
	return d.parseMetadata(md)
}

func (d *dtlsDialer) Dial(ctx context.Context, addr string, opts ...dialer.DialOption) (net.Conn, error) {
	var options dialer.DialOptions
	for _, opt := range opts {
		opt(&options)
	}

	conn, err := options.Dialer.Dial(ctx, "udp", addr)
	if err != nil {
		return nil, err
	}

	tlsCfg := d.options.TLSConfig
	if tlsCfg == nil {
		tlsCfg = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	config := dtls.Config{
		Certificates:         tlsCfg.Certificates,
		InsecureSkipVerify:   tlsCfg.InsecureSkipVerify,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		ServerName:           tlsCfg.ServerName,
		RootCAs:              tlsCfg.RootCAs,
		FlightInterval:       d.md.flightInterval,
		MTU:                  d.md.mtu,
	}

	c, err := dtls.ClientWithContext(ctx, conn, &config)
	if err != nil {
		return nil, err
	}
	return xdtls.Conn(c, d.md.bufferSize), nil
}
