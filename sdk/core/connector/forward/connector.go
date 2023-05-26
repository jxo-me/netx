package forward

import (
	"context"
	"net"

	"github.com/jxo-me/netx/sdk/core/connector"
	md "github.com/jxo-me/netx/sdk/core/metadata"
)

type forwardConnector struct {
	options connector.Options
}

func NewConnector(opts ...connector.Option) connector.IConnector {
	options := connector.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &forwardConnector{
		options: options,
	}
}

func (c *forwardConnector) Init(md md.IMetaData) (err error) {
	return nil
}

func (c *forwardConnector) Connect(ctx context.Context, conn net.Conn, network, address string, opts ...connector.ConnectOption) (net.Conn, error) {
	log := c.options.Logger.WithFields(map[string]any{
		"remote":  conn.RemoteAddr().String(),
		"local":   conn.LocalAddr().String(),
		"network": network,
		"address": address,
	})
	log.Debugf("connect %s/%s", address, network)

	return conn, nil
}
