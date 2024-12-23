package direct

import (
	"context"
	"net"

	"github.com/jxo-me/netx/core/connector"
	md "github.com/jxo-me/netx/core/metadata"
	ctxvalue "github.com/jxo-me/netx/x/ctx"
)

type directConnector struct {
	md      metadata
	options connector.Options
}

func NewConnector(opts ...connector.Option) connector.IConnector {
	options := connector.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &directConnector{
		options: options,
	}
}

func (c *directConnector) Init(md md.IMetaData) (err error) {
	return c.parseMetadata(md)
}

func (c *directConnector) Connect(ctx context.Context, _ net.Conn, network, address string, opts ...connector.ConnectOption) (net.Conn, error) {
	var cOpts connector.ConnectOptions
	for _, opt := range opts {
		opt(&cOpts)
	}

	if c.md.action == "reject" {
		return &conn{}, nil
	}

	conn, err := cOpts.Dialer.Dial(ctx, network, address)
	if err != nil {
		return nil, err
	}

	var localAddr, remoteAddr string
	if addr := conn.LocalAddr(); addr != nil {
		localAddr = addr.String()
	}
	if addr := conn.RemoteAddr(); addr != nil {
		remoteAddr = addr.String()
	}

	log := c.options.Logger.WithFields(map[string]any{
		"remote":  remoteAddr,
		"local":   localAddr,
		"network": network,
		"address": address,
		"sid":     string(ctxvalue.SidFromContext(ctx)),
	})
	log.Debugf("connect %s/%s", address, network)

	return conn, nil
}
