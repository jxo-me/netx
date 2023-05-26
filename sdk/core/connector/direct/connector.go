package forward

import (
	"context"
	"net"

	"github.com/jxo-me/netx/sdk/core/connector"
	md "github.com/jxo-me/netx/sdk/core/metadata"
)

type directConnector struct {
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
	return nil
}

func (c *directConnector) Connect(ctx context.Context, _ net.Conn, network, address string, opts ...connector.ConnectOption) (net.Conn, error) {
	var cOpts connector.ConnectOptions
	for _, opt := range opts {
		opt(&cOpts)
	}

	conn, err := cOpts.NetDialer.Dial(ctx, network, address)
	if err != nil {
		return nil, err
	}

	log := c.options.Logger.WithFields(map[string]any{
		"remote":  conn.RemoteAddr().String(),
		"local":   conn.LocalAddr().String(),
		"network": network,
		"address": address,
	})
	log.Debugf("connect %s/%s", address, network)

	return conn, nil
}
