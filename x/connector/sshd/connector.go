package sshd

import (
	"context"
	"errors"
	"net"

	"github.com/jxo-me/netx/core/connector"
	md "github.com/jxo-me/netx/core/metadata"
	ssh_util "github.com/jxo-me/netx/x/internal/util/ssh"
)

type sshdConnector struct {
	options connector.Options
}

func NewConnector(opts ...connector.Option) connector.IConnector {
	options := connector.Options{}
	for _, opt := range opts {
		opt(&options)
	}

	return &sshdConnector{
		options: options,
	}
}

func (c *sshdConnector) Init(md md.IMetaData) (err error) {
	return nil
}

func (c *sshdConnector) Connect(ctx context.Context, conn net.Conn, network, address string, opts ...connector.ConnectOption) (net.Conn, error) {
	log := c.options.Logger.WithFields(map[string]any{
		"remote":  conn.RemoteAddr().String(),
		"local":   conn.LocalAddr().String(),
		"network": network,
		"address": address,
	})
	log.Debugf("connect %s/%s", address, network)

	cc, ok := conn.(*ssh_util.ClientConn)
	if !ok {
		return nil, errors.New("ssh: invalid connection")
	}

	conn, err := cc.Client().Dial(network, address)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return conn, nil
}

// Bind implements connector.Binder.
func (c *sshdConnector) Bind(ctx context.Context, conn net.Conn, network, address string, opts ...connector.BindOption) (net.Listener, error) {
	log := c.options.Logger.WithFields(map[string]any{
		"remote":  conn.RemoteAddr().String(),
		"local":   conn.LocalAddr().String(),
		"network": network,
		"address": address,
	})
	log.Debugf("bind on %s/%s", address, network)

	cc, ok := conn.(*ssh_util.ClientConn)
	if !ok {
		return nil, errors.New("ssh: invalid connection")
	}

	if host, port, _ := net.SplitHostPort(address); host == "" {
		address = net.JoinHostPort("0.0.0.0", port)
	}

	return cc.Client().Listen(network, address)
}
